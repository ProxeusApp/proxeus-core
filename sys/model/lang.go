package model

import (
	"strings"
)

type Lang struct {
	Code    string `storm:"id"`
	Enabled bool
}

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
