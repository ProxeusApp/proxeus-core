package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/google/renameio"
)

type JSONFile struct {
	filePath string
	perm     os.FileMode
}

func NewJSONFile(filePath string, perm os.FileMode) *JSONFile {
	return &JSONFile{filePath: filePath, perm: perm}
}

func (j JSONFile) Put(d interface{}) error {
	b, err := json.MarshalIndent(d, "", "\t")
	if err != nil {
		return err
	}
	return renameio.WriteFile(j.filePath, b, j.perm)
}

func (j JSONFile) Get(d interface{}) error {
	f, err := os.Open(j.filePath)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, d)
}
