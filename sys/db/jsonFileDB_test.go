package db

import (
	"errors"
	"os"
	"testing"
)

type Flow struct {
	Test string
}

func TestJSONFileDB_IO(t *testing.T) {
	parentDir := "./embedded/"
	name := "test"
	_, err := os.Stat(parentDir)
	removeAfter := false
	if os.IsNotExist(err) {
		removeAfter = true
	}
	var jfdb *JSONFileDB
	defer func() {
		jfdb.Remove(name)
		jfdb.Close()
		if removeAfter {
			os.RemoveAll(parentDir)
		}
	}()
	jfdb, err = NewJSONFileDB(parentDir+"read", "", ".json", true)
	if err != nil {
		t.Error(err)
	}
	f := &Flow{Test: "my json"}
	err = jfdb.Put(name, f)
	if err != nil {
		t.Error(err)
	}
	var m Flow
	err = jfdb.Get(name, &m)
	if err != nil {
		t.Error(err)
	}
	if m.Test != f.Test {
		t.Error(errors.New("write value not the same as read value"))
	}
}
