package app

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ProxeusApp/proxeus-core/main/www"
	"github.com/ProxeusApp/proxeus-core/sys/eio"
	"github.com/ProxeusApp/proxeus-core/sys/file"
	"github.com/ProxeusApp/proxeus-core/sys/form"
	"github.com/ProxeusApp/proxeus-core/sys/model"
	"github.com/ProxeusApp/proxeus-core/sys/workflow"

	uuid "github.com/satori/go.uuid"
)

type ExecuteAtOnceContext struct {
	data      map[string]interface{}
	c         *www.Context
	a         model.Auth
	filePaths []string
	lang      string
}

func ExecuteWorkflowAtOnce(c *www.Context, a model.Auth, wfi *model.WorkflowItem, inputData map[string]interface{}) error {
	eaoc := &ExecuteAtOnceContext{
		data: inputData,
		a:    a,
		c:    c,
		lang: c.Lang(),
	}
	defer func() {
		for _, path := range eaoc.filePaths {
			c.System().DB.Files.Delete(path)
		}
	}()
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
			"priceretriever": {InitImplFunc: func(n *workflow.Node) (workflow.NodeIF, error) {
				return &priceRetrieverNode{ctx2: eaoc}, nil
			}, Background: true},
			"externalNode": {InitImplFunc: func(n *workflow.Node) (workflow.NodeIF, error) {
				return &externalNode{ctx2: eaoc, n: n}, nil
			}, Background: true},
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
	fileCount := len(eaoc.filePaths)
	if fileCount == 0 {
		return nil
	}
	if fileCount > 1 {
		c.Response().Header().Set("Content-Disposition", "attachment;filename=\"docs.zip\"")
		c.Response().Header().Set("Content-Type", "application/zip")
		c.Response().Committed = true
		return eaoc.WriteZIP(eaoc.filePaths, c.Response().Writer)
	}
	var buf bytes.Buffer
	err = c.System().DB.Files.Read(eaoc.filePaths[0], &buf)
	if err != nil {
		return err
	}
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=\"%s\"", filepath.Base(eaoc.filePaths[0])))
	c.Response().Header().Set("Content-Type", "application/pdf")
	c.Response().Header().Set("Content-Length", fmt.Sprintf("%d", buf.Len()))
	c.Response().Committed = true
	_, err = io.Copy(c.Response().Writer, &buf)
	return err
}

func (me *ExecuteAtOnceContext) WriteZIP(filePaths []string, writer io.Writer) error {
	zipWriter := zip.NewWriter(writer)
	defer zipWriter.Close()
	zipFile := func(p string) error {
		var buf bytes.Buffer
		err := me.c.System().DB.Files.Read(p, &buf)
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(file.InMemoryFileInfo{Path: filepath.Base(p), Len: buf.Len()})
		if err != nil {
			return err
		}
		// Change to deflate to gain better compression
		// see http://golang.org/pkg/archive/zip/#pkg-constants
		header.Method = zip.Deflate
		odtZipWriter, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(odtZipWriter, &buf)
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
	return nil, errors.New("not found")
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
		dsResp, err = me.ctx.c.System().DS.Compile(me.ctx.c.System().DB.Files,
			eio.Template{
				Data:         me.ctx.getData(),
				TemplatePath: tmpl.Path(),
				EmbedError:   false,
			})
		if err != nil {
			return false, err
		}
		me.renderedPathName = filepath.Join("wfAtOnce", uuid.NewV4().String(), tmpl.NameWithExt("pdf"))
		err = me.ctx.c.System().DB.Files.Write(me.renderedPathName, dsResp.Body)
		if err != nil {
			return false, err
		}
		me.ctx.filePaths = append(me.ctx.filePaths, me.renderedPathName)
	} else {
		return false, errors.New("not found")
	}
	return true, nil
}

func (me *EAODocTmplNodeImpl) Remove(n *workflow.Node) {}
func (me *EAODocTmplNodeImpl) Close()                  {}
