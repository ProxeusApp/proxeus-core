package www

import (
	"encoding/base64"
	"regexp"
	"strings"

	"github.com/labstack/echo"

	"github.com/ProxeusApp/proxeus-core/sys"
	"github.com/ProxeusApp/proxeus-core/sys/model"
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

func (me *Context) Session(create bool) *sys.Session {
	sess, err := getSession(me, create)
	if err != nil {
		me.Logger().Error(err)
		return nil
	}
	return sess
}

func (me *Context) SessionWithUser(usr *model.User) *sys.Session {
	sess, err := getSessionWithUser(me, true, usr)
	if err != nil {
		return nil
	}

	return sess
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

// Extract the session token from the header
func (me *Context) SessionToken() string {
	return extractSessionToken(me.Request().Header.Get("Authorization"))
}

var sessionTokenFromHeaderReg = regexp.MustCompile(`^Bearer\s([^\s]+)$`)

func extractSessionToken(headerValue string) string {
	subm := sessionTokenFromHeaderReg.FindStringSubmatch(headerValue)
	l := len(subm)
	if l != 2 {
		return ""
	}
	return subm[1]
}

// Extract the basic authentication from the header
func (me *Context) BasicAuth() (string, string) {
	return extractBasicAuth(me.Request().Header.Get("Authorization"))
}

var basicAuthFromHeaderReg = regexp.MustCompile(`^Basic\s([^\s]+)$`)

func extractBasicAuth(headerValue string) (string, string) {
	subm := basicAuthFromHeaderReg.FindStringSubmatch(headerValue)
	l := len(subm)
	if l != 2 {
		return "", ""
	}

	b, err := base64.StdEncoding.DecodeString(subm[1])
	if err != nil {
		return "", ""
	}

	fields := strings.Split(string(b), ":")

	if len(fields) != 2 {
		return "", ""
	}

	return strings.TrimSpace(fields[0]), strings.TrimSpace(fields[1])
}
