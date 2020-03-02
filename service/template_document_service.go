package service

import (
	"bytes"
	"io"
	"net/http"

	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/sys"
	"github.com/ProxeusApp/proxeus-core/sys/eio"
	"github.com/ProxeusApp/proxeus-core/sys/file"
	"github.com/ProxeusApp/proxeus-core/sys/model"
)

type (
	TemplateDocumentService interface {
		Compile(template eio.Template) (*http.Response, error)
		SaveTemplate(sess *sys.Session, id, lang, contentType, fileName string, data []byte) error
		GetTemplateVars(sess *sys.Session, contains string, settings storage.Options) ([]string, error)
		DeleteTemplate(sess *sys.Session, id, lang string) error
		AddTemplate(sess *sys.Session, id, lang, fileName string, dataReader io.Reader) error
		RemoveTemplate(sess *sys.Session, id, lang string) error
		DeleteTemplateFiles(sess *sys.Session, id string) error
		DownloadExtension(os string) (resp *http.Response, err error)
		ReadFile(path string, writer io.Writer) error
		GetTemplate(sess *sys.Session, id, lang string) (*file.IO, error)
		Get(sess *sys.Session, id string) (*model.TemplateItem, error)
		Put(sess *sys.Session, item *model.TemplateItem) error
		List(sess *sys.Session, contains string, options storage.Options) ([]*model.TemplateItem, error)
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

func ideDelete(id, lang string, sess *sys.Session) error {
	sess.Delete(id + lang)
	return sess.DeleteFile(id + lang)
}

func (me *DefaultTemplateDocumentService) SaveTemplate(sess *sys.Session, id, lang, contentType, fileName string, data []byte) error {
	tmplMeta := &file.Meta{
		Name:        fileName,
		ContentType: contentType,
	}

	fi, err := templateDB().ProvideFileInfoFor(sess, id, lang, tmplMeta) // Get and update disk info for existing file or create new, update db entry
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
		templateDB().PutVars(sess, id, lang, vars)
	}
	//remove pending file from the session
	err = ideDelete(id, lang, sess)
	if err != nil {
		return err
	}
	return nil
}

func (me *DefaultTemplateDocumentService) GetTemplateVars(sess *sys.Session, contains string, settings storage.Options) ([]string, error) {
	return templateDB().Vars(sess, contains, settings)
}

func (me *DefaultTemplateDocumentService) DeleteTemplate(sess *sys.Session, id, lang string) error {
	return templateDB().DeleteTemplate(sess, filesDB(), id, lang)
}

func (me *DefaultTemplateDocumentService) AddTemplate(sess *sys.Session, id, lang, fileName string, dataReader io.Reader) error {
	if fileName == "" {
		fileName = "unknown"
	}
	sess.Put(id+lang, fileName)
	err := sess.WriteFile(id+lang, dataReader)
	if err != nil {
		return err
	}
	sess.Put("activeTmpl"+id, lang)
	return nil
}

func (me *DefaultTemplateDocumentService) RemoveTemplate(sess *sys.Session, id, lang string) error {
	return ideDelete(id, lang, sess)
}

func (me *DefaultTemplateDocumentService) DeleteTemplateFiles(sess *sys.Session, id string) error {
	return templateDB().Delete(sess, filesDB(), id)
}

func (me *DefaultTemplateDocumentService) DownloadExtension(os string) (resp *http.Response, err error) {
	return me.ds().DownloadExtension(os)
}

func (me *DefaultTemplateDocumentService) ReadFile(path string, writer io.Writer) error {
	return filesDB().Read(path, writer)
}

func (me *DefaultTemplateDocumentService) GetTemplate(sess *sys.Session, id, lang string) (*file.IO, error) {
	return templateDB().GetTemplate(sess, id, lang)
}

func (me *DefaultTemplateDocumentService) Get(sess *sys.Session, id string) (*model.TemplateItem, error) {
	return templateDB().Get(sess, id)
}

func (me *DefaultTemplateDocumentService) Put(sess *sys.Session, item *model.TemplateItem) error {
	return templateDB().Put(sess, item)
}

func (me *DefaultTemplateDocumentService) List(sess *sys.Session, contains string, options storage.Options) ([]*model.TemplateItem, error) {
	return templateDB().List(sess, contains, options)
}

func (me *DefaultTemplateDocumentService) Exists(path string) (bool, error) {
	return filesDB().Exists(path)
}
