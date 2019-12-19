package db

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type (
	JSONFileDB struct {
		ReadDir                          string
		WriteDir                         string
		FileSuffix                       string
		openDBFiles                      map[string]*rwjfDB
		rwMutex                          sync.RWMutex
		CloseImmediatelyAfterReadOrWrite bool
		Beautify                         bool
		DisableJSONEncoding              bool
	}

	rwjfDB struct {
		read  *jfDB
		write *jfDB
	}

	jfDB struct {
		jfdbMain *JSONFileDB
		filePath string
		file     *os.File
		mutex    sync.Mutex
	}
)

// NewJSONFileDB creates an instance of JSONFileDB
// readDir cannot be empty!
// if writeDir is empty readDir is taken for writes
// fileSuffix is an extension of the end of the file name
func NewJSONFileDB(readDir, writeDir, fileSuffix string, closeImmediatelyAfterReadOrWrite bool) (*JSONFileDB, error) {
	if readDir == "" {
		return nil, errors.New("readDir cannot be empty")
	}
	return &JSONFileDB{
		ReadDir:                          readDir,
		WriteDir:                         writeDir,
		FileSuffix:                       fileSuffix,
		openDBFiles:                      make(map[string]*rwjfDB),
		CloseImmediatelyAfterReadOrWrite: closeImmediatelyAfterReadOrWrite,
	}, nil
}

func (jfdb *JSONFileDB) Put(name string, d interface{}) error {
	rwjfDB, err := jfdb.getRWDB(name, true)
	if err != nil {
		return err
	}
	if rwjfDB != nil {
		if rwjfDB.write != nil {
			err = rwjfDB.write.put(d)
		} else {
			err = rwjfDB.read.put(d)
		}
	}
	return err
}

func (jfdb *JSONFileDB) Get(name string, d interface{}) error {
	rwjfDB, err := jfdb.getRWDB(name, false)
	if err != nil {
		return err
	}
	if rwjfDB != nil {
		err = rwjfDB.read.get(d)
	}
	return err
}

func (jfdb *JSONFileDB) Remove(name string) error {
	var internalJfDB *rwjfDB
	var err error
	jfdb.rwMutex.Lock()
	internalJfDB = jfdb.openDBFiles[name]
	if internalJfDB != nil {
		internalJfDB.read.file.Close()
		if internalJfDB.write != nil {
			internalJfDB.write.file.Close()
		}
	}
	delete(jfdb.openDBFiles, name)
	err = os.Remove(filepath.Join(jfdb.ReadDir, name+jfdb.FileSuffix))
	if jfdb.WriteDir != "" {
		err = os.Remove(filepath.Join(jfdb.WriteDir, name+jfdb.FileSuffix))
	}
	jfdb.rwMutex.Unlock()
	return err
}

func (jfdb *JSONFileDB) Close() error {
	jfdb.rwMutex.Lock()
	defer jfdb.rwMutex.Unlock()
	var err error
	for _, db := range jfdb.openDBFiles {
		err = db.Close()
	}
	jfdb.openDBFiles = make(map[string]*rwjfDB)
	return err
}

func (rwDB *rwjfDB) Close() error {
	if rwDB.read != nil && rwDB.read.file != nil {
		rwDB.read.file.Close()
	}
	if rwDB.write != nil && rwDB.write.file != nil {
		rwDB.write.file.Close()
	}
	return nil
}

func (jfdb *JSONFileDB) getRWDB(name string, write bool) (*rwjfDB, error) {
	var internalJfDB *rwjfDB
	jfdb.rwMutex.RLock()
	internalJfDB = jfdb.openDBFiles[name]
	jfdb.rwMutex.RUnlock()
	if internalJfDB == nil {
		internalJfDB = &rwjfDB{}
		fp := filepath.Join(jfdb.ReadDir, name+jfdb.FileSuffix)
		f, err := mustOpenFile(fp, write)
		if write {
			if err != nil {
				return nil, err
			}
		}
		internalJfDB.read = &jfDB{jfdbMain: jfdb, filePath: fp, file: f}
		if jfdb.WriteDir != "" {
			wfp := filepath.Join(jfdb.WriteDir, name+jfdb.FileSuffix)
			f, err = mustOpenFile(wfp, write)
			if write {
				if err != nil {
					return nil, err
				}
			}
			internalJfDB.write = &jfDB{jfdbMain: jfdb, filePath: wfp, file: f}
		}
		jfdb.rwMutex.Lock()
		jfdb.openDBFiles[name] = internalJfDB
		jfdb.rwMutex.Unlock()
	}
	return internalJfDB, nil
}

func (jfdb *jfDB) put(d interface{}) error {
	//ioutil.WriteFile
	var b []byte
	var err error
	b, ok := d.([]byte)
	if !ok {
		s, ok := d.(string)
		if ok {
			b = []byte(s)
		}
	}
	if !jfdb.jfdbMain.DisableJSONEncoding {
		if !ok {
			if jfdb.jfdbMain.Beautify {
				b, err = json.MarshalIndent(d, "", "    ")
			} else {
				b, err = json.Marshal(d)
			}
			if err != nil {
				return err
			}
		} else if jfdb.jfdbMain.Beautify {
			bts := new(bytes.Buffer)
			err = json.Indent(bts, b, "", "    ")
			if err == nil {
				b = bts.Bytes()
			}
		}
	}

	err = jfdb.checkAndReopen(true)
	if err != nil {
		return err
	}
	jfdb.mutex.Lock()
	if jfdb.jfdbMain.CloseImmediatelyAfterReadOrWrite {
		defer jfdb.closeFileAndSetToNil()
	}
	defer jfdb.mutex.Unlock()
	err = jfdb.file.Truncate(0)
	if err != nil {
		return err
	}
	_, err = jfdb.file.WriteAt(b, 0)
	if err != nil {
		return err
	}
	err = jfdb.file.Sync()
	return nil
}

func (jfdb *jfDB) get(d interface{}) error {
	err := jfdb.checkAndReopen(false)
	if err != nil {
		return err
	}
	jfdb.mutex.Lock()
	_, err = jfdb.file.Seek(0, 0)
	if err != nil {
		if jfdb.jfdbMain.CloseImmediatelyAfterReadOrWrite {
			jfdb.closeFileAndSetToNil()
		}
		jfdb.mutex.Unlock()
		return err
	}
	b, err := ioutil.ReadAll(jfdb.file)
	if err != nil {
		if jfdb.jfdbMain.CloseImmediatelyAfterReadOrWrite {
			jfdb.closeFileAndSetToNil()
		}
		jfdb.mutex.Unlock()
		return err
	}
	if jfdb.jfdbMain.CloseImmediatelyAfterReadOrWrite {
		jfdb.closeFileAndSetToNil()
	}
	jfdb.mutex.Unlock()
	if jfdb.jfdbMain.DisableJSONEncoding {
		pointer, ok := d.(*[]byte)
		if ok {
			*pointer = b
		} else {
			spointer, ok := d.(*string)
			if ok {
				*spointer = string(b)
			} else {
				err = errors.New("please provide at least a *[]byte or *string")
			}
		}
	} else {
		pointer, ok := d.(*[]byte)
		if ok {
			*pointer = b
		} else {
			err = json.Unmarshal(b, d)
		}
	}
	return err
}

func (jfdb *jfDB) checkAndReopen(write bool) (err error) {
	if jfdb.file == nil {
		jfdb.file, err = mustOpenFile(jfdb.filePath, write)
		if err != nil {
			return
		}
	}
	return
}

func (jfdb *jfDB) closeFileAndSetToNil() error {
	err := jfdb.file.Close()
	jfdb.file = nil
	return err
}

func mustOpenFile(path string, create bool) (file *os.File, err error) {
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		if create {
			parentDir := filepath.Dir(path)
			_, err = os.Stat(parentDir)
			if os.IsNotExist(err) {
				err = os.MkdirAll(parentDir, 0750)
				if err != nil {
					return
				}
			}
			file, err = os.Create(path)
			if err != nil {
				return
			}
		}
	} else {
		file, err = os.OpenFile(path, os.O_RDWR, 0644)
		if err != nil {
			return
		}
	}
	return
}
