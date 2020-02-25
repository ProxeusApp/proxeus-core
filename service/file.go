package service

import (
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
	}
)

func NewFileService() *DefaultFileService {
	return &DefaultFileService{}
}

//read file in path into writer
func (me *DefaultFileService) Read(path string, w io.Writer) error {
	return filesDB().Read(path, w)
}

// Return the file by id and dataPath
func (me *DefaultFileService) GetDataFile(auth model.Auth, id, dataPath string) (*file.IO, error) {
	return userDataDB().GetDataFile(auth, id, dataPath)
}

// Return the data and files by id and dataPath
func (me *DefaultFileService) GetDataAndFiles(auth model.Auth, id, dataPath string) (interface{}, []string, error) {
	return userDataDB().GetDataAndFiles(auth, id, dataPath)
}
