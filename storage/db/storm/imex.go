package storm

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

//Imex stands for import export, short just imex.
type Imex struct {
	mainDir                 string
	dir                     string
	auth                    model.Auth
	sysDB                   *storage.DBSet
	db                      *storage.DBSet
	skipExistingOnImport    bool
	processed               map[string]map[string]string
	neededUsers             map[string]bool
	exportingAllUsersAnyway bool

	//email and ethereum address are treated as ID's
	//if this map contains items, all other entities need to update there user references on import
	//------------------------------   old ID new ID
	locatedSameUserWithDifferentID map[string]string
}

func NewImex(auth model.Auth, dbSet *storage.DBSet, dir string) (*Imex, error) {
	u := uuid.NewV4()
	imex := &Imex{
		sysDB:                          dbSet,
		db:                             &storage.DBSet{},
		processed:                      map[string]map[string]string{},
		neededUsers:                    map[string]bool{},
		locatedSameUserWithDifferentID: map[string]string{},
	}
	imex.mainDir = filepath.Join(dir, u.String())
	imex.dir = filepath.Join(imex.mainDir, "data")
	err := ensureDir(imex.dir)
	if err != nil {
		return nil, err
	}
	imex.auth = auth
	return imex, nil
}

func (me *Imex) Auth() model.Auth {
	return me.auth
}

func (me *Imex) Dir() string {
	return me.dir
}

func (me *Imex) DB() *storage.DBSet {
	return me.db
}

func (me *Imex) SysDB() *storage.DBSet {
	return me.sysDB
}

func (me *Imex) SetSkipExistingOnImport(yes bool) {
	me.skipExistingOnImport = yes
}

func (me *Imex) SkipExistingOnImport() bool {
	return me.skipExistingOnImport
}

func (me *Imex) LocatedSameUserWithDifferentID() map[string]string {
	return me.locatedSameUserWithDifferentID
}

func (me *Imex) NeededUsers() map[string]bool {
	return me.neededUsers
}
func (me *Imex) SetExportingAllUsersAnyway(b bool) {
	me.exportingAllUsersAnyway = b
}

const proxeusIdentifier = "00000_e166801d00a45901e2b3ca692a6a95e367d4a976218b485546a2da464b6c88b5"

func (me *Imex) Pack() (*os.File, error) {
	var err error
	if !me.exportingAllUsersAnyway && len(me.neededUsers) > 0 {
		ids := make([]string, len(me.neededUsers))
		i := 0
		for id := range me.neededUsers {
			ids[i] = id
			i++
		}
		err = me.sysDB.User.Export(me, ids...)
		if err != nil {
			return nil, err
		}
	}
	var f *os.File
	err = me.db.Close()
	if err != nil {
		return nil, err
	}
	err = me.writeProxeusIdentifier(&storage.ImexMeta{Version: 1})
	if err != nil {
		return nil, err
	}
	tarFile := filepath.Join(me.mainDir, "export")
	f, err = os.OpenFile(tarFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, err
	}
	err = tar.Tar(me.dir, nil, f)
	er := f.Close()
	if er != nil {
		return nil, er
	}
	if err != nil {
		return nil, err
	}
	return os.Open(tarFile)
}

func (me *Imex) writeProxeusIdentifier(imexMeta *storage.ImexMeta) error {
	bts, err := json.Marshal(imexMeta)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(me.dir, proxeusIdentifier), bts, 0600)
}

func (me *Imex) readProxeusIdentifier(imexMeta *storage.ImexMeta) error {
	f, err := os.Open(filepath.Join(me.dir, proxeusIdentifier))
	if err != nil {
		return err
	}
	defer f.Close()
	bts, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bts, imexMeta)
	return err
}

func (me *Imex) IsProcessed(kind, id string) bool {
	if m, ok := me.processed[kind]; ok {
		if _, ok := m[id]; ok {
			return true
		}
	}
	return false
}

//clean this up to make the processed keys accurate but only if there is no error
func (me *Imex) cleanI18nAllKeyMarkersWithoutErr() error {
	langs, err := me.sysDB.I18n.GetAllLangs()
	if err != nil {
		return err
	}
	for _, lang := range langs {
		if m, ok := me.processed[imexI18n]; ok {
			if errText, ok := m[lang.Code+imexAllKeysMarker]; ok && errText == "" {
				delete(m, lang.Code+imexAllKeysMarker)
			}
		}
	}
	return nil
}

func (me *Imex) ProcessedEntry(kind, id string, err error) {
	strErr := ""
	if err != nil {
		strErr = err.Error()
	}
	if m, ok := me.processed[kind]; ok {
		m[id] = strErr
	} else {
		me.processed[kind] = map[string]string{id: strErr}
	}
}

func (me *Imex) Processed() map[string]map[string]string {
	_ = me.cleanI18nAllKeyMarkersWithoutErr()
	return me.processed
}

func (me *Imex) Extract(reader io.Reader) error {
	err := tar.Untar(me.dir, reader)
	if gzip.ErrHeader == err {
		return ErrNotProxeusDB
	}
	var imexMeta storage.ImexMeta
	err = me.readProxeusIdentifier(&imexMeta)
	if err != nil {
		return ErrNotProxeusDB
	}
	return nil
}

func (me *Imex) Close() error {
	err := me.db.Close()
	if err != nil {
		return err
	}
	return os.RemoveAll(me.mainDir)
}
