// +build integration

package db

import (
	"os"
	"testing"
)

func openMongoDB(t *testing.T) DB {
	db, err := OpenDatabase("mongo", os.Getenv("PROXEUS_DATABASE_URI"), "db00")
	if err != nil {
		t.Fatal(err)
	}
	return db
}
func TestMongoCRUD(t *testing.T) {
	testCRUD(t, openMongoDB(t))
}

func TestMongoGetQuirks(t *testing.T) {
	testGetQuirks(t, openMongoDB(t))
}

func TestMongoTTL(t *testing.T) {
	testTTL(t, openMongoDB(t))
}

func TestMongoAdvancedFetching(t *testing.T) {
	testAdvancedFetching(t, openMongoDB(t))
}

func TestMongoTransactions(t *testing.T) {
	testTransactions(t, openMongoDB(t))
}
