package storm

import (
	"os"
	"regexp"
	"strings"

	"github.com/asdine/storm/q"

	"git.proxeus.com/core/central/sys/model"
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
			if l := len(e); l > 0 {
				sq.exclude = make([]interface{}, l)
				i := 0
				for k := range e {
					sq.exclude[i] = k
					i++
				}
			}
		} else if e, ok := more["exclude"].([]string); ok {
			if l := len(e); l > 0 {
				sq.exclude = make([]interface{}, l)
				for i, k := range e {
					sq.exclude[i] = k
				}
			}
		}
		if e, ok := more["include"].(map[string]interface{}); ok {
			if l := len(e); l > 0 {
				sq.include = make([]interface{}, l)
				i := 0
				for k := range e {
					sq.include[i] = k
					i++
				}
			}
		} else if e, ok := more["include"].([]string); ok {
			if l := len(e); l > 0 {
				sq.include = make([]interface{}, l)
				for i, k := range e {
					sq.include[i] = k
				}
			}
		}
	}
	if sq.limit == 0 {
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

func defaultMatcher(auth model.Authorization, contains string, params *simpleQuery, includeReadGranted bool) []q.Matcher {
	matchers := commonMatcher(auth, contains, params)
	matchers = append(matchers, q.And(IsReadGrantedFor(auth, includeReadGranted)))
	return matchers
}

func publishedMatcher(auth model.Authorization, contains string, params *simpleQuery) []q.Matcher {
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

func commonMatcher(auth model.Authorization, contains string, params *simpleQuery) []q.Matcher {
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
