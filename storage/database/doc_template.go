package database

import (
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/ProxeusApp/proxeus-core/storage"

	"github.com/ProxeusApp/proxeus-core/storage/database/db"

	uuid "github.com/satori/go.uuid"

	"github.com/ProxeusApp/proxeus-core/sys/file"
	"github.com/ProxeusApp/proxeus-core/sys/model"
)

type DocTemplateDB struct {
	db           db.DB
	baseFilePath string
}

const docTemplVersion = "doctmpl_vers"

// NewDocTemplateDB opens a handle to the Document Template database
func NewDocTemplateDB(c DBConfig) (*DocTemplateDB, error) {
	baseDir := filepath.Join(c.Dir, "document_template")
	db, err := db.OpenDatabase(c.Engine, c.URI, filepath.Join(baseDir, "document_templates"))
	if err != nil {
		return nil, err
	}
	udb := &DocTemplateDB{db: db, baseFilePath: filepath.Join(baseDir, "assets")}
	example := &model.TemplateItem{}
	udb.db.Init(example)
	initVars(udb.db)
	err = udb.db.Set(docTemplVersion, docTemplVersion, example.GetVersion())
	if err != nil {
		return nil, err
	}
	return udb, nil
}

// AssetsKeys return the base filepath of the document templates
func (me *DocTemplateDB) AssetsKey() string {
	return me.baseFilePath
}

// List fetches template items from the database based on the supplied filter options
func (me *DocTemplateDB) List(auth model.Auth, contains string, options storage.Options) ([]*model.TemplateItem, error) {
	params := makeSimpleQuery(options)
	items := make([]*model.TemplateItem, 0)
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
	if params.metaOnly {
		for _, item := range items {
			item.Data = nil
		}
	} else {
		for _, item := range items {
			for _, v := range item.Data {
				v.SetBaseDir(me.baseFilePath)
			}
		}
	}
	return items, nil
}

// Get retrieves a sinlge template item using its key
func (me *DocTemplateDB) Get(auth model.Auth, id string) (*model.TemplateItem, error) {
	var item model.TemplateItem
	err := me.db.One("ID", id, &item)
	if err != nil {
		return nil, err
	}
	itemRef := &item
	if !itemRef.IsPublishedOrReadGrantedFor(auth) {
		return nil, model.ErrAuthorityMissing
	}
	for _, v := range itemRef.Data {
		v.SetBaseDir(me.baseFilePath)
	}
	return itemRef, nil
}

// On return, n == len(buf) if and only if err == nil.
func (me *DocTemplateDB) ProvideFileInfoFor(auth model.Auth, id, lang string, fm *file.Meta) (*file.IO, error) {
	return me.getFileInfoFor(auth, id, lang, fm)
}

// PutVars inserts a new variable into the databse
func (me *DocTemplateDB) PutVars(auth model.Auth, id, lang string, vars []string) error {
	tx, err := me.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	err = updateVarsOf(auth, id+lang, vars, tx)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (me *DocTemplateDB) GetTemplate(auth model.Auth, id, lang string) (*file.IO, error) {
	return me.getFileInfoFor(auth, id, lang, nil)
}

func (me *DocTemplateDB) getFileInfoFor(auth model.Auth, id, lang string, fm *file.Meta) (*file.IO, error) {
	if id == "" || lang == "" {
		return nil, os.ErrInvalid
	}
	tmplItem, err := me.Get(auth, id)
	if err != nil {
		return nil, err
	}
	if len(tmplItem.Data) > 0 {
		if t, exists := tmplItem.Data[lang]; exists {
			if fm == nil && t.PathName() == "" {
				//if we just request the file
				return nil, os.ErrNotExist
			}
			if fm != nil && t.Name() != fm.Name {
				t.Update(fm.Name, fm.ContentType)
				err = me.put(auth, tmplItem, true)
				if err != nil {
					return nil, err
				}
			}
			return t, nil
		}
	}
	if fm != nil {
		if tmplItem.Data == nil {
			tmplItem.Data = make(map[string]*file.IO)
		}
		fi := file.New(me.baseFilePath, *fm)
		tmplItem.Data[lang] = fi
		err = me.Put(auth, tmplItem)
		if err != nil {
			return nil, err
		}
		return fi, nil
	}
	return nil, os.ErrNotExist
}

// DeleteTemplate removes a template from the database
func (me *DocTemplateDB) DeleteTemplate(auth model.Auth, files storage.FilesIF, id, lang string) error {
	if auth == nil || id == "" || lang == "" {
		return os.ErrInvalid
	}
	tmplItem, err := me.Get(auth, id)
	if err != nil {
		return err
	}
	if !tmplItem.Permissions.IsWriteGrantedFor(auth) {
		return model.ErrAuthorityMissing
	}
	if len(tmplItem.Data) > 0 {
		if t, exists := tmplItem.Data[lang]; exists {
			p := t.Path()
			if p != "" {
				delete(tmplItem.Data, lang)
				tx, err := me.db.Begin(true)
				if err != nil {
					return err
				}
				defer tx.Rollback()
				err = tx.Save(tmplItem)
				if err != nil {
					return err
				}
				remVars(auth, id+lang, tx)
				err = tx.Commit()
				if err != nil {
					return err
				}
				return files.Delete(p)
			}
		}
	}
	return os.ErrNotExist
}

func (me *DocTemplateDB) Put(auth model.Auth, item *model.TemplateItem) error {
	return me.put(auth, item, true)
}

func (me *DocTemplateDB) put(auth model.Auth, item *model.TemplateItem, updated bool) error {
	if auth == nil || item == nil {
		return os.ErrInvalid
	}
	if item.ID == "" {
		if !auth.AccessRights().AllowedToCreateEntities() {
			return model.ErrAuthorityMissing
		}
		item.ID = uuid.NewV4().String()
		item.Permissions = model.Permissions{Owner: auth.UserID()}
		item.Created = time.Now()
		item.Updated = time.Now()
		return me.db.Save(item)
	} else {
		var existing model.TemplateItem
		err := me.db.One("ID", item.ID, &existing)
		if db.NotFound(err) {
			if !auth.AccessRights().AllowedToCreateEntities() {
				return model.ErrAuthorityMissing
			}
			if item.Permissions.Owner == "" {
				item.Permissions = model.Permissions{Owner: auth.UserID()}
				item.Updated = time.Now()
			}
			return me.db.Save(item)
		}
		if err != nil {
			return err
		}
		if !existing.Permissions.IsWriteGrantedFor(auth) {
			return model.ErrAuthorityMissing
		}
		item.Permissions = *existing.Permissions.Change(auth, &item.Permissions)
		if updated {
			item.Updated = time.Now()
		}
		return me.db.Save(item)
	}
}

// Delete removes a template's files on the file system
func (me *DocTemplateDB) Delete(auth model.Auth, files storage.FilesIF, id string) error {
	tx, err := me.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var item model.TemplateItem
	err = tx.One("ID", id, &item)
	if err != nil {
		return err
	}
	uitem := &item
	if !uitem.Permissions.IsWriteGrantedFor(auth) {
		return model.ErrAuthorityMissing
	}

	//rem all tmpl docs
	for lang, tmpl := range uitem.Data {
		if tmpl != nil {
			//set base dir after loading from disk
			tmpl.SetBaseDir(me.baseFilePath)
			p := tmpl.Path()
			if p != "" {
				_ = remVars(auth, id+lang, tx)
				files.Delete(p)
			}
		}
	}

	err = tx.DeleteStruct(uitem)
	if err != nil {
		return err
	}
	return tx.Commit()
}

// Vars returns a list of variables defines for a specific template
func (me *DocTemplateDB) Vars(auth model.Auth, contains string, options storage.Options) ([]string, error) {
	contains = regexp.QuoteMeta(contains)
	params := makeSimpleQuery(options)
	tx, err := me.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	return getVars(contains, params.limit, params.index, tx)
}

func (me *DocTemplateDB) Close() error {
	return me.db.Close()
}
