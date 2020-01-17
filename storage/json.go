package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/google/renameio"
)

type JSONFile struct {
	FilePath string
}

func (j JSONFile) Put(d interface{}) error {
	b, err := json.MarshalIndent(d, "", "\t")
	if err != nil {
		return err
	}
	return renameio.WriteFile(j.FilePath, b, 0600)
}

func (j JSONFile) Get(d interface{}) error {
	f, err := os.Open(j.FilePath)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, d)
}
