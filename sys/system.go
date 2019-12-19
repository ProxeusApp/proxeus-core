package sys

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/ProxeusApp/proxeus-core/main/handlers/blockchain"
	"github.com/ProxeusApp/proxeus-core/storage"

	"log"

	cfg "github.com/ProxeusApp/proxeus-core/main/config"
	"github.com/ProxeusApp/proxeus-core/sys/cache"
	"github.com/ProxeusApp/proxeus-core/sys/email"
	"github.com/ProxeusApp/proxeus-core/sys/validate"

	"github.com/ProxeusApp/proxeus-core/storage/db/storm"
	"github.com/ProxeusApp/proxeus-core/sys/eio"
	"github.com/ProxeusApp/proxeus-core/sys/model"
	"github.com/ProxeusApp/proxeus-core/sys/session"
)

var (
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
		TestMode                    bool
		SessionMgmnt                *session.Manager
		DB                          *storage.DBSet
		DS                          *eio.DocumentServiceClient
		EmailSender                 email.EmailSender
		Cache                       *cache.UCache
		settingsDB                  storage.SettingsIF
		settingsInUse               model.Settings
		fallbackSettings            *model.Settings
		paymentListenerCancelFunc   context.CancelFunc
		signatureListenerCancelFunc context.CancelFunc
		tick                        *time.Ticker
	}
	sessionNotify struct {
		system *System
	}
)

func provideProxeusSettings() (storage.SettingsIF, error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}
	stngsDB, err := storm.NewSettingsDB(filepath.Join(u.HomeDir, ".proxeus"))
	if err != nil {
		return nil, err
	}
	return stngsDB, nil
}

func NewWithSettings(settings model.Settings) (*System, error) {
	stngsDB, err := provideProxeusSettings()
	if err != nil {
		return nil, err
	}
	me := &System{settingsDB: stngsDB, fallbackSettings: &settings}

	if strings.ToLower(settings.TestMode) == "true" {
		me.TestMode = true
	}

	err = me.init(me.GetSettings())
	if err != nil {
		return nil, err
	}
	return me, err
}

func (me *System) init(stngs *model.Settings) error {
	log.Printf("Init with settings: %#v\n", stngs)

	expiry, err := time.ParseDuration(stngs.SessionExpiry)
	if err != nil {
		expiry = time.Hour
	}

	me.DS = &eio.DocumentServiceClient{Url: stngs.DocumentServiceUrl}
	me.settingsInUse.DocumentServiceUrl = stngs.DocumentServiceUrl

	me.EmailSender, err = email.NewSparkPostEmailSender(stngs.SparkpostApiKey, stngs.EmailFrom)
	if err != nil {
		return err
	}
	me.settingsInUse.SparkpostApiKey = stngs.SparkpostApiKey

	if stngs.BlockchainNet == "ropsten" {
		cfg.Config.XESContractAddress = "0x84E0b37e8f5B4B86d5d299b0B0e33686405A3919"
		cfg.Config.EthClientURL = "https://ropsten.infura.io/v3/" + stngs.InfuraApiKey
		cfg.Config.EthWebSocketURL = "wss://ropsten.infura.io/ws/v3/" + stngs.InfuraApiKey
	} else if stngs.BlockchainNet == "mainnet" {
		cfg.Config.XESContractAddress = "0xa017ac5fac5941f95010b12570b812c974469c2c"
		cfg.Config.EthClientURL = "https://mainnet.infura.io/v3/" + stngs.InfuraApiKey
		cfg.Config.EthWebSocketURL = "wss://mainnet.infura.io/ws/v3/" + stngs.InfuraApiKey
	} else {
		log.Println("Wrong blockchain network: ", stngs.BlockchainNet)
	}

	cfg.Config.AirdropEnabled = stngs.AirdropEnabled
	cfg.Config.AirdropAmountEther = stngs.AirdropAmountEther
	cfg.Config.AirdropAmountXES = stngs.AirdropAmountXES

	me.closeDBs()
	var cacheExpiry time.Duration
	cacheExpiry, err = time.ParseDuration(stngs.CacheExpiry)
	if err != nil {
		cacheExpiry = time.Hour
	}

	_, err = os.Stat(stngs.DataDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(stngs.DataDir, 0755)
		if err != nil {
			return err
		}
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

	XESABI, err := abi.JSON(strings.NewReader(blockchain.XesMainTokenABI))
	if err != nil {
		panic(err)
	}

	log.Println("blockchain config ethURL: ", cfg.Config.EthClientURL)
	log.Println("blockchain config ethWebSocketURL: ", cfg.Config.EthWebSocketURL)

	xesAdapter := blockchain.NewAdapter(cfg.Config.XESContractAddress, XESABI)

	bcListenerPayment := blockchain.NewPaymentListener(xesAdapter, cfg.Config.EthWebSocketURL,
		cfg.Config.EthClientURL, me.DB.WorkflowPayments)
	ctxPay := context.Background()
	ctxPay, cancelPay := context.WithCancel(ctxPay)
	if me.paymentListenerCancelFunc != nil {
		me.paymentListenerCancelFunc()
	}
	me.paymentListenerCancelFunc = cancelPay
	go bcListenerPayment.Listen(ctxPay)

	ProxeusFSABI, err := abi.JSON(strings.NewReader(blockchain.ProxeusFSABI))
	if err != nil {
		panic(err)
	}

	bcListenerSignature := blockchain.NewSignatureListener(cfg.Config.EthWebSocketURL,
		cfg.Config.EthClientURL, stngs.BlockchainContractAddress, me.DB.SignatureRequests, me.DB.User, me.EmailSender, ProxeusFSABI, cfg.Config.PlatformDomain)
	ctxSig := context.Background()
	ctxSig, cancelSig := context.WithCancel(ctxPay)
	if me.signatureListenerCancelFunc != nil {
		me.signatureListenerCancelFunc()
	}
	me.signatureListenerCancelFunc = cancelSig
	go bcListenerSignature.Listen(ctxSig)

	if me.tick != nil {
		me.tick.Stop()
	}
	me.tick = time.NewTicker(time.Hour * 6)
	go me.scheduledCleanup(me.tick)

	return nil
}

func (me *System) scheduledCleanup(tick *time.Ticker) {
	for range tick.C {
		beforeTime := time.Now().AddDate(0, 0, -14)
		log.Println("[scheduler][workflowpaymentcleanup] Timing out abandoned payments from before ", beforeTime)
		err := me.DB.WorkflowPayments.SetAbandonedToTimeoutBeforeTime(beforeTime)
		if err != nil {
			log.Println("[scheduler][workflowpaymentcleanup] err: ", err.Error())
			continue
		}
		log.Println("[scheduler][workflowpaymentcleanup] Done")
	}
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

func (me *System) GetImexIFFor(fields []string) []storage.ImporterExporter {
	items := make([]storage.ImporterExporter, 0)
	for _, name := range fields {
		ex := me.DB.ImexIFByName(name)
		if ex != nil {
			items = append(items, ex)
		}
	}
	return items
}

func (me *System) Export(writer io.Writer, s *session.Session, imexIfs ...storage.ImporterExporter) (map[string]map[string]string, error) {
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

func (me *System) ExportSingle(writer io.Writer, s *session.Session, imexIfs storage.ImporterExporter, id ...string) (map[string]map[string]string, error) {
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
	imex.SetSkipExistingOnImport(skipExisting)
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
	err := me.DB.Close()
	if err != nil {
		log.Println("[system][closeDBs] err: ", err.Error())
	}
}

func (me *System) Shutdown() {
	fmt.Println("System Shutdown...")
	if me.settingsDB != nil {
		me.settingsDB.Close()
	}
	me.closeDBs()
}
