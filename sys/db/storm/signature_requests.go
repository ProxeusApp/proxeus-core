package storm

import (
	"path/filepath"
	"time"

	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/msgpack"
	"github.com/asdine/storm/q"

	"github.com/ProxeusApp/proxeus-core/sys/model"
)

type SignatureRequestsDB struct {
	db *storm.DB
}

const signatureVersion = "sig_vers"
const signatureDBDir = "signaturerequests"
const signatureDB = "signaturerequestsdb"

func NewSignatureDB(dir string) (*SignatureRequestsDB, error) {
	var err error
	var msgpackDb *storm.DB
	baseDir := filepath.Join(dir, signatureDBDir)
	err = ensureDir(baseDir)
	if err != nil {
		return nil, err
	}
	msgpackDb, err = storm.Open(filepath.Join(baseDir, signatureDB), storm.Codec(msgpack.Codec))
	if err != nil {
		return nil, err
	}
	udb := &SignatureRequestsDB{db: msgpackDb}

	example := &model.SignatureRequestItem{}
	err = udb.db.Init(example)
	if err != nil {
		return nil, err
	}
	err = udb.db.ReIndex(example)
	if err != nil {
		return nil, err
	}

	var fVersion int
	verr := udb.db.Get(signatureVersion, signatureVersion, &fVersion)
	if verr == nil && fVersion != example.GetVersion() {
		//TODO: Upgrade items in DB to match new Model
	}
	err = udb.db.Set(signatureVersion, signatureVersion, example.GetVersion())
	return udb, err
}

func (me *SignatureRequestsDB) GetBySignatory(ethAddr string) (*[]model.SignatureRequestItem, error) {
	var items []model.SignatureRequestItem
	err := me.db.Select(
		q.Eq("Signatory", ethAddr),
	).OrderBy("RequestedAt", "RevokedAt", "RejectedAt").Reverse().Find(&items)
	return &items, err
}

func (me *SignatureRequestsDB) All() (*[]model.SignatureRequestItem, error) {

	var items []model.SignatureRequestItem
	err := me.db.All(&items)
	return &items, err
}

func (me *SignatureRequestsDB) GetByID(docid string, docpath string) (*[]model.SignatureRequestItem, error) {
	var items []model.SignatureRequestItem
	err := me.db.Select(q.And(
		q.Eq("DocId", docid),
		q.Eq("DocPath", docpath),
	)).OrderBy("Requestor", "Signatory", "Revoked", "Rejected").Find(&items)
	return &items, err
}

func (me *SignatureRequestsDB) GetByHashAndSigner(hash string, signatory string) (*[]model.SignatureRequestItem, error) {
	var items []model.SignatureRequestItem
	err := me.db.Select(q.And(
		q.Eq("Hash", hash),
		q.Eq("Signatory", signatory),
	)).OrderBy("Revoked", "Rejected").Find(&items)
	return &items, err
}

func (me *SignatureRequestsDB) Add(item *model.SignatureRequestItem) error {
	err := me.db.Save(item)
	return err
}

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

func (me *SignatureRequestsDB) List(auth model.Authorization, contains string, options map[string]interface{}) ([]*model.UserDataItem, error) {
	params := makeSimpleQuery(options)
	items := make([]*model.UserDataItem, 0)
	tx, err := me.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	matchers := defaultMatcher(auth, contains, params, true)

	err = tx.Select(matchers...).
		Limit(params.limit).
		Skip(params.index).
		OrderBy("Updated").
		Reverse().
		Find(&items)

	if err != nil {
		return nil, err
	}
	if !params.metaOnly {
		for _, item := range items {
			_ = tx.Get(usrdHeavyData, item.ID, &item.Data)
		}
	}
	return items, nil
}

func (me *SignatureRequestsDB) Close() error {
	return me.db.Close()
}
