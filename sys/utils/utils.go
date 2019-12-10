package utils

import (
	"git.proxeus.com/core/central/sys"
	"git.proxeus.com/core/central/sys/form"
	"git.proxeus.com/core/central/sys/model"
	"git.proxeus.com/core/central/sys/workflow"
)

func GetAllFormFieldsWithRulesOf(wf *workflow.Workflow, a model.Authorization, s *sys.System) map[string]interface{} {
	marshaledForms := marshaledFormsOf(wf, a, s)
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

func marshaledFormsOf(wf *workflow.Workflow, a model.Authorization, s *sys.System) map[string]*model.FormItem {
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
