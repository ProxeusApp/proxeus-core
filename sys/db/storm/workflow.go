package storm

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/msgpack"
	"github.com/asdine/storm/q"
	uuid "github.com/satori/go.uuid"

	"git.proxeus.com/core/central/sys/model"
)

type WorkflowDBInterface interface {
	ListPublished(auth model.Authorization, contains string, options map[string]interface{}) ([]*model.WorkflowItem, error)
	List(auth model.Authorization, contains string, options map[string]interface{}) ([]*model.WorkflowItem, error)
	GetPublished(auth model.Authorization, id string) (*model.WorkflowItem, error)
	Get(auth model.Authorization, id string) (*model.WorkflowItem, error)
	Put(auth model.Authorization, item *model.WorkflowItem) error
	put(auth model.Authorization, item *model.WorkflowItem, updated bool) error
	getDB() *storm.DB
	updateWF(auth model.Authorization, item *model.WorkflowItem, tx storm.Node) error
	Delete(auth model.Authorization, id string) error
	saveOnly(item *model.WorkflowItem, tx storm.Node) error
	Import(imex *Imex) error
	Export(imex *Imex, id ...string) error
	Close() error
}

type WorkflowDB struct {
	db           *storm.DB
	baseFilePath string
}

const workflowHeavyData = "wh_data"
const workflowVersion = "wf_vers"

func NewWorkflowDB(dir string) (*WorkflowDB, error) {
	var err error
	var msgpackDb *storm.DB
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
	msgpackDb, err = storm.Open(filepath.Join(baseDir, "workflows"), storm.Codec(msgpack.Codec))
	if err != nil {
		return nil, err
	}
	udb := &WorkflowDB{db: msgpackDb}
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

func (me *WorkflowDB) getDB() *storm.DB {
	return me.db
}

func (me *WorkflowDB) ListPublished(auth model.Authorization, contains string, options map[string]interface{}) ([]*model.WorkflowItem, error) {
	return me.list(auth, contains, options, true)
}

func (me *WorkflowDB) List(auth model.Authorization, contains string, options map[string]interface{}) ([]*model.WorkflowItem, error) {
	return me.list(auth, contains, options, false)
}

func (me *WorkflowDB) list(auth model.Authorization, contains string, options map[string]interface{}, publishedOnly bool) ([]*model.WorkflowItem, error) {
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

func (me *WorkflowDB) GetPublished(auth model.Authorization, id string) (*model.WorkflowItem, error) {
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

func (me *WorkflowDB) Get(auth model.Authorization, id string) (*model.WorkflowItem, error) {
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

func (me *WorkflowDB) Put(auth model.Authorization, item *model.WorkflowItem) error {
	return me.put(auth, item, true)
}

func (me *WorkflowDB) put(auth model.Authorization, item *model.WorkflowItem, updated bool) error {
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

func (me *WorkflowDB) updateWF(auth model.Authorization, item *model.WorkflowItem, tx storm.Node) error {
	err := me.saveOnly(item, tx)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (me *WorkflowDB) Delete(auth model.Authorization, id string) error {
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

func (me *WorkflowDB) saveOnly(item *model.WorkflowItem, tx storm.Node) error {
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

func (me *WorkflowDB) Import(imex *Imex) error {
	err := me.init(imex)
	if err != nil {
		return err
	}
	for i := 0; true; i++ {
		items, err := imex.db.Workflow.List(imex.auth, "", map[string]interface{}{"index": i, "limit": 1000, "metaOnly": false})
		if err == nil && len(items) > 0 {
			for _, item := range items {
				if imex.skipExistingOnImport {
					_, err = imex.sysDB.Workflow.Get(imex.auth, item.ID)
					if err == nil {
						continue
					}
				}
				item.Permissions.UpdateUserID(imex.locatedSameUserWithDifferentID)

				err = imex.sysDB.Workflow.put(imex.auth, item, false)
				if err != nil {
					imex.processedEntry(imexWorkflow, item.ID, err)
					continue
				}
				imex.processedEntry(imexWorkflow, item.ID, nil)
			}
		} else {
			break
		}
	}
	return nil
}

const imexWorkflow = "Workflow"

func (me *WorkflowDB) init(imex *Imex) error {
	var err error
	if imex.db.Workflow == nil {
		imex.db.Workflow, err = NewWorkflowDB(imex.dir)
	}
	return err
}

func (me *WorkflowDB) Export(imex *Imex, id ...string) error {
	var err error
	err = me.init(imex)
	if err != nil {
		return err
	}
	specificIds := len(id) > 0
	if len(id) == 1 {
		if imex.isProcessed(imexWorkflow, id[0]) {
			return nil
		}
	}
	type NodesAfter struct {
		id    string
		store ImexIF
	}
	nodes := make(map[string]*NodesAfter)
	for i := 0; true; i++ {
		items, err := me.List(imex.auth, "", map[string]interface{}{"include": id, "index": i, "limit": 1000, "metaOnly": false})
		if err == nil && len(items) > 0 {
			var tx storm.Node
			tx, err = imex.db.Workflow.getDB().Begin(true)
			if err != nil {
				return err
			}
			for _, item := range items {
				if !imex.isProcessed(imexWorkflow, item.ID) {
					err = imex.db.Workflow.saveOnly(item, tx)
					if err != nil {
						imex.processedEntry(imexWorkflow, item.ID, err)
						continue
					}
					if item.Data != nil && item.Data.Flow != nil {
						for _, v := range item.Data.Flow.Nodes {
							if im := imex.sysDB.ImexIFByName(v.Type); im != nil {
								nodes[v.ID] = &NodesAfter{id: item.ID, store: im}
							}
						}
					}
					item.Permissions.UserIdsMap(imex.neededUsers)
					if err != nil {
						imex.processedEntry(imexWorkflow, item.ID, err)
						continue
					}
					imex.processedEntry(imexWorkflow, item.ID, nil)
				}
			}
			err = tx.Commit()
			if err != nil {
				_ = tx.Rollback()
				return err
			}
		} else {
			break
		}
		if specificIds {
			break
		}
	}
	for k, v := range nodes {
		err = v.store.Export(imex, k)
		if err != nil {
			imex.processedEntry(imexWorkflow, v.id, err)
		}
	}
	return nil
}

func (me *WorkflowDB) Close() error {
	return me.db.Close()
}
