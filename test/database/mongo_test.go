package database

import (
	"testing"

	"github.com/ProxeusApp/proxeus-core/storage/database/db"
)

func openMongoDB(t *testing.T) *db.MongoShim {
	db, err := db.OpenMongo("mongodb://localhost:27017", "db00")
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestCRUDMongo(t *testing.T) {
	testCRUD(t, openMongoDB(t))
}

func TestGetQuirksMongo(t *testing.T) {
	testGetQuirks(t, openMongoDB(t))
}

func TestAdvancedFetchingMongo(t *testing.T) {
	testAdvancedFetching(t, openMongoDB(t))
}

func TestTransactionsMongo(t *testing.T) {
	testTransactions(t, openMongoDB(t))
}
