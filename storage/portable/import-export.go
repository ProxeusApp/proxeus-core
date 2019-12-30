package portable

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/sys/model"
	"github.com/ProxeusApp/proxeus-core/sys/tar"

	uuid "github.com/satori/go.uuid"
)

var ErrNotProxeusDB = fmt.Errorf("not a Proxeus DB file")

const proxeusIdentifier = "00000_e166801d00a45901e2b3ca692a6a95e367d4a976218b485546a2da464b6c88b5"

type ImportExport struct {
	mainDir                 string
	dir                     string
	auth                    model.Auth
	sysDB                   *storage.DBSet
	db                      *storage.DBSet
	skipExistingOnImport    bool
	processed               ProcessedResults
	neededUsers             map[string]bool
	exportingAllUsersAnyway bool

	//email and ethereum address are treated as ID's
	//if this map contains items, all other entities need to update there user references on import
	//------------------------------   old ID new ID
	locatedSameUserWithDifferentID map[string]string
}

type ProcessedResults map[EntityType]map[string]string

type importExportMeta struct {
	Version int
}

func NewImportExport(auth model.Auth, dbSet *storage.DBSet, dir string) (*ImportExport, error) {
	u := uuid.NewV4()
	ie := &ImportExport{
		sysDB:                          dbSet,
		db:                             &storage.DBSet{},
		processed:                      ProcessedResults{},
		neededUsers:                    map[string]bool{},
		locatedSameUserWithDifferentID: map[string]string{},
	}
	ie.mainDir = filepath.Join(dir, u.String())
	ie.dir = filepath.Join(ie.mainDir, "data")
	err := ensureDir(ie.dir)
	if err != nil {
		return nil, err
	}
	ie.auth = auth
	return ie, nil
}

func ensureDir(dir string) error {
	var err error
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0750)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ie *ImportExport) SetSkipExistingOnImport(yes bool) {
	ie.skipExistingOnImport = yes
}

func (ie *ImportExport) Pack() (*os.File, error) {
	var err error
	if !ie.exportingAllUsersAnyway && len(ie.neededUsers) > 0 {
		ids := make([]string, len(ie.neededUsers))
		i := 0
		for id := range ie.neededUsers {
			ids[i] = id
			i++
		}
		err = ie.Export(User, ids...)
		if err != nil {
			return nil, err
		}
	}
	var f *os.File
	err = ie.db.Close()
	if err != nil {
		return nil, err
	}
	err = ie.writeProxeusIdentifier(&importExportMeta{Version: 1})
	if err != nil {
		return nil, err
	}
	tarFile := filepath.Join(ie.mainDir, "export")
	f, err = os.OpenFile(tarFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, err
	}
	err = tar.Tar(ie.dir, nil, f)
	er := f.Close()
	if er != nil {
		return nil, er
	}
	if err != nil {
		return nil, err
	}
	return os.Open(tarFile)
}

func (ie *ImportExport) writeProxeusIdentifier(meta *importExportMeta) error {
	bts, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(ie.dir, proxeusIdentifier), bts, 0600)
}

func (ie *ImportExport) readProxeusIdentifier(meta *importExportMeta) error {
	f, err := os.Open(filepath.Join(ie.dir, proxeusIdentifier))
	if err != nil {
		return err
	}
	defer f.Close()
	bts, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bts, meta)
	return err
}

func (ie *ImportExport) isProcessed(kind EntityType, id string) bool {
	if m, ok := ie.processed[kind]; ok {
		if _, ok := m[id]; ok {
			return true
		}
	}
	return false
}

//clean this up to make the processed keys accurate but only if there is no error
func (ie *ImportExport) cleanI18nAllKeyMarkersWithoutErr() error {
	langs, err := ie.sysDB.I18n.GetAllLangs()
	if err != nil {
		return err
	}
	for _, lang := range langs {
		if m, ok := ie.processed[I18n]; ok {
			if errText, ok := m[lang.Code+allKeysMarker]; ok && errText == "" {
				delete(m, lang.Code+allKeysMarker)
			}
		}
	}
	return nil
}

func (ie *ImportExport) processedEntry(kind EntityType, id string, err error) {
	strErr := ""
	if err != nil {
		strErr = err.Error()
	}
	if m, ok := ie.processed[kind]; ok {
		m[id] = strErr
	} else {
		ie.processed[kind] = map[string]string{id: strErr}
	}
}

func (ie *ImportExport) Processed() ProcessedResults {
	_ = ie.cleanI18nAllKeyMarkersWithoutErr()
	return ie.processed
}

func (ie *ImportExport) Extract(reader io.Reader) error {
	err := tar.Untar(ie.dir, reader)
	if gzip.ErrHeader == err {
		return ErrNotProxeusDB
	}
	var meta importExportMeta
	err = ie.readProxeusIdentifier(&meta)
	if err != nil {
		return ErrNotProxeusDB
	}
	return nil
}

func (ie *ImportExport) Close() error {
	err := ie.db.Close()
	if err != nil {
		return err
	}
	return os.RemoveAll(ie.mainDir)
}
