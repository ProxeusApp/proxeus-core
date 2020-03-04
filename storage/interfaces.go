// Package storage holds the storage layer for the Proxeus core
package storage

import (
	"bytes"
	"io"
	"time"

	"github.com/ProxeusApp/proxeus-core/externalnode"

	"github.com/ProxeusApp/proxeus-core/sys/file"
	"github.com/ProxeusApp/proxeus-core/sys/model"
)

// DBSet holds all references to the Proxeus Database implementations
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
	Files             FilesIF
	Session           SessionIF
}

// Options holds filter criteria for queries to the databases
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

// SettingsIF is the interface to the general Proxeus configuration file
type SettingsIF interface {
	// Put adds a new set of settings into the database / file
	Put(stngs *model.Settings) error

	// Get retrieves all settings from the database / file
	Get() (*model.Settings, error)

	// Close closes the database
	Close() error
}

// I18nIF is the interface to the translation database
type I18nIF interface {
	// Find retrieves a list of translations matching the supplied filter criteria
	Find(keyContains string, valueContains string, options Options) (map[string]map[string]string, error)

	// Get retrieves a specific translation item
	Get(lang string, key string, args ...string) (string, error)

	// GetAll returns the full set of translations
	GetAll(lang string) (map[string]string, error)

	// PutAll saves the supplied translations into the database
	PutAll(lang string, translations map[string]string) error

	// Put saves a single translation into the database
	Put(lang string, key string, text string) error

	// PutLang adds a new language to the translation database
	PutLang(lang string, enabled bool) error

	// GetLangs returns a list of languages
	GetLangs(enabled bool) ([]*model.Lang, error)

	// HasLang checks indicates if a language exists
	HasLang(lang string) bool

	// GetAllLangs returns the full list of all languages defined
	GetAllLangs() ([]*model.Lang, error)

	// PutFallback sets a specific language as default/fallback
	PutFallback(lang string) error

	// GetFallback returns the currently defined default/fallback language
	GetFallback() (string, error)

	// Close closes the database
	Close() error
}

// Form is the interface to the form database
type FormIF interface {
	// List returns the set of form items matching the supplied filter criteria
	List(auth model.Auth, contains string, options Options) ([]*model.FormItem, error)

	// Get returns one form item by its id
	Get(auth model.Auth, id string) (*model.FormItem, error)

	// Put inserts a form item
	Put(auth model.Auth, item *model.FormItem) error

	// Delete removes a form item from the database
	Delete(auth model.Auth, id string) error

	// DelComp removes a form component definition from the database
	DelComp(auth model.Auth, id string) error

	// PutComp saves a form component into the database
	PutComp(auth model.Auth, comp *model.FormComponentItem) error

	// GetComp retrieves a form component from the database
	GetComp(auth model.Auth, id string) (*model.FormComponentItem, error)

	// ListComp returns a list of form component matching the supplied filter criteria
	ListComp(auth model.Auth, contains string, options Options) (map[string]*model.FormComponentItem, error)

	// Vars returns a list of variable defined in a form
	Vars(auth model.Auth, contains string, options Options) ([]string, error)

	// Close closes the database
	Close() error
}

// WorkflowIF is the interface to the workflow database
type WorkflowIF interface {
	// ListPublished returns all workflow items matching the supplied filter options that are flagged as published
	ListPublished(auth model.Auth, contains string, options Options) ([]*model.WorkflowItem, error)

	// ListPublished returns all workflow items matching the supplied filter options
	List(auth model.Auth, contains string, options Options) ([]*model.WorkflowItem, error)

	// GetPublished returns a workflow item matching the supplied filter options that if it is flagged as published
	GetPublished(auth model.Auth, id string) (*model.WorkflowItem, error)

	// Get retrieves a worklfow item machting its id
	Get(auth model.Auth, id string) (*model.WorkflowItem, error)

	// GetList retrieves multiple workflows by matching their id
	GetList(auth model.Auth, id []string) ([]*model.WorkflowItem, error)

	//Put adds a workflow item into the database
	Put(auth model.Auth, item *model.WorkflowItem) error

	// Delete removes a workflow item from the database by matching its id
	Delete(auth model.Auth, id string) error

	// Close closes the database
	Close() error
	ExternalNodeIF
}

// ExternalNodeIF is the interface to the external node definition database
type ExternalNodeIF interface {
	// RegisterExternalNode saves an external node definition
	RegisterExternalNode(auth model.Auth, n *externalnode.ExternalNode) error

	// ListExternalNodes return a list of all external node definitions
	ListExternalNodes() []*externalnode.ExternalNode

	// DeleteExternalNode remove an external node definition
	DeleteExternalNode(auth model.Auth, id string) error

	// NodeByName retrieves an external node definition matching the supplied name
	NodeByName(auth model.Auth, name string) (*externalnode.ExternalNode, error)

	// QueryFromInstanceID return an external node instance by machting the specified id
	QueryFromInstanceID(auth model.Auth, id string) (externalnode.ExternalNodeInstance, error)

	// PutExternalNodeInstance saves an instance of an external node to the database
	PutExternalNodeInstance(auth model.Auth, i *externalnode.ExternalNodeInstance) error
}

// TemplateIF is the interface to the template database
type TemplateIF interface {

	// List fetches template items from the database based on the supplied filter options
	List(auth model.Auth, contains string, options Options) ([]*model.TemplateItem, error)

	// Get retrieves a single template item using its key
	Get(auth model.Auth, id string) (*model.TemplateItem, error)

	// ProvideFileInfoFor returns the fileinfo for a file
	ProvideFileInfoFor(auth model.Auth, id, lang string, fm *file.Meta) (*file.IO, error)

	// PutVars inserts a new variable into the databse
	PutVars(auth model.Auth, id, lang string, vars []string) error

	// GetTemplate returns the file template
	GetTemplate(auth model.Auth, id, lang string) (*file.IO, error)

	// DeleteTemplate removes a template from the database
	DeleteTemplate(auth model.Auth, files FilesIF, id, lang string) error

	//Put inserts a template item
	Put(auth model.Auth, item *model.TemplateItem) error

	// Delete removes a template's files on the file system
	Delete(auth model.Auth, files FilesIF, id string) error

	// Vars returns a list of variables defines for a specific template
	Vars(auth model.Auth, contains string, options Options) ([]string, error)

	// AssetsKeys return the base filepath of the document templates
	AssetsKey() string

	// Close closes the databases
	Close() error
}

// UserIF is the interface to the user database
type UserIF interface {

	// GetBaseFilePath returns the base file path
	GetBaseFilePath() string

	// Login tries to authenticate a user with the supplied credentials and returns the user object or an error
	Login(name, pw string) (*model.User, error)

	// Count returns the user count
	Count() (int, error)

	// List returns references to all the user object matching the supplied filter criteria
	List(auth model.Auth, contains string, options Options) ([]*model.User, error)

	// Get return a specific user object by machting its id
	Get(auth model.Auth, id string) (*model.User, error)

	// GetByBCAddress return a specific user object by matching the ethereum address
	GetByBCAddress(bcAddress string) (*model.User, error)

	// GetByEmail return a specific user object by matching the email address
	GetByEmail(email string) (*model.User, error)

	// UpdateEmail sets a new email address for a specific user id
	UpdateEmail(id, email string) error

	// Put saves a user object into the database
	Put(auth model.Auth, item *model.User) error

	// PutPw sets a new password for a specific user id
	PutPw(id, pass string) error

	// GetProfilePhoto returns a users photo
	GetProfilePhoto(auth model.Auth, id string, writer io.Writer) error

	// PutProfilePhoto sets a new photo for a specific user
	PutProfilePhoto(auth model.Auth, id string, reader io.Reader) error

	// APIKey tries to authenticate the user with the supplied API key and returns the user object or an error
	APIKey(key string) (*model.User, error)

	// CreateApiKey saves and returns a newly created random api key for a user
	CreateApiKey(auth model.Auth, userId, apiKeyName string) (string, error)

	// DeleteApiKey removes an existing API key
	DeleteApiKey(auth model.Auth, userId, hiddenApiKey string) error

	// Close closes the database
	Close() error
}

// UserDataIF is the interface to the user's workflow data database
type UserDataIF interface {
	// List returns all user data item matching the supplied filter criteria
	List(auth model.Auth, contains string, options Options, includeReadGranted bool) ([]*model.UserDataItem, error)

	// Delete removes user data and its accociates files from the database
	Delete(auth model.Auth, files FilesIF, id string) error

	// Get returns a specific user data item by matching the id
	Get(auth model.Auth, id string) (*model.UserDataItem, error)

	// GetAllFileInfosOf returns the associated file objects of a user data item
	GetAllFileInfosOf(ud *model.UserDataItem) []*file.IO

	// GetByWorkflow returns a the user data item by matching a specific workflow item
	GetByWorkflow(auth model.Auth, wf *model.WorkflowItem, finished bool) (*model.UserDataItem, bool, error)

	// GetData returns the data object for retrieving specific data from a data item
	GetData(auth model.Auth, id, dataPath string) (interface{}, error)

	// GetDataAndFiles returns the data object and associated files
	GetDataAndFiles(auth model.Auth, id, dataPath string) (interface{}, []string, error)

	// PutData inserts a data object into the data item database
	PutData(auth model.Auth, id string, dataObj map[string]interface{}) error

	// NewFile return a handle for a new data item file based on the defined base path and specific file metadata
	NewFile(auth model.Auth, meta file.Meta) *file.IO

	// GetDataFile returns the files associated with a data item
	GetDataFile(auth model.Auth, id, dataPath string) (*file.IO, error)

	// Put saves a user data item into the database
	Put(auth model.Auth, item *model.UserDataItem) error

	// AssetsKey returns the base path of the data items associated file
	AssetsKey() string

	// Close closes the database
	Close() error
}

// UserDataIF is the interface to the signature requests database
type SignatureRequestsIF interface {
	// GetBySignatory returns the list of signature requests for a specific signatory
	GetBySignatory(ethAddr string) (*[]model.SignatureRequestItem, error)

	// GetByID returns the signature request item by its id
	GetByID(docid string, docpath string) (*[]model.SignatureRequestItem, error)

	// GetByHashAndSigner returns a list of signture requests for a specific file hash and signatory
	GetByHashAndSigner(hash string, signatory string) (*[]model.SignatureRequestItem, error)

	// Add saves a signature request into the database
	Add(item *model.SignatureRequestItem) error

	// SetRejected alters the status of a signature request to rejected
	SetRejected(docid string, docpath string, signatory string) error

	// SetRevoked alters the status of a signature request to revoked
	SetRevoked(docid string, docpath string, signatory string) error

	// Close closes the database
	Close() error
}

// WorkflowPaymentsIF is the interface to the workflow payment database
type WorkflowPaymentsIF interface {
	// GetByTxHashAndStatusAndFromEthAddress returns a workflow payment item by matching the supplied filter parameters
	GetByTxHashAndStatusAndFromEthAddress(txHash, status, from string) (*model.WorkflowPaymentItem, error)

	// Get returns a specific Workflow payment item matching its id
	Get(paymentId string) (*model.WorkflowPaymentItem, error)

	// ConfirmPayment sets the status of a workflow payment item to confirmed by trying to find a matching transaction hash and searching for pending or created items matching the supplied criteria
	ConfirmPayment(txHash, from, to string, xes uint64) error

	// GetByWorkflowIdAndFromEthAddress returns a workflow payment item by matching the supplied filter parameters
	GetByWorkflowIdAndFromEthAddress(workflowID, fromEthAddr string, statuses []string) (*model.WorkflowPaymentItem, error)

	// SetAbandonedToTimeoutBeforeTime updates the status of all payment items created before the specified time to status timeout
	SetAbandonedToTimeoutBeforeTime(beforeTime time.Time) error

	// Save add a workflow payment item to the database
	Save(item *model.WorkflowPaymentItem) error

	// Update sets the status and tx hash of created workflow items matching the supplied criteria to the supplied values
	Update(paymentId, status, txHash, from string) error

	// Cancel sets the status of a workflow payment item to cancelled for the item matching the supplied id and from address
	Cancel(paymentId, from string) error

	// Redeem sets the status of a workflow payment item to redeemed for the item matching the supplied id and from address
	Redeem(workflowId, from string) error

	// Delete sets the status of a workflow payment item to deleted by matching the supplied id
	Delete(paymentId string) error

	// Remove removes a workflow payment item
	Remove(payment *model.WorkflowPaymentItem) error

	// All returns a list of all workflow payment items from the database
	All() ([]*model.WorkflowPaymentItem, error)

	// Close closes the database
	Close() error
}

// FilesIF is the interface to a generic File
type FilesIF interface {

	// Read returns a file content to the supplied writer
	Read(path string, w io.Writer) error

	// Write writes a file content from the supplied reader
	Write(path string, r io.Reader) error

	// Exists checks whether a file with a specific path exists in the file database
	Exists(path string) (bool, error)

	//Delete removes a file from the file database
	Delete(path string) error

	// Close closes the database
	Close() error
}

// SessionIF is the interface to the session database
type SessionIF interface {

	// Get returns a session
	Get(sid string) (*model.Session, error)

	// Put inserts a session
	Put(s *model.Session) error

	// Delete removes a session
	Delete(s *model.Session) error

	// GetTokenRequest returns a token request
	GetTokenRequest(t model.TokenType, id string) (*model.TokenRequest, error)

	// PutTokenRequest inserts a token request
	PutTokenRequest(r *model.TokenRequest) error

	// DeleteTokenRequest removes a token request
	DeleteTokenRequest(r *model.TokenRequest) error

	// GetValue returns the value for the provided key
	GetValue(key string, v interface{}) error

	// PutValue sets the value for a key
	PutValue(key string, v interface{}) error

	// DeleteValue removes the value of a key
	DeleteValue(key string) error

	// Close closes the database
	Close() error
}

// Close ensures proper closing of all database handles
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
		db.Files,
		db.Session,
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

// CopyFileAcross copies content between databases
func CopyFileAcross(dstDb, srcDb FilesIF, dstPath, srcPath string) (int64, error) {
	var buf bytes.Buffer
	err := srcDb.Read(srcPath, &buf)
	if err != nil {
		return 0, err
	}
	l := buf.Len()
	return int64(l), dstDb.Write(dstPath, &buf)
}

func CopyFile(db FilesIF, dstPath, srcPath string) (int64, error) {
	return CopyFileAcross(db, db, dstPath, srcPath)
}

func FileSize(db FilesIF, path string) int64 {
	var buf bytes.Buffer
	db.Read(path, &buf)
	return int64(buf.Len())
}
