package service

import (
	"bytes"
	"encoding/json"
	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/sys"
	"github.com/ProxeusApp/proxeus-core/sys/file"
	"github.com/ProxeusApp/proxeus-core/sys/form"
	"github.com/ProxeusApp/proxeus-core/sys/model"
	"github.com/ProxeusApp/proxeus-core/sys/validate"
	"io"
	"io/ioutil"
)

type (
	FormService interface {
		List(auth model.Auth, contains string, options storage.Options) ([]*model.FormItem, error)
		Get(auth model.Auth, id string) (*model.FormItem, error)
		ExportForms(auth model.Auth, id, contains string) ([]string)
		UpdateForm(auth model.Auth, id string, reader io.ReadCloser) (*model.FormItem, error)
		Delete(auth model.Auth, id string) error

		Vars(auth model.Auth, contains string, options storage.Options) ([]string, error)

		GetFormData(sess *sys.Session, id string, reset bool) (map[string]interface{}, error)
		SetFormSrc(sess *sys.Session, id string, formSrc map[string]interface{}) error
		GetFormFile(sess *sys.Session, id string, fieldname string, writer io.Writer) error
		TestFormData(sess *sys.Session, id string, input map[string]interface{}, submit bool) (validate.ErrorMap, error)
		PostFormFile(sess *sys.Session, id string, fileName string,fieldname string, reader io.ReadCloser, contentType string) (interface{},error)
	}

	DefaultFormService struct {
	}
)

func NewFormService() *DefaultFormService {
	return &DefaultFormService{}
}

func (me *DefaultFormService) UpdateForm(auth model.Auth, id string, reader io.ReadCloser) (*model.FormItem, error){
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

func (me *DefaultFormService) ExportForms(auth model.Auth, id, contains string) ([]string){
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

func (me *DefaultFormService) List(auth model.Auth, contains string, options storage.Options) ([]*model.FormItem, error) {
	return formDB().List(auth, contains, options)
}

func (me *DefaultFormService) Get(auth model.Auth, id string) (*model.FormItem, error) {
	return formDB().Get(auth, id)
}

func (me *DefaultFormService) Delete(auth model.Auth, id string) error {
	return formDB().Delete(auth, id)
}

func (me *DefaultFormService) Vars(auth model.Auth, contains string, options storage.Options) ([]string, error) {
	return formDB().Vars(auth, contains, options)
}

func (me *DefaultFormService) SetFormSrc(sess *sys.Session, id string, formSrc map[string]interface{}) error{
	dc := GetDataManager(sess)
	return dc.PutDataWithoutMerge("src"+id, formSrc)
}

func (me *DefaultFormService) TestFormData(sess *sys.Session, id string, input map[string]interface{}, submit bool) (validate.ErrorMap, error){
	dc := GetDataManager(sess)
	formSrc, _ := dc.GetData("src" + id)
	if formSrc == nil {
		item, err := formDB().Get(sess, id)
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

func (me *DefaultFormService) GetFormFile(sess *sys.Session, id string, fieldname string, writer io.Writer) error{
	dc := GetDataManager(sess)
	fi, err := dc.GetDataFile(id, fieldname)
	if err == nil {
		err = filesDB().Read(fi.Path(), writer)
	}
	return err
}

func (me *DefaultFormService) GetFormData(sess *sys.Session, id string, reset bool) (map[string]interface{}, error){
	dc := GetDataManager(sess)
	if reset {
		dc.Clear(id)
	}
	return dc.GetData(id)
}

func (me *DefaultFormService) PostFormFile(sess *sys.Session, id string, fileName string,fieldname string, reader io.ReadCloser, contentType string) (interface{},error){
	dc := GetDataManager(sess)
	formSrc, _ := dc.GetData("src" + id)
	if formSrc == nil {
		item, err := formDB().Get(sess, id)
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

func GetDataManager(sess *sys.Session) form.DataManager {
	var dc form.DataManager
	v, ok := sess.GetMemory("testDC")
	if ok {
		dc = v.(form.DataManager)
	} else {
		dc = form.NewDataManager(sess.GetSessionDir())
		sess.PutMemory("testDC", dc)
	}
	return dc
}
