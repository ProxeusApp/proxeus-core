package service

import (
	"encoding/json"
	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/sys/model"
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
