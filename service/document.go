package service

import (
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"

	"github.com/ProxeusApp/proxeus-core/main/app"
	"github.com/ProxeusApp/proxeus-core/sys/file"
	"github.com/ProxeusApp/proxeus-core/sys/form"
	"github.com/ProxeusApp/proxeus-core/sys/model"
	"github.com/ProxeusApp/proxeus-core/sys/validate"
	"github.com/ProxeusApp/proxeus-core/sys/workflow"
)

type (
	// DocumentService is an interface that provides user document functions
	DocumentService interface {
		// GetWorkflowSchema Returns the workflow by id
		GetWorkflowSchema(auth model.Auth, workflowId string) (*model.WorkflowItem, map[string]interface{}, error)

		// GetWorkflowSchema Returns the workflow by id
		Edit(auth model.Auth, userId string, formInput map[string]interface{}) error

		// GetDocApp returns the DocumentFlowInstance with the passed id
		GetDocApp(auth model.MemoryAuth, id string) *app.DocumentFlowInstance

		// Update changes the data of the workflow
		Update(auth model.MemoryAuth, id string, data map[string]interface{}) (validate.ErrorMap, error)

		// UpdateFile modifies the file of the current workflow
		UpdateFile(auth model.MemoryAuth, id, fieldName, fileName, contentType string, reader io.Reader) (*file.IO, validate.Errors, error)

		// Next proceeds to the next workflow step
		Next(auth model.MemoryAuth, id, lang string, data map[string]interface{}, final bool) (*app.DocumentFlowInstance, *app.Status, error)

		// Prev returns to the previous workflow step
		Prev(auth model.MemoryAuth, id string) (*app.Status, error)

		// GetFile returns a file by id and input name
		GetFile(auth model.MemoryAuth, id, inputName string) (*file.IO, error)

		// Preview gets a file Preview for a template
		Preview(auth model.MemoryAuth, id, templateId, lang, format string) (*app.PreviewResponse, error)

		// Delete removes the existing document with the provided id
		Delete(auth model.MemoryAuth, id string) error
	}

	DefaultDocumentService struct {
		userService UserService
		fileService FileService
	}
)

var (
	ErrDocAppNotFound = errors.New("doc app not found")
	ErrUnableToEdit   = errors.New("document edit failed")
)

func NewDocumentService(userS UserService, fileS FileService) *DefaultDocumentService {
	return &DefaultDocumentService{userService: userS, fileService: fileS}
}

// GetWorkflowSchema Returns the workflow by id
func (me *DefaultDocumentService) GetWorkflowSchema(auth model.Auth, workflowId string) (*model.WorkflowItem, map[string]interface{}, error) {

	wf, err := workflowDB().Get(auth, workflowId)
	if err != nil {
		return nil, nil, err
	}
	fieldsAndRules := me.getAllFormFieldsWithRulesOf(wf.Data, auth)

	return wf, fieldsAndRules, nil
}

// Edit modifies the document name and detail
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
	usrDataItem, err := userDataDB().Get(auth, userId)
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

	return userDataDB().Put(auth, usrDataItem)
}

// GetDocApp returns the DocumentFlowInstance with the passed id
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
		err := docApp.Init(auth, system)
		if err != nil {
			log.Println("Init err", err)
			return nil //return nil to keep existing behavior
		}
	}
	return docApp
}

// Update changes the data of the workflow
func (me *DefaultDocumentService) Update(auth model.MemoryAuth, id string, data map[string]interface{}) (validate.ErrorMap, error) {
	docApp := me.GetDocApp(auth, id)
	if docApp == nil {
		return nil, ErrDocAppNotFound
	}

	return docApp.UpdateData(data, false)
}

// UpdateFile modifies the file of the current workflow
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

// Next proceeds to the next workflow step
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
	_, _, status, err = docApp.Confirm(lang, data, userDataDB())
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
	item, err = userDataDB().Get(auth, dataID)
	if err != nil {
		return docApp, status, err
	}
	item.Finished = true

	err = userDataDB().Put(auth, item)

	return docApp, status, err
}

// Prev returns to the previous workflow step
func (me *DefaultDocumentService) Prev(auth model.MemoryAuth, id string) (*app.Status, error) {
	docApp := me.GetDocApp(auth, id)
	if docApp == nil {
		return nil, ErrDocAppNotFound
	}
	return docApp.Previous(), nil
}

// GetFile returns a file by id and input name
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

// Preview gets a file Preview for a template
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

// Delete removes the existing document with the provided id
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
				it, er := formDB().Get(a, node.ID)
				if er != nil {
					return true //continue
				}
				marshaledForms[it.ID] = it
			}
		} else if node.Type == "workflow" { // deep dive...
			it, er := workflowDB().Get(a, node.ID)
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
