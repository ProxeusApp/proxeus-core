package formbuilder

import (
	"net/http"
	"strings"

	"github.com/ProxeusApp/proxeus-core/service"

	"github.com/ProxeusApp/proxeus-core/storage/portable"

	"github.com/labstack/echo"

	"github.com/ProxeusApp/proxeus-core/main/handlers/api"

	"github.com/ProxeusApp/proxeus-core/main/handlers/helpers"
	"github.com/ProxeusApp/proxeus-core/main/www"
	"github.com/ProxeusApp/proxeus-core/sys/model"
)

var (
	formService          service.FormService
	formComponentService service.FormComponentService
)

func Init(formComponentS service.FormComponentService, formS service.FormService) {
	formComponentService = formComponentS
	formService = formS
}

// Export a form by ID
func ExportForms(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	id := c.QueryParam("id")
	contains := c.QueryParam("contains")
	exportId := formService.ExportForms(sess, id, contains)
	return api.Export(sess, []portable.EntityType{portable.Form}, c, exportId...)
}

// Returns a list of all forms
func ListHandler(e echo.Context) error {
	c := e.(*www.Context)
	contains := c.QueryParam("c")
	settings := helpers.RequestOptions(c)
	sess := c.Session(false)
	if sess != nil {
		dat, err := formService.List(sess, contains, settings)
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

// Returns a form by formID
func GetOneFormHandler(e echo.Context) error {
	c := e.(*www.Context)
	formID := c.Param("formID")
	sess := c.Session(true)
	if sess != nil {
		item, err := formService.Get(sess, formID)
		if err != nil {
			return c.String(http.StatusNotFound, err.Error())
		} else {
			return c.JSON(http.StatusOK, item)
		}
	}
	return c.NoContent(http.StatusNotFound)
}

// Update a form's formSrc
func UpdateFormHandler(e echo.Context) error {
	c := e.(*www.Context)
	ID := c.QueryParam("id")
	sess := c.Session(false)
	var err error
	if sess == nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if !strings.Contains(c.Request().Header.Get("Content-Type"), "application/json") {
		return c.String(http.StatusBadRequest, err.Error())
	}
	item, err := formService.UpdateForm(sess, ID, c.Request().Body)
	if err == nil {
		return c.JSON(http.StatusOK, item)
	}
	return c.String(http.StatusBadRequest, err.Error())
}

// Remove a form by its ID
func DeleteHandler(e echo.Context) error {
	c := e.(*www.Context)
	ID := c.Param("ID")
	sess := c.Session(false)
	if sess != nil {
		err := formService.Delete(sess, ID)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}
	return c.NoContent(http.StatusBadRequest)
}

// Returns all form-vars for using in the template IDe
func VarsHandler(e echo.Context) error {
	c := e.(*www.Context)
	containing := c.QueryParam("c")
	sess := c.Session(false)
	if sess != nil {
		resultVars, err := formService.Vars(sess, containing, helpers.RequestOptions(c))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, resultVars)
	}
	return c.NoContent(http.StatusNotFound)
}

// Return a component by ID
func GetComponentsHandler(e echo.Context) error {
	c := e.(*www.Context)
	id := c.QueryParam("id")
	contains := c.QueryParam("c")
	sess := c.Session(false)
	if sess != nil {
		var dat interface{}
		var err error
		if id != "" {
			dat, err = formComponentService.GetComp(sess, id)
		} else {
			settings := helpers.RequestOptions(c)
			dat, err = formComponentService.ListComp(sess, contains, settings)
		}
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
		return c.JSON(http.StatusOK, dat)
	}
	return c.NoContent(http.StatusNotFound)
}

// Update a component
func SetComponentHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)

	if sess == nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if !strings.Contains(c.Request().Header.Get("Content-Type"), "application/json") {
		return c.NoContent(http.StatusBadRequest)
	}

	comp, err := formComponentService.SetComp(sess, c.Request().Body)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, comp.ID)
}

// Remove a component
func DeleteComponentHandler(e echo.Context) error {
	c := e.(*www.Context)
	id := c.Param("id")
	if id != "" {
		sess := c.Session(false)
		if sess != nil {
			err := formComponentService.DelComp(sess, id)
			if err != nil {
				return err
			}
			return c.JSON(http.StatusOK, id)
		}
	}
	return c.NoContent(http.StatusBadRequest)
}
