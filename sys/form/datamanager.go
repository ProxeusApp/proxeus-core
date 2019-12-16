package form

import (
	"io"
	"os"
	"sync"

	"github.com/ProxeusApp/proxeus-core/sys/file"
)

type DataManager interface {
	OnLoad()
	GetDataFile(formID, name string) (*file.IO, error)
	GetData(formID string) (map[string]interface{}, error)
	GetDataByPath(formID, dataPath string) (interface{}, error)
	Clear(formID string) error
	GetAllData() (dat map[string]interface{}, err error)
	GetAllDataFilePathNameOnly() (dat map[string]interface{}, files []string)
	PutData(formID string, dat map[string]interface{}) error
	PutDataWithoutMerge(formID string, dat map[string]interface{}) error
	PutDataFile(formID, name string, f file.Meta, reader io.Reader) (written int64, err error)
	Close() (err error)
}

type dataManager struct {
	sync.RWMutex
	DataCluster   map[string]file.MapIO
	BaseFilePath  string
	madeFileInfos bool
}

func NewDataManager(baseFilePath string) *dataManager {
	return &dataManager{
		BaseFilePath: baseFilePath,
		DataCluster:  make(map[string]file.MapIO),
	}
}

func (me *dataManager) OnLoad() {
	me.mkFileInfos()
}

func (me *dataManager) GetDataFile(formID, name string) (*file.IO, error) {
	if formID != "" && name != "" {
		me.RLock()
		if mapIO, ok := me.DataCluster[formID]; ok {
			fi := mapIO.GetFileInfo(me.BaseFilePath, name)
			if fi != nil {
				me.RUnlock()
				return fi, nil
			}
		}
		me.RUnlock()
		return nil, os.ErrNotExist
	}
	return nil, os.ErrInvalid
}

func (me *dataManager) GetData(formID string) (map[string]interface{}, error) {
	if formID != "" {
		me.RLock()
		defer me.RUnlock()
		return me.DataCluster[formID], nil
	}
	return nil, os.ErrInvalid
}

func (me *dataManager) Clear(formID string) error {
	if formID != "" {
		me.RLock()
		defer me.RUnlock()
		delete(me.DataCluster, formID)
		return nil
	}
	return os.ErrInvalid
}

func (me *dataManager) GetDataByPath(formID, dataPath string) (interface{}, error) {
	if formID != "" {
		me.RLock()
		defer me.RUnlock()
		if mapIo, ok := me.DataCluster[formID]; ok {
			return mapIo.Get(dataPath), nil
		}
		return nil, nil
	}
	return nil, os.ErrInvalid
}

/**
returns:
{
	"key":"value",
	"fileKey": {
		"name":"myfile.pdf",
		"contentType":"application/pdf",
		"size":123,
		"path":"{{BaseFilePath}}/76ed295d-92d8-41b1-83be-7078ea9a94a2"
	}
}
*/
func (me *dataManager) GetAllData() (dat map[string]interface{}, err error) {
	me.RLock()
	defer me.RUnlock()
	if len(me.DataCluster) > 0 {
		dat = make(map[string]interface{})
		for _, formDataMap := range me.DataCluster {
			for name, maybeFile := range formDataMap {
				if maybeFile != nil {
					if fileInfo, ok := maybeFile.(*file.IO); ok {
						dat[name] = fileInfo
						continue
					} else if fileMap, ok := maybeFile.(map[string]interface{}); ok {
						if file.IsFileInfo(fileMap) {
							dat[name] = file.FromMap(me.BaseFilePath, fileMap)
							continue
						}
					}
				}
				dat[name] = maybeFile
			}
		}
	}
	return
}

/**
returns:
dat = {
	"key":"value",
	"fileKey": {
		"name":"myfile.pdf",
		"contentType":"application/pdf",
		"size":123,
		"path":"76ed295d-92d8-41b1-83be-7078ea9a94a2"
	}
}
files = [
	"{{BaseFilePath}}/76ed295d-92d8-41b1-83be-7078ea9a94a2",
	...
]
*/
func (me *dataManager) GetAllDataFilePathNameOnly() (dat map[string]interface{}, files []string) {
	me.RLock()
	defer me.RUnlock()
	if len(me.DataCluster) > 0 {
		dat = make(map[string]interface{})
		files = make([]string, 0)
		for _, formDataMap := range me.DataCluster {
			for name, maybeFile := range formDataMap {
				if maybeFile != nil {
					if fileInfo, ok := maybeFile.(*file.IO); ok {
						files = append(files, fileInfo.Path())
						dat[name] = fileInfo.ToMap(true)
						continue
					}
				}
				dat[name] = maybeFile
			}
		}
	}
	return
}

func (me *dataManager) PutData(formID string, dat map[string]interface{}) error {
	if formID != "" {
		me.Lock()
		if mapIO, ok := me.DataCluster[formID]; ok {
			mapIO.MergeWith(dat)
		} else {
			mapIO := file.MapIO{}
			mapIO.MergeWith(dat)
			me.DataCluster[formID] = mapIO
		}
		me.Unlock()
		return nil
	}
	return os.ErrInvalid
}

func (me *dataManager) PutDataWithoutMerge(formID string, dat map[string]interface{}) error {
	if formID != "" {
		me.Lock()
		me.DataCluster[formID] = dat
		me.Unlock()
		return nil
	}
	return os.ErrInvalid
}

func (me *dataManager) PutDataFile(formID, name string, f file.Meta, reader io.Reader) (written int64, err error) {
	var existingFileInfo *file.IO
	existingFileInfo, err = me.GetDataFile(formID, name)
	me.Lock()

	var formDataMap map[string]interface{}
	formDataMap = me.DataCluster[formID]
	if formDataMap == nil {
		formDataMap = make(map[string]interface{})
		me.DataCluster[formID] = formDataMap
	}
	if os.IsNotExist(err) || existingFileInfo == nil {
		existingFileInfo = file.New(me.BaseFilePath, f)
		written, err = existingFileInfo.Write(reader)
		formDataMap[name] = existingFileInfo
		me.Unlock()
	} else {
		me.Unlock()
		existingFileInfo.Update(f.Name, f.ContentType)
		written, err = existingFileInfo.Write(reader)
	}
	return
}

func (me *dataManager) mkFileInfos() {
	if !me.madeFileInfos {
		for _, mapIO := range me.DataCluster {
			mapIO.MakeFileInfos(me.BaseFilePath)
		}
		me.madeFileInfos = true
	}
}

func (me *dataManager) Close() (err error) {
	return nil
}
