package service

import (
	"bytes"
	"io"
	"net/http"

	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/sys/file"
	"github.com/ProxeusApp/proxeus-core/sys/model"

	"github.com/ProxeusApp/proxeus-core/sys/eio"
)

type (

	// TemplateDocumentService is an interface that provides template document functions
	TemplateDocumentService interface {

		// Compile compiles a template with the documentService
		Compile(template eio.Template) (*http.Response, error)

		// SaveTemplate saves the template
		SaveTemplate(auth model.Auth, id, lang, contentType, fileName string, data []byte) error

		// GetTemplateVars returns the template variables
		GetTemplateVars(auth model.Auth, contains string, settings storage.Options) ([]string, error)

		// DeleteTemplate remove a template
		DeleteTemplate(auth model.Auth, id, lang string) error

		// DeleteTemplateFiles removes template files
		DeleteTemplateFiles(auth model.Auth, id string) error

		// DownloadExtensions downloads the template assistance extension for your writer.
		DownloadExtension(os string) (resp *http.Response, err error)

		// ReadFile reads a file and writes it into the writer
		ReadFile(path string, writer io.Writer) error

		// GetTemplate returns a template file
		GetTemplate(auth model.Auth, id, lang string) (*file.IO, error)

		// Get returns a template item
		Get(auth model.Auth, id string) (*model.TemplateItem, error)

		// Put sets a template item
		Put(auth model.Auth, item *model.TemplateItem) error

		// List returns a list of template items
		List(auth model.Auth, contains string, options storage.Options) ([]*model.TemplateItem, error)

		// Exists returns whether a file under the specified path exists
		Exists(path string) (bool, error)
	}
	DefaultTemplateDocumentService struct {
	}
)

func NewTemplateDocumentService() *DefaultTemplateDocumentService {
	return &DefaultTemplateDocumentService{}
}

// Compile compiles a template with the documentService
func (me *DefaultTemplateDocumentService) Compile(template eio.Template) (*http.Response, error) {
	return me.ds().Compile(filesDB(), template)
}

// SaveTemplate saves the template
func (me *DefaultTemplateDocumentService) SaveTemplate(auth model.Auth, id, lang, contentType, fileName string, data []byte) error {
	tmplMeta := &file.Meta{
		Name:        fileName,
		ContentType: contentType,
	}

	fi, err := templateDB().ProvideFileInfoFor(auth, id, lang, tmplMeta) // Get and update disk info for existing file or create new, update db entry
	if err != nil {
		return err
	}

	err = filesDB().Write(fi.Path(), bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	var vars []string
	vars, err = me.ds().Vars(bytes.NewBuffer(data))
	if err == nil {
		//error handling not important here as keeping track of vars is not crucial
		templateDB().PutVars(auth, id, lang, vars)
	}
	return err
}

// GetTemplateVars returns the template variables
func (me *DefaultTemplateDocumentService) GetTemplateVars(auth model.Auth, contains string, settings storage.Options) ([]string, error) {
	return templateDB().Vars(auth, contains, settings)
}

// DeleteTemplate remove a template
func (me *DefaultTemplateDocumentService) DeleteTemplate(auth model.Auth, id, lang string) error {
	return templateDB().DeleteTemplate(auth, filesDB(), id, lang)
}

// DeleteTemplateFiles removes template files
func (me *DefaultTemplateDocumentService) DeleteTemplateFiles(auth model.Auth, id string) error {
	return templateDB().Delete(auth, filesDB(), id)
}

// DownloadExtensions downloads the template assistance extension for your writer.
func (me *DefaultTemplateDocumentService) DownloadExtension(os string) (resp *http.Response, err error) {
	return me.ds().DownloadExtension(os)
}

// ReadFile reads a file and writes it into the writer
func (me *DefaultTemplateDocumentService) ReadFile(path string, writer io.Writer) error {
	return filesDB().Read(path, writer)
}

// GetTemplate returns a template file
func (me *DefaultTemplateDocumentService) GetTemplate(auth model.Auth, id, lang string) (*file.IO, error) {
	return templateDB().GetTemplate(auth, id, lang)
}

// Get returns a template item
func (me *DefaultTemplateDocumentService) Get(auth model.Auth, id string) (*model.TemplateItem, error) {
	return templateDB().Get(auth, id)
}

// Put sets a template item
func (me *DefaultTemplateDocumentService) Put(auth model.Auth, item *model.TemplateItem) error {
	return templateDB().Put(auth, item)
}

// List returns a list of template items
func (me *DefaultTemplateDocumentService) List(auth model.Auth, contains string, options storage.Options) ([]*model.TemplateItem, error) {
	return templateDB().List(auth, contains, options)
}

// Exists returns whether a file under the specified path exists
func (me *DefaultTemplateDocumentService) Exists(path string) (bool, error) {
	return filesDB().Exists(path)
}

func (me *DefaultTemplateDocumentService) ds() *eio.DocumentServiceClient {
	return system.DS
}
