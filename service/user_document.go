package service

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/sys/eio"
	"github.com/ProxeusApp/proxeus-core/sys/model"
)

type (
	UserDocumentService interface {
		List(auth model.Auth, contains string, settings storage.Options) ([]*model.UserDataItem, error)
		Get(auth model.Auth, id string) (*model.UserDataItem, error)
		Put(auth model.Auth, userDataItem *model.UserDataItem) error
		GetDocFile(auth model.Auth, id, dataPath, inlineOrAttachment string) (*FileHeaderResponse, string, error)
		GetTemplateWithFormatFile(auth model.Auth, id, dataPath, format, inlineOrAttachment string) (*FileHeaderResponse, io.ReadCloser, error)
		GetByWorkflow(auth model.Auth, wf *model.WorkflowItem, finished bool) (*model.UserDataItem, bool, error)
	}

	DefaultUserDocumentService struct {
		userService             UserService
		fileService             FileService
		templateDocumentService TemplateDocumentService
	}

	FileHeaderResponse struct {
		ContentType        string
		ContentDisposition string
		ContentLength      string
		ContentPages       string
	}
)

func NewUserDocumentService(userS UserService, fileS FileService, templateDocumentS TemplateDocumentService) *DefaultUserDocumentService {
	return &DefaultUserDocumentService{userService: userS, fileService: fileS, templateDocumentService: templateDocumentS}
}

// List returns a list of UserDataItem that contain the string passed in contain.
// settings are used to modify the result list
func (me *DefaultUserDocumentService) List(auth model.Auth, contains string, settings storage.Options) ([]*model.UserDataItem, error) {
	return userDataDB().List(auth, contains, settings, false)
}

// Get returns the UserDataItem with the id
func (me *DefaultUserDocumentService) Get(auth model.Auth, id string) (*model.UserDataItem, error) {
	return userDataDB().Get(auth, id)
}

//Returns the userDataItem with the provided workflow
func (me *DefaultUserDocumentService) GetByWorkflow(auth model.Auth, wf *model.WorkflowItem, finished bool) (*model.UserDataItem, bool, error) {
	return userDataDB().GetByWorkflow(auth, wf, finished)
}

func (me *DefaultUserDocumentService) Put(auth model.Auth, userDataItem *model.UserDataItem) error {
	return userDataDB().Put(auth, userDataItem)
}

// GetDocFile returns file info for a pdf file that was generated from a workflow
func (me *DefaultUserDocumentService) GetDocFile(auth model.Auth, id, dataPath, inlineOrAttachment string) (*FileHeaderResponse, string, error) {
	fileInfo, err := me.fileService.GetDataFile(auth, id, dataPath)
	if err != nil {
		return nil, "", os.ErrNotExist
	}

	fileName := fileInfo.NameWithExt("pdf")
	contentDisposition := fmt.Sprintf(`%s; filename="%s"`, inlineOrAttachment, fileName)
	contentLength := strconv.FormatInt(fileInfo.Size(), 10)

	headerResponse := &FileHeaderResponse{
		ContentType:        fileInfo.ContentType(),
		ContentDisposition: contentDisposition,
		ContentLength:      contentLength,
		ContentPages:       "",
	}

	return headerResponse, fileInfo.Path(), nil
}

// GetTemplateWithFormatFile returns file info for a docx file that was generated from a workflow
func (me *DefaultUserDocumentService) GetTemplateWithFormatFile(auth model.Auth, id, dataPath, format, inlineOrAttachment string) (*FileHeaderResponse, io.ReadCloser, error) {
	fileInfo, err := me.fileService.GetDataFile(auth, id, dataPath)
	if err != nil {
		return nil, nil, err
	}
	dat, files, err := me.fileService.GetDataAndFiles(auth, id, "input")
	if err != nil {
		log.Println("[fileService][GetFile] GetDataAndFiles err: ", err.Error())
	}
	formt := eio.Format(format)
	template := eio.Template{
		Format:       formt,
		Data:         map[string]interface{}{"input": dat},
		TemplatePath: fileInfo.Path(),
		Assets:       files,
	}
	dsResp, err := me.templateDocumentService.Compile(template)
	if err != nil {
		return nil, nil, err
	}

	contentDisposition := fmt.Sprintf("%s;filename=\"%s\"", inlineOrAttachment, fileInfo.NameWithExt(formt.String()))

	headerResponse := &FileHeaderResponse{
		ContentType:        dsResp.Header.Get("Content-Type"),
		ContentDisposition: contentDisposition,
		ContentLength:      dsResp.Header.Get("Content-Length"),
		ContentPages:       dsResp.Header.Get("Content-Pages"),
	}

	return headerResponse, dsResp.Body, nil
}
