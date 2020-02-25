package service

import (
	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/sys"
)

var system *sys.System

func Init(s *sys.System) {
	system = s
}

func paymentsDB() storage.WorkflowPaymentsIF {
	return system.DB.WorkflowPayments
}

func workflowDB() storage.WorkflowIF {
	return system.DB.Workflow
}

func userDB() storage.UserIF {
	return system.DB.User
}

func userDataDB() storage.UserDataIF {
	return system.DB.UserData
}

func formDB() storage.FormIF {
	return system.DB.Form
}

func templateDB() storage.TemplateIF {
	return system.DB.Template
}

func settingsDB() storage.SettingsIF {
	return system.DB.Settings
}

func filesDB() storage.FilesIF {
	return system.DB.Files
}
