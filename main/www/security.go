package www

import (
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo"

	"git.proxeus.com/core/central/sys/model"
)

type Security struct {
}

func NewSecurity() *Security {
	return &Security{}
}

func (a *Security) With(role model.Role) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) error {
			c := e.(*Context)
			if sess := c.Session(false); sess != nil && sess.AccessRights().IsGrantedFor(role) {
				return next(c)
			}
			return redirectToLogin(c)
		}
	}
}

func redirectToLogin(c echo.Context) error {
	if c.IsWebSocket() {
		return echo.ErrUnauthorized
	}
	req := c.Request()
	isAjax := req.Header.Get("X-Requested-With") == "XMLHttpRequest"
	var referer string
	if isAjax {
		referer = getURI(req.Header.Get("Host"), req.Header.Get("Referer"))
	} else {
		referer = getURI(req.Header.Get("Host"), req.URL.String())
	}
	if referer != "" {
		c.SetCookie(&http.Cookie{
			Name:    "R",
			Value:   base64.RawURLEncoding.EncodeToString([]byte(referer)),
			Path:    "/",
			Expires: time.Now().Add(time.Hour * 24),
		})
	}

	if isAjax {
		return echo.ErrUnauthorized
	}
	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}

func getURI(host, url string) string {
	if host == "" {
		return ""
	}
	i := strings.Index(url, host)
	if i != -1 {
		return url[strings.Index(url, host)+len(host):]
	}
	return ""
}
