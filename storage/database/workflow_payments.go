package database

import (
	"errors"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/asdine/storm/q"

	"github.com/ProxeusApp/proxeus-core/storage/database/db"
	"github.com/ProxeusApp/proxeus-core/sys/model"
)

type WorkflowPaymentsDB struct {
	db db.DB
}

const workflowPaymentVersion = "payment_vers"
const WorkflowPaymentDBDir = "workflowpayments"
const WorkflowPaymentDB = "workflowpaymentsdb"

func NewWorkflowPaymentDB(c DBConfig) (*WorkflowPaymentsDB, error) {
	var err error

	baseDir := filepath.Join(c.Dir, WorkflowPaymentDBDir)
	err = ensureDir(baseDir)
	if err != nil {
		return nil, err
	}
	db, err := db.OpenDatabase(c.Engine, c.URI, filepath.Join(baseDir, WorkflowPaymentDB))
	if err != nil {
		return nil, err
	}
	udb := &WorkflowPaymentsDB{db: db}

	example := &model.WorkflowPaymentItem{}
	err = udb.db.Init(example)
	if err != nil {
		return nil, err
	}
	err = udb.db.ReIndex(example)
	if err != nil {
		return nil, err
	}

	err = udb.db.Set(workflowPaymentVersion, workflowPaymentVersion, example.GetVersion())
	if err != nil {
		return nil, err
	}

	return udb, nil
}

func (me *WorkflowPaymentsDB) All() ([]*model.WorkflowPaymentItem, error) {
	var items []*model.WorkflowPaymentItem

	err := me.db.All(&items)
	return items, err
}

func (me *WorkflowPaymentsDB) Get(paymentId string) (*model.WorkflowPaymentItem, error) {
	var item model.WorkflowPaymentItem

	query := me.db.Select(
		q.Eq("ID", paymentId),
	).OrderBy("CreatedAt").Reverse().Limit(1) //always match newest entry

	err := query.First(&item)
	return &item, err
}
func (me *WorkflowPaymentsDB) GetByTxHashAndStatusAndFromEthAddress(txHash, status,
	fromEthAddr string) (*model.WorkflowPaymentItem, error) {

	var item model.WorkflowPaymentItem

	query := me.db.Select(
		q.Eq("TxHash", txHash),
		q.Eq("Status", status),
		q.Eq("From", fromEthAddr),
	).OrderBy("CreatedAt").Reverse()

	err := query.First(&item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (me *WorkflowPaymentsDB) GetByWorkflowIdAndFromEthAddress(workflowID, fromEthAddr string,
	statuses []string) (*model.WorkflowPaymentItem, error) {

	var (
		item  model.WorkflowPaymentItem
		query db.Query
	)

	if len(statuses) == 0 {
		query = me.db.Select(
			q.Eq("WorkflowID", workflowID),
			q.Eq("From", fromEthAddr),
		)
	} else {
		query = me.db.Select(
			q.Eq("WorkflowID", workflowID),
			q.Eq("From", fromEthAddr),
			q.In("Status", statuses),
		)
	}

	query.OrderBy("CreatedAt").Reverse()

	err := query.First(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (me *WorkflowPaymentsDB) SetAbandonedToTimeoutBeforeTime(beforeTime time.Time) error {
	query := me.db.Select(
		q.Or(
			q.Eq("Status", model.PaymentStatusCreated),
			q.Eq("Status", model.PaymentStatusPending),
		),
		q.Lt("CreatedAt", beforeTime),
	)

	return query.Each(new(model.WorkflowPaymentItem), func(record interface{}) error {
		u := record.(*model.WorkflowPaymentItem)
		u.Status = model.PaymentStatusTimeout
		return me.Save(u)
	})
}

func (me *WorkflowPaymentsDB) Save(item *model.WorkflowPaymentItem) error {
	if item.CreatedAt.IsZero() {
		item.CreatedAt = time.Now()
	}
	return me.db.Save(item)
}

func (me *WorkflowPaymentsDB) ConfirmPayment(txHash, from, to string, xes uint64) error {
	tx, err := me.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var item model.WorkflowPaymentItem

	// Initially try to get payment by TxHash
	query := tx.Select(
		q.Eq("TxHash", txHash),
		q.Eq("From", from),
		q.Eq("To", to),
		q.Eq("Xes", xes),
		q.In("Status", []string{model.PaymentStatusPending, model.PaymentStatusCreated}),
	).OrderBy("CreatedAt").Reverse().Limit(1) //always match newest entry

	err = query.First(&item)
	if err != nil {
		if !db.NotFound(err) {
			return err
		}

		// prioritize PaymentStatusPending over PaymentStatusCreated
		query := tx.Select(
			q.Eq("From", from),
			q.Eq("To", to),
			q.Eq("Xes", xes),
			q.Eq("Status", model.PaymentStatusPending),
		).OrderBy("CreatedAt").Reverse().Limit(1) //always match newest entry

		err = query.First(&item)
		if err != nil {
			if !db.NotFound(err) {
				return err
			}

			query = tx.Select(
				q.Eq("From", from),
				q.Eq("To", to),
				q.Eq("Xes", xes),
				q.Eq("Status", model.PaymentStatusCreated),
			).OrderBy("CreatedAt").Reverse().Limit(1) //always match newest entry

			err = query.First(&item)
			if err != nil {
				return err
			}
		}
	}

	item.Status = model.PaymentStatusConfirmed
	if item.TxHash == "" {
		item.TxHash = txHash
	}

	err = tx.Update(&item)
	if err != nil {
		log.Println("[ConfirmPayment] tx.Update err: ", err.Error())
		return err
	}

	return tx.Commit()
}

func (me *WorkflowPaymentsDB) Redeem(workflowId, from string) error {
	tx, err := me.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var item model.WorkflowPaymentItem

	query := tx.Select(
		q.Eq("WorkflowID", workflowId),
		q.Eq("From", from),
		q.Eq("Status", model.PaymentStatusConfirmed),
	).OrderBy("CreatedAt").Reverse().Limit(1) //always match newest entry

	err = query.First(&item)
	if err != nil {
		return err
	}

	item.Status = model.PaymentStatusRedeemed

	err = tx.Update(&item)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (me *WorkflowPaymentsDB) Cancel(paymentId, from string) error {
	tx, err := me.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var item model.WorkflowPaymentItem

	query := tx.Select(
		q.Eq("ID", paymentId),
		q.Eq("From", from),
		q.Eq("Status", model.PaymentStatusCreated),
	).OrderBy("CreatedAt").Reverse() //always match newest entry

	err = query.First(&item)
	if err != nil {
		return err
	}

	item.Status = model.PaymentStatusCancelled

	err = tx.Update(&item)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (me *WorkflowPaymentsDB) Delete(paymentId string) error {
	tx, err := me.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var item model.WorkflowPaymentItem

	err = tx.One("ID", paymentId, &item)
	if err != nil {
		return err
	}

	item.Status = model.PaymentStatusDeleted

	err = tx.Update(&item)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (me *WorkflowPaymentsDB) Remove(payment *model.WorkflowPaymentItem) error {
	return me.db.DeleteStruct(payment)
}

var errNothingToUpdate = errors.New("nothing to update")

func (me *WorkflowPaymentsDB) Update(paymentId, status, txHash, from string) error {
	tx, err := me.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var item model.WorkflowPaymentItem
	query := tx.Select(
		q.Eq("ID", paymentId),
		q.Eq("From", from),
		q.Eq("Status", model.PaymentStatusCreated),
	).OrderBy("CreatedAt").Reverse().Limit(1) //always match newest entry

	err = query.First(&item)
	if err != nil {
		return err
	}

	if strings.TrimSpace(status) == "" && strings.TrimSpace(txHash) == "" {
		return errNothingToUpdate
	}

	if strings.TrimSpace(status) != "" {
		item.Status = status
	}
	if strings.TrimSpace(txHash) != "" {
		item.TxHash = txHash
	}

	err = tx.Update(&item)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (me *WorkflowPaymentsDB) Close() error {
	return me.db.Close()
}
