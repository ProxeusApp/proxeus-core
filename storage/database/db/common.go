package db

import (
	"fmt"
	"time"

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

func ttlRefresh(opts []*GetOptions) bool {
	ttlRefresh := true
	for _, o := range opts {
		if o.NoTTLRefresh {
			ttlRefresh = false
		}
	}
	return ttlRefresh
}

func ttlDuration(opts []*SetOptions) time.Duration {
	var ttl time.Duration
	for _, o := range opts {
		if o.TTL > 0 {
			ttl = o.TTL
		}
	}
	return ttl
}
