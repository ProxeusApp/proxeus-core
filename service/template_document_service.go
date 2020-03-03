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
	TemplateDocumentService interface {
		Compile(template eio.Template) (*http.Response, error)
		SaveTemplate(auth model.Auth, id, lang, contentType, fileName string, data []byte) error
		GetTemplateVars(auth model.Auth, contains string, settings storage.Options) ([]string, error)
		DeleteTemplate(auth model.Auth, id, lang string) error
		DeleteTemplateFiles(auth model.Auth, id string) error
		DownloadExtension(os string) (resp *http.Response, err error)
		ReadFile(path string, writer io.Writer) error
		GetTemplate(auth model.Auth, id, lang string) (*file.IO, error)
		Get(auth model.Auth, id string) (*model.TemplateItem, error)
		Put(auth model.Auth, item *model.TemplateItem) error
		List(auth model.Auth, contains string, options storage.Options) ([]*model.TemplateItem, error)
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

func (me *DefaultTemplateDocumentService) ds() *eio.DocumentServiceClient {
	return system.DS
}

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

func (me *DefaultTemplateDocumentService) GetTemplateVars(auth model.Auth, contains string, settings storage.Options) ([]string, error) {
	return templateDB().Vars(auth, contains, settings)
}

func (me *DefaultTemplateDocumentService) DeleteTemplate(auth model.Auth, id, lang string) error {
	return templateDB().DeleteTemplate(auth, filesDB(), id, lang)
}

func (me *DefaultTemplateDocumentService) DeleteTemplateFiles(auth model.Auth, id string) error {
	return templateDB().Delete(auth, filesDB(), id)
}

func (me *DefaultTemplateDocumentService) DownloadExtension(os string) (resp *http.Response, err error) {
	return me.ds().DownloadExtension(os)
}

func (me *DefaultTemplateDocumentService) ReadFile(path string, writer io.Writer) error {
	return filesDB().Read(path, writer)
}

func (me *DefaultTemplateDocumentService) GetTemplate(auth model.Auth, id, lang string) (*file.IO, error) {
	return templateDB().GetTemplate(auth, id, lang)
}

func (me *DefaultTemplateDocumentService) Get(auth model.Auth, id string) (*model.TemplateItem, error) {
	return templateDB().Get(auth, id)
}

func (me *DefaultTemplateDocumentService) Put(auth model.Auth, item *model.TemplateItem) error {
	return templateDB().Put(auth, item)
}

func (me *DefaultTemplateDocumentService) List(auth model.Auth, contains string, options storage.Options) ([]*model.TemplateItem, error) {
	return templateDB().List(auth, contains, options)
}

func (me *DefaultTemplateDocumentService) Exists(path string) (bool, error) {
	return filesDB().Exists(path)
}
