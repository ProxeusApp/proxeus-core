package storm

import (
	"path/filepath"

	"github.com/ProxeusApp/proxeus-core/storage/db"

	"github.com/ProxeusApp/proxeus-core/sys/model"
	"github.com/ProxeusApp/proxeus-core/sys/validate"
)

type SettingsDB struct {
	jf db.JSONFile
}

func NewSettingsDB(baseDir string) (*SettingsDB, error) {
	baseDir = filepath.Join(baseDir, "settings")
	err := ensureDir(baseDir)
	if err != nil {
		return nil, err
	}
	return &SettingsDB{
		jf: db.JSONFile{FilePath: filepath.Join(baseDir, "main.json")},
	}, nil
}

func (se *SettingsDB) Put(s *model.Settings) error {
	err := validate.Struct(s)
	if err != nil {
		return err
	}
	return se.jf.Put(s)
}

func (se *SettingsDB) Get() (*model.Settings, error) {
	var s model.Settings
	err := se.jf.Get(&s)
	return &s, err
}

func (se *SettingsDB) Close() error { return nil }
