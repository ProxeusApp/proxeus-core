package db

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/asdine/storm"
)

func OpenDatabase(engine, uri, name string) (DB, error) {
	switch engine {
	case "mongo":
		return OpenMongo(uri, name)
	case "storm", "":
		return OpenStorm(name)
	}
	return nil, fmt.Errorf("unknown db engine '%s'", engine)
}

func NotFound(err error) bool {
	return err == mongo.ErrNoDocuments || err == storm.ErrNotFound
}
