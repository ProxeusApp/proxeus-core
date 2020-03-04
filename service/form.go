package service

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/sys/file"
	"github.com/ProxeusApp/proxeus-core/sys/form"
	"github.com/ProxeusApp/proxeus-core/sys/model"
	"github.com/ProxeusApp/proxeus-core/sys/validate"
)

type (

	// FormService is an interface that provides form functions
	FormService interface {
		// List returns a list of all forms
		List(auth model.Auth, contains string, options storage.Options) ([]*model.FormItem, error)

		// Get returns a form by id
		Get(auth model.Auth, id string) (*model.FormItem, error)

		// ExportForms returns the id of all forms
		ExportForms(auth model.Auth, id, contains string) []string

		// UpdateForm updates the form with the data from the reader
		UpdateForm(auth model.Auth, id string, reader io.ReadCloser) (*model.FormItem, error)

		// Delete a form by id
		Delete(auth model.Auth, id string) error

		// Returns a list of variable defined in a form
		Vars(auth model.Auth, contains string, options storage.Options) ([]string, error)

		// GetFormData returns the data of a form
		// If the reset flag is set the data will be cleared
		GetFormData(auth model.MemoryAuth, id string, reset bool) (map[string]interface{}, error)

		// SetFormSrc sets the data from formSrc to the form
		SetFormSrc(auth model.MemoryAuth, id string, formSrc map[string]interface{}) error

		// GetFormFile reads the contents of a form file into the writer
		GetFormFile(auth model.MemoryAuth, id string, fieldname string, writer io.Writer) error

		// TestFormData validates and sets data
		TestFormData(auth model.MemoryAuth, id string, input map[string]interface{}, submit bool) (validate.ErrorMap, error)

		// PostFormFile sets the form file from the reader
		PostFormFile(auth model.MemoryAuth, id string, fileName string, fieldname string, reader io.ReadCloser, contentType string) (interface{}, error)
	}

	DefaultFormService struct {
	}
)

func NewFormService() *DefaultFormService {
	return &DefaultFormService{}
}

// UpdateForm updates the form with the data from the reader
func (me *DefaultFormService) UpdateForm(auth model.Auth, id string, reader io.ReadCloser) (*model.FormItem, error) {
	body, _ := ioutil.ReadAll(reader)
	item := model.FormItem{}
	err := json.Unmarshal(body, &item)
	if err == nil {
		item.ID = id
		err = formDB().Put(auth, &item)
		return &item, err
	}
	return nil, nil
}

// ExportForms returns the id of all forms
func (me *DefaultFormService) ExportForms(auth model.Auth, id, contains string) []string {
	var exportId []string
	if id != "" {
		exportId = []string{id}
	} else if contains != "" {
		items, _ := formDB().List(auth, contains, storage.Options{Limit: 1000})
		if len(items) > 0 {
			exportId = make([]string, len(items))
			for i, item := range items {
				exportId[i] = item.ID
			}
		}
	}
	return exportId
}

// List returns a list of all forms
func (me *DefaultFormService) List(auth model.Auth, contains string, options storage.Options) ([]*model.FormItem, error) {
	return formDB().List(auth, contains, options)
}

// Get returns a form by id
func (me *DefaultFormService) Get(auth model.Auth, id string) (*model.FormItem, error) {
	return formDB().Get(auth, id)
}

// Delete a form by id
func (me *DefaultFormService) Delete(auth model.Auth, id string) error {
	return formDB().Delete(auth, id)
}

// Returns a list of variable defined in a form
func (me *DefaultFormService) Vars(auth model.Auth, contains string, options storage.Options) ([]string, error) {
	return formDB().Vars(auth, contains, options)
}

// SetFormSrc sets the data from formSrc to the form
func (me *DefaultFormService) SetFormSrc(auth model.MemoryAuth, id string, formSrc map[string]interface{}) error {
	dc := GetDataManager(auth)
	return dc.PutDataWithoutMerge("src"+id, formSrc)
}

// TestFormData validates and sets data
func (me *DefaultFormService) TestFormData(auth model.MemoryAuth, id string, input map[string]interface{}, submit bool) (validate.ErrorMap, error) {
	dc := GetDataManager(auth)
	formSrc, _ := dc.GetData("src" + id)
	if formSrc == nil {
		item, err := formDB().Get(auth, id)
		if err == nil && item != nil {
			formSrc = item.Data
		}
	}
	presistedData, err := dc.GetData(id)
	if err != nil {
		return nil, err
	}
	pd := file.MapIO(presistedData)
	pd.MergeWith(input)
	errors, err := form.Validate(pd, formSrc, submit)
	if err != nil {
		return errors, nil
	}
	if len(errors) > 0 {
		return errors, nil
	}
	err = dc.PutData(id, input)
	if err == nil {
		return nil, nil
	}
	return nil, err
}

// GetFormFile reads the contents of a form file into the writer
func (me *DefaultFormService) GetFormFile(auth model.MemoryAuth, id string, fieldname string, writer io.Writer) error {
	dc := GetDataManager(auth)
	fi, err := dc.GetDataFile(id, fieldname)
	if err == nil {
		err = filesDB().Read(fi.Path(), writer)
	}
	return err
}

// GetFormData returns the data of a form
// If the reset flag is set the data will be cleared
func (me *DefaultFormService) GetFormData(auth model.MemoryAuth, id string, reset bool) (map[string]interface{}, error) {
	dc := GetDataManager(auth)
	if reset {
		dc.Clear(id)
	}
	return dc.GetData(id)
}

// PostFormFile sets the form file from the reader
func (me *DefaultFormService) PostFormFile(auth model.MemoryAuth, id string, fileName string, fieldname string, reader io.ReadCloser, contentType string) (interface{}, error) {
	dc := GetDataManager(auth)
	formSrc, _ := dc.GetData("src" + id)
	if formSrc == nil {
		item, err := formDB().Get(auth, id)
		if err == nil && item != nil {
			formSrc = item.Data
		}
	}
	defer reader.Close()
	buf, err := form.ValidateFile(reader, formSrc, fieldname)
	if err != nil {
		return nil, err
	}
	err = dc.PutDataFile(filesDB(), id, fieldname,
		file.Meta{
			Name:        fileName,
			ContentType: contentType,
		},
		bytes.NewBuffer(buf),
	)
	if err == nil {
		fData, err := dc.GetDataByPath(id, fieldname)
		return fData, err
	}
	return nil, err
}

// GetDataManager returns the form.DataManager
func GetDataManager(auth model.MemoryAuth) form.DataManager {
	var dc form.DataManager
	v, ok := auth.GetMemory("testDC")
	if ok {
		dc = v.(form.DataManager)
	} else {
		dc = form.NewDataManager(auth.GetSessionDir())
		auth.PutMemory("testDC", dc)
	}
	return dc
}
