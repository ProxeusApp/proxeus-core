package service

import (
	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/sys"
)

type (
	baseService struct {
		system *sys.System
	}
)

func (me *baseService) paymentsDB() storage.WorkflowPaymentsIF {
	return me.system.DB.WorkflowPayments
}

func (me *baseService) workflowDB() storage.WorkflowIF {
	return me.system.DB.Workflow
}

func (me *baseService) userDB() storage.UserIF {
	return me.system.DB.User
}

func (me *baseService) userDataDB() storage.UserDataIF {
	return me.system.DB.UserData
}

func (me *baseService) formDB() storage.FormIF {
	return me.system.DB.Form
}

func (me *baseService) templateDB() storage.TemplateIF {
	return me.system.DB.Template
}
