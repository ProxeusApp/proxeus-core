package service

import (
	"github.com/ProxeusApp/proxeus-core/sys"
	"github.com/ProxeusApp/proxeus-core/sys/eio"
	"net/http"
)

type (
	TemplateDocumentService interface {
		Compile(template eio.Template) (*http.Response, error)
	}
	DefaultTemplateDocumentService struct {
		*baseService
	}
)

func NewTemplateDocumentService(system *sys.System) *DefaultTemplateDocumentService {
	return &DefaultTemplateDocumentService{&baseService{system: system}}
}

func (me *DefaultTemplateDocumentService) Compile(template eio.Template) (*http.Response, error) {
	return me.ds().Compile(me.filesDB(), template)
}

func (me *DefaultTemplateDocumentService) ds() *eio.DocumentServiceClient {
	return me.system.DS
}
