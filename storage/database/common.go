package database

import (
	"fmt"

	"github.com/ProxeusApp/proxeus-core/storage/database/db"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/asdine/storm"
)

func OpenDatabase(engine, uri, name string) (db.DB, error) {
	switch engine {
	case "mongo":
		return db.OpenMongo(uri, name)
	case "storm", "":
		return db.OpenStorm(name)
	}
	return nil, fmt.Errorf("unknown db engine '%s'", engine)
}

func NotFound(err error) bool {
	return err == mongo.ErrNoDocuments || err == storm.ErrNotFound
}
