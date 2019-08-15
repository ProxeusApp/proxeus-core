package mysql

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type (
	Config struct {
		User     string
		Password string
		Database string
	}
	simpleQuery struct {
		metaOnly bool
		index    int
		limit    int
		exclude  map[string]interface{}
		include  map[string]interface{}
	}
)

func CreateDBConnection(conf Config) (*sql.DB, error) {
	return sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", conf.User, conf.Password, conf.Database))
}

func CreateDBConnectionByDataSource(dataSourceName string) (*sql.DB, error) {
	log.Println("conn", dataSourceName)
	return sql.Open("mysql", dataSourceName)
}

func makeSimpleQuery(more map[string]interface{}) *simpleQuery {
	const Limit = 1000
	sq := &simpleQuery{metaOnly: true, index: 0, limit: 20}
	if len(more) > 0 {
		if mo, ok := more["metaOnly"].(bool); ok {
			sq.metaOnly = mo
		}
		if i, ok := more["index"].(int); ok {
			sq.index = i
		} else if i, ok := more["index"].(float64); ok {
			sq.index = int(i)
		}
		if l, ok := more["limit"].(int); ok {
			sq.limit = l
			if sq.limit > Limit {
				sq.limit = Limit
			}
		} else if l, ok := more["limit"].(float64); ok {
			sq.limit = int(l)
			if sq.limit > Limit {
				sq.limit = Limit
			}
		}
		if e, ok := more["exclude"].(map[string]interface{}); ok {
			sq.exclude = e
		}
		if e, ok := more["include"].(map[string]interface{}); ok {
			sq.include = e
		}
	}
	if sq.limit == 0 {
		sq.limit = Limit
	}
	sq.index = sq.index * sq.limit
	return sq
}
