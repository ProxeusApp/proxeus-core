package storm

import "github.com/ProxeusApp/proxeus-core/storage"

func NewDBSet(sDB storage.SettingsIF, folderPath string) (me *storage.DBSet, err error) {
	me = &storage.DBSet{Settings: sDB}
	settings, err := me.Settings.Get()
	if err != nil {
		return nil, err
	}
	conf := DBConfig{
		Dir:    folderPath,
		Engine: settings.DatabaseEngine,
		URI:    settings.DatabaseURI,
	}
	me.I18n, err = NewI18nDB(conf)
	if err != nil {
		return nil, err
	}
	me.Form, err = NewFormDB(conf)
	if err != nil {
		return nil, err
	}
	me.Template, err = NewDocTemplateDB(conf)
	if err != nil {
		return nil, err
	}
	me.Workflow, err = NewWorkflowDB(conf)
	if err != nil {
		return nil, err
	}
	me.User, err = NewUserDB(conf)
	if err != nil {
		return nil, err
	}
	me.UserData, err = NewUserDataDB(conf)
	if err != nil {
		return nil, err
	}
	me.SignatureRequests, err = NewSignatureDB(conf)
	if err != nil {
		return nil, err
	}
	me.WorkflowPayments, err = NewWorkflowPaymentDB(conf)
	if err != nil {
		return nil, err
	}
	return
}

type DBConfig struct {
	Dir    string
	Engine string
	URI    string
}
