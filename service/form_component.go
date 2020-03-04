package service

import (
	"encoding/json"
	"io"

	"github.com/ProxeusApp/proxeus-core/test/assets"

	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/storage/database/db"
	"github.com/ProxeusApp/proxeus-core/sys/model"

	"io/ioutil"
	"log"
	"path/filepath"
)

type (

	// FormComponentService is an interface that provides form component functions
	FormComponentService interface {

		//EnsureDefaultFormComponents creates all default form components
		EnsureDefaultFormComponents(auth model.Auth) error

		// DelComp removes a form component
		DelComp(auth model.Auth, id string) error

		// SetComp sets a form components
		SetComp(auth model.Auth, reader io.ReadCloser) (*model.FormComponentItem, error)

		// GetComp returns a form component by id
		GetComp(auth model.Auth, id string) (*model.FormComponentItem, error)

		// ListComp lists all components that match contains and storage.options
		ListComp(auth model.Auth, contains string, options storage.Options) (map[string]*model.FormComponentItem, error)
	}

	DefaultFormComponentService struct {
	}
)

func NewFormComponentService() *DefaultFormComponentService {
	return &DefaultFormComponentService{}
}

//EnsureDefaultFormComponents creates all default form components
func (me *DefaultFormComponentService) EnsureDefaultFormComponents(auth model.Auth) error {
	_, err := formDB().ListComp(auth, "", storage.Options{})
	if err != nil && !db.NotFound(err) {
		return err
	}
	defaultFormcomponentents := []string{"HC1", "HC2", "HC3", "HC5", "HC7", "HC8", "HC9", "HC10", "HC11", "HC12"}
	for _, formCompId := range defaultFormcomponentents {
		jsonFile, err := assets.Asset(filepath.Join("test", "assets", "components", formCompId+".json"))
		if err != nil {
			log.Println(err)
			return err
		}
		var comp model.FormComponentItem
		err = json.Unmarshal(jsonFile, &comp)
		if err != nil {
			log.Println(err)
			return err
		}

		err = formDB().PutComp(auth, &comp)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

// DelComp removes a form component
func (me *DefaultFormComponentService) DelComp(auth model.Auth, id string) error {
	return formDB().DelComp(auth, id)
}

// SetComp sets a form components
func (me *DefaultFormComponentService) SetComp(auth model.Auth, reader io.ReadCloser) (*model.FormComponentItem, error) {
	body, _ := ioutil.ReadAll(reader)
	var comp model.FormComponentItem
	err := json.Unmarshal(body, &comp)
	if err != nil {
		return nil, err
	}
	return &comp, formDB().PutComp(auth, &comp)
}

// GetComp returns a form component by id
func (me *DefaultFormComponentService) GetComp(auth model.Auth, id string) (*model.FormComponentItem, error) {
	return formDB().GetComp(auth, id)
}

// ListComp lists all components that match contains and storage.options
func (me *DefaultFormComponentService) ListComp(auth model.Auth, contains string, options storage.Options) (map[string]*model.FormComponentItem, error) {
	return formDB().ListComp(auth, contains, options)
}
