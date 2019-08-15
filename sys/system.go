package sys

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"git.proxeus.com/core/central/sys/cache"
	"git.proxeus.com/core/central/sys/email"
	"git.proxeus.com/core/central/sys/validate"

	"log"

	"git.proxeus.com/core/central/sys/db/storm"
	"git.proxeus.com/core/central/sys/eio"
	"git.proxeus.com/core/central/sys/model"
	"git.proxeus.com/core/central/sys/session"
)

var (
	ErrAccessDenied = fmt.Errorf("access denied")

	ReadAllFile = func(path string) ([]byte, error) {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		return ioutil.ReadAll(f)
	}
)

type (
	System struct {
		Debug            bool
		SessionMgmnt     *session.Manager
		DB               *storm.DBSet
		DS               *eio.DocumentServiceClient
		EmailSender      email.EmailSender
		Cache            *cache.UCache
		settingsDB       *storm.SettingsDB
		settingsInUse    model.Settings
		fallbackSettings *model.Settings
	}
	sessionNotify struct {
		system *System
	}
)

func provideProxeusSettings() (*storm.SettingsDB, error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}
	var stngsDB *storm.SettingsDB
	stngsDB, err = storm.NewSettingsDB(filepath.Join(u.HomeDir, ".proxeus"))
	if err != nil {
		return nil, err
	}
	return stngsDB, nil
}

func New() (*System, error) {
	stngsDB, err := provideProxeusSettings()
	if err != nil {
		return nil, err
	}
	me := &System{settingsDB: stngsDB}

	err = me.init(me.GetSettings())
	if err != nil {
		return nil, err
	}
	return me, err
}

func NewWithSettings(settings model.Settings) (*System, error) {
	stngsDB, err := provideProxeusSettings()
	if err != nil {
		return nil, err
	}
	me := &System{settingsDB: stngsDB, fallbackSettings: &settings}

	err = me.init(me.GetSettings())
	if err != nil {
		return nil, err
	}
	return me, err
}

func (me *System) init(stngs *model.Settings) error {
	//err := validate.Struct(stngs)
	//if err != nil {
	//	return err
	//}
	var err error
	var expiry time.Duration
	expiry, err = time.ParseDuration(stngs.SessionExpiry)
	if err != nil {
		expiry = time.Hour
	}

	me.Debug = true //TODO read system settings from db
	if me.DS == nil || me.settingsInUse.DocumentServiceUrl != stngs.DocumentServiceUrl {
		me.DS = &eio.DocumentServiceClient{Url: stngs.DocumentServiceUrl}
		me.settingsInUse.DocumentServiceUrl = stngs.DocumentServiceUrl
	}

	if me.EmailSender == nil || me.settingsInUse.SparkpostApiKey != stngs.SparkpostApiKey {
		me.EmailSender, err = email.NewSparkPostEmailSender(stngs.SparkpostApiKey)
		if err != nil {
			return err
		}
		me.settingsInUse.SparkpostApiKey = stngs.SparkpostApiKey
	}
	if me.Cache == nil || me.settingsInUse.DataDir != stngs.DataDir || me.settingsInUse.CacheExpiry != stngs.CacheExpiry {
		me.closeDBs()
		var cacheExpiry time.Duration
		cacheExpiry, err = time.ParseDuration(stngs.CacheExpiry)
		if err != nil {
			cacheExpiry = time.Hour
		}
		config := cache.UCacheConfig{
			DiskStorePath: filepath.Join(stngs.DataDir, "cache"),
			StoreType:     cache.DiskCache,
			ExtendExpiry:  false,
			DefaultExpiry: cacheExpiry,
		}
		me.Cache, err = cache.NewUCache(config)
		if err != nil {
			return err
		}

		me.DB, err = storm.NewDBSet(me.settingsDB, stngs.DataDir)
		if err != nil {
			return err
		}
		sessionNotify := &sessionNotify{system: me}
		me.SessionMgmnt, err = session.NewManagerWithNotify(stngs.DataDir, expiry, sessionNotify)
		if err != nil {
			return err
		}
		me.settingsInUse.DataDir = stngs.DataDir
	}
	me.DB.SignatureRequestsDB, err = storm.NewSignatureDB(stngs.DataDir)
	me.DB.WorkflowPaymentsDB, err = storm.NewWorkflowPaymentDB(stngs.DataDir)
	return nil
}

func (me *System) Configured() (bool, error) {
	count, err := me.DB.User.Count()
	if err != nil {
		return false, err
	}
	var s *model.Settings
	s, err = me.settingsDB.Get()
	if err != nil {
		return false, err
	}
	//validate to ensure settings loaded from the disk are still valid otherwise force configuration
	err = validate.Struct(s)
	if err != nil {
		return false, nil
	}
	return s != nil && count > 0, nil
}

func (me *System) GetSettings() *model.Settings {
	stngs, _ := me.settingsDB.Get()
	if stngs == nil {
		if me.fallbackSettings != nil {
			stngs = me.fallbackSettings
		} else {
			stngs = model.NewDefaultSettings()
		}
	}
	return stngs
}

func (me *System) PutSettings(stngs *model.Settings) error {
	err := validate.Struct(stngs)
	if err != nil {
		return err
	}
	err = me.settingsDB.Put(stngs)
	if err != nil {
		return err
	}
	return me.init(stngs)
}

func (me *sessionNotify) OnSessionCreated(id string, s *session.Session) {
	log.Println("OnSessionCreated", s)
}

func (me *sessionNotify) OnSessionLoaded(id string, s *session.Session) {
	log.Println("OnSessionLoaded", s)
}

func (me *sessionNotify) OnSessionExpired(id string, s *session.Session) {
	log.Println("OnSessionExpired", s, id)
}

func (me *sessionNotify) OnSessionRemoved(id string) {
	log.Println("OnSessionRemoved", id)
}

func (me *System) GetImexIFFor(fields []string) []storm.ImexIF {
	items := make([]storm.ImexIF, 0)
	for _, name := range fields {
		ex := me.DB.ImexIFByName(name)
		if ex != nil {
			items = append(items, ex)
		}
	}
	return items
}

func (me *System) Export(writer io.Writer, s *session.Session, imexIfs ...storm.ImexIF) (map[string]map[string]string, error) {
	imex, err := storm.NewImex(s, me.DB, s.SessionDir())
	if err != nil {
		return nil, err
	}
	defer imex.Close()
	for _, ex := range imexIfs {
		err = ex.Export(imex)
		if err != nil {
			return nil, err
		}
	}
	var f *os.File
	f, err = imex.Pack()
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(writer, f)
	if err != nil {
		return nil, err
	}
	err = f.Close()
	if err != nil {
		return nil, err
	}
	return imex.Processed(), nil
}

func (me *System) ExportSingle(writer io.Writer, s *session.Session, imexIfs storm.ImexIF, id ...string) (map[string]map[string]string, error) {
	imex, err := storm.NewImex(s, me.DB, s.SessionDir())
	if err != nil {
		return nil, err
	}
	defer imex.Close()
	err = imexIfs.Export(imex, id...)
	if err != nil {
		return nil, err
	}
	var f *os.File
	f, err = imex.Pack()
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(writer, f)
	if err != nil {
		return nil, err
	}
	err = f.Close()
	if err != nil {
		return nil, err
	}
	return imex.Processed(), nil
}

func (me *System) Import(reader io.Reader, s *session.Session, skipExisting bool) (map[string]map[string]string, error) {
	imex, err := storm.NewImex(s, me.DB, s.SessionDir())
	if err != nil {
		return nil, err
	}
	defer imex.Close()
	imex.SkipExistingOnImport(skipExisting)
	err = imex.Extract(reader)
	if err != nil {
		return nil, err
	}
	err = me.DB.Settings.Import(imex)
	if err != nil {
		return nil, err
	}
	//User must be imported before entities with permissions
	err = me.DB.User.Import(imex)
	if err != nil {
		return nil, err
	}
	err = me.DB.I18n.Import(imex)
	if err != nil {
		return nil, err
	}

	err = me.DB.Template.Import(imex)
	if err != nil {
		return nil, err
	}
	err = me.DB.Form.Import(imex)
	if err != nil {
		return nil, err
	}
	err = me.DB.Workflow.Import(imex)
	if err != nil {
		return nil, err
	}
	err = me.DB.UserData.Import(imex)
	if err != nil {
		return nil, err
	}
	return imex.Processed(), nil
}

func (me *System) closeDBs() {
	if me.Cache != nil {
		me.Cache.Close()
		me.Cache = nil
	}
	if me.SessionMgmnt != nil {
		_ = me.SessionMgmnt.Close()
		me.SessionMgmnt = nil
	}
	if me.DB == nil {
		return
	}
	_ = me.DB.Close()
}

func (me *System) Shutdown() {
	fmt.Println("System Shutdown...")
	if me.settingsDB != nil {
		me.settingsDB.Close()
	}
	me.closeDBs()
}
