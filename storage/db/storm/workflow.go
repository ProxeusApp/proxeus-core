package storm

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/storage/database"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	uuid "github.com/satori/go.uuid"

	"github.com/ProxeusApp/proxeus-core/sys/model"
)

type WorkflowDB struct {
	db           database.Shim
	baseFilePath string
}

const workflowHeavyData = "wh_data"
const workflowVersion = "wf_vers"

func NewWorkflowDB(dir string) (*WorkflowDB, error) {
	var err error

	baseDir := filepath.Join(dir, "workflow")
	err = ensureDir(baseDir)
	if err != nil {
		return nil, err
	}
	assetDir := filepath.Join(baseDir, "assets")
	err = ensureDir(assetDir)
	if err != nil {
		return nil, err
	}
	db, err := database.OpenStorm(filepath.Join(baseDir, "workflows"))
	if err != nil {
		return nil, err
	}
	udb := &WorkflowDB{db: db}
	udb.baseFilePath = assetDir

	example := &model.WorkflowItem{}
	udb.db.Init(example)

	var fVersion int
	verr := udb.db.Get(workflowVersion, workflowVersion, &fVersion)
	if verr == nil && fVersion != example.GetVersion() {
		log.Println("upgrade db", fVersion, "mem", example.GetVersion())
	}
	err = udb.db.Set(workflowVersion, workflowVersion, example.GetVersion())
	if err != nil {
		return nil, err
	}
	return udb, nil
}

func (me *WorkflowDB) AssetsKey() string {
	return me.baseFilePath
}

func (me *WorkflowDB) ListPublished(auth model.Auth, contains string, options storage.Options) ([]*model.WorkflowItem, error) {
	return me.list(auth, contains, options, true)
}

func (me *WorkflowDB) List(auth model.Auth, contains string, options storage.Options) ([]*model.WorkflowItem, error) {
	return me.list(auth, contains, options, false)
}

func (me *WorkflowDB) list(auth model.Auth, contains string, options storage.Options, publishedOnly bool) ([]*model.WorkflowItem, error) {
	params := makeSimpleQuery(options)
	items := make([]*model.WorkflowItem, 0)
	tx, err := me.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var matchers []q.Matcher
	if publishedOnly {
		matchers = publishedMatcher(auth, contains, params)
	} else {
		matchers = defaultMatcher(auth, contains, params, true)
	}

	//when user account is deleted the users workslows will be set to deactivated
	m := q.And(
		q.Eq("Deactivated", false),
	)

	matchers = append(matchers, m)

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
			_ = tx.Get(workflowHeavyData, item.ID, &item.Data)
		}
	}
	return items, nil
}

func (me *WorkflowDB) GetPublished(auth model.Auth, id string) (*model.WorkflowItem, error) {
	var item model.WorkflowItem
	err := me.db.One("ID", id, &item)
	if err != nil {
		return nil, err
	}
	itemRef := &item
	if !(itemRef.OwnedBy(auth) || itemRef.IsPublishedFor(auth)) {
		return nil, model.ErrAuthorityMissing
	}
	me.db.Get(workflowHeavyData, itemRef.ID, &itemRef.Data)
	return itemRef, nil
}

func (me *WorkflowDB) Get(auth model.Auth, id string) (*model.WorkflowItem, error) {
	var item model.WorkflowItem
	err := me.db.One("ID", id, &item)
	if err != nil {
		return nil, err
	}
	itemRef := &item
	if !itemRef.IsPublishedOrReadGrantedFor(auth) {
		return nil, model.ErrAuthorityMissing
	}
	me.db.Get(workflowHeavyData, itemRef.ID, &itemRef.Data)
	return itemRef, nil
}

// Retrieve multiple workflows. If one is not found, an error is returned
func (me *WorkflowDB) GetList(auth model.Auth, ids []string) ([]*model.WorkflowItem, error) {
	var workflows []*model.WorkflowItem
	for _, id := range ids {
		workflow, err := me.Get(auth, id)
		if err != nil {
			return nil, err
		}
		workflows = append(workflows, workflow)
	}
	return workflows, nil
}

func (me *WorkflowDB) Put(auth model.Auth, item *model.WorkflowItem) error {
	return me.put(auth, item, true)
}

func (me *WorkflowDB) put(auth model.Auth, item *model.WorkflowItem, updated bool) error {
	if item == nil {
		return os.ErrInvalid
	}
	if item.ID == "" {
		if !auth.AccessRights().AllowedToCreateEntities() {
			return model.ErrAuthorityMissing
		}
		u2 := uuid.NewV4()
		item.ID = u2.String()
		item.Permissions = model.Permissions{Owner: auth.UserID()}
		tx, err := me.db.Begin(true)
		if err != nil {
			return err
		}
		defer tx.Rollback()
		item.Created = time.Now()
		item.Updated = time.Now()
		return me.updateWF(auth, item, tx)
	} else {
		tx, err := me.db.Begin(true)
		if err != nil {
			return err
		}
		defer tx.Rollback()
		var existing model.WorkflowItem
		err = tx.One("ID", item.ID, &existing)
		if err == storm.ErrNotFound {
			if !auth.AccessRights().AllowedToCreateEntities() {
				return model.ErrAuthorityMissing
			}
			if item.Permissions.Owner == "" {
				item.Permissions = model.Permissions{Owner: auth.UserID()}
				item.Updated = time.Now()
			}
			return me.updateWF(auth, item, tx)
		}
		if err != nil {
			return err
		}
		if existing.Permissions.IsWriteGrantedFor(auth) {
			item.Permissions = *existing.Permissions.Change(auth, &item.Permissions)
			if updated {
				item.Updated = time.Now()
			}
			return me.updateWF(auth, item, tx)
		} else {
			return model.ErrAuthorityMissing
		}
	}
}

func (me *WorkflowDB) updateWF(auth model.Auth, item *model.WorkflowItem, tx database.Shim) error {
	err := me.saveOnly(item, tx)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (me *WorkflowDB) Delete(auth model.Auth, id string) error {
	tx, err := me.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var item model.WorkflowItem
	err = tx.One("ID", id, &item)
	if err != nil {
		return err
	}
	uitem := &item
	if !uitem.Permissions.IsWriteGrantedFor(auth) {
		return model.ErrAuthorityMissing
	}
	//error handling not important here
	tx.Delete(workflowHeavyData, uitem.ID)

	err = tx.DeleteStruct(&item)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (me *WorkflowDB) saveOnly(item *model.WorkflowItem, tx database.Shim) error {
	if item.Data != nil {
		err := tx.Set(workflowHeavyData, item.ID, item.Data)
		if err != nil {
			return err
		}
	}
	cp := *item
	cp.Data = nil
	return tx.Save(&cp)
}

func (me *WorkflowDB) Close() error {
	return me.db.Close()
}
