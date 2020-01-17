package database

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/ProxeusApp/proxeus-core/storage/database/db"
)

func openStormDB(t *testing.T, path string) *db.StormShim {
	db, err := db.OpenStorm(path)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestCRUDStorm(t *testing.T) {
	f, _ := ioutil.TempFile("", "test_db_")
	defer os.Remove(f.Name())

	testCRUD(t, openStormDB(t, f.Name()))
}

func TestGetQuirksStorm(t *testing.T) {
	f, _ := ioutil.TempFile("", "test_db_")
	defer os.Remove(f.Name())

	testGetQuirks(t, openStormDB(t, f.Name()))
}

func TestAdvancedFetchingStorm(t *testing.T) {
	f, _ := ioutil.TempFile("", "test_db_")
	defer os.Remove(f.Name())

	testAdvancedFetching(t, openStormDB(t, f.Name()))
}

func TestTransactionsStorm(t *testing.T) {
	f, _ := ioutil.TempFile("", "test_db_")
	defer os.Remove(f.Name())

	testTransactions(t, openStormDB(t, f.Name()))
}
