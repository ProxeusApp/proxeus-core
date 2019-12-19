package storm

import "github.com/ProxeusApp/proxeus-core/storage"

func NewDBSet(sDB storage.SettingsIF, folderPath string) (me *storage.DBSet, err error) {
	me = &storage.DBSet{Settings: sDB}
	me.I18n, err = NewI18nDB(folderPath)
	if err != nil {
		return nil, err
	}
	me.Form, err = NewFormDB(folderPath)
	if err != nil {
		return nil, err
	}
	me.Template, err = NewDocTemplateDB(folderPath)
	if err != nil {
		return nil, err
	}
	me.Workflow, err = NewWorkflowDB(folderPath)
	if err != nil {
		return nil, err
	}
	me.User, err = NewUserDB(folderPath)
	if err != nil {
		return nil, err
	}
	me.UserData, err = NewUserDataDB(folderPath)
	if err != nil {
		return nil, err
	}
	me.SignatureRequests, err = NewSignatureDB(folderPath)
	if err != nil {
		return nil, err
	}
	me.WorkflowPayments, err = NewWorkflowPaymentDB(folderPath)
	if err != nil {
		return nil, err
	}
	return
}
