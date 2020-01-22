package www

import (
	"encoding/gob"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"

	"github.com/ProxeusApp/proxeus-core/sys"
	"github.com/ProxeusApp/proxeus-core/sys/model"
)

func SessionMiddleware() echo.MiddlewareFunc {
	gob.Register(map[string]interface{}{})
	gob.Register(map[string]map[string]string{})
	sessionStore := sessions.NewCookieStore([]byte("secret_Dummy_1234"), []byte("12345678901234567890123456789012"))
	return session.Middleware(sessionStore)
}

var anonymousUser = &model.User{Role: model.PUBLIC}

func init() {
	//init here because fields belong to super struct
	anonymousUser.ID = ""
	anonymousUser.Name = "anonymous"
}

func getSessionWithUser(c *Context, create bool, usr *model.User) (currentSession *sys.Session, err error) {
	if !create || usr == nil {
		if csess := c.Get("sys.session"); csess != nil {
			var ok bool
			if currentSession, ok = csess.(*sys.Session); ok {
				return
			}
		}
	}

	sess, err := session.Get("s", c)
	if sess == nil || err != nil {
		return
	}
	if sid, ok := sess.Values["id"]; ok {
		//session exists
		if sidStr, ok := sid.(string); ok {
			if c.System() != nil {
				currentSession, err = c.System().GetSession(sidStr)
			}
		}
	}
	if create && usr != nil && currentSession != nil {
		if currentSession.S.UsrID != usr.ID || currentSession.AccessRights() != usr.Role || currentSession.S.UserName != usr.Name {
			currentSession.DeleteAll()
			currentSession = nil
		}
	}
	if currentSession == nil && create {
		if usr == nil {
			usr = anonymousUser
		}
		currentSession, err = c.System().NewSession(usr)
		if err != nil {
			return
		}
		if currentSession != nil {
			sess.Values["id"] = currentSession.S.ID
			options := sessions.Options{
				Path:     "/",
				MaxAge:   60 * 30, // 30 minutes,
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
			}
			settings := c.System().GetSettings()
			if strings.ToLower(settings.TestMode) != "true" {
				options.Secure = true
			}
			sess.Options = &options
			c.Set("sys.session", currentSession)
			err = sess.Save(c.Request(), c.Response())
		}
	}
	return
}
func getSession(c *Context, create bool) (currentSession *sys.Session, err error) {
	return getSessionWithUser(c, create, nil)
}

func delSession(c *Context) (err error) {
	var sess *sys.Session
	sess, err = getSession(c, false)
	if err != nil {
		return err
	}
	if sess == nil {
		return os.ErrNotExist
	}
	err = sess.DeleteAll()
	DeleteCookie(c, "s")
	return err
}

var pastTime = time.Unix(0, 0)

func DeleteCookie(c echo.Context, name string) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = ""
	cookie.Expires = pastTime
	c.SetCookie(cookie)
}
