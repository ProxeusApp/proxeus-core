package i18n

import (
	"bytes"
	"regexp"
	"strconv"
)

const (
	phStart = "\\{"
	phEnd   = "\\}"
)

type I18nResolver struct {
	RegexCache map[int]*regexp.Regexp
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
	newPath.Write([]byte(phStart))
	newPath.Write([]byte(strconv.Itoa(i)))
	newPath.Write([]byte(phEnd))
	by := newPath.String()
	return &by
}
