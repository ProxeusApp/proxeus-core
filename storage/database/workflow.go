package database

import (
	"github.com/ProxeusApp/proxeus-core/externalnode"
	"os"
	"path/filepath"
	"time"

	"github.com/ProxeusApp/proxeus-core/storage"

	"github.com/ProxeusApp/proxeus-core/storage/database/db"

	"github.com/asdine/storm/q"
	uuid "github.com/satori/go.uuid"

	"github.com/ProxeusApp/proxeus-core/sys/model"
)

type WorkflowDB struct {
	db           db.DB
	baseFilePath string
}

const workflowHeavyData = "wh_data"
const workflowVersion = "wf_vers"
const externalNodeInstance = "ExternalNodeInstance"

// NewWorkflowDB returns a handle to the workflow database
func NewWorkflowDB(c DBConfig) (*WorkflowDB, error) {
	baseDir := filepath.Join(c.Dir, "workflow")
	db, err := db.OpenDatabase(c.Engine, c.URI, filepath.Join(baseDir, "workflows"))
	if err != nil {
		return nil, err
	}
	udb := &WorkflowDB{db: db, baseFilePath: filepath.Join(baseDir, "assets")}

	udb.db.Init(workflowHeavyData)

	udb.db.Init(externalNodeInstance)
	udb.db.Init(new(externalnode.ExternalNode))

	example := &model.WorkflowItem{}
	udb.db.Init(example)

	err = udb.db.Set(workflowVersion, workflowVersion, example.GetVersion())
	if err != nil {
		return nil, err
	}
	return udb, nil
}

// ListPublished returns all workflow items matching the supplied filter options that are flagged as published
func (me *WorkflowDB) ListPublished(auth model.Auth, contains string, options storage.Options) ([]*model.WorkflowItem, error) {
	return me.list(auth, contains, options, true)
}

// ListPublished returns all workflow items matching the supplied filter options
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

// ListPublished returns a workflow item matching the supplied filter options that if it is flagged as published
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

// Get retrieves a worklfow item machting its id
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

// GetList retrieves multiple workflows by matching their id
func (me *WorkflowDB) GetList(auth model.Auth, ids []string) ([]*model.WorkflowItem, error) {
	var workflows []*model.WorkflowItem
	for _, id := range ids {
		if id == "" {
			continue
		}
		workflow, err := me.Get(auth, id)
		if err != nil {
			return nil, err
		}
		workflows = append(workflows, workflow)
	}
	return workflows, nil
}

//Put adds a workflow item into the database
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
		item.ID = uuid.NewV4().String()
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
		if db.NotFound(err) {
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

func (me *WorkflowDB) updateWF(auth model.Auth, item *model.WorkflowItem, tx db.DB) error {
	err := me.saveOnly(item, tx)
	if err != nil {
		return err
	}
	return tx.Commit()
}

// Delete removes a workflow item from the database by matching its id
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

func (me *WorkflowDB) saveOnly(item *model.WorkflowItem, tx db.DB) error {
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

// RegisterExternalNode saves an external node definition
func (me *WorkflowDB) RegisterExternalNode(auth model.Auth, n *externalnode.ExternalNode) error {
	return me.db.Save(n)
}

// NodeByName retrieves an external node definition matching the supplied name
func (me *WorkflowDB) NodeByName(auth model.Auth, name string) (*externalnode.ExternalNode, error) {
	var i externalnode.ExternalNode
	err := me.db.One("Name", name, &i)
	return &i, err
}

// QueryFromInstanceID return an external node instance by machting the specified id
func (me *WorkflowDB) QueryFromInstanceID(auth model.Auth, id string) (externalnode.ExternalQuery, error) {
	var i externalnode.ExternalNodeInstance
	err := me.db.Get(externalNodeInstance, id, &i)
	if err != nil {
		return externalnode.ExternalQuery{}, err
	}
	n, err := me.NodeByName(auth, i.NodeName)
	if err != nil {
		return externalnode.ExternalQuery{}, err
	}
	return externalnode.ExternalQuery{
		ExternalNode:         n,
		ExternalNodeInstance: &i,
	}, nil
}

// ListExternalNodes return a list of all external node definitions
func (me *WorkflowDB) ListExternalNodes() []*externalnode.ExternalNode {
	var l []*externalnode.ExternalNode
	me.db.Select().Each(new(externalnode.ExternalNode), func(record interface{}) error {
		item := record.(*externalnode.ExternalNode)
		l = append(l, item)
		return nil
	})
	return l
}

// DeleteExternalNode remove an external node definition
func (me *WorkflowDB) DeleteExternalNode(auth model.Auth, id string) error {
	return me.db.DeleteStruct(&externalnode.ExternalNode{ID: id})
}

// PutExternalNodeInstance saves an instance of an external node to the database
func (me *WorkflowDB) PutExternalNodeInstance(auth model.Auth, item *externalnode.ExternalNodeInstance) error {
	tx, err := me.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var i externalnode.ExternalNodeInstance
	err = tx.Get(externalNodeInstance, item.ID, &i)
	if db.NotFound(err) {
		if !auth.AccessRights().AllowedToCreateEntities() {
			return model.ErrAuthorityMissing
		}
	} else if err != nil {
		return err
	}
	err = tx.Set(externalNodeInstance, item.ID, item)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (me *WorkflowDB) Close() error {
	return me.db.Close()
}
