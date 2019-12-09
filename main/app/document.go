package app

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/labstack/echo"

	"github.com/ProxeusApp/proxeus-core/sys"
	"github.com/ProxeusApp/proxeus-core/sys/db/storm"
	"github.com/ProxeusApp/proxeus-core/sys/eio"
	"github.com/ProxeusApp/proxeus-core/sys/file"
	"github.com/ProxeusApp/proxeus-core/sys/form"
	"github.com/ProxeusApp/proxeus-core/sys/model"
	"github.com/ProxeusApp/proxeus-core/sys/validate"
	"github.com/ProxeusApp/proxeus-core/sys/workflow"
)

//TODO replace DataCluster with file.MapIO. DataCluster was meant to be used for guest users only to prevent from storing data that will be never used again after the session is expired
type (
	DocumentFlowInstance struct {
		WFID             string            `json:"WFID"`
		DataID           string            `json:"dataID"`
		DataCluster      *form.DataManager `json:"dataCluster"`
		wfContext        *workflow.Engine
		system           *sys.System
		statusResult     *Status
		dirty            bool
		templateMap      map[string]*model.TemplateItem `json:"-"`
		templateMapLock  sync.RWMutex                   `json:"-"`
		Steps            []workflow.Step                `json:"steps"`
		initLock         sync.Mutex                     `json:"-"`
		wfItem           *model.WorkflowItem            `json:"-"`
		Started          bool
		startedID        string
		auth             model.Authorization
		SelectedDocLangs map[string]string
		Confirmed        bool
		unorderedData    map[string]interface{}
		allFormFields    []string
	}

	Status struct {
		TargetName  string          `json:"targetName"`
		Steps       []workflow.Step `json:"steps"`
		HasNext     bool            `json:"hasNext"`
		HasPrev     bool            `json:"hasPrev"`
		Docs        []interface{}   `json:"docs"`
		CurrentType string          `json:"currentType"`
		UserData    interface{}     `json:"userData,omitempty"`
		Data        interface{}     `json:"data,omitempty"`
	}

	FormNodeImpl struct {
		ctx       *DocumentFlowInstance
		presented bool
	}
	DocTmplNodeImpl struct {
		ctx *DocumentFlowInstance
	}
)

func NewDocumentApp(usrData *model.UserDataItem, auth model.Authorization, system *sys.System, wfid, baseFilePath string) (*DocumentFlowInstance, error) {
	if system == nil || wfid == "" {
		return nil, os.ErrInvalid
	}

	workflowItem, err := system.DB.Workflow.Get(auth, wfid)

	if err == nil && workflowItem != nil && workflowItem.Data != nil {
		doc := &DocumentFlowInstance{
			DataID:           usrData.ID,
			SelectedDocLangs: map[string]string{},
			system:           system,
			WFID:             wfid,
			DataCluster:      form.NewDataManager(baseFilePath),
			Started:          true,
			auth:             auth,
		}
		if len(usrData.Data) > 0 { // TODO REMOVE COMMENT
			doc.unorderedData = usrData.Data["input"].(map[string]interface{})
		}
		return doc, doc.init(workflowItem.Data, nil)
	}
	return nil, os.ErrNotExist
}

func (me *DocumentFlowInstance) Init(auth model.Authorization, system *sys.System) error {
	me.auth = auth
	me.system = system
	workflowItem, err := system.DB.Workflow.Get(auth, me.WFID)
	if err == nil && workflowItem != nil && workflowItem.Data != nil {
		return me.init(workflowItem.Data, me.Steps)
	}
	return os.ErrNotExist
}

func (me *DocumentFlowInstance) NeedToBeInitialized() bool {
	return me.dirty
}

func (me *DocumentFlowInstance) OnLoad() {
	me.dirty = true
	if me.DataCluster != nil {
		me.DataCluster.OnLoad()
	}
}

func (me *DocumentFlowInstance) init(wfd *workflow.Workflow, state []workflow.Step) error {
	me.dirty = false
	me.Confirmed = false
	me.statusResult = &Status{}
	me.templateClear()

	conf := workflow.Config{
		GetWorkflow: me.getWorkflow,
		State:       state,
		GetData:     me.getData, // for condition execution
		NodeImpl: map[string]*workflow.NodeDef{
			"mailsender":     {InitImplFunc: me.newMailSender, Background: true},
			"priceretriever": {InitImplFunc: me.newPriceRetriever, Background: true},
			"form":           {InitImplFunc: me.newFormNodeImpl, Background: false},
			"template":       {InitImplFunc: me.newDocTmplNodeImpl, Background: true},
		},
	}
	var err error
	me.wfContext, err = workflow.New(wfd, conf)
	return err
}

func (me *DocumentFlowInstance) isLangAvailable(lang string) (bool, error) {
	langs, err := me.system.DB.I18n.GetLangs(true)
	if err != nil {
		return false, err
	}
	for _, l := range langs {
		if l.Code == lang {
			return true, nil
		}
	}
	return false, nil
}

func (me *DocumentFlowInstance) getWorkflow(id string) (*workflow.Workflow, error) {
	item, err := me.system.DB.Workflow.Get(me.auth, id)
	if err != nil {
		return nil, err
	}
	if item.Data != nil {
		return item.Data, nil
	}
	return nil, os.ErrNotExist
}

func (me *DocumentFlowInstance) getData() interface{} {
	dd, _ := me.DataCluster.GetAllData()
	return map[string]interface{}{"input": dd}
}

func (me *DocumentFlowInstance) getDataFor(id string) (map[string]interface{}, error) {
	m, err := me.DataCluster.GetData(id)
	if m == nil && len(me.unorderedData) > 0 {
		f, err := me.system.DB.Form.Get(me.auth, id)
		if err == nil {
			frmData := map[string]interface{}{}
			for _, name := range form.Vars(f.Data) {
				if v, ok := me.unorderedData[name]; ok {
					frmData[name] = v
				}
			}
			err = me.DataCluster.PutData(id, frmData)
			if err != nil {
				return nil, err
			}
			return frmData, nil
		}
	}
	return m, err
}

func (me *DocumentFlowInstance) getDataByPath(dataPath string) (interface{}, error) {
	return me.system.DB.UserData.GetData(me.auth, me.DataID, dataPath)
}

func (me *DocumentFlowInstance) newMailSender(n *workflow.Node) (workflow.NodeIF, error) {
	return &mailSenderNode{ctx: me}, nil
}

func (me *DocumentFlowInstance) newPriceRetriever(n *workflow.Node) (workflow.NodeIF, error) {
	return &priceRetrieverNode{ctx: me}, nil
}

func (me *DocumentFlowInstance) newFormNodeImpl(n *workflow.Node) (workflow.NodeIF, error) {
	return &FormNodeImpl{ctx: me}, nil
}

func (me *DocumentFlowInstance) newDocTmplNodeImpl(n *workflow.Node) (workflow.NodeIF, error) {
	return &DocTmplNodeImpl{ctx: me}, nil
}

func (me *DocumentFlowInstance) getDataWithFiles() (d map[string]interface{}, files []string) {
	d, files = me.DataCluster.GetAllDataFilePathNameOnly()
	if d == nil {
		d = map[string]interface{}{}
	}
	d = map[string]interface{}{"input": d}
	return
}

func (me *DocumentFlowInstance) GetFile(name string) (*file.IO, error) {
	n, err := me.wfContext.Current()
	if err != nil {
		return nil, err
	}
	if n != nil && n.Type == "form" {
		return me.DataCluster.GetDataFile(n.ID, name)
	}
	return nil, os.ErrNotExist
}

func (me *DocumentFlowInstance) UpdateData(d map[string]interface{}, submit bool) (verrs validate.ErrorMap, err error) {
	var n *workflow.Node
	n, err = me.wfContext.Current()
	if n != nil && n.Type == "form" {
		verrs, err = form.Validate(d, me.statusResult.Data, false)
		if err == nil && len(verrs) == 0 {
			err = me.writeData(n, d)
		}
	}
	return
}

func (me *DocumentFlowInstance) readData(dataPath string) (interface{}, error) {
	return me.system.DB.UserData.GetData(me.auth, me.DataID, dataPath)
}

func (me *DocumentFlowInstance) writeField(n *workflow.Node, fname string, val interface{}) error {
	return me.writeData(n, map[string]interface{}{fname: val})
}

func (me *DocumentFlowInstance) writeData(n *workflow.Node, d map[string]interface{}) error {
	if n != nil {
		err := me.DataCluster.PutData(n.ID, d)
		if err == nil {
			err = me.system.DB.UserData.PutData(me.auth, me.DataID, map[string]interface{}{"input": d})
		}
		return err
	}
	return nil
}

func (me *DocumentFlowInstance) UpdateFile(name string, fm file.Meta, reader io.Reader) (verrs validate.Errors, err error) {
	var n *workflow.Node
	n, err = me.wfContext.Current()
	if err != nil {
		return
	}
	var tmpFile *validate.TmpFile
	tmpFile, err = form.ValidateFile(reader, me.statusResult.Data, name)
	if err != nil {
		if er, ok := err.(validate.Errors); ok {
			verrs = er
			err = nil
		}
	}
	defer func() {
		if tmpFile != nil {
			_ = tmpFile.Close()
		}
	}()
	if len(verrs) == 0 {
		_, err = me.DataCluster.PutDataFile(n.ID, name, fm, tmpFile)
		if err != nil {
			return
		}
		var f *file.IO //TODO better storage handling on session and persistent db, right now we just duplicate it
		var dbF *file.IO
		f, err = me.DataCluster.GetDataFile(n.ID, name)
		if err != nil {
			return
		}
		dbF, err = me.system.DB.UserData.GetDataFile(me.auth, me.DataID, "input."+name)
		if err == os.ErrNotExist {
			dbF = me.system.DB.UserData.NewFile(me.auth, f.Meta())
			err = me.system.DB.UserData.PutData(me.auth, me.DataID, map[string]interface{}{"input": map[string]interface{}{name: dbF}})
		}
		if err != nil {
			return
		}
		var of *os.File
		of, err = os.Open(f.Path())
		if err != nil {
			return
		}
		defer of.Close()
		_, err = dbF.Write(of)
	}
	return
}

func (me *DocumentFlowInstance) WF() *model.WorkflowItem {
	if me.wfItem == nil {
		me.wfItem, _ = me.system.DB.Workflow.Get(me.auth, me.WFID)
		me.wfItem.Data = nil
	}
	return me.wfItem
}

func (me *DocumentFlowInstance) Confirm(currentAppLang string, tmpls map[string]interface{}, store *storm.UserDataDB) (map[string]interface{}, map[string]*file.IO, *Status, error) {
	if me.Confirmed {
		return nil, nil, me.statusResult, nil
	}
	var err error
	defer func() {
		if err != nil {
			me.Confirmed = false
		}
	}()
	me.Confirmed = true
	var documentHex string

	dat, files := me.getDataWithFiles()
	me.statusResult.Docs = make([]interface{}, 0)
	finalDocs := make([]interface{}, 0)
	finalDocRefs := make([]*file.IO, 0)
	if len(me.templateMap) > 0 {
		var tlang string
		for id, tmplItem := range me.templateMap {
			if tmplItem != nil && tmplItem.Data != nil {
				tlang = me.SelectedDocLangs[id]
				if tlang == "" && tmpls != nil {
					if l, ok := tmpls[id]; ok {
						if ls, ok := l.(string); ok {
							if ls != "" {
								tlang = ls
							}
						}
					}
				}
				if tlang == "" {
					tlang = currentAppLang
				}
				var available bool
				if available, err = me.isLangAvailable(tlang); !available || err != nil {
					err = os.ErrNotExist
					return nil, nil, me.statusResult, err
				}
				if tmpl, ok := tmplItem.Data[tlang]; ok {
					var dsResp *http.Response
					dsResp, err = me.system.DS.Compile(eio.Template{
						Data:         dat,
						TemplatePath: tmpl.Path(),
						Assets:       files,
						EmbedError:   false,
					})
					if err != nil {
						return nil, nil, me.statusResult, err
					}

					// get the pdf into a []byte
					var documentBytes []byte
					documentBytes, err = ioutil.ReadAll(dsResp.Body)
					if err != nil {
						return nil, nil, me.statusResult, err
					}
					// Restore ds response body because the buffer reader clears the buffer
					dsResp.Body = ioutil.NopCloser(bytes.NewBuffer(documentBytes))
					// hash the document
					documentHex = crypto.Keccak256Hash(documentBytes).String()
					me.statusResult.addDoc(id, documentHex, tmplItem)
					contentType := dsResp.Header.Get("Content-Type")
					i64, _ := strconv.ParseInt(dsResp.Header.Get("Content-Length"), 10, 0)
					finalDoc := store.NewFile(
						me.auth,
						file.Meta{
							Name:        tmpl.NameWithExt("pdf"),
							ContentType: contentType,
							Size:        i64,
						},
					)
					tmplRef := store.NewFile(me.auth, tmpl.Meta())
					var tfr *os.File
					tfr, err = os.Open(tmpl.Path())
					if err != nil {
						return nil, nil, nil, err
					}
					_, err = tmplRef.Write(tfr)
					tfr.Close()
					_, err = finalDoc.Write(dsResp.Body)
					dsResp.Body.Close()
					finalDoc.Hash = documentHex
					finalDocRefs = append(finalDocRefs, tmplRef)
					finalDocs = append(finalDocs, finalDoc)
					finalDoc.SetRef(fmt.Sprintf("tmpls[%d]", len(finalDocs)-1))
					if err != nil {
						return nil, nil, me.statusResult, err
					}
				}
			}
		}
	}
	dat["tmpls"] = finalDocRefs
	dat["docs"] = finalDocs
	var item *model.UserDataItem
	item, err = store.Get(me.auth, me.DataID)
	item.Data = dat
	if err == nil {
		err = store.Put(me.auth, item)
	}
	return dat, nil, me.statusResult, err
}

func (me *DocumentFlowInstance) Next(d map[string]interface{}) (*Status, error) {
	n, _ := me.wfContext.Current()
	if n != nil && len(d) > 0 {
		err := me.DataCluster.PutData(n.ID, d)
		if err != nil {
			return me.statusResult, err
		}
		//TODO could check for guest or higher
		err = me.system.DB.UserData.PutData(me.auth, me.DataID, map[string]interface{}{"input": d})
		if err != nil {
			return me.statusResult, err
		}
	}
	return me.next(d)
}

func (me *DocumentFlowInstance) Current(d map[string]interface{}) (*Status, error) {
	err := me.Init(me.auth, me.system)
	if err != nil {
		return me.statusResult, err
	}
	n, err := me.wfContext.Current()
	if err != nil {
		return me.statusResult, err
	}
	err = me.ensureWeBeginWithAForegroundNode(n, d)
	if err != nil {
		return me.statusResult, err
	}
	err = me.currentStatus(n)
	if err != nil {
		return me.statusResult, err
	}
	if n != nil && n.Type == "form" {
		formItem, err := me.system.DB.Form.Get(me.auth, n.ID)
		if err == nil {
			me.statusResult.Data = form.GetFormSrc(formItem.Data)
		}
		me.statusResult.UserData, _ = me.getDataFor(n.ID)
	}
	return me.statusResult, nil
}

func (me *DocumentFlowInstance) ensureWeBeginWithAForegroundNode(n *workflow.Node, d map[string]interface{}) error {
	if me.Started || !me.wfContext.HasPrevious() {
		_, err := me.next(d)
		if err == nil {
			me.Started = false
		}
		return err
	}
	return nil
}

func (me *DocumentFlowInstance) next(d map[string]interface{}) (*Status, error) {
	if me.statusResult.HasNext {
		me.Confirmed = false
	}
	var err error
	me.statusResult.HasNext, err = me.wfContext.Next()
	if err != nil {
		return me.statusResult, err
	}
	me.statusResult.Steps = me.wfContext.State()
	me.Steps = me.statusResult.Steps
	me.statusResult.HasPrev = me.hasPrev()
	return me.statusResult, err
}

func (me *DocumentFlowInstance) hasPrev() bool {
	formsCount := 0
	for _, a := range me.statusResult.Steps {
		if a.Type == "form" {
			formsCount++
		}
	}
	return formsCount > 1
}

func (me *Status) addDoc(id, hash string, tmpls *model.TemplateItem) {
	if me.Docs == nil {
		me.Docs = make([]interface{}, 0)
	}
	langs := make([]string, 0)
	if tmpls != nil {
		for k := range tmpls.Data {
			langs = append(langs, k)
		}
	}
	me.Docs = append(me.Docs, map[string]interface{}{"id": id, "hash": hash, "langs": langs, "name": tmpls.Name, "detail": tmpls.Detail})
}

func (me *DocumentFlowInstance) templateAdd(id string, tmplItem *model.TemplateItem) {
	me.templateMapLock.Lock()
	defer me.templateMapLock.Unlock()
	if me.templateMap == nil {
		me.templateMap = map[string]*model.TemplateItem{}
	}
	me.templateMap[id] = tmplItem
}
func (me *DocumentFlowInstance) templateClear() {
	me.templateMapLock.Lock()
	defer me.templateMapLock.Unlock()
	me.templateMap = map[string]*model.TemplateItem{}
}

func (me *DocumentFlowInstance) templateExists(id string) bool {
	me.templateMapLock.RLock()
	defer me.templateMapLock.RUnlock()
	if me.templateMap != nil {
		if _, ok := me.templateMap[id]; ok {
			return true
		}
	}
	return false
}

func (me *DocumentFlowInstance) Previous() *Status {
	me.Confirmed = false
	me.statusResult.HasPrev, _ = me.wfContext.Previous(true)
	if !me.statusResult.HasPrev {
		me.Started = true
	}
	me.currentStatus(nil)
	return me.statusResult
}

func (me *DocumentFlowInstance) getTemplate(id, lang string) (*file.IO, error) {
	me.templateMapLock.RLock()
	defer me.templateMapLock.RUnlock()
	if me.templateMap != nil {
		if tmplm, ok := me.templateMap[id]; ok {
			if tmplm != nil {
				return tmplm.GetTemplate(lang)
			}
		}
	}
	return nil, os.ErrNotExist
}

func (me *DocumentFlowInstance) Preview(id, lang, strFormat string, resp *echo.Response) error {
	if available, err := me.isLangAvailable(lang); !available || err != nil {
		return os.ErrNotExist
	}
	tmpl, err := me.getTemplate(id, lang)
	if err != nil {
		return err
	}
	format := eio.Format(strFormat)
	resp.Header().Set("Content-Disposition", fmt.Sprintf("%s;filename=\"%s\"", "attachment", tmpl.NameWithExt(format.String())))
	dataMap, files := me.getDataWithFiles()
	dsResp, err := me.system.DS.Compile(eio.Template{
		TemplatePath: tmpl.Path(),
		Data:         dataMap,
		Format:       format,
		Assets:       files,
	})
	if err != nil {
		return err
	}
	me.SelectedDocLangs[id] = lang
	resp.Header().Set("Content-Type", dsResp.Header.Get("Content-Type"))
	resp.Header().Set("Content-Length", dsResp.Header.Get("Content-Length"))
	defer dsResp.Body.Close()
	resp.Committed = true //prevents from-> http: multiple response.WriteHeader calls
	_, err = io.Copy(resp.Writer, dsResp.Body)
	return nil
}

func (me *DocumentFlowInstance) currentStatus(n *workflow.Node) error {
	var err error
	if n == nil {
		n, err = me.wfContext.Current()
		if err != nil {
			return err
		}
	}
	me.statusResult.CurrentType = n.Type
	me.statusResult.Steps = me.wfContext.State()
	me.Steps = me.statusResult.Steps
	me.statusResult.HasNext = me.wfContext.HasNext()
	me.statusResult.HasPrev = me.hasPrev()
	return err
}

func (me *DocumentFlowInstance) updateDocs() {
	me.statusResult.Docs = make([]interface{}, 0)
	for k, v := range me.templateMap {
		me.statusResult.addDoc(k, "", v)
	}
}

/**
***** Node Implementation type "form" *******
***************** START *********************
 */

func (me *FormNodeImpl) Execute(n *workflow.Node) (proceed bool, err error) {
	//formData, err := me.ctx.getDataFor(n.ID)
	formDataIf, _ := me.ctx.getDataByPath("input")
	var formData map[string]interface{}
	if v, ok := formDataIf.(map[string]interface{}); ok {
		formData = v
	}
	me.ctx.statusResult.CurrentType = n.Type

	if !me.presented {
		formItem, err := me.ctx.system.DB.Form.Get(me.ctx.auth, n.ID)
		if err != nil {
			return false, err
		}
		me.ctx.statusResult.TargetName = formItem.Name
		formItem.Data = form.GetFormSrc(formItem.Data) // compatibility func
		if formItem.Data == nil {
			return false, fmt.Errorf("form empty")
		}
		me.ctx.statusResult.Data = formItem.Data
		me.ctx.statusResult.UserData = formData
		me.presented = true
		return false, nil
	}
	formItem, err := me.ctx.system.DB.Form.Get(me.ctx.auth, n.ID)
	if err != nil {
		return false, err
	}
	me.ctx.statusResult.Data = form.GetFormSrc(formItem.Data)
	me.ctx.statusResult.UserData = formData
	errs, err := form.Validate(formData, me.ctx.statusResult.Data, true)
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

func (me *FormNodeImpl) Remove(n *workflow.Node) {}
func (me *FormNodeImpl) Close() {
	me.ctx = nil
}

/**
***************** END ***********************
***** Node Implementation type "form" *******
 */

/**
***** Node Implementation type "template" *******
***************** START *********************
 */

func (me *DocTmplNodeImpl) Execute(n *workflow.Node) (proceed bool, err error) {
	me.ctx.statusResult.CurrentType = n.Type
	uniqueNodeId := n.WFUniqueID()
	if !me.ctx.templateExists(uniqueNodeId) {
		var tmplItem *model.TemplateItem
		tmplItem, err = me.ctx.system.DB.Template.Get(me.ctx.auth, n.ID)
		if err != nil {
			return false, err
		}
		if tmplItem != nil && tmplItem.Data != nil {
			if tmplItem.Data != nil {
				me.ctx.templateAdd(uniqueNodeId, tmplItem)
				me.ctx.updateDocs()
			}
		}
	}
	return true, nil
}

func (me *DocTmplNodeImpl) Remove(n *workflow.Node) {
	//remove the template from the app collection as it is not part of the path anymore
	delete(me.ctx.templateMap, n.WFUniqueID())
	me.ctx.updateDocs()
}

func (me *DocTmplNodeImpl) Close() {
	me.ctx = nil
}

/**
***************** END ***********************
***** Node Implementation type "template" *******
 */

func (me *DocumentFlowInstance) Close() error {
	if me.DataCluster != nil {
		return me.DataCluster.Close()
	}
	return nil
}
