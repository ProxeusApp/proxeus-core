package mysql

import (
	"bytes"
	"database/sql"
	"os"
	"regexp"
	"sync"
	"time"

	"git.proxeus.com/core/central/sys/i18n"
	"git.proxeus.com/core/central/sys/model"
)

type MysqlI18n struct {
	db           *sql.DB
	resolver     *i18n.I18nResolver
	allCache     map[string]map[string]string
	allCacheLock sync.RWMutex
	re           *regexp.Regexp
}

func NewI18nStore(db *sql.DB) (*MysqlI18n, error) {
	i := &MysqlI18n{db: db, resolver: &i18n.I18nResolver{}, allCache: make(map[string]map[string]string), re: regexp.MustCompile(`^[A-Za-z_]{2,6}$`)}
	i.setup()
	return i, nil
}

func (me *MysqlI18n) setup() (err error) {
	_, err = me.db.Exec("UPDATE internationalization f SET f.loc = 'en' WHERE f.loc = '*';")
	return
}

func (me *MysqlI18n) Find(keyContains string, valueContains string, options map[string]interface{}) (res map[string]map[string]string, err error) {
	sQuery := makeSimpleQuery(options)
	var r *sql.Rows
	if keyContains != "" {
		r, err = me.db.Query("SELECT f.`loc`,f.`code`,f.`text` FROM `internationalization` f WHERE f.code LIKE ? order by f.code LIMIT ?,?;", "%"+keyContains+"%", sQuery.index, sQuery.limit)
	} else if valueContains != "" {
		r, err = me.db.Query("SELECT f.`loc`,f.`code`,f.`text` FROM `internationalization` f WHERE f.text LIKE ? order by f.code LIMIT ?,?;", "%"+valueContains+"%", sQuery.index, sQuery.limit)
	} else {
		return
	}

	if err != nil {
		return
	}
	defer r.Close()
	var (
		loc  sql.NullString
		code sql.NullString
		text sql.NullString
	)
	queryAgainFull := make([]interface{}, 0)
	inQuery := &bytes.Buffer{}
	inQuery.WriteString("SELECT f.`loc`,f.`code`,f.`text` FROM `internationalization` f WHERE f.code in(")
	i := 0
	if r.Err() != nil {
		err = r.Err()
		return
	}
	for r.Next() {
		err = r.Scan(&loc, &code, &text)
		if err != nil {
			return
		}
		if code.Valid {
			queryAgainFull = append(queryAgainFull, code.String)
			if i > 0 {
				inQuery.WriteString(",?")
			} else {
				inQuery.WriteString("?")
			}
			i++
		}
	}
	if i == 0 {
		return nil, nil
	}
	inQuery.WriteString(")  order by f.`code`,f.`text`")
	r, err = me.db.Query(inQuery.String(), queryAgainFull...)
	if err != nil {
		return
	}
	res = make(map[string]map[string]string)
	if r.Err() != nil {
		err = r.Err()
		return
	}
	for r.Next() {
		err = r.Scan(&loc, &code, &text)
		if err != nil {
			return
		}
		if code.Valid {
			queryAgainFull = append(queryAgainFull)
		}

		if loc.Valid && code.Valid && text.Valid {
			if res[code.String] == nil {
				res[code.String] = make(map[string]string)
			}
			res[code.String][loc.String] = text.String
		}
	}
	return res, nil
}

func (me *MysqlI18n) Get(lang string, key string, args ...string) (resolvedText string, err error) {
	var r *sql.Rows
	r, err = me.db.Query("SELECT f.`text` FROM `internationalization` f WHERE f.`loc` = ? AND f.`code` = ? LIMIT 0,1;", lang, key)
	if err != nil {
		return
	}
	defer r.Close()
	var (
		text sql.NullString
	)
	for r.Next() {
		err = r.Scan(&text)
		if err != nil {
			return
		}
		if text.Valid {
			resolvedText = text.String
			resolvedText = me.resolver.Resolve(resolvedText, args...)
			return
		}
	}
	resolvedText = key
	return
}

func (me *MysqlI18n) GetAll(lang string) (res map[string]string, err error) {
	me.allCacheLock.RLock()
	if m, ok := me.allCache[lang]; ok {
		me.allCacheLock.RUnlock()
		return m, nil
	}
	me.allCacheLock.RUnlock()
	var r *sql.Rows
	r, err = me.db.Query("SELECT f.`code`, f.`text` FROM `internationalization` f WHERE f.`loc` = ?;", lang)
	if err != nil {
		return
	}
	defer r.Close()
	var (
		code sql.NullString
		text sql.NullString
	)
	res = make(map[string]string)
	for r.Next() {
		err = r.Scan(&code, &text)
		if err != nil {
			return
		}
		if code.Valid && text.Valid {
			res[code.String] = text.String
		}
	}
	me.allCacheLock.Lock()
	me.allCache[lang] = res
	me.allCacheLock.Unlock()
	return
}

func (me *MysqlI18n) Put(lang string, key string, text string) (err error) {
	if lang != "" && key != "" {
		if !me.re.MatchString(lang) {
			return os.ErrInvalid
		}
		me.allCacheLock.Lock()
		if m, ok := me.allCache[lang]; ok {
			m[key] = text
		}
		me.allCacheLock.Unlock()
		var result sql.Result
		result, err = me.db.Exec("UPDATE internationalization f SET f.text = ? WHERE f.loc = ? AND f.code = ?;",
			text,
			lang,
			key)
		if err == nil {
			var count int64
			count, err = result.RowsAffected()
			if count == 0 {
				now := time.Now()
				result, _ = me.db.Exec("INSERT INTO internationalization (`version`, `code`, `date_created`, `last_updated`, `loc`, `relevance`, `text`) VALUES (1,?,?,?,?,1,?);",
					key,
					now,
					now,
					lang,
					text)
			}
		}
	}
	return
}

func (me *MysqlI18n) Delete(keyContains string) error {
	return nil
}

func (me *MysqlI18n) PutLang(code string, enabled bool) (err error) {
	if !me.re.MatchString(code) {
		return os.ErrInvalid
	}
	me.db.Exec("INSERT INTO `language` (`code`, `label`, `enabled`) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE `code`=?, `label`=?, `enabled`=?", code, code, enabled, code, code, enabled)
	return
}

func (me *MysqlI18n) langCount(code string) (count int) {
	rows, err := me.db.Query("SELECT COUNT(*) as count FROM `language` where `code`=?;", code)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&count)
	}
	return
}

func (me *MysqlI18n) GetLangs(enabled bool) (langs []*model.Lang, err error) {
	var r *sql.Rows
	r, err = me.db.Query("SELECT `code`, `enabled` FROM `language` f WHERE f.enabled = ?;", enabled)
	if err != nil {
		return
	}
	defer r.Close()
	var (
		code     sql.NullString
		senabled sql.NullString
	)
	langs = make([]*model.Lang, 0)
	for r.Next() {
		err = r.Scan(&code, &senabled)
		if err != nil {
			return
		}
		lang := &model.Lang{}
		if senabled.Valid {
			lang.Enabled = senabled.String == "\x01"
		}
		if code.Valid {
			lang.Code = code.String
		}
		langs = append(langs, lang)
	}
	return
}

func (me *MysqlI18n) GetAllLangs() (langs []*model.Lang, err error) {
	var r *sql.Rows
	r, err = me.db.Query("SELECT `code`, `enabled` FROM `language` f;")
	if err != nil {
		return
	}
	defer r.Close()
	var (
		code    sql.NullString
		enabled sql.NullString
	)

	langs = make([]*model.Lang, 0)
	for r.Next() {
		err = r.Scan(&code, &enabled)
		if err != nil {
			return
		}
		lang := &model.Lang{}
		if enabled.Valid {
			lang.Enabled = enabled.String == "\x01"
		}
		if code.Valid {
			lang.Code = code.String
		}
		langs = append(langs, lang)
	}
	return
}

func (me *MysqlI18n) PutFallback(l string) error {
	return nil
}

func (me *MysqlI18n) GetFallback() (string, error) {
	return "en", nil
}

func (me *MysqlI18n) Close() error {
	return nil
}
