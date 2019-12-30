package storm

import (
	"os"
	"path/filepath"

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

func (me *SettingsDB) Close() error {
	return me.jfdb.Close()
}
