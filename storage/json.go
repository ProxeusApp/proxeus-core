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

// NewJSONFile creates a representation for a File including its path and file permission
func NewJSONFile(filePath string, perm os.FileMode) *JSONFile {
	return &JSONFile{filePath: filePath, perm: perm}
}

// Put inserts a new line into the JSON File
func (j *JSONFile) Put(d interface{}) error {
	b, err := json.MarshalIndent(d, "", "\t")
	if err != nil {
		return err
	}
	return renameio.WriteFile(j.filePath, b, j.perm)
}

// Get retrieves the current content of the JSON File
func (j *JSONFile) Get(d interface{}) error {
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
