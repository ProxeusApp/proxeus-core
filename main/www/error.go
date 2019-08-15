package www

import (
	"log"
	"net/http"

	"github.com/labstack/echo"

	"git.proxeus.com/core/central/sys"
)

// DefaultHTTPErrorHandler is the default HTTP error handler
func DefaultHTTPErrorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		msg  interface{}
		bts  []byte
	)
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
	} else {
		msg = err.Error()
	}
	if _, ok := msg.(string); ok {
		msg = echo.Map{"message": msg}
	}
	if !c.Response().Committed {
		if c.Request().Header.Get("X-Requested-With") == "XMLHttpRequest" {
			if err = c.JSON(code, msg); err != nil {
				goto ERROR
			}
		} else {
			bts, err = sys.ReadAllFile("frontend.html")
			if err == nil {
				if err = c.HTMLBlob(code, bts); err != nil {
					goto ERROR
				}
				return
			}

		}
	}
ERROR:
	log.Println(err)
}
