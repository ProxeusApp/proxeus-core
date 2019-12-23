package storm

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ProxeusApp/proxeus-core/storage"

	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/msgpack"
	"github.com/asdine/storm/q"
	uuid "github.com/satori/go.uuid"

	"github.com/ProxeusApp/proxeus-core/sys/model"
)

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

func (me *WorkflowDB) updateWF(auth model.Auth, item *model.WorkflowItem, tx storm.Node) error {
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

func (me *WorkflowDB) Import(imex storage.ImexIF) error {
	err := me.init(imex)
	if err != nil {
		return err
	}
	for i := 0; true; i++ {
		items, err := imex.DB().Workflow.List(imex.Auth(), "", storage.IndexOptions(i))
		if err == nil && len(items) > 0 {
			for _, item := range items {
				if imex.SkipExistingOnImport() {
					_, err = imex.SysDB().Workflow.Get(imex.Auth(), item.ID)
					if err == nil {
						continue
					}
				}
				item.Permissions.UpdateUserID(imex.LocatedSameUserWithDifferentID())

				err = me.put(imex.Auth(), item, false)
				if err != nil {
					imex.ProcessedEntry(imexWorkflow, item.ID, err)
					continue
				}
				imex.ProcessedEntry(imexWorkflow, item.ID, nil)
			}
		} else {
			break
		}
	}
	return nil
}

const imexWorkflow = "Workflow"

func (me *WorkflowDB) init(imex storage.ImexIF) error {
	var err error
	if imex.DB().Workflow == nil {
		imex.DB().Workflow, err = NewWorkflowDB(imex.Dir())
	}
	return err
}

func (me *WorkflowDB) Export(imex storage.ImexIF, id ...string) error {
	var err error
	err = me.init(imex)
	if err != nil {
		return err
	}
	specificIds := len(id) > 0
	if len(id) == 1 {
		if imex.IsProcessed(imexWorkflow, id[0]) {
			return nil
		}
	}
	type NodesAfter struct {
		id    string
		store storage.ImporterExporter
	}
	nodes := make(map[string]*NodesAfter)
	for i := 0; true; i++ {
		items, err := me.List(imex.Auth(), "", storage.IndexOptions(i).WithInclude(id))
		if err == nil && len(items) > 0 {
			var tx storm.Node
			tx, err = me.getDB().Begin(true)
			if err != nil {
				return err
			}
			for _, item := range items {
				if !imex.IsProcessed(imexWorkflow, item.ID) {
					err = me.saveOnly(item, tx)
					if err != nil {
						imex.ProcessedEntry(imexWorkflow, item.ID, err)
						continue
					}
					if item.Data != nil && item.Data.Flow != nil {
						for _, v := range item.Data.Flow.Nodes {
							if im := imex.SysDB().ImexIFByName(v.Type); im != nil {
								nodes[v.ID] = &NodesAfter{id: item.ID, store: im}
							}
						}
					}
					item.Permissions.UserIdsMap(imex.NeededUsers())
					if err != nil {
						imex.ProcessedEntry(imexWorkflow, item.ID, err)
						continue
					}
					imex.ProcessedEntry(imexWorkflow, item.ID, nil)
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
			imex.ProcessedEntry(imexWorkflow, v.id, err)
		}
	}
	return nil
}

func (me *WorkflowDB) Close() error {
	return me.db.Close()
}
