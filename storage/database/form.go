package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/ProxeusApp/proxeus-core/storage"

	"github.com/ProxeusApp/proxeus-core/storage/database/db"

	"github.com/asdine/storm/q"
	uuid "github.com/satori/go.uuid"

	"github.com/ProxeusApp/proxeus-core/sys/form"
	"github.com/ProxeusApp/proxeus-core/sys/model"
)

type FormDB struct {
	db           db.DB
	baseFilePath string
}

//formHeavyData helps us to load the data of the model.Form entity when it is requested by metaOnly = false
const formHeavyData = "form_data"

//formVersion helps us to keep track of the structure version between persistence and memory
const formVersion = "form_version"

const formCompVersion = "formComp_version"

func NewFormDB(c DBConfig) (*FormDB, error) {
	baseDir := filepath.Join(c.Dir, "form")
	db, err := db.OpenDatabase(c.Engine, c.URI, filepath.Join(baseDir, "forms"))
	if err != nil {
		return nil, err
	}
	udb := &FormDB{db: db, baseFilePath: filepath.Join(baseDir, "assets")}

	udb.db.Init(formHeavyData)

	example := &model.FormItem{}
	exampleComp := &model.FormComponentItem{}
	udb.db.Init(example)
	udb.db.Init(exampleComp)
	initVars(udb.db)

	err = udb.db.Set(formVersion, formVersion, example.GetVersion())
	if err != nil {
		return nil, err
	}
	err = udb.db.Set(formCompVersion, formCompVersion, exampleComp.GetVersion())
	if err != nil {
		return nil, err
	}
	return udb, nil
}

func (me *FormDB) List(auth model.Auth, contains string, options storage.Options) ([]*model.FormItem, error) {
	params := makeSimpleQuery(options)
	var items []*model.FormItem
	tx, err := me.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	matcher := defaultMatcher(auth, contains, params, true)
	err = tx.Select(matcher...).
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
			_ = tx.Get(formHeavyData, item.ID, &item.Data)
		}
	}
	return items, nil
}

func (me *FormDB) Get(auth model.Auth, id string) (*model.FormItem, error) {
	var item model.FormItem
	err := me.db.One("ID", id, &item)
	if err != nil {
		return nil, err
	}
	itemRef := &item
	if !itemRef.IsPublishedOrReadGrantedFor(auth) {
		return nil, model.ErrAuthorityMissing
	}
	me.db.Get(formHeavyData, itemRef.ID, &itemRef.Data)
	return itemRef, err
}

func (me *FormDB) Put(auth model.Auth, item *model.FormItem) error {
	return me.put(auth, item, true)
}

func (me *FormDB) put(auth model.Auth, item *model.FormItem, updated bool) error {
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
		return me.updateForm(auth, item, tx)
	} else {
		tx, err := me.db.Begin(true)
		if err != nil {
			return err
		}
		defer tx.Rollback()
		var existing model.FormItem
		err = tx.One("ID", item.ID, &existing)
		if db.NotFound(err) {
			err = nil
			if !auth.AccessRights().AllowedToCreateEntities() {
				return model.ErrAuthorityMissing
			}
			if item.Permissions.Owner == "" {
				item.Permissions = model.Permissions{Owner: auth.UserID()}
				item.Updated = time.Now()
			}
			return me.updateForm(auth, item, tx)
		}
		if err != nil {
			return err
		}
		existingRef := &existing
		if existingRef.Permissions.IsWriteGrantedFor(auth) {
			item.Permissions = *existingRef.Permissions.Change(auth, &item.Permissions)
			if updated {
				item.Updated = time.Now()
			}
			return me.updateForm(auth, item, tx)
		} else {
			return model.ErrAuthorityMissing
		}
	}
}

func (me *FormDB) Delete(auth model.Auth, id string) error {
	tx, err := me.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var item model.FormItem
	err = tx.One("ID", id, &item)
	if err != nil {
		return err
	}
	uitem := &item
	if !uitem.Permissions.IsWriteGrantedFor(auth) {
		return model.ErrAuthorityMissing
	}
	//error handling not important here
	tx.Delete(formHeavyData, uitem.ID)

	err = tx.DeleteStruct(&item)
	if err != nil {
		return err
	}

	err = remVars(auth, id, tx)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (me *FormDB) updateForm(auth model.Auth, item *model.FormItem, tx db.DB) error {
	err := me.updateVars(auth, item, tx)
	if err != nil {
		return err
	}
	if len(item.Data) > 0 {
		err = tx.Set(formHeavyData, item.ID, item.Data)
		if err != nil {
			return err
		}
	}
	cp := *item
	cp.Data = nil
	err = tx.Save(&cp)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (me *FormDB) updateVars(auth model.Auth, item *model.FormItem, tx db.DB) error {
	formVars := form.Vars(item.Data)
	if len(formVars) > 0 {
		err := updateVarsOf(auth, item.ID, formVars, tx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (me *FormDB) DelComp(auth model.Auth, id string) error {
	if !auth.AccessRights().IsGrantedForUserModifications() {
		return model.ErrAuthorityMissing
	}
	tx, err := me.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var item model.FormComponentItem
	err = tx.One("ID", id, &item)
	if err != nil {
		return err
	}
	err = tx.DeleteStruct(&item)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (me *FormDB) PutComp(auth model.Auth, comp *model.FormComponentItem) error {
	if comp == nil {
		return errors.New("item cannot be nil")
	}
	if !auth.AccessRights().AllowedToCreateEntities() {
		return model.ErrAuthorityMissing
	}

	if comp.ID == "" { //insert new item
		comp.ID = uuid.NewV4().String()
	}

	return me.db.Save(comp)
}

func (me *FormDB) GetComp(auth model.Auth, id string) (*model.FormComponentItem, error) {
	if id == "" {
		return nil, fmt.Errorf("id can't be empty")
	}
	var err error
	var fbi model.FormComponentItem
	err = me.db.One("ID", id, &fbi)
	if err != nil {
		return nil, err
	}
	return &fbi, nil
}

//makes it possible to search to content no matter what field type it is
type formComponentSearchMatcher struct {
	Contains     string
	re           *regexp.Regexp
	reCompileErr bool
}

func (me *formComponentSearchMatcher) MatchField(v interface{}) (bool, error) {
	if me.reCompileErr {
		return false, nil
	}
	bts, err := json.Marshal(v)
	if err != nil {
		return false, nil
	}
	if me.re == nil {
		me.re, err = regexp.Compile(`(?i)` + me.Contains)
		if err != nil {
			me.reCompileErr = true
			return false, nil
		}
	}
	return len(me.re.FindIndex(bts)) > 0, nil
}

func (me *FormDB) ListComp(auth model.Auth, contains string, options storage.Options) (map[string]*model.FormComponentItem, error) {
	params := makeSimpleQuery(options)
	items := make(map[string]*model.FormComponentItem)
	tx, err := me.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	contains = containsCaseInsensitiveReg(contains)
	matchers := make([]q.Matcher, 0)
	if contains != "" {
		fcSearchMatcher := &formComponentSearchMatcher{Contains: contains}
		matchers = append(matchers,
			q.And(
				q.Or(
					q.Re("Name", contains),
					q.Re("Detail", contains),
					q.NewFieldMatcher("Template", fcSearchMatcher),
					q.NewFieldMatcher("Settings", fcSearchMatcher),
				),
			),
		)
	}
	if len(params.exclude) > 0 {
		matchers = append(matchers,
			q.And(
				q.Not(q.In("ID", params.exclude)),
			),
		)
	}
	if len(params.include) > 0 {
		matchers = append(matchers,
			q.And(
				q.In("ID", params.include),
			),
		)
	}
	err = tx.Select(matchers...).
		Limit(params.limit).
		Skip(params.index).
		OrderBy("Updated").
		Reverse().
		Each(new(model.FormComponentItem), func(record interface{}) error {
			item := record.(*model.FormComponentItem)
			items[item.ID] = item
			return nil
		})
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (me *FormDB) Vars(auth model.Auth, contains string, options storage.Options) ([]string, error) {
	contains = regexp.QuoteMeta(contains)
	params := makeSimpleQuery(options)
	tx, err := me.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	return getVars(contains, params.limit, params.index, tx)
}

func (me *FormDB) Close() error {
	return me.db.Close()
}
