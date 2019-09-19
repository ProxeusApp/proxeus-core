package storm

import (
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/msgpack"
	"github.com/asdine/storm/q"
	uuid "github.com/satori/go.uuid"

	"git.proxeus.com/core/central/sys/file"
	"git.proxeus.com/core/central/sys/model"
)

type UserDataDB struct {
	db           *storm.DB
	baseFilePath string
	mainDir      string
}

const usrdHeavyData = "usrd_data"
const usrdVersion = "usrd_vers"
const usrdMainDir = "userdata"

func NewUserDataDB(dir string) (*UserDataDB, error) {
	var err error
	var msgpackDb *storm.DB
	baseDir := filepath.Join(dir, usrdMainDir)
	err = ensureDir(baseDir)
	if err != nil {
		return nil, err
	}
	assetDir := filepath.Join(baseDir, "assets")
	err = ensureDir(assetDir)
	if err != nil {
		return nil, err
	}
	msgpackDb, err = storm.Open(filepath.Join(baseDir, "usrdb"), storm.Codec(msgpack.Codec))
	if err != nil {
		return nil, err
	}
	udb := &UserDataDB{db: msgpackDb, mainDir: baseDir}
	udb.baseFilePath = assetDir

	example := &model.UserDataItem{}
	err = udb.db.Init(example)
	if err != nil {
		return nil, err
	}
	err = udb.db.ReIndex(example)
	if err != nil {
		return nil, err
	}
	var fVersion int
	verr := udb.db.Get(usrdVersion, usrdVersion, &fVersion)
	if verr == nil && fVersion != example.GetVersion() {
		log.Println("upgrade db", fVersion, "mem", example.GetVersion())
	}
	err = udb.db.Set(usrdVersion, usrdVersion, example.GetVersion())
	if err != nil {
		return nil, err
	}
	return udb, nil
}

func (me *UserDataDB) List(auth model.Authorization, contains string, options map[string]interface{}, includeReadGranted bool) ([]*model.UserDataItem, error) {
	params := makeSimpleQuery(options)
	items := make([]*model.UserDataItem, 0)
	tx, err := me.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	matchers := defaultMatcher(auth, contains, params, includeReadGranted)

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

func (me *UserDataDB) Delete(auth model.Authorization, id string) error {
	tx, err := me.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var item model.UserDataItem
	err = tx.One("ID", id, &item)
	if err != nil {
		return err
	}
	uitem := &item
	if !uitem.Permissions.IsWriteGrantedFor(auth) {
		return model.ErrAuthorityMissing
	}
	//error handling not important here
	tx.Get(usrdHeavyData, uitem.ID, &uitem.Data)
	if len(uitem.Data) > 0 {
		_, files := file.MapIO(uitem.Data).GetAllDataAndFiles(me.baseFilePath)
		for _, v := range files {
			err = os.Remove(v)
			if err != nil {
				return err
			}
		}
	}
	//error handling not important here
	tx.Delete(usrdHeavyData, uitem.ID)

	err = tx.DeleteStruct(&item)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (me *UserDataDB) Get(auth model.Authorization, id string) (*model.UserDataItem, error) {
	var item model.UserDataItem
	err := me.db.One("ID", id, &item)
	if err != nil {
		return nil, err
	}
	itemRef := &item
	if !itemRef.Permissions.IsReadGrantedFor(auth) {
		return nil, model.ErrAuthorityMissing
	}
	//error handling not important as we don't care whether data exists etc..
	me.db.Get(usrdHeavyData, itemRef.ID, &itemRef.Data)
	return itemRef, nil
}

func (me *UserDataDB) GetAllFileInfosOf(ud *model.UserDataItem) []*file.IO {
	m := file.MapIO(ud.Data)
	return m.GetAllFileInfos(me.baseFilePath)
}

func (me *UserDataDB) GetByWorkflow(auth model.Authorization, wf *model.WorkflowItem, finished bool) (*model.UserDataItem, bool, error) {
	var item model.UserDataItem
	matchers := defaultMatcher(auth, "", nil, true)
	matchers = append(matchers, q.And(q.Eq("WorkflowID", wf.ID), q.Eq("Finished", finished)))
	alreadyStarted := false
	err := me.db.Select(matchers...).OrderBy("Created").Reverse().First(&item)
	if err != nil {
		return nil, alreadyStarted, err
	}
	alreadyStarted = true
	if !item.Permissions.IsReadGrantedFor(auth) {
		return nil, alreadyStarted, model.ErrAuthorityMissing
	}
	err = me.db.Get(usrdHeavyData, item.ID, &item.Data)
	return &item, alreadyStarted, err
}

func (me *UserDataDB) GetData(auth model.Authorization, id, dataPath string) (interface{}, error) {
	tx, err := me.db.Begin(false)
	if err != nil {
		return nil, err
	}
	var item model.UserDataItem
	defer tx.Rollback()
	err = tx.One("ID", id, &item)
	if err != nil {
		return nil, err
	}
	if !item.Permissions.IsReadGrantedFor(auth) {
		return nil, model.ErrAuthorityMissing
	}
	err = tx.Get(usrdHeavyData, item.ID, &item.Data)
	if err != nil {
		return nil, err
	}
	if item.Data != nil {
		mapIO := file.MapIO(item.Data)
		mapIO.MakeFileInfos(me.baseFilePath)
		return mapIO.Get(dataPath), nil
	}
	return nil, os.ErrNotExist
}

func (me *UserDataDB) GetDataAndFiles(auth model.Authorization, id, dataPath string) (interface{}, []string, error) {
	tx, err := me.db.Begin(false)
	if err != nil {
		return nil, nil, err
	}
	var item model.UserDataItem
	defer tx.Rollback()
	err = tx.One("ID", id, &item)
	if err != nil {
		return nil, nil, err
	}
	if !item.Permissions.IsReadGrantedFor(auth) {
		return nil, nil, model.ErrAuthorityMissing
	}
	err = tx.Get(usrdHeavyData, item.ID, &item.Data)
	if err != nil {
		return nil, nil, err
	}
	if item.Data != nil {
		d := file.MapIO(item.Data).Get(dataPath)
		if d != nil {
			if dm, ok := d.(map[string]interface{}); ok {
				dat, files := file.MapIO(dm).GetAllDataAndFiles(me.baseFilePath)
				return dat, files, nil
			}
		}
		return d, nil, nil
	}
	return nil, nil, os.ErrNotExist
}

func (me *UserDataDB) PutData(auth model.Authorization, id string, dataObj map[string]interface{}) error {
	tx, err := me.db.Begin(true)
	if err != nil {
		return err
	}
	var item model.UserDataItem
	defer tx.Rollback()
	err = tx.One("ID", id, &item)
	if err != nil {
		return err
	}
	if !item.Permissions.IsWriteGrantedFor(auth) {
		return model.ErrAuthorityMissing
	}
	err = tx.Get(usrdHeavyData, item.ID, &item.Data)
	if err != nil && err != storm.ErrNotFound {
		return err
	}
	if item.Data == nil {
		item.Data = map[string]interface{}{}
	}
	file.MapIO(item.Data).MergeWith(dataObj)
	err = tx.Set(usrdHeavyData, item.ID, item.Data)
	if err != nil {
		return err
	}
	item.Data = nil
	item.Updated = time.Now()
	err = tx.Save(&item)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (me *UserDataDB) NewFile(auth model.Authorization, meta file.Meta) *file.IO {
	return file.New(me.baseFilePath, meta)
}

func (me *UserDataDB) GetDataFile(auth model.Authorization, id, dataPath string) (*file.IO, error) {
	tx, err := me.db.Begin(false)
	if err != nil {
		return nil, err
	}
	var item model.UserDataItem
	defer tx.Rollback()
	err = tx.One("ID", id, &item)
	if err != nil {
		return nil, err
	}
	if !item.Permissions.IsReadGrantedFor(auth) {
		return nil, model.ErrAuthorityMissing
	}
	err = tx.Get(usrdHeavyData, item.ID, &item.Data)
	if err != nil {
		return nil, err
	}
	fi := file.MapIO(item.Data).GetFileInfo(me.baseFilePath, dataPath)
	if fi != nil {
		return fi, nil
	}
	return nil, os.ErrNotExist
}

func (me *UserDataDB) Put(auth model.Authorization, item *model.UserDataItem) error {
	return me.put(auth, item, true)
}

func (me *UserDataDB) put(auth model.Authorization, item *model.UserDataItem, updated bool) error {
	if item == nil {
		return os.ErrInvalid
	}
	if item.ID == "" {
		if !auth.AccessRights().AllowedToCreateUserData() {
			return model.ErrAuthorityMissing
		}
		u2 := uuid.NewV4()
		item.ID = u2.String()
		item.Permissions = model.Permissions{Owner: auth.UserID()}
		item.Created = time.Now()
		item.Updated = item.Created
		tx, err := me.db.Begin(true)
		if err != nil {
			return err
		}
		defer tx.Rollback()
		return me.updateItem(auth, item, tx)
	} else {
		tx, err := me.db.Begin(true)
		if err != nil {
			return err
		}
		defer tx.Rollback()
		var existing model.UserDataItem
		err = tx.One("ID", item.ID, &existing)
		if err == storm.ErrNotFound {
			if !auth.AccessRights().AllowedToCreateUserData() {
				return model.ErrAuthorityMissing
			}
			if item.Permissions.Owner == "" {
				item.Permissions = model.Permissions{Owner: auth.UserID()}
				item.Updated = time.Now()
			}
			return me.updateItem(auth, item, tx)
		}
		if err != nil {
			return err
		}
		if existing.Permissions.IsWriteGrantedFor(auth) {
			item.Permissions = *existing.Permissions.Change(auth, &item.Permissions)
			if updated {
				item.Updated = time.Now()
			}
			return me.updateItem(auth, item, tx)
		} else {
			return model.ErrAuthorityMissing
		}
	}
}

func (me *UserDataDB) updateItem(auth model.Authorization, item *model.UserDataItem, tx storm.Node) error {
	err := me.saveOnly(item, tx)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (me *UserDataDB) saveOnly(item *model.UserDataItem, tx storm.Node) error {
	if item.Data != nil {
		err := tx.Set(usrdHeavyData, item.ID, item.Data)
		if err != nil {
			return err
		}
	}
	cp := *item
	cp.Data = nil
	return tx.Save(&cp)
}

func (me *UserDataDB) Import(imex *Imex) error {
	err := me.init(imex)
	if err != nil {
		return err
	}

	for i := 0; true; i++ {
		items, err := imex.db.UserData.List(imex.auth, "", map[string]interface{}{"index": i, "limit": 1000, "metaOnly": false}, true)
		if err == nil && len(items) > 0 {
			for _, item := range items {
				if imex.skipExistingOnImport {
					_, err = imex.sysDB.UserData.Get(imex.auth, item.ID)
					if err == nil {
						continue
					}
				}
				item.Permissions.UpdateUserID(imex.locatedSameUserWithDifferentID)

				err = imex.sysDB.UserData.put(imex.auth, item, false)
				if err != nil {
					imex.processedEntry(imexUserData, item.ID, err)
					continue
				}
				hadError := false
				fios := imex.db.UserData.GetAllFileInfosOf(item)
				for _, fio := range fios {
					f, err := os.OpenFile(filepath.Join(imex.sysDB.UserData.baseFilePath, fio.PathName()), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
					if err != nil {
						hadError = true
						imex.processedEntry(imexUserData, item.ID, err)
						continue
					}
					_, err = fio.Read(f)
					if err != nil {
						hadError = true
						imex.processedEntry(imexUserData, item.ID, err)
					}
					err = f.Close()
					if err != nil {
						hadError = true
						imex.processedEntry(imexUserData, item.ID, err)
					}
				}
				if !hadError {
					imex.processedEntry(imexUserData, item.ID, nil)
				}
			}
		} else {
			break
		}
	}
	return nil
}

const imexUserData = "UserData"

func (me *UserDataDB) init(imex *Imex) error {
	var err error
	if imex.db.UserData == nil {
		imex.db.UserData, err = NewUserDataDB(imex.dir)
	}
	return err
}

func (me *UserDataDB) Export(imex *Imex, id ...string) error {
	err := me.init(imex)
	if err != nil {
		return err
	}

	specificIds := len(id) > 0
	for i := 0; true; i++ {
		items, err := me.List(imex.auth, "", map[string]interface{}{"include": id, "index": i, "limit": 1000, "metaOnly": false}, true)
		if err == nil && len(items) > 0 {
			var tx storm.Node
			tx, err = imex.db.UserData.db.Begin(true)
			if err != nil {
				return err
			}
			workflows := map[string]bool{}
			wg := &sync.WaitGroup{}
			wg.Add(1)
			fileCopyErrs := map[string]error{}
			go func() {
				for _, item := range items {
					fios := me.GetAllFileInfosOf(item)
					for _, fio := range fios {
						f, err := os.OpenFile(filepath.Join(imex.db.UserData.baseFilePath, fio.PathName()), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
						if err != nil {
							fileCopyErrs[item.ID] = err
							continue
						}
						_, err = fio.Read(f)
						if err != nil {
							fileCopyErrs[item.ID] = err
						}
						err = f.Close()
						if err != nil {
							fileCopyErrs[item.ID] = err
						}
					}
				}
				wg.Done()
			}()
			for _, item := range items {
				if !imex.isProcessed(imexUserData, item.ID) {
					err = imex.db.UserData.saveOnly(item, tx)
					if err != nil {
						imex.processedEntry(imexUserData, item.ID, err)
						continue
					}
					if item.WorkflowID != "" {
						workflows[item.WorkflowID] = true
					}
					item.Permissions.UserIdsMap(imex.neededUsers)
					if err != nil {
						imex.processedEntry(imexUserData, item.ID, err)
						continue
					}
					imex.processedEntry(imexUserData, item.ID, nil)
				}
			}
			err = tx.Commit()
			if err != nil {
				_ = tx.Rollback()
				return err
			}
			if len(workflows) > 0 {
				wfIds := make([]string, len(workflows))
				i := 0
				for wfID := range workflows {
					wfIds[i] = wfID
					i++
				}
				err = imex.sysDB.Workflow.Export(imex, wfIds...)
				if err != nil {
					return err
				}
			}
			wg.Wait()
			for k, v := range fileCopyErrs {
				imex.processedEntry(imexUserData, k, v)
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

func (me *UserDataDB) Close() error {
	return me.db.Close()
}

func (me *UserDataDB) Remove() error {
	return os.RemoveAll(me.mainDir)
}
