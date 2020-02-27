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

// NewWorkflowPaymentDB return a handle to the workflow payment database
func NewWorkflowPaymentDB(c DBConfig) (*WorkflowPaymentsDB, error) {
	baseDir := filepath.Join(c.Dir, WorkflowPaymentDBDir)
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

// All returns a list of all workflow payment items from the database
func (me *WorkflowPaymentsDB) All() ([]*model.WorkflowPaymentItem, error) {
	var items []*model.WorkflowPaymentItem

	err := me.db.All(&items)
	return items, err
}

// Get returns a specific Workflow payment item matching its id
func (me *WorkflowPaymentsDB) Get(paymentId string) (*model.WorkflowPaymentItem, error) {
	var item model.WorkflowPaymentItem

	query := me.db.Select(
		q.Eq("ID", paymentId),
	).OrderBy("CreatedAt").Reverse().Limit(1) //always match newest entry

	err := query.First(&item)
	return &item, err
}

// GetByTxHashAndStatusAndFromEthAddress returns a workflow payment item by matching the supplied filter parameters
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

// GetByWorkflowIdAndFromEthAddress returns a workflow payment item by matching the supplied filter parameters
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

// SetAbandonedToTimeoutBeforeTime updates the status of all payment items created before the specified time to status timeout
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

// Save add a workflow payment item to the database
func (me *WorkflowPaymentsDB) Save(item *model.WorkflowPaymentItem) error {
	if item.CreatedAt.IsZero() {
		item.CreatedAt = time.Now()
	}
	return me.db.Save(item)
}

// ConfirmPayment sets the status of a workflow payment item to confirmed by trying to find a matching transaction hash and searching for pending or created items matching the supplied criteria
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

// Redeem sets the status of a workflow payment item to redeemed for the item matching the supplied id and from address
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

// Cancel sets the status of a workflow payment item to cancelled for the item matching the supplied id and from address
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

// Delete sets the status of a workflow payment item to deleted by matching the supplied id
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

// Update sets the status and tx hash of created workflow items matching the supplied criteria to the supplied values
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
