package www

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/ProxeusApp/proxeus-core/storage"

	"sort"
	"strconv"

	"github.com/labstack/echo"
)

var (
	defaultLang = ""
	cookieName  = "lang"
)

type WebI18n struct {
	i18nStore storage.I18nIF
	Lang      string
}

func NewI18n(i18nStore storage.I18nIF, c echo.Context) *WebI18n {
	wi18n := WebI18n{i18nStore, ""}
	wi18n.Lang = c.Param(cookieName)
	//write cookie only when reading from Accept-Language
	if wi18n.Lang == "" {
		wi18n.Lang = c.QueryParam(cookieName)
		if wi18n.Lang == "" {
			wi18n.Lang = c.FormValue(cookieName)
			if wi18n.Lang == "" {
				cookie, _ := c.Cookie(cookieName)
				if cookie != nil {
					wi18n.Lang = cookie.Value
					if wi18n.Lang == "" {
						wi18n.Lang = wi18n.findAcceptableLang(c.Request().Header.Get("Accept-Language"))
						writeLangCookie(c, wi18n.Lang)
						//cookieSet = true
					}
				} else {
					wi18n.Lang = wi18n.findAcceptableLang(c.Request().Header.Get("Accept-Language"))
					writeLangCookie(c, wi18n.Lang)
					//cookieSet = true
				}
			}
		}
	}
	return &wi18n
}

func writeLangCookie(c echo.Context, lang string) {
	cookie := new(http.Cookie)
	cookie.Path = "/"
	cookie.Name = cookieName
	cookie.Value = lang
	cookie.Expires = time.Now().Add(24 * 10 * time.Hour)
	c.SetCookie(cookie)
}

func (wi18n *WebI18n) findAcceptableLang(acceptLangHeader string) string {
	if acceptLangHeader != "" {
		langs := ParseAcceptLanguageHeader(acceptLangHeader)
		if len(langs) > 0 {
			possibleLangs, _ := wi18n.i18nStore.GetLangs(true)
			for _, httpAcceptLang := range langs {
				for _, lang := range possibleLangs {
					if lang.Matches(httpAcceptLang.Lang) {
						return lang.Code
					}
				}
			}
		}
	}
	return wi18n.getDefaultLang()
}

type HttpAceptLanguage struct {
	Lang   string
	Weight float64
}

var regexParseAcceptLanguageWeight = regexp.MustCompile(`q=([0-9\.]+)`)

//ParseAcceptLanguageHeader reads the weighted header values
//like de-CH,de;q=0.9,en;q=0.8,en-US;q=0.7
//into [{Lang:de-CH, Weight:1} {Lang:de Weight:0.9} {Lang:en Weight:0.8} {Lang:en-US Weight:0.7}]
func ParseAcceptLanguageHeader(headerValue string) []*HttpAceptLanguage {
	splittedList := strings.Split(headerValue, ",")
	langList := make([]*HttpAceptLanguage, 0)
	for i, a := range splittedList {
		b := strings.Split(strings.TrimSpace(a), ";")
		if len(b) > 0 {
			if len(b[0]) > 1 {
				lang := &HttpAceptLanguage{Lang: b[0], Weight: 0}
				//try to read the weight
				if len(b) > 1 {
					abb := regexParseAcceptLanguageWeight.FindAllStringSubmatch(b[1], -1)
					if len(abb) > 0 && len(abb[0]) > 1 {
						if s, err := strconv.ParseFloat(abb[0][1], 64); err == nil {
							lang.Weight = s
						}
					}
				}
				if lang.Weight == 0 && i == 0 {
					lang.Weight = 1
				}
				langList = append(langList, lang)
			}
		}
	}
	sort.Slice(langList, func(i, j int) bool {
		return langList[i].Weight > langList[j].Weight
	})
	return langList
}

//Translate
func (wi18n *WebI18n) T(b ...interface{}) string {
	lang := wi18n.Lang
	l := len(b)
	d := make([]string, l)
	wi18n.makeStrArray(l, d, b)
	var key string
	if l > 1 && strings.HasPrefix(d[0], "lang:") {
		lang = strings.Replace(d[0], "lang:", "", 1)
		if strings.HasPrefix(d[1], "key:") {
			key = strings.Replace(d[1], "key:", "", 1)
		} else {
			key = d[1]
		}
		if l > 2 {
			return wi18n.resolveKey(lang, key, d[2:]...)
		}
	} else {
		key = d[0]
		if l > 1 {
			return wi18n.resolveKey(lang, key, d[1:]...)
		}
	}
	t := wi18n.resolveKey(lang, key)
	return t
}

func (wi18n *WebI18n) resolveKey(lang, key string, args ...string) string {
	if wi18n.i18nStore != nil {
		key, _ = wi18n.i18nStore.Get(lang, key, args...)
	}
	return key
}
func (wi18n *WebI18n) getDefaultLang() string {
	if defaultLang == "" {
		if wi18n.i18nStore != nil {
			defaultLang, _ = wi18n.i18nStore.GetFallback()
			return defaultLang
		}
		defaultLang = "en"
	}
	return defaultLang
}

func (wi18n *WebI18n) GetAll() map[string]string {
	if wi18n.i18nStore != nil {
		lngs, _ := wi18n.i18nStore.GetAll(wi18n.Lang)
		return lngs
	}
	return nil
}

func (wi18n *WebI18n) makeStrArray(l int, strArray []string, d []interface{}) []string {
	i := 0
	for ; i < l; i++ {
		strArray[i] = fmt.Sprintf("%v", d[i])
	}
	return strArray
}
