package database

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/ProxeusApp/proxeus-core/storage/database/db"
)

type FileDB struct {
	db db.DB
}

type storedFile struct {
	ID   string `storm:"id"`
	Data []byte
}

func NewFileDB(c DBConfig) (*FileDB, error) {
	baseDir := path.Join(c.Dir, "file")
	err := ensureDir(baseDir)
	if err != nil {
		return nil, err
	}
	db, err := OpenDatabase(c.Engine, c.URI, filepath.Join(baseDir, "files"))
	if err != nil {
		return nil, err
	}
	return &FileDB{db: db}, nil
}

func (d *FileDB) Read(path string, w io.Writer) error {
	if path == "" {
		return os.ErrNotExist
	}
	var f storedFile
	err := d.db.Get("storedFile", path, &f)
	if err != nil {
		return err
	}
	_, err = w.Write(f.Data)
	return err
}

func (d *FileDB) Write(path string, r io.Reader) error {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	f := storedFile{ID: path, Data: buf}
	return d.db.Set("storedFile", f.ID, &f)
}

func (d *FileDB) Exists(path string) (bool, error) {
	var f storedFile
	err := d.db.Get("storedFile", path, &f)
	if NotFound(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
func (d *FileDB) Delete(path string) error {
	return d.db.Delete("storedFile", path)
}
func (d *FileDB) Close() error { return d.db.Close() }
