package www

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/ProxeusApp/proxeus-core/sys/session"
)

// SessionAuthToken create a request session if a valid API Key is found
func SessionTokenAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(e echo.Context) error {
		c := e.(*Context)

		if sess := c.Session(false); sess != nil {
			return next(c)
		}

		// We first check if we can authenticate with an API key
		sess, err := sessionFromSessionToken(c)
		if err != nil {
			// We had an session token but it not valid
			return c.NoContent(http.StatusUnauthorized)
		}

		var removeCookie bool
		if sess != nil {
			c.Set("sys.session", sess)
			removeCookie = true
		}

		if err = next(c); err != nil {
			c.Error(err)
		}

		if removeCookie {
			c.Response().Header().Del("Set-Cookie")
		}

		return nil
	}
}

func sessionFromSessionToken(c *Context) (*session.Session, error) {
	token := c.SessionToken()
	if token == "" {
		return nil, nil
	}

	sess, err := c.System().SessionMgmnt.Get(token)
	if err != nil {
		return nil, err
	}
	return sess, nil
}
