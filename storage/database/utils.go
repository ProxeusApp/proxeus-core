package database

import (
	"os"
	"regexp"
	"strings"

	"github.com/ProxeusApp/proxeus-core/storage"

	"github.com/asdine/storm/q"

	"github.com/ProxeusApp/proxeus-core/sys/model"
)

type (
	simpleQuery struct {
		metaOnly bool
		rawIndex int
		index    int
		limit    int
		exclude  []interface{}
		include  []interface{}
	}
)

func containsCaseInsensitiveReg(contains string) string {
	contains = strings.TrimSpace(contains)
	if contains == "" {
		return ""
	}
	return "(?i)" + regexp.QuoteMeta(contains)
}

func makeSimpleQuery(o storage.Options) *simpleQuery {
	const Limit = 1000
	sq := &simpleQuery{metaOnly: o.MetaOnly, index: o.Index, limit: o.Limit}

	if l := len(o.Exclude); l > 0 {
		sq.exclude = make([]interface{}, l)
		i := 0
		for k := range o.Exclude {
			sq.exclude[i] = k
			i++
		}
	}

	if l := len(o.Include); l > 0 {
		sq.include = make([]interface{}, l)
		i := 0
		for k := range o.Include {
			sq.include[i] = k
			i++
		}
	}

	if sq.limit == 0 || sq.limit > Limit {
		sq.limit = Limit
	}
	sq.rawIndex = sq.index
	sq.index = sq.index * sq.limit
	return sq
}

func ensureDir(dir string) error {
	var err error
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0750)
		if err != nil {
			return err
		}
	}
	return nil
}

func defaultMatcher(auth model.Auth, contains string, params *simpleQuery, includeReadGranted bool) []q.Matcher {
	matchers := commonMatcher(auth, contains, params)
	matchers = append(matchers, q.And(IsReadGrantedFor(auth, includeReadGranted)))
	return matchers
}

func publishedMatcher(auth model.Auth, contains string, params *simpleQuery) []q.Matcher {
	matchers := commonMatcher(auth, contains, params)
	var m q.Matcher
	if auth == nil {
		m = q.Eq("Published", true)
	} else {
		m = q.Or(
			q.Eq("Owner", auth.UserID()),
			q.Eq("Published", true),
		)
	}
	matchers = append(matchers, q.And(m))
	return matchers
}

func commonMatcher(auth model.Auth, contains string, params *simpleQuery) []q.Matcher {
	contains = containsCaseInsensitiveReg(contains)
	matchers := make([]q.Matcher, 0)
	if contains != "" {
		matchers = append(matchers,
			q.And(
				q.Or(
					q.Re("Name", contains),
					q.Re("Detail", contains),
				),
			),
		)
	}
	if params != nil {
		if len(params.exclude) > 0 {
			matchers = append(matchers,
				q.And(
					q.Not(q.In("ID", params.exclude)),
				),
			)
		}
		if len(params.include) > 0 {
			matchers = append(matchers,
				q.And(
					q.In("ID", params.include),
				),
			)
		}
	}
	return matchers
}
