package database

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/sys/model"
	"github.com/ProxeusApp/proxeus-core/sys/validate"
)

type SettingsDB struct {
	jf *storage.JSONFile
}

func resolveHomeDirectory(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}

	u, err := user.Current()
	if err != nil {
		return path
	}

	return filepath.Clean(filepath.Join(u.HomeDir, path[1:]))
}

// NewSettingsDB returns a handle to access the settings database / file
func NewSettingsDB(settingsFile string, initialSettings *model.Settings) (*SettingsDB, error) {
	path := resolveHomeDirectory(settingsFile)
	err := os.MkdirAll(filepath.Dir(path), 0750)
	if err != nil {
		return nil, err
	}

	sdb := &SettingsDB{
		jf: storage.NewJSONFile(path, 0600),
	}

	_, err = sdb.Get()
	if err != nil {
		if initialSettings == nil {
			initialSettings = model.NewDefaultSettings()
		}
		err = sdb.Put(initialSettings)
		if err != nil {
			return nil, err
		}
	}

	return sdb, nil
}

// Put adds a new set of settings into the database / file
func (se *SettingsDB) Put(s *model.Settings) error {
	err := validate.Struct(s)
	if err != nil {
		return err
	}
	return se.jf.Put(s)
}

// Get retrieves all settings from the database / file
func (se *SettingsDB) Get() (*model.Settings, error) {
	var s model.Settings
	err := se.jf.Get(&s)
	return &s, err
}

func (se *SettingsDB) Close() error { return nil }
