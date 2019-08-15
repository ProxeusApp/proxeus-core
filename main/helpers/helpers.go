package helpers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path"
	"strconv"
	"strings"

	"github.com/labstack/echo"
)

func AbsoluteURL(c echo.Context, uri ...string) string {
	return c.Scheme() + "://" + path.Join(c.Request().Host, path.Join(uri...))
}

func ReadReqSettings(c echo.Context) map[string]interface{} {
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
		settings := map[string]interface{}{
			"limit": ilimit,
			"index": iindex,
		}
		if include != "" {
			var includeMap map[string]interface{}
			json.Unmarshal([]byte(include), &includeMap)
			settings["include"] = includeMap
		}
		if excludes != "" {
			var excludesMap map[string]interface{}
			json.Unmarshal([]byte(excludes), &excludesMap)
			settings["exclude"] = excludesMap
		}
		if metaOnly != "" {
			metaOnlyBool, err := strconv.ParseBool(metaOnly)
			if err == nil {
				settings["metaOnly"] = metaOnlyBool
			}
		}
		return settings
	} else {
		settngs := &struct {
			MetaOnly bool                   `json:"metaOnly"`
			Index    int                    `json:"index"`
			Limit    int                    `json:"limit"`
			Exclude  map[string]interface{} `json:"exclude"`
			Include  map[string]interface{} `json:"include"`
		}{}
		bts, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return nil
		}
		err = json.Unmarshal(bts, &settngs)
		if err != nil {
			return nil
		}
		return map[string]interface{}{
			"limit":    settngs.Limit,
			"index":    settngs.Index,
			"exclude":  settngs.Exclude,
			"include":  settngs.Include,
			"metaOnly": settngs.MetaOnly,
		}
	}
	return nil
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
