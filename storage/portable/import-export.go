package portable

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/storage/database"
	"github.com/ProxeusApp/proxeus-core/sys/model"
	"github.com/ProxeusApp/proxeus-core/sys/tar"
)

var ErrNotProxeusDB = fmt.Errorf("not a Proxeus DB file")

const proxeusIdentifier = "00000_e166801d00a45901e2b3ca692a6a95e367d4a976218b485546a2da464b6c88b5"

type ImportExport struct {
	mainDir                 string
	dataDir                 string
	settingsFile            string
	auth                    model.Auth
	sysDB                   *storage.DBSet
	db                      *storage.DBSet
	dbConfig                database.DBConfig
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

// NewImportExport creates a configuration to import and export database files
func NewImportExport(auth model.Auth, dbSet *storage.DBSet, dir string) (*ImportExport, error) {
	ie := &ImportExport{
		sysDB:                          dbSet,
		db:                             &storage.DBSet{},
		processed:                      ProcessedResults{},
		neededUsers:                    map[string]bool{},
		locatedSameUserWithDifferentID: map[string]string{},
	}
	ie.mainDir = filepath.Join(dir, "export")
	ie.dataDir = filepath.Join(ie.mainDir, "data")
	err := os.MkdirAll(ie.dataDir, 0750)
	if err != nil {
		return nil, err
	}
	ie.settingsFile = filepath.Join(ie.dataDir, "settings", "main.json")
	ie.auth = auth
	ie.dbConfig = database.DBConfig{Engine: "storm", Dir: ie.dataDir}
	return ie, err
}

// Initialized the Files DB
// Make sure the path exists before opening the db. E.g. in import first unpack the zip to the wanted path then open FileDB
func (ie *ImportExport) InitFilesDB() error {
	var err error
	ie.db.Files, err = database.NewFileDB(ie.dbConfig)
	return err
}

func (ie *ImportExport) SetSkipExistingOnImport(yes bool) {
	ie.skipExistingOnImport = yes
}

// Pack produces and return the handle to a tar file containing the export of the database files
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
	err = tar.Tar(ie.dataDir, f)
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
	return ioutil.WriteFile(filepath.Join(ie.dataDir, proxeusIdentifier), bts, 0600)
}

func (ie *ImportExport) readProxeusIdentifier(meta *importExportMeta) error {
	f, err := os.Open(filepath.Join(ie.dataDir, proxeusIdentifier))
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

// clean this up to make the processed keys accurate but only if there is no error
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
	err := tar.Untar(ie.dataDir, reader)
	if gzip.ErrHeader == err {
		return ErrNotProxeusDB
	}
	var meta importExportMeta
	err = ie.readProxeusIdentifier(&meta)
	if err != nil {
		log.Print("[ImportExport] Extract err: ", err.Error())
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
