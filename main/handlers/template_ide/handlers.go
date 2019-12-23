package template_ide

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/labstack/echo"

	"github.com/ProxeusApp/proxeus-core/main/handlers/api"
	"github.com/ProxeusApp/proxeus-core/main/handlers/formbuilder"
	"github.com/ProxeusApp/proxeus-core/main/helpers"
	"github.com/ProxeusApp/proxeus-core/main/www"
	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/sys/eio"
	"github.com/ProxeusApp/proxeus-core/sys/file"
	"github.com/ProxeusApp/proxeus-core/sys/model"
	"github.com/ProxeusApp/proxeus-core/sys/session"
)

var rendererHelper = func(e echo.Context, tmplPath, fileName string) error {
	c := e.(*www.Context)
	inlineOrAttachment := c.QueryParam("inline")
	if inlineOrAttachment == "" {
		inlineOrAttachment = "attachment"
	} else {
		inlineOrAttachment = "inline"
	}
	sess := c.Session(false)
	if sess == nil || tmplPath == "" {
		return os.ErrInvalid
	}
	dc := formbuilder.GetDataManager(sess)
	dataMap, files := dc.GetAllDataFilePathNameOnly()
	if dataMap == nil {
		dataMap = map[string]interface{}{}
	}
	dataMap = map[string]interface{}{"input": dataMap}
	format := eio.Format(c.QueryParam("format"))
	resp := c.Response()
	resp.Header().Set("Content-Disposition", fmt.Sprintf("%s;filename=\"%s\"", inlineOrAttachment, file.NameWithExt(fileName, format.String())))
	dsResp, err := c.System().DS.Compile(
		eio.Template{
			Format:       format,
			Data:         dataMap,
			TemplatePath: tmplPath,
			Assets:       files,
			EmbedError:   true,
		})
	if err != nil {
		return err
	}
	resp.Header().Set("Content-Type", dsResp.Header.Get("Content-Type"))
	resp.Header().Set("Content-Length", dsResp.Header.Get("Content-Length"))

	defer dsResp.Body.Close()
	resp.Committed = true //prevents from-> http: multiple response.WriteHeader calls
	_, err = io.Copy(resp.Writer, dsResp.Body)
	return err
}

func ExportTemplate(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	var id []string
	if c.QueryParam("id") != "" {
		id = []string{c.QueryParam("id")}
	} else if c.QueryParam("contains") != "" {
		items, _ := c.System().DB.Template.List(sess, c.QueryParam("contains"), storage.Options{Limit: 1000})
		if len(items) > 0 {
			id = make([]string, len(items))
			for i, item := range items {
				id[i] = item.ID
			}
		}
	}
	return api.Export(sess, []storage.ImporterExporter{c.System().DB.Template}, c, id...)
}

func IdeFormHandler(e echo.Context) error {
	c := e.(*www.Context)
	contains := c.QueryParam("c")
	sess := c.Session(false)
	if sess != nil {
		settings := helpers.RequestOptions(c)
		settings.MetaOnly = false
		dat, err := c.System().DB.Form.List(sess, contains, settings)
		if err == nil && dat != nil {
			dc := formbuilder.GetDataManager(sess)
			for _, it := range dat {
				if it.Data != nil {
					it.Data["data"], _ = dc.GetData(it.ID)
				}
			}
			return c.JSON(http.StatusOK, dat)
		}
		return c.NoContent(http.StatusNotFound)
	}
	return c.NoContent(http.StatusBadRequest)
}

func ListHandler(e echo.Context) error {
	c := e.(*www.Context)
	contains := c.QueryParam("c")
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	settings := helpers.RequestOptions(c)
	dat, err := c.System().DB.Template.List(sess, contains, settings)
	if err != nil || dat == nil {
		if err == model.ErrAuthorityMissing {
			return c.NoContent(http.StatusUnauthorized)
		}
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, dat)
}

func UpdateHandler(e echo.Context) error {
	c := e.(*www.Context)
	ID := c.QueryParam("id")
	sess := c.Session(false)
	if sess != nil {
		if strings.Contains(c.Request().Header.Get("Content-Type"), "application/json") {
			body, err := ioutil.ReadAll(c.Request().Body)
			if err == nil {
				item := model.TemplateItem{}
				err = json.Unmarshal(body, &item)
				if err == nil {
					item.ID = ID
					err = c.System().DB.Template.Put(sess, &item)
					if err != nil {
						return err
					}
					return c.JSON(http.StatusOK, item)
				}
			}
		}
	}
	return c.NoContent(http.StatusBadRequest)
}

func OneTmplHandler(e echo.Context) error {
	c := e.(*www.Context)
	id := c.Param("id")
	sess := c.Session(false)
	if sess != nil {
		item, err := c.System().DB.Template.Get(sess, id)
		if err != nil {
			return c.NoContent(http.StatusNotFound)
		}
		return c.JSON(http.StatusOK, item)
	}
	return c.NoContent(http.StatusBadRequest)
}

func DownloadTemplateHandler(e echo.Context) error {
	c := e.(*www.Context)
	id := c.Param("id")
	lang := c.Param("lang")
	sess := c.Session(false)
	if sess != nil {
		fi, err := c.System().DB.Template.GetTemplate(sess, id, lang)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
		if _, ok := c.QueryParams()["raw"]; ok {
			c.Response().Header().Set("Content-Disposition", fmt.Sprintf("%s;filename=\"%s\"", "attachment", fi.Name()))
			c.Response().Committed = true //prevents from-> http: multiple response.WriteHeader calls
			_, err = fi.Read(c.Response().Writer)
			if err == nil {
				return c.NoContent(http.StatusOK)
			}
		} else {
			err = rendererHelper(c, fi.Path(), fi.Name())
			if err != nil {
				if err, ok := err.(net.Error); ok && err.Timeout() {
					return c.NoContent(http.StatusServiceUnavailable)
				}
				return c.NoContent(http.StatusBadRequest)
			}
		}
	}
	return c.NoContent(http.StatusNotFound)
}

func IdeGetTmpAssDownload(e echo.Context) error {
	c := e.(*www.Context)
	dsResp, err := c.System().DS.DownloadExtension(c.QueryParam("os"))
	fmt.Println("after download", err)
	if err != nil {
		return err
	}
	resp := c.Response()
	resp.Header().Set("Content-Type", dsResp.Header.Get("Content-Type"))
	resp.Header().Set("Content-Length", dsResp.Header.Get("Content-Length"))
	resp.Header().Set("Content-Disposition", dsResp.Header.Get("Content-Disposition"))
	_, err = io.Copy(resp.Writer, dsResp.Body)
	return err
}

func DeleteHandler(e echo.Context) error {
	c := e.(*www.Context)
	ID := c.Param("ID")
	sess := c.Session(false)
	if sess != nil {
		err := c.System().DB.Template.Delete(sess, ID)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}
	return c.NoContent(http.StatusBadRequest)
}

func IdeGetDownloadHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess != nil {
		var (
			tmplPath string
			tmplName string
			exists   bool
		)
		id := c.Param("id")
		var lang string
		sess.Get("activeTmpl"+id, &lang)
		sess.Get(id+lang, &tmplName)
		if lang == "" {
			return c.NoContent(http.StatusNotFound)
		}
		//the file wasn't deleted, don't provide it without the name
		if tmplName != "" {
			tmplPath, exists = sess.FilePath(id + lang)
		}
		if !exists {
			//load it from the persistent store
			fi, err := c.System().DB.Template.GetTemplate(sess, id, lang)
			if err != nil {
				return c.NoContent(http.StatusNotFound)
			}
			tmplPath = fi.Path()
			tmplName = fi.Name()
		}
		if _, ok := c.QueryParams()["raw"]; ok { //if confirm var exists
			if tmplPath != "" {
				c.Response().Header().Set("Content-Disposition", fmt.Sprintf("%s;filename=\"%s\"", "attachment", tmplName))
				return c.File(tmplPath)
			}
		} else {
			err := rendererHelper(c, tmplPath, tmplName)
			if err != nil {
				if err, ok := err.(net.Error); ok && err.Timeout() {
					return c.String(http.StatusServiceUnavailable, err.Error())
				}
				return c.String(http.StatusBadRequest, err.Error())
			}
			return nil
		}
	}
	return c.NoContent(http.StatusNotFound)
}

func IdeGetDeleteHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess != nil {
		id := c.Param("id")
		lang := c.Param("lang")
		err := ideDelete(id, lang, sess)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}
	return c.NoContent(http.StatusBadRequest)
}

func ideDelete(id, lang string, sess *session.Session) error {
	sess.Delete(id + lang)
	return sess.DeleteFile(id + lang)
}

func IdePostUploadHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess != nil {
		id := c.Param("id")
		lang := c.Param("lang")
		fileName, _ := url.QueryUnescape(c.Request().Header.Get("File-Name"))
		if fileName == "" {
			fileName = "unknown"
		}
		sess.Put(id+lang, fileName)
		_, err := sess.WriteFile(id+lang, c.Request().Body)
		c.Request().Body.Close()
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		sess.Put("activeTmpl"+id, lang)
		return c.NoContent(http.StatusOK)
	}
	return c.NoContent(http.StatusBadRequest)
}

func DeleteTemplateHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess != nil {
		id := c.Param("id")
		lang := c.Param("lang")
		err := c.System().DB.Template.DeleteTemplate(sess, id, lang)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}
	return c.NoContent(http.StatusBadRequest)
}

func IdeSetActiveHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess != nil {
		id := c.Param("id")
		lang := c.Param("lang")
		sess.Put("activeTmpl"+id, lang)
		return c.NoContent(http.StatusOK)
	}
	return c.NoContent(http.StatusBadRequest)
}

func VarsTemplateHandler(e echo.Context) error {
	c := e.(*www.Context)
	contains := c.QueryParam("c")
	sess := c.Session(false)
	if sess != nil {
		settings := helpers.RequestOptions(c)
		dat, err := c.System().DB.Template.Vars(sess, contains, settings)
		if err != nil || len(dat) == 0 {
			return c.NoContent(http.StatusNotFound)
		}
		for i, a := range dat {
			dat[i] = strings.TrimPrefix(a, "input.")
		}
		return c.JSON(http.StatusOK, dat)
	}
	return c.NoContent(http.StatusBadRequest)
}

func UploadTemplateHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess != nil {
		id := c.Param("id")
		lang := c.Param("lang")
		fileName, _ := url.QueryUnescape(c.Request().Header.Get("File-Name"))
		tmplMeta := &file.Meta{
			Name:        fileName,
			ContentType: c.Request().Header.Get("Content-Type"),
		}
		fi, err := c.System().DB.Template.ProvideFileInfoFor(sess, id, lang, tmplMeta)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		_, err = fi.Write(c.Request().Body)
		c.Request().Body.Close()
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		var vars []string
		vars, err = c.System().DS.Vars("input.", fi.Path())
		log.Println(vars, err)
		if err == nil {
			//error handling not important here as keeping track of vars is not crucial
			c.System().DB.Template.PutVars(sess, id, lang, vars)
		}
		//remove pending file from the session
		err = ideDelete(id, lang, sess)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
		return c.NoContent(http.StatusOK)
	}
	return c.NoContent(http.StatusBadRequest)
}
