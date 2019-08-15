package model

import (
	"strings"
)

type (
	Lang struct {
		Code    string `storm:"id"`
		Enabled bool
	}

	I18nStoreIF interface {
		Find(keyContains string, valueContains string, options map[string]interface{}) (map[string]map[string]string, error)
		Get(lang string, key string, args ...string) (string, error)
		GetInsert(lang string, key string, args ...string) (string, error)
		GetAll(lang string) (map[string]string, error)
		Put(lang string, key string, text string) error
		Delete(keyContains string) error

		PutLang(code string, enabled bool) error
		GetLangs(enabled bool) ([]*Lang, error)
		GetAllLangs() ([]*Lang, error)
		PutFallback(l string) error
		GetFallback() (string, error)
		Close() error
	}
)

func (l *Lang) Matches(code string) bool {
	if l.Code == code {
		return true
	}
	if len(code) > 2 {
		i := strings.Index(code, "-")
		if i > -1 && strings.Contains(l.Code, code[0:i]) {
			return true
		}
		i = strings.Index(code, "_")
		if i > -1 && strings.Contains(l.Code, code[0:i]) {
			return true
		}
	}
	return false
}
