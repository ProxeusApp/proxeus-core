package db

import (
	"os"
	"path/filepath"

	"github.com/cznic/ql"
)

func NewSQLFileDB(filePath string) (*ql.DB, error) {
	filePath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		dirPath := filepath.Dir(filePath)
		_, err := os.Stat(dirPath)
		if os.IsNotExist(err) {
			err = os.MkdirAll(filepath.Dir(filePath), 0750)
			if err != nil {
				return nil, err
			}
		}
		//file, err := os.Create(filePath)
		//if err != nil {return nil, err}
		//file.Close()
	}

	opts := &ql.Options{
		CanCreate: true,
	}
	db, err := ql.OpenFile(filePath, opts)
	return db, err
}
