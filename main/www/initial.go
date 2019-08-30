package www

import (
	"net/http"
	"strings"

	"github.com/labstack/echo"

	"git.proxeus.com/core/central/sys"
	"git.proxeus.com/core/central/sys/model"
)

type InitialHandler struct {
	configured      bool
	cleanOnNextCall bool
}

func NewInitialHandler(configured bool) *InitialHandler {
	return &InitialHandler{configured: configured}
}

func (me *InitialHandler) Handler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(e echo.Context) error {
		c := e.(*Context)
		if !me.configured {
			sess := c.Session(false)
			if sess == nil || sess.AccessRights() != model.ROOT { //ensure we have a tmp root session to power up
				sess = c.SessionWithUser(&model.User{ID: "XYZ", Role: model.ROOT})
			}

			if me.cleanOnNextCall && c.Request().RequestURI != "/api/import/results" && c.Request().RequestURI != "/api/init" {
				me.configured, _ = c.System().Configured()
				me.cleanOnNextCall = false
				er2 := c.System().SessionMgmnt.Clean()
				if er2 != nil {
					return c.NoContent(http.StatusInternalServerError)
				}
				return next(c)
			}
			if strings.ToLower(c.Request().Method) == "get" {
				if !strings.HasPrefix(c.Request().RequestURI, "/api/") &&
					!strings.HasPrefix(c.Request().RequestURI, "/static/") &&
					!strings.HasPrefix(c.Request().RequestURI, "/favicon.ico") {
					bts, err := sys.ReadAllFile("initial.html")
					if err != nil {
						return c.NoContent(http.StatusNotFound)
					}
					return c.HTMLBlob(http.StatusOK, bts)
				}
				return next(c)
			} else {
				if strings.HasPrefix(c.Request().RequestURI, "/api/init") || strings.HasPrefix(c.Request().RequestURI, "/api/import") {
					er := next(c)
					me.configured, _ = c.System().Configured()
					if me.configured {
						me.configured = false
						//to let /api/import/results through as all sessions will be deleted afterwards, this makes it possible
						//to view the results before they are gone
						me.cleanOnNextCall = true
					}
					return er
				}
				return c.NoContent(http.StatusBadRequest)
			}

		}
		return next(c)
	}
}
