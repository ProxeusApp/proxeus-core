package utils

import (
	"strings"

	"github.com/ProxeusApp/proxeus-core/sys"
	"github.com/ProxeusApp/proxeus-core/sys/form"
	"github.com/ProxeusApp/proxeus-core/sys/model"
	"github.com/ProxeusApp/proxeus-core/sys/workflow"
)

// FindFieldNameContaining finds field names containing all the provided strings
// to make it possible to use take names like:
// lastName_pr with lastname
// CTXFirstName with firstname
// ConsumerFirstName with firstname
// ProducerFirstName with firstname
func FindFieldNameContaining(fields []string, containing ...string) []string {
	var fieldNames []string
	containingCount := 0
	for i := range containing {
		containing[i] = strings.ToLower(containing[i])
	}
	for _, v := range fields {
		fnameLower := strings.ToLower(v)
		containingCount = 0
		for _, c := range containing {
			if strings.Contains(fnameLower, c) {
				containingCount++
			}
		}
		if containingCount == len(containing) {
			fieldNames = append(fieldNames, v)
		}
	}
	return fieldNames
}

func GetAllFormFieldsOf(wf *workflow.Workflow, a model.Authorization, s *sys.System) []string {
	marshaledForms := MarshaledFormsOf(wf, a, s)
	allFieldsMap := map[string]bool{}
	//collect all form fields
	for _, formItem := range marshaledForms {
		vars := form.Vars(formItem.Data)
		for _, v := range vars {
			allFieldsMap[v] = true
		}
	}
	allFields := make([]string, len(allFieldsMap))
	i := 0
	for k := range allFieldsMap {
		allFields[i] = k
		i++
	}
	return allFields
}

func GetAllFormFieldsWithRulesOf(wf *workflow.Workflow, a model.Authorization, s *sys.System) map[string]interface{} {
	marshaledForms := MarshaledFormsOf(wf, a, s)
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

func MarshaledFormsOf(wf *workflow.Workflow, a model.Authorization, s *sys.System) map[string]*model.FormItem {
	if wf == nil {
		return nil
	}
	marshaledForms := map[string]*model.FormItem{}
	loop := &workflow.Looper{}
	//loop recursively and collect all forms
	wf.Loop(loop, func(l *workflow.Looper, node *workflow.Node) bool {
		if node.Type == "form" {
			if _, ok := marshaledForms[node.ID]; !ok {
				it, er := s.DB.Form.Get(a, node.ID)
				if er != nil {
					return true //continue
				}
				marshaledForms[it.ID] = it
			}
		} else if node.Type == "workflow" { // deep dive...
			it, er := s.DB.Workflow.Get(a, node.ID)
			if er != nil {
				return true //continue
			}
			it.LoopNodes(l, nil)
		}
		return true //continue
	})
	return marshaledForms
}
