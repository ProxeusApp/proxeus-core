package www

import (
	"log"
	"net/http"

	"github.com/labstack/echo"

	"github.com/ProxeusApp/proxeus-core/sys"
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

	if c.Response().Committed {
		log.Println(err)
		return
	}

	if c.Request().Header.Get("X-Requested-With") == "XMLHttpRequest" {
		err := c.JSON(code, msg)
		if err != nil {
			log.Println(err)
		}
		return
	}

	bts, err = sys.ReadAllFile("frontend.html")
	if err != nil {
		log.Println(err)
		return
	}

	err = c.HTMLBlob(code, bts)
	if err != nil {
		log.Println(err)
	}
}
