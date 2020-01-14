package file

import (
	"encoding/json"
	"log"
	"path/filepath"
	"testing"
)

type f struct {
	Path        string `json:"path"`
	ContentType string `json:"contentType"`
	Size        int64  `json:"size"`
}

func fromJSONBytes(baseFilePath string, bts []byte) (*IO, error) {
	f := &IO{baseFileDir: baseFilePath}
	err := f.UnmarshalJSON(bts)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func TestFileIO(t *testing.T) {
	baseDir := "./testdir"
	var err error

	//write
	fi := New(baseDir, Meta{Name: "Abc", ContentType: "SomeCT", Size: 123})
	var bts []byte
	bts, err = fi.MarshalJSON()
	if err != nil {
		t.Error(err)
		return
	}
	testInF := f{}
	err = json.Unmarshal(bts, &testInF)
	if err != nil {
		t.Error(err)
		return
	}

	//check New path
	if testInF.Path != fi.PathName() {
		t.Error(err, testInF.Path, fi.PathName())
	}

	bts, err = fi.MarshalJSON()
	if err != nil {
		t.Error(err)
		return
	}
	testInF = f{}
	err = json.Unmarshal(bts, &testInF)
	if err != nil {
		t.Error(err)
		return
	}

	//check Output path
	if testInF.Path != fi.PathName() {
		t.Error(err, testInF.Path)
	}

	//read
	log.Println(baseDir)
	fi2, err := fromJSONBytes(baseDir, bts)
	if err != nil {
		t.Error(err)
		return
	}

	//overwrite
	bts, err = fi2.MarshalJSON()
	testF := f{}
	err = json.Unmarshal(bts, &testF)
	if err != nil {
		t.Error(err)
		return
	}

	//check Output path
	if testF.Path != fi2.PathName() {
		t.Error(err, testF.Path)
	}

	//test switch from Output to New
	bts, err = fi2.MarshalJSON()
	testF = f{}
	err = json.Unmarshal(bts, &testF)
	if err != nil {
		t.Error(err)
		return
	}

	//check New path
	if testF.Path != fi2.PathName() {
		t.Error(err, testF.Path, filepath.Join(baseDir, fi2.PathName()))
	}
}
