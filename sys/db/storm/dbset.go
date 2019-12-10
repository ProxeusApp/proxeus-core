package storm

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	uuid "github.com/satori/go.uuid"

	"github.com/ProxeusApp/proxeus-core/sys/model"
	"github.com/ProxeusApp/proxeus-core/sys/tar"
)

var ErrNotProxeusDB = fmt.Errorf("not a Proxeus DB file")

type DBSet struct {
	Settings            *SettingsDB
	I18n                *I18nDB
	Form                *FormDB
	Workflow            WorkflowDBInterface
	Template            *DocTemplateDB
	User                UserDBInterface
	UserData            *UserDataDB
	SignatureRequestsDB *SignatureRequestsDB
	WorkflowPaymentsDB  WorkflowPaymentsDBInterface
	reflCache           map[string]ImexIF
	reflCacheLock       sync.Mutex
}

//Imex stands for import export, short just imex.
type Imex struct {
	mainDir                 string
	dir                     string
	auth                    model.Authorization
	sysDB                   *DBSet
	db                      *DBSet
	skipExistingOnImport    bool
	processed               map[string]map[string]string
	neededUsers             map[string]bool
	exportingAllUsersAnyway bool

	//email and ethereum address are treated as ID's
	//if this map contains items, all other entities need to update there user references on import
	//------------------------------   old ID new ID
	locatedSameUserWithDifferentID map[string]string
}

type ImexIF interface {
	Export(imex *Imex, id ...string) error
	Import(imex *Imex) error
}

func NewImex(auth model.Authorization, dbSet *DBSet, dir string) (*Imex, error) {
	u := uuid.NewV4()
	imex := &Imex{
		sysDB:                          dbSet,
		db:                             &DBSet{},
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

func (me *Imex) SkipExistingOnImport(yes bool) {
	me.skipExistingOnImport = yes
}

const proxeusIdentifier = "00000_e166801d00a45901e2b3ca692a6a95e367d4a976218b485546a2da464b6c88b5"

type ImexMeta struct {
	Version int
}

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
	err = me.writeProxeusIdentifier(&ImexMeta{Version: 1})
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

func (me *Imex) writeProxeusIdentifier(imexMeta *ImexMeta) error {
	bts, err := json.Marshal(imexMeta)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(me.dir, proxeusIdentifier), bts, 0600)
}

func (me *Imex) readProxeusIdentifier(imexMeta *ImexMeta) error {
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

func (me *Imex) isProcessed(kind, id string) bool {
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

func (me *Imex) processedEntry(kind, id string, err error) {
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
	var imexMeta ImexMeta
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

func NewDBSet(sDB *SettingsDB, folderPath string) (me *DBSet, err error) {
	me = &DBSet{Settings: sDB}
	me.I18n, err = NewI18nDB(folderPath)
	if err != nil {
		return nil, err
	}
	me.Form, err = NewFormDB(folderPath)
	if err != nil {
		return nil, err
	}
	me.Template, err = NewDocTemplateDB(folderPath)
	if err != nil {
		return nil, err
	}
	me.Workflow, err = NewWorkflowDB(folderPath)
	if err != nil {
		return nil, err
	}
	me.User, err = NewUserDB(folderPath)
	if err != nil {
		return nil, err
	}
	me.UserData, err = NewUserDataDB(folderPath)
	if err != nil {
		return nil, err
	}
	me.SignatureRequestsDB, err = NewSignatureDB(folderPath)
	if err != nil {
		return nil, err
	}
	me.WorkflowPaymentsDB, err = NewWorkflowPaymentDB(folderPath)
	if err != nil {
		return nil, err
	}
	return
}

func (me *DBSet) ImexIFByName(name string) ImexIF {
	name = strings.ToLower(name)
	me.reflCacheLock.Lock()
	defer me.reflCacheLock.Unlock()
	if me.reflCache == nil {
		me.reflCache = map[string]ImexIF{}
	}
	if im, ok := me.reflCache[name]; ok {
		return im
	}
	v := reflect.Indirect(reflect.ValueOf(me))
	strType := v.Type()
	size := v.NumField()
	for i := 0; i < size; i++ {
		tf := strType.Field(i)
		if strings.ToLower(tf.Name) == name {
			fv := v.Field(i)
			if fv.IsValid() && fv.CanInterface() {
				if imexIf, ok := fv.Interface().(ImexIF); ok {
					me.reflCache[name] = imexIf
					return imexIf
				}
			}
		}
	}
	me.reflCache[name] = nil
	return nil
}

func (me *DBSet) Close() error {
	//no need to close settings
	var err error
	if me.I18n != nil {
		err = me.I18n.Close()
		if err != nil {
			return err
		}
		me.I18n = nil
	}
	if me.Workflow != nil {
		err = me.Workflow.Close()
		if err != nil {
			return err
		}
		me.Workflow = nil
	}
	if me.Template != nil {
		err = me.Template.Close()
		if err != nil {
			return err
		}
		me.Template = nil
	}
	if me.Form != nil {
		err = me.Form.Close()
		if err != nil {
			return err
		}
		me.Form = nil
	}
	if me.User != nil {
		err = me.User.Close()
		if err != nil {
			return err
		}
		me.User = nil
	}
	if me.UserData != nil {
		err = me.UserData.Close()
		if err != nil {
			return err
		}
		me.UserData = nil
	}
	if me.SignatureRequestsDB != nil {
		err = me.SignatureRequestsDB.Close()
		if err != nil {
			return err
		}
		me.SignatureRequestsDB = nil
	}
	if me.WorkflowPaymentsDB != nil {
		err = me.WorkflowPaymentsDB.Close()
		if err != nil {
			return err
		}
		me.WorkflowPaymentsDB = nil
	}
	return nil
}
