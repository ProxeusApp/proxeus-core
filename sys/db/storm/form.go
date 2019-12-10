package storm

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/msgpack"
	"github.com/asdine/storm/q"
	uuid "github.com/satori/go.uuid"

	"github.com/ProxeusApp/proxeus-core/sys/form"
	"github.com/ProxeusApp/proxeus-core/sys/model"
)

type FormDB struct {
	db           *storm.DB
	baseFilePath string
}

//formHeavyData helps us to load the data of the model.Form entity when it is requested by metaOnly = false
const formHeavyData = "form_data"

//formVersion helps us to keep track of the structure version between persistence and memory
const formVersion = "form_version"

const formCompVersion = "formComp_version"

func NewFormDB(dir string) (*FormDB, error) {
	var err error
	var msgpackDb *storm.DB
	baseDir := filepath.Join(dir, "form")
	err = ensureDir(baseDir)
	if err != nil {
		return nil, err
	}
	assetDir := filepath.Join(baseDir, "assets")
	err = ensureDir(assetDir)
	if err != nil {
		return nil, err
	}
	msgpackDb, err = storm.Open(filepath.Join(baseDir, "forms"), storm.Codec(msgpack.Codec))
	if err != nil {
		return nil, err
	}
	udb := &FormDB{db: msgpackDb}
	udb.baseFilePath = assetDir

	example := &model.FormItem{}
	exampleComp := &model.FormComponentItem{}
	udb.db.Init(example)
	udb.db.Init(exampleComp)
	initVars(udb.db)
	var fVersion int
	verr := udb.db.Get(formVersion, formVersion, &fVersion)
	if verr == nil && fVersion != example.GetVersion() {
		log.Println("upgrade db", fVersion, "mem", example.GetVersion())
	}
	err = udb.db.Set(formVersion, formVersion, example.GetVersion())
	var fcmpVersion int

	verr = udb.db.Get(formCompVersion, formCompVersion, &fcmpVersion)
	if verr == nil && fcmpVersion != exampleComp.GetVersion() {
		log.Println("upgrade db", fcmpVersion, "mem", exampleComp.GetVersion())
	}
	err = udb.db.Set(formCompVersion, formCompVersion, exampleComp.GetVersion())
	if err != nil {
		return nil, err
	}
	return udb, nil
}

func (me *FormDB) List(auth model.Authorization, contains string, options map[string]interface{}) ([]*model.FormItem, error) {
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

func (me *FormDB) Get(auth model.Authorization, id string) (*model.FormItem, error) {
	var item model.FormItem
	err := me.db.One("ID", id, &item)
	if err != nil {
		return nil, err
	}
	itemRef := &item
	if !itemRef.IsPublishedOrReadGrantedFor(auth) {
		return nil, model.ErrAuthorityMissing
	}
	//error handling not important
	_ = me.db.Get(formHeavyData, itemRef.ID, &itemRef.Data)
	return itemRef, nil
}

func (me *FormDB) Put(auth model.Authorization, item *model.FormItem) error {
	return me.put(auth, item, true)
}

func (me *FormDB) put(auth model.Authorization, item *model.FormItem, updated bool) error {
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
		return me.updateForm(auth, item, tx)
	} else {
		tx, err := me.db.Begin(true)
		if err != nil {
			return err
		}
		defer tx.Rollback()
		var existing model.FormItem
		err = tx.One("ID", item.ID, &existing)
		if err == storm.ErrNotFound {
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

func (me *FormDB) Delete(auth model.Authorization, id string) error {
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

func (me *FormDB) updateForm(auth model.Authorization, item *model.FormItem, tx storm.Node) error {
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

func (me *FormDB) updateVars(auth model.Authorization, item *model.FormItem, tx storm.Node) error {
	formVars := form.Vars(item.Data)
	if len(formVars) > 0 {
		err := updateVarsOf(auth, item.ID, formVars, tx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (me *FormDB) DelComp(auth model.Authorization, id string) error {
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

func (me *FormDB) PutComp(auth model.Authorization, comp *model.FormComponentItem) error {
	if comp == nil {
		return errors.New("item cannot be nil")
	}
	if !auth.AccessRights().AllowedToCreateEntities() {
		return model.ErrAuthorityMissing
	}

	if comp.ID == "" { //insert new item
		u2 := uuid.NewV4()
		comp.ID = u2.String()
	}

	return me.db.Save(comp)
}

func (me *FormDB) GetComp(auth model.Authorization, id string) (*model.FormComponentItem, error) {
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

func (me *FormDB) ListComp(auth model.Authorization, contains string, options map[string]interface{}) (map[string]*model.FormComponentItem, error) {
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

func (me *FormDB) Vars(auth model.Authorization, contains string, options map[string]interface{}) ([]string, error) {
	contains = regexp.QuoteMeta(contains)
	params := makeSimpleQuery(options)
	tx, err := me.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	return getVars(contains, params.limit, params.index, tx)
}

func (me *FormDB) Import(imex *Imex) error {
	err := me.init(imex)
	if err != nil {
		return err
	}
	for i := 0; true; i++ {
		items, err := imex.db.Form.ListComp(imex.auth, "", map[string]interface{}{"index": i, "limit": 1000, "metaOnly": false})
		if err == nil && len(items) > 0 {
			for _, item := range items {
				if imex.skipExistingOnImport {
					_, err = imex.sysDB.Form.GetComp(imex.auth, item.ID)
					if err == nil {
						continue
					}
				}
				err = imex.sysDB.Form.PutComp(imex.auth, item)
				if err != nil {
					imex.processedEntry(imexFormComponent, item.ID, err)
					continue
				}
				imex.processedEntry(imexFormComponent, item.ID, nil)
			}
		} else {
			break
		}
	}
	for i := 0; true; i++ {
		items, err := imex.db.Form.List(imex.auth, "", map[string]interface{}{"index": i, "limit": 1000, "metaOnly": false})
		if err == nil && len(items) > 0 {
			for _, item := range items {
				if imex.skipExistingOnImport {
					_, err = imex.sysDB.Form.Get(imex.auth, item.ID)
					if err == nil {
						continue
					}
				}

				item.Permissions.UpdateUserID(imex.locatedSameUserWithDifferentID)

				err = imex.sysDB.Form.put(imex.auth, item, false)
				if err != nil {
					imex.processedEntry(imexForm, item.ID, err)
					continue
				}
				imex.processedEntry(imexForm, item.ID, nil)
				//TODO export i18n of comps
			}
		} else {
			break
		}
	}
	return nil
}

const imexForm = "Form"
const imexFormComponent = "FormComponent"

func (me *FormDB) init(imex *Imex) error {
	var err error
	if imex.db.Form == nil {
		imex.db.Form, err = NewFormDB(imex.dir)
	}
	return err
}

func (me *FormDB) Export(imex *Imex, id ...string) error {
	var err error
	err = me.init(imex)
	if err != nil {
		return err
	}
	specificIds := len(id) > 0
	if len(id) == 1 {
		if imex.isProcessed(imexForm, id[0]) {
			return nil
		}
	}
	formComps := make(map[string]bool)

	for i := 0; true; i++ {
		items, err := imex.sysDB.Form.List(imex.auth, "", map[string]interface{}{"include": id, "index": i, "limit": 1000, "metaOnly": false})
		if err == nil && len(items) > 0 {
			var tx storm.Node
			tx, err = imex.db.Form.db.Begin(true)
			if err != nil {
				return err
			}
			for _, item := range items {
				if imex.isProcessed(imexForm, item.ID) {
					continue
				}
				if len(item.Data) > 0 {
					err = tx.Set(formHeavyData, item.ID, item.Data)
					if err != nil {
						imex.processedEntry(imexForm, item.ID, err)
						continue
					}
				}
				formSrc := form.GetFormSrc(item.Data)
				item.Data = nil //
				err = tx.Save(item)
				if err != nil {
					imex.processedEntry(imexForm, item.ID, err)
					continue
				}
				form.LoopComponents(formSrc, func(compId, compInstId string, compMain interface{}, comp map[string]interface{}) bool {
					if compId == "" {
						imex.processedEntry(imexForm, item.ID, fmt.Errorf("form contains a component with an empty id"))
					}
					formComps[compId] = true
					return true
				})
				item.Permissions.UserIdsMap(imex.neededUsers)
				imex.processedEntry(imexForm, item.ID, nil)
				//TODO export i18n of comps
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
	if len(formComps) > 0 {
		var tx storm.Node
		tx, err = imex.db.Form.db.Begin(true)
		if err != nil {
			return err
		}
		for compId := range formComps {
			var fbi model.FormComponentItem
			err = imex.sysDB.Form.db.One("ID", compId, &fbi)
			if err != nil {
				imex.processedEntry(imexFormComponent, compId, err)
				continue
			}
			err = tx.Save(&fbi)
			imex.processedEntry(imexFormComponent, compId, err)
		}
		err = tx.Commit()
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return nil
}

func (me *FormDB) Close() error {
	return me.db.Close()
}
