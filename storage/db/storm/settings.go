package storm

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ProxeusApp/proxeus-core/storage"

	"github.com/ProxeusApp/proxeus-core/storage/db"

	"github.com/ProxeusApp/proxeus-core/sys/model"
	"github.com/ProxeusApp/proxeus-core/sys/validate"
)

type SettingsDB struct {
	jfdb *db.JSONFileDB
}

func NewSettingsDB(baseDir string) (*SettingsDB, error) {
	baseDir = filepath.Join(baseDir, "settings")
	err := ensureDir(baseDir)
	if err != nil {
		return nil, err
	}
	jfdb, err := db.NewJSONFileDB(baseDir, "", ".json", true)
	if err != nil {
		return nil, err
	}
	jfdb.Beautify = true
	return &SettingsDB{jfdb: jfdb}, nil
}

func (me *SettingsDB) Put(stngs *model.Settings) error {
	err := validate.Struct(stngs)
	if err != nil {
		return err
	}
	err = me.jfdb.Put("main", stngs)
	if err != nil {
		return err
	}
	return nil
}

func (me *SettingsDB) Get() (*model.Settings, error) {
	var stngs model.Settings
	err := me.jfdb.Get("main", &stngs)
	if os.IsNotExist(err) {
		err = ensureDir(me.jfdb.ReadDir)
		if err != nil {
			return nil, err
		}
		err = me.jfdb.Get("main", &stngs)
	}
	if err != nil {
		return nil, err
	}
	return &stngs, nil
}

func (me *SettingsDB) Import(imex storage.ImexIF) error {
	var err error
	imex.DB().Settings, err = NewSettingsDB(imex.Dir())
	if err != nil {
		imex.ProcessedEntry(imexSettings, imexSettings, err)
		return nil
	}
	var s *model.Settings
	s, err = imex.DB().Settings.Get()
	if err != nil { //does not exist
		return nil
	}
	if !imex.Auth().AccessRights().IsGrantedFor(model.ROOT) {
		imex.ProcessedEntry(imexSettings, imexSettings, fmt.Errorf("no authority to import"))
		//return nil to not break the export, just ignore the call
		return nil
	}
	err = imex.SysDB().Settings.Put(s)
	if err != nil {
		imex.ProcessedEntry(imexSettings, imexSettings, err)
		return nil
	}
	imex.ProcessedEntry(imexSettings, imexSettings, nil)
	return nil
}

const imexSettings = "Settings"

func (me *SettingsDB) Export(imex storage.ImexIF, id ...string) error {
	if !imex.Auth().AccessRights().IsGrantedFor(model.ROOT) {
		imex.ProcessedEntry(imexSettings, imexSettings, fmt.Errorf("no authority to export"))
		//return nil to not break the export, just ignore the call
		return nil
	}
	settings, err := me.Get()
	if err != nil {
		imex.ProcessedEntry(imexSettings, imexSettings, err)
		return nil
	}
	var sdb *SettingsDB
	sdb, err = NewSettingsDB(imex.Dir())
	if err != nil {
		imex.ProcessedEntry(imexSettings, imexSettings, err)
		return nil
	}
	err = sdb.Put(settings)
	if err != nil {
		imex.ProcessedEntry(imexSettings, imexSettings, err)
		return nil
	}
	sdb.Close()
	imex.ProcessedEntry(imexSettings, imexSettings, nil)
	return nil
}

func (me *SettingsDB) Close() error {
	return me.jfdb.Close()
}
