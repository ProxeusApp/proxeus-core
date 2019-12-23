package formbuilder

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/labstack/echo"

	"github.com/ProxeusApp/proxeus-core/main/handlers/api"
	"github.com/ProxeusApp/proxeus-core/storage"

	"io/ioutil"

	"github.com/ProxeusApp/proxeus-core/main/helpers"
	"github.com/ProxeusApp/proxeus-core/main/www"
	"github.com/ProxeusApp/proxeus-core/sys/model"
)

func ExportForms(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	var id []string
	if c.QueryParam("id") != "" {
		id = []string{c.QueryParam("id")}
	} else if c.QueryParam("contains") != "" {
		items, _ := c.System().DB.Form.List(sess, c.QueryParam("contains"), storage.Options{Limit: 1000})
		if len(items) > 0 {
			id = make([]string, len(items))
			for i, item := range items {
				id[i] = item.ID
			}
		}
	}
	return api.Export(sess, []storage.ImporterExporter{c.System().DB.Form}, c, id...)
}

func ListHandler(e echo.Context) error {
	c := e.(*www.Context)
	contains := c.QueryParam("c")
	settings := helpers.RequestOptions(c)
	sess := c.Session(false)
	if sess != nil {
		dat, err := c.System().DB.Form.List(sess, contains, settings)
		if err != nil || dat == nil {
			if err == model.ErrAuthorityMissing {
				return c.NoContent(http.StatusUnauthorized)
			}
			return c.NoContent(http.StatusNotFound)
		}
		return c.JSON(http.StatusOK, dat)
	}
	return c.NoContent(http.StatusUnauthorized)
}

func GetOneFormHandler(e echo.Context) error {
	c := e.(*www.Context)
	formID := c.Param("formID")
	sess := c.Session(true)
	if sess != nil {
		item, err := c.System().DB.Form.Get(sess, formID)
		if err == nil {
			return c.JSON(http.StatusOK, item)
		}
	}
	return c.NoContent(http.StatusNotFound)
}

func UpdateFormHandler(e echo.Context) error {
	c := e.(*www.Context)
	ID := c.QueryParam("id")
	sess := c.Session(false)
	if sess != nil {
		if strings.Contains(c.Request().Header.Get("Content-Type"), "application/json") {
			body, _ := ioutil.ReadAll(c.Request().Body)
			item := model.FormItem{}
			err := json.Unmarshal(body, &item)
			if err == nil {
				item.ID = ID
				err = c.System().DB.Form.Put(sess, &item)
				if err == nil {
					return c.JSON(http.StatusOK, item)
				}
			}
		}
	}
	return c.NoContent(http.StatusBadRequest)
}

func GetComponentsHandler(e echo.Context) error {
	c := e.(*www.Context)
	id := c.QueryParam("id")
	contains := c.QueryParam("c")
	sess := c.Session(false)
	if sess != nil {
		var dat interface{}
		var err error
		if id != "" {
			dat, err = c.System().DB.Form.GetComp(sess, id)
		} else {
			settings := helpers.RequestOptions(c)
			dat, err = c.System().DB.Form.ListComp(sess, contains, settings)
		}
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
		return c.JSON(http.StatusOK, dat)
	}
	return c.NoContent(http.StatusNotFound)
}

func SetComponentHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess != nil {
		if strings.Contains(c.Request().Header.Get("Content-Type"), "application/json") {
			body, _ := ioutil.ReadAll(c.Request().Body)
			var comp model.FormComponentItem
			err := json.Unmarshal(body, &comp)
			if err == nil {
				err = c.System().DB.Form.PutComp(sess, &comp)
				if err != nil {
					return c.NoContent(http.StatusInternalServerError)
				}
				return c.JSON(http.StatusOK, comp.ID)
			}
		}
	}
	return c.NoContent(http.StatusBadRequest)
}

func DeleteHandler(e echo.Context) error {
	c := e.(*www.Context)
	ID := c.Param("ID")
	sess := c.Session(false)
	if sess != nil {
		err := c.System().DB.Form.Delete(sess, ID)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}
	return c.NoContent(http.StatusBadRequest)
}

func DeleteComponentHandler(e echo.Context) error {
	c := e.(*www.Context)
	id := c.Param("id")
	if id != "" {
		sess := c.Session(false)
		if sess != nil {
			err := c.System().DB.Form.DelComp(sess, id)
			if err != nil {
				return err
			}
			return c.JSON(http.StatusOK, id)
		}
	}
	return c.NoContent(http.StatusBadRequest)
}

func VarsHandler(e echo.Context) error {
	c := e.(*www.Context)
	containing := c.QueryParam("c")
	sess := c.Session(false)
	if sess != nil {
		resultVars, err := c.System().DB.Form.Vars(sess, containing, helpers.RequestOptions(c))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, resultVars)
	}
	return c.NoContent(http.StatusNotFound)
}
