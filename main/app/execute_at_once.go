package app

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/asdine/storm"

	"git.proxeus.com/core/central/main/www"
	"git.proxeus.com/core/central/sys/eio"
	"git.proxeus.com/core/central/sys/form"
	"git.proxeus.com/core/central/sys/model"
	"git.proxeus.com/core/central/sys/workflow"
)

type ExecuteAtOnceContext struct {
	data       map[string]interface{}
	c          *www.Context
	a          model.Authorization
	tmpDirPath string
	lang       string
}

func ExecuteWorkflowAtOnce(c *www.Context, a model.Authorization, wfi *model.WorkflowItem, inputData map[string]interface{}) error {
	tmpDirPath, err := ioutil.TempDir(os.TempDir(), "wfAtOnce")
	if err != nil {
		return err
	}
	defer func() {
		_ = os.RemoveAll(tmpDirPath)
	}()
	eaoc := &ExecuteAtOnceContext{
		data:       inputData,
		a:          a,
		c:          c,
		lang:       c.Lang(),
		tmpDirPath: tmpDirPath,
	}
	conf := workflow.Config{
		GetWorkflow: func(id string) (*workflow.Workflow, error) {
			item, err := c.System().DB.Workflow.Get(a, id)
			if err != nil {
				return nil, err
			}
			if item.Data != nil {
				return item.Data, nil
			}
			return nil, os.ErrNotExist
		},
		GetData: eaoc.getData,
		NodeImpl: map[string]*workflow.NodeDef{
			"form": {InitImplFunc: func(n *workflow.Node) (workflow.NodeIF, error) {
				return &EAOFormNodeImpl{ctx: eaoc}, nil
			}, Background: false},
			"template": {InitImplFunc: func(n *workflow.Node) (workflow.NodeIF, error) {
				return &EAODocTmplNodeImpl{ctx: eaoc}, nil
			}, Background: true},
		},
	}
	wfContext, err := workflow.New(wfi.Data, conf)
	if err != nil {
		return err
	}
	for wfContext.LoopNext() {
	}
	_, err = wfContext.Current()
	if err != nil {
		return err
	}
	filePaths, err := listFiles(eaoc.tmpDirPath)
	if err != nil {
		return err
	}
	fileCount := len(filePaths)
	if fileCount == 0 {
		return nil
	}
	if fileCount > 1 {
		c.Response().Header().Set("Content-Disposition", "attachment;filename=\"docs.zip\"")
		c.Response().Header().Set("Content-Type", "application/zip")
		c.Response().Committed = true
		return eaoc.WriteZIP(filePaths, c.Response().Writer)
	}
	renderedTmpl, err := os.Open(filePaths[0])
	if err != nil {
		return err
	}
	fstat, err := renderedTmpl.Stat()
	if err != nil {
		return err
	}
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=\"%s\"", filepath.Base(filePaths[0])))
	c.Response().Header().Set("Content-Type", "application/pdf")
	c.Response().Header().Set("Content-Length", fmt.Sprintf("%d", fstat.Size()))
	c.Response().Committed = true
	_, err = io.Copy(c.Response().Writer, renderedTmpl)
	return err
}

func (me *ExecuteAtOnceContext) WriteZIP(filePaths []string, writer io.Writer) error {
	zipWriter := zip.NewWriter(writer)
	defer func() {
		if zipWriter != nil && zipWriter.Close() != nil {
			return
		}
	}()
	zipFile := func(p string) error {
		renderedTmpl, err := os.Open(p)
		if err != nil {
			return err
		}
		defer renderedTmpl.Close()
		fstat, err := renderedTmpl.Stat()
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(fstat)
		if err != nil {
			return err
		}
		header.Name = filepath.Base(p)

		// Change to deflate to gain better compression
		// see http://golang.org/pkg/archive/zip/#pkg-constants
		header.Method = zip.Deflate
		odtZipWriter, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.CopyN(odtZipWriter, renderedTmpl, fstat.Size())
		if err != nil {
			return err
		}
		return nil
	}
	for _, p := range filePaths {
		err := zipFile(p)
		if err != nil {
			return err
		}
	}
	return nil
}

func (me *ExecuteAtOnceContext) getData() interface{} {
	return map[string]interface{}{"input": me.data}
}

type EAOFormNodeImpl struct {
	ctx *ExecuteAtOnceContext
}

type EAODocTmplNodeImpl struct {
	ctx              *ExecuteAtOnceContext
	renderedPathName string
}

func (me *EAOFormNodeImpl) getDataFor(formItem *model.FormItem, id string) (map[string]interface{}, error) {
	if formItem != nil {
		formData := map[string]interface{}{}
		vars := form.Vars(formItem.Data)
		for _, k := range vars {
			formData[k] = me.ctx.data[k]
		}
		return formData, nil
	}
	return nil, storm.ErrNotFound
}

func (me *EAOFormNodeImpl) Execute(n *workflow.Node) (proceed bool, err error) {
	formItem, err := me.ctx.c.System().DB.Form.Get(me.ctx.a, n.ID)
	if err != nil {
		return false, err
	}
	formData, err := me.getDataFor(formItem, n.ID)
	if err != nil {
		return false, err
	}
	errs, err := form.Validate(formData, form.GetFormSrc(formItem.Data), true)
	if err != nil {
		//definition error
		return false, err
	}
	if len(errs) > 0 {
		//validation errors
		return false, errs
	}
	return true, nil
}

func (me *EAOFormNodeImpl) Remove(n *workflow.Node) {}
func (me *EAOFormNodeImpl) Close()                  {}

func (me *EAODocTmplNodeImpl) Execute(n *workflow.Node) (proceed bool, err error) {
	var tmplItem *model.TemplateItem
	tmplItem, err = me.ctx.c.System().DB.Template.Get(me.ctx.a, n.ID)
	if err != nil {
		return false, err
	}
	if tmpl, ok := tmplItem.Data[me.ctx.lang]; ok {
		var dsResp *http.Response
		dsResp, err = me.ctx.c.System().DS.Compile(eio.Template{
			Data:         me.ctx.getData(),
			TemplatePath: tmpl.Path(),
			EmbedError:   false,
		})
		if err != nil {
			return false, err
		}
		me.renderedPathName = filepath.Join(me.ctx.tmpDirPath, tmpl.NameWithExt("pdf"))
		f, err := os.OpenFile(me.renderedPathName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return false, err
		}
		_, err = io.Copy(f, dsResp.Body)
		if err != nil {
			return false, err
		}
		err = f.Close()
		if err != nil {
			return false, err
		}
	} else {
		return false, storm.ErrNotFound
	}
	return true, nil
}

func (me *EAODocTmplNodeImpl) Remove(n *workflow.Node) {
	//remove the template from the app collection as it is not part of the path anymore
	if me.renderedPathName != "" {
		_ = os.Remove(me.renderedPathName)
	}
}

func (me *EAODocTmplNodeImpl) Close() {}

func listFiles(dirPath string) ([]string, error) {
	serializedFiles := make([]string, 0)
	err := rec(&serializedFiles, dirPath)
	if err != nil {
		return nil, err
	}
	return serializedFiles, nil
}

func rec(i *[]string, name string) error {
	fi, err := os.Stat(name)
	if err != nil {
		return err
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		files, err := ioutil.ReadDir(name)
		if err != nil {
			return err
		}
		for _, f := range files {
			err = rec(i, filepath.Join(name, f.Name()))
			if err != nil {
				return err
			}
		}
	case mode.IsRegular():
		*i = append(*i, name)
	}
	return nil
}
