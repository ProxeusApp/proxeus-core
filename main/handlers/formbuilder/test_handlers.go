package formbuilder

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ProxeusApp/proxeus-core/sys"

	"github.com/labstack/echo"

	"io/ioutil"
	"strings"

	"github.com/ProxeusApp/proxeus-core/main/handlers/helpers"
	"github.com/ProxeusApp/proxeus-core/main/www"
	"github.com/ProxeusApp/proxeus-core/sys/file"
	"github.com/ProxeusApp/proxeus-core/sys/form"
	"github.com/ProxeusApp/proxeus-core/sys/validate"
)

func GetDataId(e echo.Context) error {
	c := e.(*www.Context)
	reset, _ := strconv.ParseBool(c.QueryParam("reset"))
	id := c.Param("id")
	if id != "" {
		sess := c.Session(false)
		if sess != nil {
			dc := GetDataManager(sess)
			if reset {
				dc.Clear(id)
			}
			newData, err := dc.GetData(id)
			if err == nil {
				return c.JSON(http.StatusOK, newData)
			}
		}
	}
	return c.NoContent(http.StatusBadRequest)
}

func GetFileTypes(e echo.Context) error {
	return e.JSON(http.StatusOK, validate.FileTypes())
}

func SetFormSrcHandler(e echo.Context) error {
	c := e.(*www.Context)
	id := c.Param("id")
	if id != "" {
		if strings.Contains(c.Request().Header.Get("Content-Type"), "application/json") {
			body, _ := ioutil.ReadAll(c.Request().Body)
			formSrc := map[string]interface{}{}
			err := json.Unmarshal(body, &formSrc)
			if err == nil && len(formSrc) > 0 {
				sess := c.Session(false)
				if sess != nil {
					dc := GetDataManager(sess)
					err = dc.PutDataWithoutMerge("src"+id, formSrc)
					if err == nil {
						return c.NoContent(http.StatusOK)
					}
				}
			}
		}
	}
	return c.NoContent(http.StatusBadRequest)
}

func TestFormDataHandler(e echo.Context) error {
	c := e.(*www.Context)
	submit, _ := strconv.ParseBool(c.QueryParam("s"))
	id := c.Param("id")
	if id != "" {
		input, err := helpers.ParseDataFromReq(c)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
		sess := c.Session(false)
		if sess != nil {
			dc := GetDataManager(sess)
			formSrc, _ := dc.GetData("src" + id)
			if formSrc == nil {
				item, err := c.System().DB.Form.Get(sess, id)
				if err == nil && item != nil {
					formSrc = item.Data
				}
			}
			presistedData, err := dc.GetData(id)
			if err != nil {
				return c.NoContent(http.StatusBadRequest)
			}
			pd := file.MapIO(presistedData)
			pd.MergeWith(input)
			verrs, err := form.Validate(pd, formSrc, submit)
			if err == nil {
				if len(verrs) > 0 {
					verrs.Translate(func(key string, args ...string) string {
						return c.I18n().T(key, args)
					})
					return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"errors": verrs})
				} else {
					err = dc.PutData(id, input)
					if err == nil {
						return c.NoContent(http.StatusOK)
					}
				}
			}
		}
	}
	return c.NoContent(http.StatusBadRequest)
}

func GetFileIdFieldName(e echo.Context) error {
	c := e.(*www.Context)
	fieldname := c.Param("fieldname")
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	id := c.Param("id")
	if id != "" {
		dc := GetDataManager(sess)
		fi, err := dc.GetDataFile(id, fieldname)
		if err == nil {
			c.Response().Committed = true //prevents from-> http: multiple response.WriteHeader calls
			err = c.System().DB.Files.Read(fi.Path(), c.Response().Writer)
			if err == nil {
				return c.NoContent(http.StatusOK)
			}
		}
	}
	return c.NoContent(http.StatusNotFound)
}

func PostFileIdFieldName(e echo.Context) error {
	c := e.(*www.Context)
	id := c.Param("id")
	fieldname := c.Param("fieldname")
	if fieldname != "" && id != "" {
		sess := c.Session(false)
		if sess != nil {
			dc := GetDataManager(sess)
			formSrc, _ := dc.GetData("src" + id)
			if formSrc == nil {
				item, err := c.System().DB.Form.Get(sess, id)
				if err == nil && item != nil {
					formSrc = item.Data
				}
			}
			fileName, _ := url.QueryUnescape(c.Request().Header.Get("File-Name"))
			if fileName == "" {
				fileName = "unknown"
			}
			defer c.Request().Body.Close()
			buf, err := form.ValidateFile(c.Request().Body, formSrc, fieldname)
			if err != nil {
				if er, ok := err.(validate.ErrorMap); ok {
					er.Translate(func(key string, args ...string) string {
						return c.I18n().T(key, args)
					})
					return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"errors": er})
				}
				return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"errors": err})
			}
			//TODO improve efficienty of file move
			err = dc.PutDataFile(c.System().DB.Files, id, fieldname,
				file.Meta{
					Name:        fileName,
					ContentType: c.Request().Header.Get("Content-Type"),
				},
				bytes.NewBuffer(buf),
			)
			if err == nil {
				fData, _ := dc.GetDataByPath(id, fieldname)
				return c.JSON(http.StatusOK, fData)
			}
		}
	}
	return c.NoContent(http.StatusBadRequest)
}

func GetDataManager(sess *sys.Session) form.DataManager {
	var dc form.DataManager
	v, ok := sess.GetMemory("testDC")
	if ok {
		dc = v.(form.DataManager)
	} else {
		dc = form.NewDataManager(sess.GetSessionDir())
		sess.PutMemory("testDC", dc)
	}
	return dc
}
