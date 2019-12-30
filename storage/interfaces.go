package storage

import (
	"io"
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
}

type Options struct {
	Limit    int                    `json:"limit"`
	Index    int                    `json:"index"`
	Include  map[string]interface{} `json:"include"`
	Exclude  map[string]interface{} `json:"exclude"`
	MetaOnly bool                   `json:"metaOnly"`
}

func IndexOptions(i int) Options {
	return Options{
		Index: i,
		Limit: 1000,
	}
}

func (o Options) WithInclude(in []string) Options {
	o.Include = map[string]interface{}{}
	for i, v := range in {
		o.Include[v] = i
	}
	return o
}

type SettingsIF interface {
	Put(stngs *model.Settings) error
	Get() (*model.Settings, error)
	Close() error
}

type I18nIF interface {
	Find(keyContains string, valueContains string, options Options) (map[string]map[string]string, error)
	Get(lang string, key string, args ...string) (string, error)
	GetAll(lang string) (map[string]string, error)
	PutAll(lang string, translations map[string]string) error
	Put(lang string, key string, text string) error
	PutLang(lang string, enabled bool) error
	GetLangs(enabled bool) ([]*model.Lang, error)
	HasLang(lang string) bool
	GetAllLangs() ([]*model.Lang, error)
	PutFallback(lang string) error
	GetFallback() (string, error)
	Close() error
}

type FormIF interface {
	List(auth model.Auth, contains string, options Options) ([]*model.FormItem, error)
	Get(auth model.Auth, id string) (*model.FormItem, error)
	Put(auth model.Auth, item *model.FormItem) error
	Delete(auth model.Auth, id string) error
	DelComp(auth model.Auth, id string) error
	PutComp(auth model.Auth, comp *model.FormComponentItem) error
	GetComp(auth model.Auth, id string) (*model.FormComponentItem, error)
	ListComp(auth model.Auth, contains string, options Options) (map[string]*model.FormComponentItem, error)
	Vars(auth model.Auth, contains string, options Options) ([]string, error)
	AssetsKey() string
	Close() error
}

type WorkflowIF interface {
	ListPublished(auth model.Auth, contains string, options Options) ([]*model.WorkflowItem, error)
	List(auth model.Auth, contains string, options Options) ([]*model.WorkflowItem, error)
	GetPublished(auth model.Auth, id string) (*model.WorkflowItem, error)
	Get(auth model.Auth, id string) (*model.WorkflowItem, error)
	GetList(auth model.Auth, id []string) ([]*model.WorkflowItem, error)
	Put(auth model.Auth, item *model.WorkflowItem) error
	Delete(auth model.Auth, id string) error
	AssetsKey() string
	Close() error
}

type TemplateIF interface {
	List(auth model.Auth, contains string, options Options) ([]*model.TemplateItem, error)
	Get(auth model.Auth, id string) (*model.TemplateItem, error)
	ProvideFileInfoFor(auth model.Auth, id, lang string, fm *file.Meta) (*file.IO, error)
	PutVars(auth model.Auth, id, lang string, vars []string) error
	GetTemplate(auth model.Auth, id, lang string) (*file.IO, error)
	DeleteTemplate(auth model.Auth, id, lang string) error
	Put(auth model.Auth, item *model.TemplateItem) error
	Delete(auth model.Auth, id string) error
	Vars(auth model.Auth, contains string, options Options) ([]string, error)
	AssetsKey() string
	Close() error
}

type UserIF interface {
	GetBaseFilePath() string
	Login(name, pw string) (*model.User, error)
	Count() (int, error)
	List(auth model.Auth, contains string, options Options) ([]*model.User, error)
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
	AssetsKey() string
	Close() error
}

type UserDataIF interface {
	List(auth model.Auth, contains string, options Options, includeReadGranted bool) ([]*model.UserDataItem, error)
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
	AssetsKey() string
	Close() error
}

type SignatureRequestsIF interface {
	GetBySignatory(ethAddr string) (*[]model.SignatureRequestItem, error)
	GetByID(docid string, docpath string) (*[]model.SignatureRequestItem, error)
	GetByHashAndSigner(hash string, signatory string) (*[]model.SignatureRequestItem, error)
	Add(item *model.SignatureRequestItem) error
	SetRejected(docid string, docpath string, signatory string) error
	SetRevoked(docid string, docpath string, signatory string) error
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
