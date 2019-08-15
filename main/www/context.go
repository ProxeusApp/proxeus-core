package www

import (
	"regexp"

	"github.com/labstack/echo"

	"git.proxeus.com/core/central/sys"
	"git.proxeus.com/core/central/sys/model"
	"git.proxeus.com/core/central/sys/session"
)

var singleSystem *sys.System

func SetSystem(s *sys.System) {
	singleSystem = s
}

type Context struct {
	echo.Context
	webI18n *WebI18n
}

func (me *Context) Lang() string {
	if me.webI18n != nil {
		return me.webI18n.Lang
	}
	me.webI18n = NewI18n(me.System().DB.I18n, me)
	return me.webI18n.Lang
}

func (me *Context) Session(create bool) *session.Session {
	sess, err := getSession(me, create)
	if err != nil {
		return nil
	}
	return sess
}

func (me *Context) SessionWithUser(usr *model.User) *session.Session {
	sess, err := getSessionWithUser(me, true, usr)
	if err != nil {
		return nil
	}
	return sess
}

//Auth checks if there is a session available otherwise it retrieves the api key if possible
func (me *Context) Auth() (model.Authorization, error) {
	sess := me.Session(false)
	if sess == nil {
		u, err := useApiKeyAsUserAuth(me)
		if err != nil {
			return nil, err
		}
		return u, nil
	}
	return sess, nil
}

func (me *Context) ApiKey() (string, error) {
	_, apiKey := readApiKeyFromHeader(me.Request().Header.Get("Authorization"))
	if len(apiKey) > 0 {
		return apiKey, nil
	}
	return "", echo.ErrNotFound
}

func (me *Context) EndSession() {
	_ = delSession(me)
}

func (me *Context) System() *sys.System {
	return singleSystem
}

func (me *Context) I18n() *WebI18n {
	if me.webI18n != nil {
		return me.webI18n
	}
	me.webI18n = NewI18n(me.System().DB.I18n, me)
	return me.webI18n
}

var apiKeyFromHeaderReg = regexp.MustCompile(`\s*(\w+)?\s*([^\s]*)`)

//Returns type and key as string.
func readApiKeyFromHeader(headerValue string) (string, string) {
	subm := apiKeyFromHeaderReg.FindAllStringSubmatch(headerValue, 1)
	l := len(subm)
	if l == 1 {
		l = len(subm[0])
		if l == 3 {
			if len(subm[0][2]) == 0 {
				return "", subm[0][1]
			} else {
				return subm[0][1], subm[0][2]
			}
		}
	}
	return "", ""
}

func useApiKeyAsUserAuth(c *Context) (model.Authorization, error) {
	apiKey, err := c.ApiKey()
	if err != nil {
		return nil, err
	}
	u, err := c.System().DB.User.APIKey(apiKey)
	if err != nil {
		return nil, err
	}
	return u, nil
}
