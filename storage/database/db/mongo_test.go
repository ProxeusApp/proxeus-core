package db

import (
	"testing"
)

func openMongoDB(t *testing.T) DB {
	t.Skip("needs mongo setup")
	db, err := OpenDatabase("mongo", "mongodb://localhost:27017", "db00")
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
