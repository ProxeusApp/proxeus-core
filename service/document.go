package service

import (
	"errors"
	"fmt"
	"github.com/ProxeusApp/proxeus-core/main/app"
	"github.com/ProxeusApp/proxeus-core/sys"
	"github.com/ProxeusApp/proxeus-core/sys/file"
	"github.com/ProxeusApp/proxeus-core/sys/form"
	"github.com/ProxeusApp/proxeus-core/sys/model"
	"github.com/ProxeusApp/proxeus-core/sys/validate"
	"github.com/ProxeusApp/proxeus-core/sys/workflow"
	"io"
	"log"
	"regexp"
)

type (
	DocumentService interface {
		GetWorkflowSchema(auth model.Auth, workflowId string) (*model.WorkflowItem, map[string]interface{}, error)
		Edit(auth model.Auth, userId string, formInput map[string]interface{}) error
		GetDocApp(auth model.MemoryAuth, id string) *app.DocumentFlowInstance
		Update(auth model.MemoryAuth, id string, data map[string]interface{}) (validate.ErrorMap, error)
		UpdateFile(auth model.MemoryAuth, id, fieldName, fileName, contentType string, reader io.Reader) (*file.IO, validate.Errors, error)
		Next(auth model.MemoryAuth, id, lang string, data map[string]interface{}, final bool) (*app.DocumentFlowInstance, *app.Status, error)
		Prev(auth model.MemoryAuth, id string) (*app.Status, error)
		GetFile(auth model.MemoryAuth, id, inputName string) (*file.IO, error)
		Preview(auth model.MemoryAuth, id, templateId, lang, format string) (*app.PreviewResponse, error)
		Delete(auth model.MemoryAuth, id string) error
	}

	DefaultDocumentService struct {
		userService UserService
		fileService FileService
		*baseService
	}
)

var (
	ErrDocAppNotFound = errors.New("doc app not found")
	ErrUnableToEdit   = errors.New("document edit failed")
)

func NewDocumentService(system *sys.System, userS UserService, fileS FileService) *DefaultDocumentService {
	return &DefaultDocumentService{baseService: &baseService{system: system}, userService: userS, fileService: fileS}
}

// Return the workflow by id
func (me *DefaultDocumentService) GetWorkflowSchema(auth model.Auth, workflowId string) (*model.WorkflowItem, map[string]interface{}, error) {

	wf, err := me.workflowDB().Get(auth, workflowId)
	if err != nil {
		return nil, nil, err
	}
	fieldsAndRules := me.getAllFormFieldsWithRulesOf(wf.Data, auth)

	return wf, fieldsAndRules, nil
}

// Edit the document name and detail
func (me *DefaultDocumentService) Edit(auth model.Auth, userId string, formInput map[string]interface{}) error {
	filenameRegex := regexp.MustCompile(`^[^\s][\p{L}\d.,_\-&: ]{3,}[^\s]$`)
	name, ok := formInput["name"]
	if !ok {
		return ErrUnableToEdit
	}
	fName, ok := name.(string)
	if !ok {
		return ErrUnableToEdit
	}
	if len(fName) >= 80 || !filenameRegex.MatchString(fName) {
		return ErrUnableToEdit
	}
	usrDataItem, err := me.userDataDB().Get(auth, userId)
	if err != nil {
		return err
	}
	detail, ok := formInput["detail"]
	if !ok {
		return ErrUnableToEdit
	}
	fDetail, ok := detail.(string)
	if !ok {
		return ErrUnableToEdit
	}
	usrDataItem.Name = fName
	usrDataItem.Detail = fDetail

	return me.userDataDB().Put(auth, usrDataItem)
}

// Returns the DocumentFlowInstance with the passed id
func (me *DefaultDocumentService) GetDocApp(auth model.MemoryAuth, id string) *app.DocumentFlowInstance {
	if auth == nil {
		return nil
	}
	var docApp *app.DocumentFlowInstance
	sessDocAppID := fmt.Sprintf("docApp_%s", id)
	v, ok := auth.GetMemory(sessDocAppID)
	if !ok {
		return nil
	}
	docApp = v.(*app.DocumentFlowInstance)
	if docApp != nil && docApp.NeedToBeInitialized() {
		err := docApp.Init(auth, me.system)
		if err != nil {
			log.Println("Init err", err)
			return nil //return nil to keep existing behavior
		}
	}
	return docApp
}

// Update data of the workflow
func (me *DefaultDocumentService) Update(auth model.MemoryAuth, id string, data map[string]interface{}) (validate.ErrorMap, error) {
	docApp := me.GetDocApp(auth, id)
	if docApp == nil {
		return nil, ErrDocAppNotFound
	}

	return docApp.UpdateData(data, false)
}

// Update the file of the current workflow
func (me *DefaultDocumentService) UpdateFile(auth model.MemoryAuth, id, fieldName, fileName, contentType string, reader io.Reader) (*file.IO, validate.Errors, error) {
	docApp := me.GetDocApp(auth, id)
	if docApp == nil {
		return nil, nil, ErrDocAppNotFound
	}
	verrs, err := docApp.UpdateFile(fieldName, file.Meta{Name: fileName, ContentType: contentType, Size: 0}, reader)
	if err != nil || len(verrs) > 0 {
		return nil, verrs, err
	}

	finfo, err := docApp.GetFile(fieldName)
	if err != nil {
		return nil, nil, err
	}

	return finfo, verrs, err

}

// Go to next workflow step
func (me *DefaultDocumentService) Next(auth model.MemoryAuth, id, lang string, data map[string]interface{}, final bool) (*app.DocumentFlowInstance, *app.Status, error) {
	docApp := me.GetDocApp(auth, id)
	if docApp == nil {
		return nil, nil, ErrDocAppNotFound
	}

	status, err := docApp.Next(data)

	if err != nil || status.HasNext {
		return docApp, status, err
	}

	//done
	_, _, status, err = docApp.Confirm(lang, data, me.userDataDB())
	if err != nil {
		return docApp, status, err
	}
	//after tx success
	if !final {
		return docApp, status, nil
	}

	dataID := docApp.DataID
	me.removeWorkflowDocumentFromSession(auth, id)
	var item *model.UserDataItem
	item, err = me.userDataDB().Get(auth, dataID)
	if err != nil {
		return docApp, status, err
	}
	item.Finished = true

	err = me.userDataDB().Put(auth, item)

	return docApp, status, err
}

// Go to preview workflow step
func (me *DefaultDocumentService) Prev(auth model.MemoryAuth, id string) (*app.Status, error) {
	docApp := me.GetDocApp(auth, id)
	if docApp == nil {
		return nil, ErrDocAppNotFound
	}
	return docApp.Previous(), nil
}

// Return a file by id and input name
func (me *DefaultDocumentService) GetFile(auth model.MemoryAuth, id, inputName string) (*file.IO, error) {
	docApp := me.GetDocApp(auth, id)
	if docApp == nil {
		return nil, ErrDocAppNotFound
	}

	finfo, err := docApp.GetFile(inputName)
	if err == nil && finfo != nil {
		return finfo, err //if file is found return here
	}
	if docApp.DataID == "" {
		return finfo, err
	}

	dataPath := fmt.Sprintf("input.%s", inputName)
	return me.fileService.GetDataFile(auth, docApp.DataID, dataPath)
}

// Get a file Preview for a template
func (me *DefaultDocumentService) Preview(auth model.MemoryAuth, id, templateId, lang, format string) (*app.PreviewResponse, error) {
	if id == "" || templateId == "" || lang == "" || auth == nil {
		return nil, ErrDocAppNotFound
	}

	docApp := me.GetDocApp(auth, id)
	if docApp == nil {
		return nil, ErrDocAppNotFound
	}

	return docApp.Preview(templateId, lang, format)
}

func (me *DefaultDocumentService) Delete(auth model.MemoryAuth, id string) error {
	userDataItem, err := me.userService.GetUserDataById(auth, id)
	if err != nil {
		return err
	}
	me.removeWorkflowDocumentFromSession(auth, userDataItem.WorkflowID)
	return me.userService.DeleteUserData(auth, id)
}

func (me *DefaultDocumentService) getAllFormFieldsWithRulesOf(wf *workflow.Workflow, auth model.Auth) map[string]interface{} {
	marshaledForms := me.marshaledFormsOf(wf, auth)
	fieldsAndRules := map[string]interface{}{}
	//collect all form fields
	for _, formItem := range marshaledForms {
		vars := form.Vars(formItem.Data)
		for _, v := range vars {
			fieldsAndRules[v] = form.RulesOf(formItem.Data, v)
		}
	}
	return fieldsAndRules
}

func (me *DefaultDocumentService) marshaledFormsOf(wf *workflow.Workflow, a model.Auth) map[string]*model.FormItem {
	if wf == nil {
		return nil
	}
	marshaledForms := map[string]*model.FormItem{}
	loop := &workflow.Looper{}
	//loop recursively and collect all forms
	wf.Loop(loop, func(l *workflow.Looper, node *workflow.Node) bool {
		if node.Type == "form" {
			if _, ok := marshaledForms[node.ID]; !ok {
				it, er := me.formDB().Get(a, node.ID)
				if er != nil {
					return true //continue
				}
				marshaledForms[it.ID] = it
			}
		} else if node.Type == "workflow" { // deep dive...
			it, er := me.workflowDB().Get(a, node.ID)
			if er != nil {
				return true //continue
			}
			it.LoopNodes(l, nil)
		}
		return true //continue
	})
	return marshaledForms
}

func (me *DefaultDocumentService) removeWorkflowDocumentFromSession(auth model.MemoryAuth, id string) {
	auth.DeleteMemory(fmt.Sprintf("docApp_%s", id))
}
