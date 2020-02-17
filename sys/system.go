package sys

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/patrickmn/go-cache"

	cfg "github.com/ProxeusApp/proxeus-core/main/config"
	"github.com/ProxeusApp/proxeus-core/main/handlers/blockchain"
	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/storage/database"
	"github.com/ProxeusApp/proxeus-core/storage/portable"
	"github.com/ProxeusApp/proxeus-core/sys/eio"
	"github.com/ProxeusApp/proxeus-core/sys/email"
	"github.com/ProxeusApp/proxeus-core/sys/model"
	"github.com/ProxeusApp/proxeus-core/sys/validate"
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
		TestMode    bool
		AllowHttp   bool
		DB          *storage.DBSet
		DS          *eio.DocumentServiceClient
		EmailSender email.EmailSender

		settingsDB                  storage.SettingsIF
		settingsInUse               model.Settings
		paymentListenerCancelFunc   context.CancelFunc
		signatureListenerCancelFunc context.CancelFunc
		tick                        *time.Ticker
		cache                       *cache.Cache
	}
)

func NewWithSettings(settingsFile string, initialSettings *model.Settings) (*System, error) {
	stngsDB, err := database.NewSettingsDB(settingsFile, initialSettings)
	if err != nil {
		return nil, err
	}
	me := &System{settingsDB: stngsDB}

	if strings.ToLower(initialSettings.TestMode) == "true" {
		me.TestMode = true
		email.TestMode = true
	}

	if strings.ToLower(initialSettings.AllowHttp) == "true" {
		me.AllowHttp = true
	}

	err = me.init(me.GetSettings())
	if err != nil {
		return nil, err
	}
	return me, err
}

func (me *System) init(stngs *model.Settings) error {
	stngs.DataDir, _ = filepath.Abs(stngs.DataDir)

	me.DS = &eio.DocumentServiceClient{Url: stngs.DocumentServiceUrl}
	me.settingsInUse.DocumentServiceUrl = stngs.DocumentServiceUrl

	var err error
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

	err = os.MkdirAll(stngs.DataDir, 0755)
	if err != nil {
		return err
	}

	me.DB, err = database.NewDBSet(me.settingsDB, stngs.DataDir)
	if err != nil {
		return err
	}
	me.settingsInUse.DataDir = stngs.DataDir

	cacheExpiry, err := time.ParseDuration(stngs.SessionExpiry)
	if err != nil {
		return err
	}
	me.cache = cache.New(cacheExpiry, 10*time.Minute)

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
	return stngs
}

func (me *System) PutSettings(stngs *model.Settings) error {
	err := me.settingsDB.Put(stngs)
	if err != nil {
		return err
	}
	return me.init(stngs)
}

func (me *System) Export(writer io.Writer, s *Session, entities []portable.EntityType) (portable.ProcessedResults, error) {
	dir := filepath.Join(os.TempDir(), s.GetSessionDir())
	defer os.RemoveAll(dir)
	ie, err := portable.NewImportExport(s, me.DB, dir)
	if err != nil {
		return nil, err
	}
	defer ie.Close()
	for _, entity := range entities {
		err = ie.Export(entity)
		if err != nil {
			return nil, err
		}
	}
	var f *os.File
	f, err = ie.Pack()
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
	return ie.Processed(), nil
}

func (me *System) ExportSingle(writer io.Writer, s *Session, entity portable.EntityType, id ...string) (portable.ProcessedResults, error) {
	dir := filepath.Join(os.TempDir(), s.GetSessionDir())
	defer os.RemoveAll(dir)
	ie, err := portable.NewImportExport(s, me.DB, dir)
	if err != nil {
		return nil, err
	}
	defer ie.Close()
	err = ie.Export(entity, id...)
	if err != nil {
		return nil, err
	}
	var f *os.File
	f, err = ie.Pack()
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
	return ie.Processed(), nil
}

func (me *System) Import(reader io.Reader, s *Session, skipExisting bool) (portable.ProcessedResults, error) {
	dir := filepath.Join(os.TempDir(), s.GetSessionDir())
	defer os.RemoveAll(dir)
	ie, err := portable.NewImportExport(s, me.DB, dir)
	if err != nil {
		return nil, err
	}
	defer ie.Close()
	ie.SetSkipExistingOnImport(skipExisting)
	err = ie.Extract(reader)
	if err != nil {
		return nil, err
	}
	toImport := []portable.EntityType{
		portable.Settings,
		portable.User, //User must be imported before entities with permissions
		portable.I18n,
		portable.Template,
		portable.Form,
		portable.Workflow,
		portable.UserData,
	}
	for _, entityType := range toImport {
		err = ie.Import(entityType)
		if err != nil {
			return nil, err
		}
	}
	return ie.Processed(), nil
}

func (me *System) GetSession(sid string) (*Session, error) {
	s, err := me.DB.Session.Get(sid)
	if err != nil {
		return nil, err
	}
	return &Session{S: s, db: me.DB, cache: me.cache}, err
}

func (me *System) NewSession(usr *model.User) (*Session, error) {
	return NewSession(me, usr)
}

func (me *System) closeDBs() {
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
