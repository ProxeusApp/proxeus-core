package service

import (
	"github.com/ProxeusApp/proxeus-core/sys"
	"io"
)

type (
	FileService interface {
		Read(path string, w io.Writer) error
	}

	DefaultFileService struct {
		*baseService
	}
)

func NewFileService(system *sys.System) *DefaultFileService {
	return &DefaultFileService{&baseService{system: system}}
}

func (me *DefaultFileService) Read(path string, w io.Writer) error {
	return me.filesDB().Read(path, w)
}
