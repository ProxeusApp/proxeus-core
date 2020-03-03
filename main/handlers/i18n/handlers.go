package i18n

import (
	"net/http"

	"github.com/ProxeusApp/proxeus-core/storage/portable"

	"github.com/labstack/echo"

	"github.com/ProxeusApp/proxeus-core/main/handlers/api"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/ProxeusApp/proxeus-core/main/handlers/helpers"
	"github.com/ProxeusApp/proxeus-core/main/www"
	"github.com/ProxeusApp/proxeus-core/sys"
)

// Render app.html
func IndexHandler(e echo.Context) error {
	c := e.(*www.Context)
	bts, err := sys.ReadAllFile("app.html")
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	return c.HTMLBlob(http.StatusOK, bts)
}

// Returns a json containing i18n metadata {
//		"langListSize",
//		"langList",
//		"activeLangs",
//		"langFallback",
//		"fallbackTranslations",
//	}
func MetaHandler(e echo.Context) error {
	c := e.(*www.Context)
	fallback, _ := c.System().DB.I18n.GetFallback()
	langs, _ := c.System().DB.I18n.GetAllLangs()
	activeLangs, _ := c.System().DB.I18n.GetLangs(true)
	fallbackTranslations, _ := c.System().DB.I18n.GetAll(fallback)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"langListSize":         len(langs),
		"langList":             langs,
		"activeLangs":          activeLangs,
		"langFallback":         fallback,
		"fallbackTranslations": fallbackTranslations,
	})
}

// Export the i18n definitions of the platform
func ExportI18n(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	key := c.QueryParam("k")
	value := c.QueryParam("v")
	var id []string
	if c.QueryParam("id") != "" {
		id = []string{c.QueryParam("id")}
	} else if key != "" || value != "" {
		//TODO specific translations could be added
	}
	return api.Export(sess, []portable.EntityType{portable.I18n}, c, id...)
}

// Returns all language keys
func AllHandler(e echo.Context) error {
	c := e.(*www.Context)
	wi18n := www.NewI18n(c.System().DB.I18n, c)
	return c.JSON(http.StatusOK, wi18n.GetAll())
}

// Search the key of a given key or a given translation
func FindHandler(e echo.Context) error {
	c := e.(*www.Context)
	key := c.QueryParam("k")
	value := c.QueryParam("v")
	settings := helpers.RequestOptions(c)
	da, err := c.System().DB.I18n.Find(key, value, settings)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, da)
}

// Returns a list of matching i18n keys for a given form field
func FormBuilderI18nSearchHandler(e echo.Context) error {
	c := e.(*www.Context)
	containing := c.QueryParam("c")
	settings := helpers.RequestOptions(c)
	da, err := c.System().DB.I18n.Find(containing, containing, settings)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, da)
}

// Update the i18n definitions
func UpdateHandler(e echo.Context) error {
	c := e.(*www.Context)
	if strings.Contains(c.Request().Header.Get("Content-Type"), "application/json") {
		body, _ := ioutil.ReadAll(c.Request().Body)
		trans := make(map[string]map[string]string)
		err := json.Unmarshal(body, &trans)
		if err == nil && len(trans) > 0 {
			for key, item := range trans {
				for lang, text := range item {
					err = c.System().DB.I18n.Put(lang, key, text)
					if err != nil {
						break
					}
				}
			}
			if err == nil {
				return c.JSON(http.StatusOK, map[string]interface{}{"msg": c.I18n().T("updated")})
			} else {
				fmt.Println(err)
			}
		}
	}
	return c.NoContent(http.StatusBadRequest)
}

// Set the language fallback
// @param lang => string
func SetFallbackHandler(e echo.Context) error {
	c := e.(*www.Context)
	lang := c.QueryParam("lang")
	if lang == "" {
		err := c.System().DB.I18n.PutFallback(lang)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
	}
	return c.NoContent(http.StatusOK)
}

// Switch the language
func LangSwitchHandler(e echo.Context) error {
	c := e.(*www.Context)
	if strings.Contains(c.Request().Header.Get("Content-Type"), "application/json") {
		body, _ := ioutil.ReadAll(c.Request().Body)
		params := make(map[string]bool)
		err := json.Unmarshal(body, &params)
		if err == nil && params != nil {
			for code, active := range params {
				err = c.System().DB.I18n.PutLang(code, active)
				if err != nil {
					return c.NoContent(http.StatusInternalServerError)
				}
			}
			respData := map[string]interface{}{"msg": c.I18n().T("updated")}
			return c.JSON(http.StatusOK, respData)
		}
	}
	return c.NoContent(http.StatusBadRequest)
}

// Translate a given set of language keys
func TranslateHandler(e echo.Context) error {
	c := e.(*www.Context)
	if strings.Contains(c.Request().Header.Get("Content-Type"), "application/json") {
		body, _ := ioutil.ReadAll(c.Request().Body)
		var arrOrStr interface{}
		err := json.Unmarshal(body, &arrOrStr)
		if err == nil {
			wi18n := c.I18n()
			arr, ok := arrOrStr.([]interface{})
			if ok {
				l := len(arr)
				resArr := make([]string, l)
				if l > 0 {
					for i, key := range arr {
						resArr[i] = wi18n.T(key)
					}
				}
				return c.JSON(http.StatusOK, resArr)
			} else {
				i18nKeyStr, ok := arrOrStr.(string)
				if ok {
					return c.JSON(http.StatusOK, wi18n.T(i18nKeyStr))
				}
			}
		}
	}
	return c.JSON(http.StatusOK, "")
}
