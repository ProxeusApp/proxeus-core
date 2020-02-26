package service

import (
	"github.com/ProxeusApp/proxeus-core/sys/eio"
	"net/http"
)

type (
	TemplateDocumentService interface {
		Compile(template eio.Template) (*http.Response, error)
	}
	DefaultTemplateDocumentService struct {
	}
)

func NewTemplateDocumentService() *DefaultTemplateDocumentService {
	return &DefaultTemplateDocumentService{}
}

func (me *DefaultTemplateDocumentService) Compile(template eio.Template) (*http.Response, error) {
	return me.ds().Compile(filesDB(), template)
}

func (me *DefaultTemplateDocumentService) ds() *eio.DocumentServiceClient {
	return system.DS
}
