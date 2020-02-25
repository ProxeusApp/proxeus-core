package service

import (
	"encoding/json"
	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/storage/database/db"
	"github.com/ProxeusApp/proxeus-core/sys"
	"github.com/ProxeusApp/proxeus-core/sys/model"
	"io"

	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type (
	FormComponentService interface {
		EnsureDefaultFormComponents(auth model.Auth)
		DelComp(auth model.Auth, id string) error
		SetComp(auth model.Auth, reader io.ReadCloser) (*model.FormComponentItem, error)
		GetComp(auth model.Auth, id string) (*model.FormComponentItem, error)
		ListComp(auth model.Auth, contains string, options storage.Options) (map[string]*model.FormComponentItem, error)
	}

	DefaultFormComponentService struct {
		*baseService
	}
)

func NewFormComponentService(system *sys.System) *DefaultFormComponentService {
	return &DefaultFormComponentService{&baseService{system: system}}
}


func (me *DefaultFormComponentService) EnsureDefaultFormComponents(auth model.Auth){
	dat, err := me.system.DB.Form.ListComp(auth, "", storage.Options{})
	if db.NotFound(err) || (err == nil && dat == nil) {
		defaultFormcomponentents := []string{"HC1", "HC2", "HC3", "HC5", "HC7", "HC8", "HC9", "HC10", "HC11", "HC12"}
		for _, formCompId := range defaultFormcomponentents {
			jsonFile, err := os.Open(filepath.Join("test", "assets", "components", formCompId+".json"))
			if err != nil {
				log.Println(err)
			}

			body, _ := ioutil.ReadAll(jsonFile)
			var comp model.FormComponentItem
			err = json.Unmarshal(body, &comp)
			if err != nil {
				log.Println(err)
				jsonFile.Close()
				continue
			}

			err = me.system.DB.Form.PutComp(auth, &comp)
			if err != nil {
				log.Println(err)
				jsonFile.Close()
				continue
			}
		}
	}
}

func (me *DefaultFormComponentService) DelComp(auth model.Auth, id string) error {
	return me.system.DB.Form.DelComp(auth, id)
}

func (me *DefaultFormComponentService) SetComp(auth model.Auth, reader io.ReadCloser) (*model.FormComponentItem, error){
	body, _ := ioutil.ReadAll(reader)
	var comp model.FormComponentItem
	err := json.Unmarshal(body, &comp)
	if err != nil {
		return nil, err
	}
	return &comp, me.system.DB.Form.PutComp(auth, &comp)
}

func (me *DefaultFormComponentService) GetComp(auth model.Auth, id string) (*model.FormComponentItem, error) {
	return me.system.DB.Form.GetComp(auth, id)
}

func (me *DefaultFormComponentService) ListComp(auth model.Auth, contains string, options storage.Options) (map[string]*model.FormComponentItem, error) {
	return me.system.DB.Form.ListComp(auth, contains, options)
}
