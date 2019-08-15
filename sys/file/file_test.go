package file

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

type (
	F struct {
		Path        string `json:"path"`
		ContentType string `json:"contentType"`
		Size        int64  `json:"size"`
	}
)

func TestFileIO(t *testing.T) {
	baseDir := "./testdir"
	fileData := "hello"
	file2Data := "hello2"
	err := ensureDir(baseDir)
	if err != nil {
		t.Error(err)
		return
	}
	err = ioutil.WriteFile(filepath.Join(baseDir, "abc"), []byte(fileData), 0600)
	if err != nil {
		t.Error(err)
		return
	}
	err = ioutil.WriteFile(filepath.Join(baseDir, "abc2"), []byte(file2Data), 0600)
	if err != nil {
		t.Error(err)
		return
	}
	var f *os.File
	var n int64
	f, err = os.Open(filepath.Join(baseDir, "abc"))
	//write
	fi := New(baseDir, Meta{Name: "Abc", ContentType: "SomeCT", Size: 123})
	_, err = fi.Write(f)
	var bts []byte
	bts, err = fi.MarshalJSON()
	if err != nil {
		t.Error(err)
		return
	}
	testInF := F{}
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
	testInF = F{}
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
	fi2, err := FromJSONBytes(baseDir, bts)
	if err != nil {
		t.Error(err)
		return
	}
	buf := &bytes.Buffer{}
	n, err = fi2.Read(buf)
	if err != nil {
		t.Error(n, err)
		return
	}
	if buf.String() != fileData {
		t.Error(err)
	}

	//overwrite
	bts, err = fi2.MarshalJSON()
	testF := F{}
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
	testF = F{}
	err = json.Unmarshal(bts, &testF)
	if err != nil {
		t.Error(err)
		return
	}

	//check New path
	if testF.Path != fi2.PathName() {
		t.Error(err, testF.Path, filepath.Join(baseDir, fi2.PathName()))
	}

	f, err = os.Open(filepath.Join(baseDir, "abc2"))
	if err != nil {
		t.Error(err)
		return
	}
	_, err = fi.Write(f)
	if err != nil {
		t.Error(err)
		return
	}

	buf = &bytes.Buffer{}
	n, err = fi.Read(buf)
	if err != nil {
		t.Error(n, err)
		return
	}
	if buf.String() != file2Data {
		t.Error(err, buf.String())
	}
	os.RemoveAll(baseDir)
}

func ensureDir(dir string) error {
	var err error
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0750)
		if err != nil {
			return err
		}
	}
	return nil
}
