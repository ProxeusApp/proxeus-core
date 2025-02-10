package formbuilder

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/labstack/echo/v4"

	"io/ioutil"
	"strings"

	"github.com/ProxeusApp/proxeus-core/main/handlers/helpers"
	"github.com/ProxeusApp/proxeus-core/main/www"
	"github.com/ProxeusApp/proxeus-core/sys/validate"
)

func GetDataId(e echo.Context) error {
	c := e.(*www.Context)
	reset, _ := strconv.ParseBool(c.QueryParam("reset"))
	id := c.Param("id")
	if id == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusBadRequest)
	}
	newData, err := formService.GetFormData(sess, id, reset)
	if err == nil {
		return c.JSON(http.StatusOK, newData)
	}
	return c.NoContent(http.StatusBadRequest)
}

func GetFileTypes(e echo.Context) error {
	return e.JSON(http.StatusOK, validate.FileTypes())
}

func SetFormSrcHandler(e echo.Context) error {
	c := e.(*www.Context)
	id := c.Param("id")
	if id == "" {
		return c.NoContent(http.StatusBadRequest)
	}
	if !strings.Contains(c.Request().Header.Get("Content-Type"), "application/json") {
		return c.NoContent(http.StatusBadRequest)
	}
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusBadRequest)
	}
	body, _ := ioutil.ReadAll(c.Request().Body)
	formSrc := map[string]interface{}{}
	err := json.Unmarshal(body, &formSrc)
	if err != nil || len(formSrc) <= 0 {
		return c.NoContent(http.StatusBadRequest)
	}
	err = formService.SetFormSrc(sess, id, formSrc)
	if err == nil {
		return c.NoContent(http.StatusOK)
	}
	return c.NoContent(http.StatusBadRequest)
}

func TestFormDataHandler(e echo.Context) error {
	c := e.(*www.Context)
	submit, _ := strconv.ParseBool(c.QueryParam("s"))
	id := c.Param("id")
	if id == "" {
		return c.NoContent(http.StatusBadRequest)
	}
	input, err := helpers.ParseDataFromReq(c)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusBadRequest)
	}
	errors, err := formService.TestFormData(sess, id, input, submit)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	if len(errors) > 0 {
		errors.Translate(func(key string, args ...string) string {
			return c.I18n().T(key, args)
		})
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"errors": errors})
	}
	return c.NoContent(http.StatusOK)
}

func GetFileIdFieldName(e echo.Context) error {
	c := e.(*www.Context)
	fieldname := c.Param("fieldname")
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	id := c.Param("id")
	if id == "" {
		return c.NoContent(http.StatusNotFound)
	}
	err := formService.GetFormFile(sess, id, fieldname, c.Response().Writer)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	c.Response().Committed = true //prevents from-> http: multiple response.WriteHeader calls
	return c.NoContent(http.StatusOK)
}

func PostFileIdFieldName(e echo.Context) error {
	c := e.(*www.Context)
	id := c.Param("id")
	fieldname := c.Param("fieldname")
	if fieldname == "" || id == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusBadRequest)
	}

	fileName, _ := url.QueryUnescape(c.Request().Header.Get("File-Name"))
	if fileName == "" {
		fileName = "unknown"
	}

	fData, err := formService.PostFormFile(sess, id, fileName, fieldname, c.Request().Body, c.Request().Header.Get("Content-Type"))
	if err != nil {
		if er, ok := err.(validate.ErrorMap); ok {
			er.Translate(func(key string, args ...string) string {
				return c.I18n().T(key, args)
			})
			return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"errors": er})
		}
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"errors": err})
	}
	return c.JSON(http.StatusOK, fData)
}
