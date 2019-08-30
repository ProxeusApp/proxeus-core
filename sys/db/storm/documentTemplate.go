package storm

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/msgpack"
	uuid "github.com/satori/go.uuid"

	"git.proxeus.com/core/central/sys/file"
	"git.proxeus.com/core/central/sys/model"
)

type DocTemplateDB struct {
	db           *storm.DB
	baseFilePath string
}

const docTemplVersion = "doctmpl_vers"

func NewDocTemplateDB(dir string) (*DocTemplateDB, error) {
	var err error
	var msgpackDb *storm.DB
	baseDir := filepath.Join(dir, "document_template")
	err = ensureDir(baseDir)
	if err != nil {
		return nil, err
	}
	assetDir := filepath.Join(baseDir, "assets")
	err = ensureDir(assetDir)
	if err != nil {
		return nil, err
	}
	msgpackDb, err = storm.Open(filepath.Join(baseDir, "document_templates"), storm.Codec(msgpack.Codec))
	if err != nil {
		return nil, err
	}
	udb := &DocTemplateDB{db: msgpackDb}
	udb.baseFilePath = assetDir

	example := &model.TemplateItem{}

	udb.db.Init(example)
	initVars(udb.db)
	//udb.db.ReIndex(example)
	var fVersion int
	verr := udb.db.Get(docTemplVersion, docTemplVersion, &fVersion)
	if verr == nil && fVersion == 0 {
		//do the checks... and upgrade
		log.Println("upgrade document template structure", fVersion, "->", example.GetVersion())
		//udb.db.ReIndex(example)
		type Template struct {
			Lang        string `json:"lang,omitempty"`
			Filename    string `json:"name,omitempty"`
			Size        int64  `json:"size,omitempty"`
			ContentType string `json:"contentType,omitempty"`
			Path        string `json:"path,omitempty"`
		}
		type TemplateItem struct {
			model.Permissions
			ID     string `json:"id" storm:"id"`
			Name   string `json:"name" storm:"index"`
			Detail string `json:"detail"`
			//Permissions Permissions `json:"permissions"`
			Updated time.Time            `json:"updated" storm:"index"`
			Created time.Time            `json:"created" storm:"index"`
			Data    map[string]*Template `json:"data"`
		}
		var oldItems []*TemplateItem
		err = udb.db.All(&oldItems)
		if err != nil && err != storm.ErrNotFound {
			return nil, err
		}
		for _, v := range oldItems {
			item := &model.TemplateItem{ID: v.ID, Name: v.Name, Detail: v.Detail, Updated: v.Updated, Created: v.Created}
			item.Permissions = v.Permissions
			if v.Data != nil {
				item.Data = model.TemplateLangMap{}
			}
			for lang, tmpl := range v.Data {
				fi := file.FromMap(udb.baseFilePath, map[string]interface{}{
					"path":        tmpl.Path,
					"name":        tmpl.Filename,
					"contentType": tmpl.ContentType,
					"size":        tmpl.Size,
				})
				item.Data[lang] = fi
			}
			err = udb.db.Save(item)
			if err != nil {
				return nil, err
			}

		}
		err = udb.db.Set(docTemplVersion, docTemplVersion, example.GetVersion())
		if err != nil {
			return nil, err
		}
	}

	return udb, nil
}

func (me *DocTemplateDB) List(auth model.Authorization, contains string, options map[string]interface{}) ([]*model.TemplateItem, error) {
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

func (me *DocTemplateDB) Get(auth model.Authorization, id string) (*model.TemplateItem, error) {
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
func (me *DocTemplateDB) ProvideFileInfoFor(auth model.Authorization, id, lang string, fm *file.Meta) (*file.IO, error) {
	return me.getFileInfoFor(auth, id, lang, fm)
}

func (me *DocTemplateDB) PutVars(auth model.Authorization, id, lang string, vars []string) error {
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

func (me *DocTemplateDB) GetTemplate(auth model.Authorization, id, lang string) (*file.IO, error) {
	return me.getFileInfoFor(auth, id, lang, nil)
}

func (me *DocTemplateDB) getFileInfoFor(auth model.Authorization, id, lang string, fm *file.Meta) (*file.IO, error) {
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

func (me *DocTemplateDB) DeleteTemplate(auth model.Authorization, id, lang string) error {
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
				return os.Remove(p)
			}
		}
	}
	return os.ErrNotExist
}

func (me *DocTemplateDB) Put(auth model.Authorization, item *model.TemplateItem) error {
	return me.put(auth, item, true)
}

func (me *DocTemplateDB) put(auth model.Authorization, item *model.TemplateItem, updated bool) error {
	if auth == nil || item == nil {
		return os.ErrInvalid
	}
	if item.ID == "" {
		if !auth.AccessRights().AllowedToCreateEntities() {
			return model.ErrAuthorityMissing
		}
		u2 := uuid.NewV4()
		item.ID = u2.String()
		item.Permissions = model.Permissions{Owner: auth.UserID()}
		item.Created = time.Now()
		item.Updated = time.Now()
		return me.db.Save(item)
	} else {
		var existing model.TemplateItem
		err := me.db.One("ID", item.ID, &existing)
		if err == storm.ErrNotFound {
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

func (me *DocTemplateDB) Delete(auth model.Authorization, id string) error {
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
				_ = os.Remove(p)
			}
		}
	}

	err = tx.DeleteStruct(uitem)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (me *DocTemplateDB) Vars(auth model.Authorization, contains string, options map[string]interface{}) ([]string, error) {
	contains = regexp.QuoteMeta(contains)
	params := makeSimpleQuery(options)
	tx, err := me.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	return getVars(contains, params.limit, params.index, tx)
}

func (me *DocTemplateDB) Import(imex *Imex) error {
	err := me.init(imex)
	if err != nil {
		return err
	}
	for i := 0; true; i++ {
		items, err := imex.db.Template.List(imex.auth, "", map[string]interface{}{"index": i, "limit": 1000, "metaOnly": false})
		if err == nil && len(items) > 0 {
			for _, item := range items {
				if imex.skipExistingOnImport {
					_, err = imex.sysDB.Template.Get(imex.auth, item.ID)
					if err == nil {
						continue
					}
				}

				item.Permissions.UpdateUserID(imex.locatedSameUserWithDifferentID)

				err = imex.sysDB.Template.put(imex.auth, item, false)
				if err != nil {
					imex.processedEntry(imexTemplate, item.ID, err)
					continue
				}
				var fi *file.IO
				hadError := false
				for lang, tmplItem := range item.Data {
					fi, err = imex.sysDB.Template.GetTemplate(imex.auth, item.ID, lang)
					if err != nil {
						hadError = true
						imex.processedEntry(imexTemplate, item.ID, err)
						continue
					}
					_, err = tmplItem.CpTo(fi)
					if err != nil {
						hadError = true
						imex.processedEntry(imexTemplate, item.ID, err)
					}
				}
				if !hadError {
					imex.processedEntry(imexTemplate, item.ID, nil)
				}
			}
		} else {
			break
		}
	}
	return nil
}

const imexTemplate = "Template"

func (me *DocTemplateDB) init(imex *Imex) error {
	var err error
	if imex.db.Template == nil {
		imex.db.Template, err = NewDocTemplateDB(imex.dir)
	}
	return err
}

func (me *DocTemplateDB) Export(imex *Imex, id ...string) error {
	err := me.init(imex)
	if err != nil {
		return err
	}
	specificIds := len(id) > 0
	if len(id) == 1 {
		if imex.isProcessed(imexTemplate, id[0]) {
			return nil
		}
	}
	var tx storm.Node
	for i := 0; true; i++ {
		items, err := imex.sysDB.Template.List(imex.auth, "", map[string]interface{}{"include": id, "index": i, "limit": 1000, "metaOnly": false})
		if err == nil && len(items) > 0 {
			tx, err = imex.db.Template.db.Begin(true)
			if err != nil {
				return err
			}
			wg := &sync.WaitGroup{}
			wg.Add(1)
			fileCopyErrs := map[string]error{}
			go func() {
				for _, item := range items {
					var exportFile *os.File
					for _, tmplItem := range item.Data {
						exportFile, err = os.OpenFile(filepath.Join(imex.db.Template.baseFilePath, tmplItem.PathName()), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
						if os.IsExist(err) {
							err = nil
						}
						if err != nil {
							if exportFile != nil {
								_ = exportFile.Close()
							}
							fileCopyErrs[item.ID] = err
							continue
						}
						_, err = tmplItem.Read(exportFile)
						if err != nil {
							fileCopyErrs[item.ID] = err
						}

					}
				}
				wg.Done()
			}()

			for _, item := range items {
				if !imex.isProcessed(imexTemplate, item.ID) {
					err = tx.Save(item)
					if err != nil {
						imex.processedEntry(imexTemplate, item.ID, err)
						continue
					}
					item.Permissions.UserIdsMap(imex.neededUsers)
					if err != nil {
						imex.processedEntry(imexTemplate, item.ID, err)
						continue
					}
					imex.processedEntry(imexTemplate, item.ID, nil)
				}
			}
			err = tx.Commit()
			if err != nil {
				_ = tx.Rollback()
				return err
			}
			wg.Wait()
			for k, v := range fileCopyErrs {
				imex.processedEntry(imexTemplate, k, v)
			}
		} else {
			break
		}
		if specificIds {
			break
		}
	}
	return nil
}

func (me *DocTemplateDB) Close() error {
	return me.db.Close()
}
