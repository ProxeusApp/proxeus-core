package database

import (
	"os"
	"path/filepath"
	"time"

	"github.com/ProxeusApp/proxeus-core/storage"

	"github.com/asdine/storm/q"
	uuid "github.com/satori/go.uuid"

	"github.com/ProxeusApp/proxeus-core/storage/database/db"
	"github.com/ProxeusApp/proxeus-core/sys/file"
	"github.com/ProxeusApp/proxeus-core/sys/model"
)

type UserDataDB struct {
	db           db.DB
	baseFilePath string
	mainDir      string
}

const usrdHeavyData = "usrd_data"
const usrdVersion = "usrd_vers"
const usrdMainDir = "userdata"

// NewUserDataDB returns a handle to the user data database
func NewUserDataDB(c DBConfig) (*UserDataDB, error) {
	baseDir := filepath.Join(c.Dir, usrdMainDir)
	db, err := db.OpenDatabase(c.Engine, c.URI, filepath.Join(baseDir, "usrdb"))
	if err != nil {
		return nil, err
	}
	udb := &UserDataDB{db: db, mainDir: baseDir,
		baseFilePath: filepath.Join(baseDir, "assets")}

	udb.db.Init(usrdHeavyData)

	example := &model.UserDataItem{}
	err = udb.db.Init(example)
	if err != nil {
		return nil, err
	}
	err = udb.db.ReIndex(example)
	if err != nil {
		return nil, err
	}
	err = udb.db.Set(usrdVersion, usrdVersion, example.GetVersion())
	if err != nil {
		return nil, err
	}
	return udb, nil
}

// List returns all user data item matching the supplied filter criteria
func (me *UserDataDB) List(auth model.Auth, contains string, options storage.Options, includeReadGranted bool) ([]*model.UserDataItem, error) {
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

// Delete removes user data and its accociates files from the database
func (me *UserDataDB) Delete(auth model.Auth, filesDB storage.FilesIF, id string) error {
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
			err = filesDB.Delete(v)
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

// Get returns a specific user data item by matching the id
func (me *UserDataDB) Get(auth model.Auth, id string) (*model.UserDataItem, error) {
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

// GetAllFileInfosOf returns the associated file objects of a user data item
func (me *UserDataDB) GetAllFileInfosOf(ud *model.UserDataItem) []*file.IO {
	m := file.MapIO(ud.Data)
	return m.GetAllFileInfos(me.baseFilePath)
}

// GetByWorkflow returns a the user data item by matching a specific workflow item
func (me *UserDataDB) GetByWorkflow(auth model.Auth, wf *model.WorkflowItem, finished bool) (*model.UserDataItem, bool, error) {
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

// GetData returns the data object for retrieving specific data from a data item
func (me *UserDataDB) GetData(auth model.Auth, id, dataPath string) (interface{}, error) {
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

// GetDataAndFiles returns the data object and associated files
func (me *UserDataDB) GetDataAndFiles(auth model.Auth, id, dataPath string) (interface{}, []string, error) {
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

// PutData inserts a data object into the data item database
func (me *UserDataDB) PutData(auth model.Auth, id string, dataObj map[string]interface{}) error {
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
	if err != nil && !db.NotFound(err) {
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

// NewFile return a handle for a new data item file based on the defined base path and specific file metadata
func (me *UserDataDB) NewFile(auth model.Auth, meta file.Meta) *file.IO {
	return file.New(me.baseFilePath, meta)
}

// GetDataAndFiles returns the files associated with a data item
func (me *UserDataDB) GetDataFile(auth model.Auth, id, dataPath string) (*file.IO, error) {
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

// Put saves a user data item into the database
func (me *UserDataDB) Put(auth model.Auth, item *model.UserDataItem) error {
	return me.put(auth, item, true)
}

func (me *UserDataDB) put(auth model.Auth, item *model.UserDataItem, updated bool) error {
	if item == nil {
		return os.ErrInvalid
	}
	if item.ID == "" {
		if !auth.AccessRights().AllowedToCreateUserData() {
			return model.ErrAuthorityMissing
		}
		item.ID = uuid.NewV4().String()
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
		if db.NotFound(err) {
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

func (me *UserDataDB) updateItem(auth model.Auth, item *model.UserDataItem, tx db.DB) error {
	err := me.saveOnly(item, tx)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (me *UserDataDB) saveOnly(item *model.UserDataItem, tx db.DB) error {
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

// AssetsKey returns the base path of the data items associated file
func (me *UserDataDB) AssetsKey() string {
	return me.baseFilePath
}

func (me *UserDataDB) Close() error {
	return me.db.Close()
}
