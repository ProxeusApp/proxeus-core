package db

import (
	"io/ioutil"
	"os"
	"testing"
)

func openStormDB(t *testing.T, path string) DB {
	db, err := OpenDatabase("storm", "", path)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestStormCRUD(t *testing.T) {
	f, _ := ioutil.TempFile("", "test_db_")
	defer os.Remove(f.Name())

	testCRUD(t, openStormDB(t, f.Name()))
}

func TestStormGetQuirks(t *testing.T) {
	f, _ := ioutil.TempFile("", "test_db_")
	defer os.Remove(f.Name())

	testGetQuirks(t, openStormDB(t, f.Name()))
}

func TestStormTTL(t *testing.T) {
	f, _ := ioutil.TempFile("", "test_db_")
	defer os.Remove(f.Name())

	testTTL(t, openStormDB(t, f.Name()))
}

func TestStormAdvancedFetching(t *testing.T) {
	f, _ := ioutil.TempFile("", "test_db_")
	defer os.Remove(f.Name())

	testAdvancedFetching(t, openStormDB(t, f.Name()))
}

func TestStormTransactions(t *testing.T) {
	f, _ := ioutil.TempFile("", "test_db_")
	defer os.Remove(f.Name())

	testTransactions(t, openStormDB(t, f.Name()))
}
