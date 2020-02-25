package service

import (
	"github.com/ProxeusApp/proxeus-core/sys"
	"github.com/ProxeusApp/proxeus-core/sys/file"
	"github.com/ProxeusApp/proxeus-core/sys/model"
	"io"
)

type (
	FileService interface {
		Read(path string, w io.Writer) error
		GetDataFile(auth model.Auth, id, dataPath string) (*file.IO, error)
		GetDataAndFiles(auth model.Auth, id, dataPath string) (interface{}, []string, error)
	}

	DefaultFileService struct {
		*baseService
	}
)

func NewFileService(system *sys.System) *DefaultFileService {
	return &DefaultFileService{&baseService{system: system}}
}

//read file in path into writer
func (me *DefaultFileService) Read(path string, w io.Writer) error {
	return me.filesDB().Read(path, w)
}

// Return the file by id and dataPath
func (me *DefaultFileService) GetDataFile(auth model.Auth, id, dataPath string) (*file.IO, error) {
	return me.userDataDB().GetDataFile(auth, id, dataPath)
}

// Return the data and files by id and dataPath
func (me *DefaultFileService) GetDataAndFiles(auth model.Auth, id, dataPath string) (interface{}, []string, error) {
	return me.userDataDB().GetDataAndFiles(auth, id, dataPath)
}
