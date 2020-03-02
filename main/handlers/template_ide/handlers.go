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

	"github.com/ProxeusApp/proxeus-core/sys"

	"github.com/ProxeusApp/proxeus-core/service"

	"github.com/ProxeusApp/proxeus-core/storage/portable"

	"github.com/labstack/echo"

	"github.com/ProxeusApp/proxeus-core/main/handlers/api"
	"github.com/ProxeusApp/proxeus-core/main/handlers/helpers"
	"github.com/ProxeusApp/proxeus-core/main/www"
	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/sys/eio"
	"github.com/ProxeusApp/proxeus-core/sys/file"
	"github.com/ProxeusApp/proxeus-core/sys/model"
)

var (
	templateS service.TemplateDocumentService
	formS     service.FormService
)

func Init(templateService service.TemplateDocumentService, formService service.FormService) {
	templateS = templateService
	formS = formService
}

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
	dc := service.GetDataManager(sess)
	dataMap, files := dc.GetAllDataFilePathNameOnly()
	if dataMap == nil {
		dataMap = map[string]interface{}{}
	}
	dataMap = map[string]interface{}{"input": dataMap}
	format := eio.Format(c.QueryParam("format"))
	resp := c.Response()
	resp.Header().Set("Content-Disposition", fmt.Sprintf("%s;filename=\"%s\"", inlineOrAttachment, file.NameWithExt(fileName, format.String())))
	dsResp, err := c.System().DS.Compile(c.System().DB.Files,
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
		items, _ := templateS.List(sess, c.QueryParam("contains"), storage.Options{Limit: 1000})
		if len(items) > 0 {
			id = make([]string, len(items))
			for i, item := range items {
				id[i] = item.ID
			}
		}
	}
	return api.Export(sess, []portable.EntityType{portable.Template}, c, id...)
}

func IdeFormHandler(e echo.Context) error {
	c := e.(*www.Context)
	contains := c.QueryParam("c")
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusBadRequest)
	}
	settings := helpers.RequestOptions(c)
	settings.MetaOnly = false
	dat, err := formS.List(sess, contains, settings)
	if err != nil || dat == nil {
		return c.NoContent(http.StatusNotFound)
	}
	dc := service.GetDataManager(sess)
	for _, it := range dat {
		if it.Data != nil {
			it.Data["data"], _ = dc.GetData(it.ID)
		}
	}
	return c.JSON(http.StatusOK, dat)
}

func ListHandler(e echo.Context) error {
	c := e.(*www.Context)
	contains := c.QueryParam("c")
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	settings := helpers.RequestOptions(c)
	dat, err := templateS.List(sess, contains, settings)
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
	if sess == nil {
		return c.NoContent(http.StatusBadRequest)
	}
	if !strings.Contains(c.Request().Header.Get("Content-Type"), "application/json") {
		return c.NoContent(http.StatusBadRequest)
	}
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	item := model.TemplateItem{}
	err = json.Unmarshal(body, &item)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	item.ID = ID
	err = templateS.Put(sess, &item)
	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}
	return c.JSON(http.StatusOK, item)

}

func OneTmplHandler(e echo.Context) error {
	c := e.(*www.Context)
	id := c.Param("id")
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusBadRequest)
	}
	item, err := templateS.Get(sess, id)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, item)

}

func DownloadTemplateHandler(e echo.Context) error {
	c := e.(*www.Context)
	id := c.Param("id")
	lang := c.Param("lang")
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusNotFound)
	}
	_, raw := c.QueryParams()["raw"]
	var err error
	fi, err := templateS.GetTemplate(sess, id, lang)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	if raw {
		c.Response().Header().Set("Content-Disposition", fmt.Sprintf("%s;filename=\"%s\"", "attachment", fi.Name()))
		c.Response().Committed = true //prevents from-> http: multiple response.WriteHeader calls
		err = templateS.ReadFile(fi.Path(), c.Response().Writer)
	} else {
		err = rendererHelper(c, fi.Path(), fi.Name())
	}
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			return c.NoContent(http.StatusServiceUnavailable)
		}
		return c.NoContent(http.StatusBadRequest)
	}
	return c.NoContent(http.StatusOK)
}

func IdeGetTmpAssDownload(e echo.Context) error {
	c := e.(*www.Context)
	dsResp, err := templateS.DownloadExtension(c.QueryParam("os"))
	if err != nil {
		fmt.Println("after download", err)
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
	id := c.Param("ID")
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusBadRequest)
	}
	err := templateS.DeleteTemplateFiles(sess, id)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.NoContent(http.StatusOK)
}

func IdeGetDownloadHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusNotFound)
	}
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
		var err error
		exists, err = templateS.Exists(sess.FilePath(id + lang))
		tmplPath = sess.FilePath(id + lang)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}
	if !exists {
		//load it from the persistent store
		fi, err := templateS.GetTemplate(sess, id, lang)
		if err != nil {
			return c.NoContent(http.StatusNotFound)
		}
		tmplPath = fi.Path()
		tmplName = fi.Name()
	}

	if _, ok := c.QueryParams()["raw"]; ok { //if confirm var exists
		if tmplPath != "" {
			c.Response().Header().Set("Content-Disposition", fmt.Sprintf("%s;filename=\"%s\"", "attachment", tmplName))
			return templateS.ReadFile(tmplPath, c.Response().Writer)
		}
		return c.NoContent(http.StatusBadRequest)
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
	return c.NoContent(http.StatusOK)
}

func IdeGetDeleteHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusBadRequest)
	}
	id := c.Param("id")
	lang := c.Param("lang")
	err := ideDelete(id, lang, sess)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.NoContent(http.StatusOK)
}

func IdePostUploadHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusBadRequest)
	}
	id := c.Param("id")
	lang := c.Param("lang")
	fileName, _ := url.QueryUnescape(c.Request().Header.Get("File-Name"))
	dataReader := c.Request().Body
	if fileName == "" {
		fileName = "unknown"
	}
	sess.Put(id+lang, fileName)
	err := sess.WriteFile(id+lang, dataReader)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	sess.Put("activeTmpl"+id, lang)
	c.Request().Body.Close()
	return c.NoContent(http.StatusOK)
}

func DeleteTemplateHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusBadRequest)
	}
	id := c.Param("id")
	lang := c.Param("lang")
	err := templateS.DeleteTemplate(sess, id, lang)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.NoContent(http.StatusOK)
}

func IdeSetActiveHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusBadRequest)
	}
	id := c.Param("id")
	lang := c.Param("lang")
	sess.Put("activeTmpl"+id, lang)
	return c.NoContent(http.StatusOK)
}

func VarsTemplateHandler(e echo.Context) error {
	c := e.(*www.Context)
	contains := c.QueryParam("c")
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusBadRequest)
	}
	settings := helpers.RequestOptions(c)
	dat, err := templateS.GetTemplateVars(sess, contains, settings)
	if err != nil || len(dat) == 0 {
		return c.NoContent(http.StatusNotFound)
	}
	for i, a := range dat {
		dat[i] = strings.TrimPrefix(a, "input.")
	}
	return c.JSON(http.StatusOK, dat)

}

func UploadTemplateHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusBadRequest)
	}
	id := c.Param("id")
	lang := c.Param("lang")
	contentType := c.Request().Header.Get("Content-Type")
	fileName, _ := url.QueryUnescape(c.Request().Header.Get("File-Name"))
	data, err := ioutil.ReadAll(c.Request().Body)

	err = templateS.SaveTemplate(sess, id, lang, contentType, fileName, data)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	//remove pending file from the session
	err = ideDelete(id, lang, sess)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	return c.NoContent(http.StatusOK)
}

func ideDelete(id, lang string, sess *sys.Session) error {
	sess.Delete(id + lang)
	return sess.DeleteFile(id + lang)
}
