package helpers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path"
	"strconv"
	"strings"

	"github.com/ProxeusApp/proxeus-core/storage"

	"github.com/labstack/echo"
)

// Gets the absolute url of a path (uri)
func AbsoluteURL(c echo.Context, uri ...string) string {
	return AbsoluteURLWithScheme(c.Scheme(), c.Request().Host, uri...)
}

// Gets the absolute url of a path (uri)
func AbsoluteURLWithScheme(scheme, host string, uri ...string) string {
	return scheme + "://" + path.Join(host, path.Join(uri...))
}

func RequestOptions(c echo.Context) storage.Options {
	if strings.ToLower(c.Request().Method) == "get" {
		metaOnly := c.QueryParam("m")
		limit := c.QueryParam("l")
		excludes := c.QueryParam("e")
		include := c.QueryParam("in")
		index := c.QueryParam("i")

		ilimit := 20
		if limit != "" {
			i, err := strconv.Atoi(limit)
			if err == nil {
				ilimit = i
			}
		}
		iindex := 0
		if index != "" {
			i, err := strconv.Atoi(index)
			if err == nil {
				iindex = i
			}
		}
		settings := storage.Options{
			Limit: ilimit,
			Index: iindex,
		}
		if include != "" {
			var includeMap map[string]interface{}
			json.Unmarshal([]byte(include), &includeMap)
			settings.Include = includeMap
		}
		if excludes != "" {
			var excludesMap map[string]interface{}
			json.Unmarshal([]byte(excludes), &excludesMap)
			settings.Exclude = excludesMap
		}
		if metaOnly != "" {
			metaOnlyBool, err := strconv.ParseBool(metaOnly)
			if err == nil {
				settings.MetaOnly = metaOnlyBool
			}
		} else {
			// default is meta only
			settings.MetaOnly = true
		}
		return settings
	}

	var settings storage.Options
	bts, _ := ioutil.ReadAll(c.Request().Body)
	json.Unmarshal(bts, &settings)
	return settings
}

func ParseDataFromReq(c echo.Context) (map[string]interface{}, error) {
	req := c.Request()
	formInput := make(map[string]interface{})
	var err error
	if strings.HasPrefix(req.Header.Get("Content-Type"), "application/json") {
		err = json.NewDecoder(req.Body).Decode(&formInput)
		if err != nil {
			return formInput, err
		}
	} else {
		if err := req.ParseForm(); err != nil {
			return formInput, err
		}
		delete(req.Form, "s") // go includes query params here as well
		for k, v := range req.Form {
			if len(v) == 1 && v[0] == "" {
				continue
			}
			formInput[k] = v
		}
	}
	if _, ok := formInput[""]; ok {
		return map[string]interface{}{}, errors.New("can't post data with empty key, fix your form definition")
	}
	return formInput, err
}
