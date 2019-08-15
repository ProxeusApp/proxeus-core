package storm

import (
	"fmt"
	"os"
	"path/filepath"

	"git.proxeus.com/core/central/sys/db"
	"git.proxeus.com/core/central/sys/model"
	"git.proxeus.com/core/central/sys/validate"
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

func (me *SettingsDB) Import(imex *Imex) error {
	var err error
	imex.db.Settings, err = NewSettingsDB(imex.dir)
	if err != nil {
		imex.processedEntry(imexSettings, imexSettings, err)
		return nil
	}
	var s *model.Settings
	s, err = imex.db.Settings.Get()
	if err != nil { //does not exist
		return nil
	}
	if !imex.auth.AccessRights().IsGrantedFor(model.ROOT) {
		imex.processedEntry(imexSettings, imexSettings, fmt.Errorf("no authority to import"))
		//return nil to not break the export, just ignore the call
		return nil
	}
	err = imex.sysDB.Settings.Put(s)
	if err != nil {
		imex.processedEntry(imexSettings, imexSettings, err)
		return nil
	}
	imex.processedEntry(imexSettings, imexSettings, nil)
	return nil
}

const imexSettings = "Settings"

func (me *SettingsDB) Export(imex *Imex, id ...string) error {
	if !imex.auth.AccessRights().IsGrantedFor(model.ROOT) {
		imex.processedEntry(imexSettings, imexSettings, fmt.Errorf("no authority to export"))
		//return nil to not break the export, just ignore the call
		return nil
	}
	settings, err := me.Get()
	if err != nil {
		imex.processedEntry(imexSettings, imexSettings, err)
		return nil
	}
	var sdb *SettingsDB
	sdb, err = NewSettingsDB(imex.dir)
	if err != nil {
		imex.processedEntry(imexSettings, imexSettings, err)
		return nil
	}
	err = sdb.Put(settings)
	if err != nil {
		imex.processedEntry(imexSettings, imexSettings, err)
		return nil
	}
	sdb.Close()
	imex.processedEntry(imexSettings, imexSettings, nil)
	return nil
}

func (me *SettingsDB) Close() {
	_ = me.jfdb.Close()
}
