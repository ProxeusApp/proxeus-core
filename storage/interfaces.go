package storage

import (
	"io"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/ProxeusApp/proxeus-core/sys/file"
	"github.com/ProxeusApp/proxeus-core/sys/model"
)

type DBSet struct {
	Settings          SettingsIF
	I18n              I18nIF
	Form              FormIF
	Workflow          WorkflowIF
	Template          TemplateIF
	User              UserIF
	UserData          UserDataIF
	SignatureRequests SignatureRequestsIF
	WorkflowPayments  WorkflowPaymentsIF

	//TODO(mmal): remove
	reflCache     map[string]ImporterExporter
	reflCacheLock sync.Mutex
}

type SettingsIF interface {
	ImporterExporter
	Put(stngs *model.Settings) error
	Get() (*model.Settings, error)
	Close() error
}

type I18nIF interface {
	ImporterExporter
	Find(keyContains string, valueContains string, options map[string]interface{}) (map[string]map[string]string, error)
	Get(lang string, key string, args ...string) (string, error)
	GetInsert(lang string, key string, args ...string) (string, error)
	GetAll(lang string) (map[string]string, error)
	PutAll(lang string, translations map[string]string) error
	Put(lang string, key string, text string) error
	Delete(keyContains string) error
	PutLang(lang string, enabled bool) error
	GetLangs(enabled bool) ([]*model.Lang, error)
	HasLang(lang string) bool
	GetAllLangs() ([]*model.Lang, error)
	PutFallback(lang string) error
	GetFallback() (string, error)
	Close() error
}

type FormIF interface {
	ImporterExporter
	List(auth model.Auth, contains string, options map[string]interface{}) ([]*model.FormItem, error)
	Get(auth model.Auth, id string) (*model.FormItem, error)
	Put(auth model.Auth, item *model.FormItem) error
	Delete(auth model.Auth, id string) error
	DelComp(auth model.Auth, id string) error
	PutComp(auth model.Auth, comp *model.FormComponentItem) error
	GetComp(auth model.Auth, id string) (*model.FormComponentItem, error)
	ListComp(auth model.Auth, contains string, options map[string]interface{}) (map[string]*model.FormComponentItem, error)
	Vars(auth model.Auth, contains string, options map[string]interface{}) ([]string, error)
	Close() error
}

type WorkflowIF interface {
	ImporterExporter
	ListPublished(auth model.Auth, contains string, options map[string]interface{}) ([]*model.WorkflowItem, error)
	List(auth model.Auth, contains string, options map[string]interface{}) ([]*model.WorkflowItem, error)
	GetPublished(auth model.Auth, id string) (*model.WorkflowItem, error)
	Get(auth model.Auth, id string) (*model.WorkflowItem, error)
	GetList(auth model.Auth, id []string) ([]*model.WorkflowItem, error)
	Put(auth model.Auth, item *model.WorkflowItem) error
	Delete(auth model.Auth, id string) error
	Close() error
}

type TemplateIF interface {
	ImporterExporter
	List(auth model.Auth, contains string, options map[string]interface{}) ([]*model.TemplateItem, error)
	Get(auth model.Auth, id string) (*model.TemplateItem, error)
	ProvideFileInfoFor(auth model.Auth, id, lang string, fm *file.Meta) (*file.IO, error)
	PutVars(auth model.Auth, id, lang string, vars []string) error
	GetTemplate(auth model.Auth, id, lang string) (*file.IO, error)
	DeleteTemplate(auth model.Auth, id, lang string) error
	Put(auth model.Auth, item *model.TemplateItem) error
	Delete(auth model.Auth, id string) error
	Vars(auth model.Auth, contains string, options map[string]interface{}) ([]string, error)
	Close() error
}

type UserIF interface {
	ImporterExporter
	GetBaseFilePath() string
	Login(name, pw string) (*model.User, error)
	Count() (int, error)
	List(auth model.Auth, contains string, options map[string]interface{}) ([]*model.User, error)
	Get(auth model.Auth, id string) (*model.User, error)
	GetByBCAddress(bcAddress string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	UpdateEmail(id, email string) error
	Put(auth model.Auth, item *model.User) error
	PutPw(id, pass string) error
	GetProfilePhoto(auth model.Auth, id string, writer io.Writer) (n int64, err error)
	PutProfilePhoto(auth model.Auth, id string, reader io.Reader) (written int64, err error)
	APIKey(key string) (*model.User, error)
	CreateApiKey(auth model.Auth, userId, apiKeyName string) (string, error)
	DeleteApiKey(auth model.Auth, userId, hiddenApiKey string) error
	Close() error
}

type UserDataIF interface {
	ImporterExporter
	List(auth model.Auth, contains string, options map[string]interface{}, includeReadGranted bool) ([]*model.UserDataItem, error)
	Delete(auth model.Auth, id string) error
	Get(auth model.Auth, id string) (*model.UserDataItem, error)
	GetAllFileInfosOf(ud *model.UserDataItem) []*file.IO
	GetByWorkflow(auth model.Auth, wf *model.WorkflowItem, finished bool) (*model.UserDataItem, bool, error)
	GetData(auth model.Auth, id, dataPath string) (interface{}, error)
	GetDataAndFiles(auth model.Auth, id, dataPath string) (interface{}, []string, error)
	PutData(auth model.Auth, id string, dataObj map[string]interface{}) error
	NewFile(auth model.Auth, meta file.Meta) *file.IO
	GetDataFile(auth model.Auth, id, dataPath string) (*file.IO, error)
	Put(auth model.Auth, item *model.UserDataItem) error
	Close() error
	Remove() error
}

type SignatureRequestsIF interface {
	GetBySignatory(ethAddr string) (*[]model.SignatureRequestItem, error)
	All() (*[]model.SignatureRequestItem, error)
	GetByID(docid string, docpath string) (*[]model.SignatureRequestItem, error)
	GetByHashAndSigner(hash string, signatory string) (*[]model.SignatureRequestItem, error)
	Add(item *model.SignatureRequestItem) error
	SetRejected(docid string, docpath string, signatory string) error
	SetRevoked(docid string, docpath string, signatory string) error
	List(auth model.Auth, contains string, options map[string]interface{}) ([]*model.UserDataItem, error)
	Close() error
}

type WorkflowPaymentsIF interface {
	GetByTxHashAndStatusAndFromEthAddress(txHash, status, from string) (*model.WorkflowPaymentItem, error)
	Get(paymentId string) (*model.WorkflowPaymentItem, error)
	ConfirmPayment(txHash, from, to string, xes uint64) error
	GetByWorkflowIdAndFromEthAddress(workflowID, fromEthAddr string, statuses []string) (*model.WorkflowPaymentItem, error)
	SetAbandonedToTimeoutBeforeTime(beforeTime time.Time) error
	Save(item *model.WorkflowPaymentItem) error
	Update(paymentId, status, txHash, from string) error
	Cancel(paymentId, from string) error
	Redeem(workflowId, from string) error
	Delete(paymentId string) error
	Remove(payment *model.WorkflowPaymentItem) error
	All() ([]*model.WorkflowPaymentItem, error)
	Close() error
}

type ImexIF interface {
	Auth() model.Auth
	SetSkipExistingOnImport(yes bool)
	SkipExistingOnImport() bool
	Pack() (*os.File, error)
	ProcessedEntry(kind, id string, err error)
	Processed() map[string]map[string]string
	IsProcessed(kind, id string) bool
	LocatedSameUserWithDifferentID() map[string]string
	NeededUsers() map[string]bool
	SetExportingAllUsersAnyway(b bool)
	Extract(reader io.Reader) error
	Dir() string
	DB() *DBSet
	SysDB() *DBSet
	Close() error
}

type ImporterExporter interface {
	Export(imex ImexIF, id ...string) error
	Import(imex ImexIF) error
}

type ImexMeta struct {
	Version int
}

func (db *DBSet) Close() error {
	cs := []io.Closer{
		db.Settings,
		db.I18n,
		db.Form,
		db.Workflow,
		db.Template,
		db.User,
		db.UserData,
		db.SignatureRequests,
		db.WorkflowPayments,
	}
	for _, c := range cs {
		if c == nil {
			continue
		}
		err := c.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (me *DBSet) ImexIFByName(name string) ImporterExporter {
	name = strings.ToLower(name)
	me.reflCacheLock.Lock()
	defer me.reflCacheLock.Unlock()
	if me.reflCache == nil {
		me.reflCache = map[string]ImporterExporter{}
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
				if imexIf, ok := fv.Interface().(ImporterExporter); ok {
					me.reflCache[name] = imexIf
					return imexIf
				}
			}
		}
	}
	me.reflCache[name] = nil
	return nil
}
