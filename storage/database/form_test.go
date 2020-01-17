package database

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ProxeusApp/proxeus-core/storage"

	"github.com/ProxeusApp/proxeus-core/sys/model"

	. "github.com/onsi/gomega"
)

func TestForm(t *testing.T) {
	RegisterTestingT(t)
	fo := testDBSet.Form

	options := storage.IndexOptions(0)

	item := &model.FormItem{
		ID:   "1",
		Name: "some name",
		Data: simpleFormData("some_field"),
	}
	item2 := &model.FormItem{
		Name: "some name 2",
	}

	comp := &model.FormComponentItem{
		ID:       "2",
		Name:     "some comp",
		Template: "tpl",
	}

	// put
	Expect(fo.Put(dummy, item)).To(Succeed())
	item.Name = "some name 1"
	Expect(fo.Put(dummy, item)).To(Succeed())
	Expect(fo.Put(dummy, item2)).To(Succeed())
	li, _ := fo.List(dummy, item.Name, options)
	item2.ID = li[0].ID

	// get
	Expect(fo.Get(dummy, item.ID)).To(equalJSON(item))
	Expect(fo.Get(dummy, item2.ID)).To(equalJSON(item))
	Expect(fo.Vars(dummy, item.ID, options)).To(Equal([]string{}))

	// components
	Expect(fo.PutComp(dummy, comp)).To(Succeed())
	Expect(fo.GetComp(dummy, comp.ID)).To(Equal(comp))
	Expect(fo.ListComp(dummy, comp.Template, options)).
		To(equalJSON(map[string]*model.FormComponentItem{comp.ID: comp}))

	// delete
	Expect(fo.DelComp(dummySuperAdmin, comp.ID)).To(Succeed())
	Expect(fo.Delete(dummy, item.ID)).To(Succeed())
}

func simpleFormData(fieldName string) map[string]interface{} {
	j := fmt.Sprintf(`{
    "formSrc": {
      "components": {
        "5zvr98w21yynozx60nhmc": {
          "_compId": "HC2",
          "_order": 0,
          "autocomplete": "on",
          "help": "test-help",
          "label": "test-label",
          "name": "%s",
          "placeholder": "test-placeholder",
          "validate": {
            "required": true
          }
        }
      },
      "v": 2
    }
  }`, fieldName)
	var result map[string]interface{}
	err := json.Unmarshal([]byte(j), &result)
	if err != nil {
		return nil
	}
	return result
}
