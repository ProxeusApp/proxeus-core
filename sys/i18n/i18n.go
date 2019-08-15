package i18n

//"../db"

import (
	"bytes"
	"regexp"
	"strconv"

	"github.com/asdine/storm"

	"git.proxeus.com/core/central/sys/model"
)

var (
	phStart       = []byte("\\{")
	phEnd         = []byte("\\}")
	keySplitRegex = regexp.MustCompile(`^([^\.]+)\.(.*)`)
)

const (
	i18nBucket = "_i18n_"
)

type I18nSettings struct {
	ID       string `storm:"id"`
	Fallback string
}

type I18n struct {
	db          *storm.DB
	resolver    *I18nResolver
	settings    *I18nSettings
	activeLangs []*model.Lang
}

type I18nResolver struct {
	RegexCache map[int]*regexp.Regexp
}

func (i18n *I18n) setup(db *storm.DB) error {
	i18n.db = db
	i18n.GetFallback()
	var err error
	i18n.activeLangs, err = i18n.getLangs(true)
	return err
}

func (i18n *I18n) Find(keyContains string, valueContains string, limit int) (map[string]map[string]string, error) {
	tx, err := i18n.db.Bolt.Begin(false)
	res := make(map[string]map[string]string)
	if err != nil {
		return res, err
	}
	defer tx.Rollback()
	index := 0
	b := tx.Bucket([]byte(i18nBucket))
	if b == nil {
		return res, nil
	}
	c := b.Cursor()
	var l int
	var key string
	var lang string
	keyRegex := regexp.MustCompile(`^([^\.]+)\.(.*(` + regexp.QuoteMeta(keyContains) + `).*)`)
	valueRegex := regexp.MustCompile(`(.*(` + regexp.QuoteMeta(valueContains) + `).*)`)
	if keyContains != "" && valueContains == "" {
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if keyRegex.Match(k) {
				pieces := keySplitRegex.FindAllSubmatch([]byte(k), -1)
				if len(pieces) > 0 {
					l = len(v) - 1
					lang = string(pieces[0][1])
					key = string(pieces[0][2])
					if res[key] == nil {
						res[key] = make(map[string]string)
						index++
					}
					res[key][lang] = string(v[1:l])
					if index == limit {
						break
					}
				}
			}
		}
	} else if keyContains == "" && valueContains != "" {
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if valueRegex.Match(v) {
				pieces := keySplitRegex.FindAllSubmatch([]byte(k), -1)
				if len(pieces) > 0 {
					l = len(v) - 1
					lang = string(pieces[0][1])
					key = string(pieces[0][2])
					if res[key] == nil {
						res[key] = make(map[string]string)
						index++
					}
					res[key][lang] = string(v[1:l])
					if index == limit {
						break
					}
				}
			}
		}
	} else if keyContains != "" && valueContains != "" {
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if keyRegex.Match(k) && valueRegex.Match(v) {
				pieces := keySplitRegex.FindAllSubmatch([]byte(k), -1)
				if len(pieces) > 0 {
					l = len(v) - 1
					lang = string(pieces[0][1])
					key = string(pieces[0][2])
					if res[key] == nil {
						res[key] = make(map[string]string)
						index++
					}
					res[key][lang] = string(v[1:l])
					if index == limit {
						break
					}
				}
			}
		}
	} else if keyContains == "" && valueContains == "" {
		for k, v := c.First(); k != nil; k, v = c.Next() {
			pieces := keySplitRegex.FindAllSubmatch([]byte(k), -1)
			if len(pieces) > 0 {
				l = len(v) - 1
				lang = string(pieces[0][1])
				key = string(pieces[0][2])
				if res[key] == nil {
					res[key] = make(map[string]string)
					index++
				}
				res[key][lang] = string(v[1:l])
				if index == limit {
					break
				}
			}
		}
	}
	return res, nil
}

func (i18n *I18n) Delete(keyContains string) error {
	if keyContains == "" {
		return nil
	}
	tx, err := i18n.db.Bolt.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	b := tx.Bucket([]byte(i18nBucket))
	if b == nil {
		return nil
	}
	c := b.Cursor()
	keyRegex := regexp.MustCompile(`^([^\.]+)\.(.*(` + regexp.QuoteMeta(keyContains) + `).*)`)
	for k, _ := c.First(); k != nil; k, _ = c.Next() {
		if keyRegex.Match(k) {
			pieces := keySplitRegex.FindAllSubmatch([]byte(k), -1)
			if len(pieces) > 0 {
				c.Delete()
			}
		}
	}
	return nil
}

func (i18n *I18n) Get(lang string, key string, args ...string) (string, error) {
	sel := make([]string, 1)
	sel[0] = lang + "." + key
	var text string

	err := i18n.db.Get(i18nBucket, lang+"."+key, &text)
	if err == nil && text != "" {
		if i18n.resolver == nil {
			i18n.resolver = &I18nResolver{}
		}
		return i18n.resolver.Resolve(text, args...), nil
	}
	if i18n.settings.Fallback != "" {
		if lang != i18n.settings.Fallback {
			//fallback
			return i18n.Get(i18n.settings.Fallback, key, args...)
		}
	}
	return key, nil
}

func (i18n *I18n) Put(lang string, key string, text string) error {
	return i18n.db.Set(i18nBucket, lang+"."+key, text)
}

func (i18n *I18n) PutLang(code string, enabled bool) error {
	var err error
	err = i18n.db.Save(&model.Lang{Code: code, Enabled: enabled})
	if err != nil {
		return err
	}
	i18n.activeLangs, err = i18n.getLangs(true)
	return err
}

func (i18n *I18n) GetLangs(enabled bool) ([]*model.Lang, error) {
	var err error
	if enabled {
		if i18n.activeLangs == nil || len(i18n.activeLangs) == 0 {
			i18n.activeLangs, err = i18n.getLangs(enabled)
		}
		return i18n.activeLangs, err
	}
	return i18n.getLangs(enabled)
}

func (i18n *I18n) getLangs(enabled bool) ([]*model.Lang, error) {
	var langs []model.Lang
	i18n.db.Find("Enabled", enabled, &langs)
	langCodes := make([]*model.Lang, len(langs))
	for i := range langs {
		langCodes[i] = &langs[i]
	}
	return langCodes, nil
}

func (i18n *I18n) GetAllLangs() ([]*model.Lang, error) {
	var langs []model.Lang
	i18n.db.All(&langs)
	langRefs := make([]*model.Lang, len(langs))
	for i := range langs {
		langRefs[i] = &langs[i]
	}
	return langRefs, nil
}

func (i18n *I18n) PutFallback(l string) error {
	if l != "" {
		i18n.settings.Fallback = l
		s := *i18n.settings
		return i18n.db.Save(&s)
	}
	return nil
}

func (i18n *I18n) GetFallback() (string, error) {
	if i18n.settings == nil {
		var s I18nSettings
		i18n.db.One("ID", "main", &s)
		if s == (I18nSettings{}) {
			i18n.settings = &I18nSettings{ID: "main"}
		} else {
			i18n.settings = &s
		}
	}
	return i18n.settings.Fallback, nil
}

func (i18n *I18n) Close() error {
	return nil
}

func (i18nr *I18nResolver) Resolve(text string, args ...string) string {
	if len(args) > 0 {
		for i, match := range args {
			text = i18nr.rexp(i).ReplaceAllString(text, match)
		}
	}
	return text
}

func (i18nr *I18nResolver) rexp(i int) *regexp.Regexp {
	if i18nr.RegexCache == nil {
		i18nr.RegexCache = make(map[int]*regexp.Regexp)
	}
	r := i18nr.RegexCache[i]
	if r != nil {
		return r
	}
	r = regexp.MustCompile(*i18nr.msgPh(i))
	i18nr.RegexCache[i] = r
	return r
}

func (i18nr *I18nResolver) msgPh(i int) *string {
	newPath := new(bytes.Buffer)
	newPath.Write(phStart)
	newPath.Write([]byte(strconv.Itoa(i)))
	newPath.Write(phEnd)
	by := newPath.String()
	return &by
}
