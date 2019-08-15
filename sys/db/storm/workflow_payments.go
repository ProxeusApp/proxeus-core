package storm

import (
	"path/filepath"

	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/msgpack"
	"github.com/asdine/storm/q"

	"git.proxeus.com/core/central/sys/model"
)

type WorkflowPaymentsDBInterface interface {
	GetByTxHash(txHash string) (*model.WorkflowPaymentItem, error)
	GetByWorkflowId(workflowID string) (*model.WorkflowPaymentItem, error)
	GetByWorkflowIdAndFromEthAddress(workflowID, ethAddr string) (*model.WorkflowPaymentItem, error)
	Add(item *model.WorkflowPaymentItem) error
	Delete(txHash string) error
	Close() error
}

type WorkflowPaymentsDB struct {
	db *storm.DB
}

const workflowPaymentVersion = "sig_vers"
const workflowPaymentDBDir = "workflowpayments"
const workflowPaymentDB = "workflowpaymentsdb"

func NewWorkflowPaymentDB(dir string) (*WorkflowPaymentsDB, error) {
	var err error
	var msgpackDb *storm.DB
	baseDir := filepath.Join(dir, workflowPaymentDBDir)
	err = ensureDir(baseDir)
	if err != nil {
		return nil, err
	}
	msgpackDb, err = storm.Open(filepath.Join(baseDir, workflowPaymentDB), storm.Codec(msgpack.Codec))
	if err != nil {
		return nil, err
	}
	udb := &WorkflowPaymentsDB{db: msgpackDb}

	example := &model.WorkflowPaymentItem{}
	udb.db.Init(example)
	udb.db.ReIndex(example)

	err = udb.db.Set(workflowPaymentVersion, workflowPaymentVersion, example.GetVersion())
	if err != nil {
		return nil, err
	}

	return udb, nil
}

func (me *WorkflowPaymentsDB) GetByTxHash(txHash string) (*model.WorkflowPaymentItem, error) {
	tx, err := me.db.Begin(false)
	if err != nil {
		return nil, err
	}
	var item model.WorkflowPaymentItem
	defer tx.Rollback()
	err = tx.One("Hash", txHash, &item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (me *WorkflowPaymentsDB) GetByWorkflowId(workflowID string) (*model.WorkflowPaymentItem, error) {
	tx, err := me.db.Begin(false)
	if err != nil {
		return nil, err
	}
	var item model.WorkflowPaymentItem
	defer tx.Rollback()
	err = tx.One("WorkflowID", workflowID, &item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (me *WorkflowPaymentsDB) GetByWorkflowIdAndFromEthAddress(workflowID, ethAddr string) (*model.WorkflowPaymentItem, error) {
	var item model.WorkflowPaymentItem

	query := me.db.Select(q.Eq("WorkflowID", workflowID), q.Eq("From", ethAddr))

	err := query.First(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (me *WorkflowPaymentsDB) Add(item *model.WorkflowPaymentItem) error {
	tx, err := me.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	tx.Save(item)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (me *WorkflowPaymentsDB) Delete(txHash string) error {
	tx, err := me.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var item model.WorkflowPaymentItem
	err = tx.One("Hash", txHash, &item)
	if err != nil {
		return err
	}

	err = tx.DeleteStruct(&item)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (me *WorkflowPaymentsDB) Close() error {
	return me.db.Close()
}
