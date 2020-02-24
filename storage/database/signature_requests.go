package database

import (
	"path/filepath"
	"time"

	"github.com/asdine/storm/q"

	"github.com/ProxeusApp/proxeus-core/storage/database/db"
	"github.com/ProxeusApp/proxeus-core/sys/model"
)

type SignatureRequestsDB struct {
	db db.DB
}

const signatureVersion = "sig_vers"
const signatureDBDir = "signaturerequests"
const signatureDB = "signaturerequestsdb"

// NewSignatureDB return a handle to the signature request database
func NewSignatureDB(c DBConfig) (*SignatureRequestsDB, error) {
	baseDir := filepath.Join(c.Dir, signatureDBDir)
	db, err := db.OpenDatabase(c.Engine, c.URI, filepath.Join(baseDir, signatureDB))
	if err != nil {
		return nil, err
	}
	udb := &SignatureRequestsDB{db: db}

	example := &model.SignatureRequestItem{}
	err = udb.db.Init(example)
	if err != nil {
		return nil, err
	}
	err = udb.db.ReIndex(example)
	if err != nil {
		return nil, err
	}
	err = udb.db.Set(signatureVersion, signatureVersion, example.GetVersion())
	return udb, err
}

// GetBySignatory returns the list of signature requests for a specific signatory
func (me *SignatureRequestsDB) GetBySignatory(ethAddr string) (*[]model.SignatureRequestItem, error) {
	var items []model.SignatureRequestItem
	err := me.db.Select(
		q.Eq("Signatory", ethAddr),
	).OrderBy("RequestedAt", "RevokedAt", "RejectedAt").Reverse().Find(&items)
	return &items, err
}

// GetByID returns the signature request item by its id
func (me *SignatureRequestsDB) GetByID(docid string, docpath string) (*[]model.SignatureRequestItem, error) {
	var items []model.SignatureRequestItem
	err := me.db.Select(q.And(
		q.Eq("DocId", docid),
		q.Eq("DocPath", docpath),
	)).OrderBy("Requestor", "Signatory", "Revoked", "Rejected").Find(&items)
	return &items, err
}

// GetByHashAndSigner returns a list of signture requests for a specific file hash and signatory
func (me *SignatureRequestsDB) GetByHashAndSigner(hash string, signatory string) (*[]model.SignatureRequestItem, error) {
	var items []model.SignatureRequestItem
	err := me.db.Select(q.And(
		q.Eq("Hash", hash),
		q.Eq("Signatory", signatory),
	)).OrderBy("Revoked", "Rejected").Find(&items)
	return &items, err
}

// Add saves a signature request into the database
func (me *SignatureRequestsDB) Add(item *model.SignatureRequestItem) error {
	err := me.db.Save(item)
	return err
}

// SetRejected alters the status of a signature request to rejected
func (me *SignatureRequestsDB) SetRejected(docid string, docpath string, signatory string) error {
	var items []model.SignatureRequestItem
	err := me.db.Select(q.And(
		q.Eq("DocId", docid),
		q.Eq("DocPath", docpath),
		q.Eq("Signatory", signatory),
		q.Eq("Rejected", false),
		q.Eq("Revoked", false),
	)).Find(&items)
	if err != nil {
		return err
	}
	items[0].Rejected = true
	items[0].RejectedAt = time.Now()
	err = me.db.Save(&items[0])
	return err
}

// SetRevoked alters the status of a signature request to revoked
func (me *SignatureRequestsDB) SetRevoked(docid string, docpath string, signatory string) error {
	var items []model.SignatureRequestItem
	err := me.db.Select(q.And(
		q.Eq("DocId", docid),
		q.Eq("DocPath", docpath),
		q.Eq("Signatory", signatory),
		q.Eq("Rejected", false),
		q.Eq("Revoked", false),
	)).Find(&items)
	if err != nil {
		return err
	}
	items[0].Revoked = true
	items[0].RevokedAt = time.Now()
	err = me.db.Save(&items[0])
	return err
}

func (me *SignatureRequestsDB) Close() error {
	return me.db.Close()
}
