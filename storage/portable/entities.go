package portable

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/storage/db/storm"
	"github.com/ProxeusApp/proxeus-core/sys/file"
	"github.com/ProxeusApp/proxeus-core/sys/form"
	"github.com/ProxeusApp/proxeus-core/sys/model"
)

type EntityType string

const (
	Settings EntityType = "Settings"
	I18n     EntityType = "I18n"
	Form     EntityType = "Form"
	Workflow EntityType = "Workflow"
	Template EntityType = "Template"
	User     EntityType = "User"
	UserData EntityType = "UserData"
)

// embedded entities
const formComponent = "FormComponent"

func StringToEntityType(str string) EntityType {
	switch strings.ToLower(str) {
	case "settings":
		return Settings
	case "i18n":
		return I18n
	case "form":
		return Form
	case "workflow":
		return Workflow
	case "template":
		return Template
	case "user":
		return User
	case "userdata":
		return UserData
	}
	return ""
}

func (ie *ImportExport) Export(t EntityType, ids ...string) error {
	switch t {
	case Settings:
		return ie.exportSettings(ids...)
	case I18n:
		return ie.exportI18n(ids...)
	case Form:
		return ie.exportForm(ids...)
	case Workflow:
		return ie.exportWorkflow(ids...)
	case Template:
		return ie.exportTemplate(ids...)
	case User:
		return ie.exportUser(ids...)
	case UserData:
		return ie.exportUserData(ids...)
	}
	return fmt.Errorf("type '%s' Not Found", t)
}

func (ie *ImportExport) Import(t EntityType) error {
	switch t {
	case Settings:
		return ie.importSettings()
	case I18n:
		return ie.importI18n()
	case Form:
		return ie.importForm()
	case Workflow:
		return ie.importWorkflow()
	case Template:
		return ie.importTemplate()
	case User:
		return ie.importUser()
	case UserData:
		return ie.importUserData()
	}
	return fmt.Errorf("type '%s' Not Found", t)
}

func (ie *ImportExport) exportTemplate(id ...string) error {
	if ie.db.Template == nil {
		var err error
		ie.db.Template, err = storm.NewDocTemplateDB(ie.dbConfig)
		if err != nil {
			return err
		}
	}
	specificIds := len(id) > 0
	if len(id) == 1 {
		if ie.isProcessed(Template, id[0]) {
			return nil
		}
	}
	for i := 0; true; i++ {
		items, err := ie.sysDB.Template.List(ie.auth, "", storage.IndexOptions(i).WithInclude(id))
		if err == nil && len(items) > 0 {
			wg := &sync.WaitGroup{}
			wg.Add(1)
			fileCopyErrs := map[string]error{}
			go func() {
				for _, item := range items {
					for _, tmplItem := range item.Data {
						path := filepath.Join(ie.db.Template.AssetsKey(), tmplItem.PathName())
						_, err = storage.CopyFileAcross(ie.db.Files, ie.sysDB.Files, path, tmplItem.Path())
						if err != nil {
							fileCopyErrs[item.ID] = err
						}
					}
				}
				wg.Done()
			}()

			for _, item := range items {
				if !ie.isProcessed(Template, item.ID) {
					item.Permissions.UserIdsMap(ie.neededUsers)
					if err != nil {
						ie.processedEntry(Template, item.ID, err)
						continue
					}
					ie.processedEntry(Template, item.ID, nil)
				}
			}
			wg.Wait()
			for k, v := range fileCopyErrs {
				ie.processedEntry(Template, k, v)
			}
		} else {
			break
		}
		if specificIds {
			break
		}
	}
	return nil
}

func (ie *ImportExport) importTemplate() error {
	if ie.db.Template == nil {
		var err error
		ie.db.Template, err = storm.NewDocTemplateDB(ie.dbConfig)
		if err != nil {
			return err
		}
	}
	for i := 0; true; i++ {
		items, err := ie.db.Template.List(ie.auth, "", storage.IndexOptions(i))
		if err == nil && len(items) > 0 {
			for _, item := range items {
				if ie.skipExistingOnImport {
					_, err = ie.sysDB.Template.Get(ie.auth, item.ID)
					if err == nil {
						continue
					}
				}

				item.Permissions.UpdateUserID(ie.locatedSameUserWithDifferentID)

				err = ie.sysDB.Template.Put(ie.auth, item)
				if err != nil {
					ie.processedEntry(Template, item.ID, err)
					continue
				}
				var fi *file.IO
				hadError := false
				for lang, tmplItem := range item.Data {
					fi, err = ie.sysDB.Template.GetTemplate(ie.auth, item.ID, lang)
					if err != nil {
						hadError = true
						ie.processedEntry(Template, item.ID, err)
						continue
					}
					n, err := storage.CopyFileAcross(ie.sysDB.Files, ie.db.Files, fi.Path(), tmplItem.Path())
					fi.SetSize(n)
					if err != nil {
						hadError = true
						ie.processedEntry(Template, item.ID, err)
					}
				}
				if !hadError {
					ie.processedEntry(Template, item.ID, nil)
				}
			}
		} else {
			break
		}
	}
	return nil
}

func (ie *ImportExport) exportSettings(id ...string) error {
	if ie.db.Settings == nil {
		var err error
		ie.db.Settings, err = storm.NewSettingsDB(ie.dir)
		if err != nil {
			ie.processedEntry(Settings, string(Settings), err)
			return err
		}
	}
	if !ie.auth.AccessRights().IsGrantedFor(model.ROOT) {
		ie.processedEntry(Settings, string(Settings), fmt.Errorf("no authority to export"))
		//return nil to not break the export, just ignore the call
		return nil
	}
	settings, err := ie.sysDB.Settings.Get()
	if err != nil {
		ie.processedEntry(Settings, string(Settings), err)
		return nil
	}
	err = ie.db.Settings.Put(settings)
	if err != nil {
		ie.processedEntry(Settings, string(Settings), err)
		return nil
	}
	ie.db.Settings.Close()
	ie.processedEntry(Settings, string(Settings), nil)
	return nil
}

const allKeysMarker = "_all"

func (ie *ImportExport) exportI18n(id ...string) error {
	if ie.db.I18n == nil {
		var err error
		ie.db.I18n, err = storm.NewI18nDB(ie.dbConfig)
		if err != nil {
			return err
		}
	}
	specificIds := len(id) > 0
	langs, err := ie.sysDB.I18n.GetAllLangs()
	if err == nil && len(langs) > 0 {
		for _, lang := range langs {
			if !ie.isProcessed(I18n, lang.Code) {
				err = ie.db.I18n.PutLang(lang.Code, lang.Enabled)
				ie.processedEntry(I18n, lang.Code, err)
			}
			if specificIds {
				var trans map[string]string
				trans, err = ie.sysDB.I18n.GetAll(lang.Code)
				if err != nil {
					ie.processedEntry(I18n, lang.Code+allKeysMarker, err)
					continue
				}
				for _, v := range id {
					if !ie.isProcessed(I18n, lang.Code+"_"+v) {
						if text, ok := trans[v]; ok {
							err = ie.db.I18n.Put(lang.Code, v, text)
							ie.processedEntry(I18n, lang.Code+"_"+v, err)
						}
					}
				}
				ie.processedEntry(I18n, lang.Code+allKeysMarker, err)
			} else if !ie.isProcessed(I18n, lang.Code+allKeysMarker) {
				var trans map[string]string
				trans, err = ie.sysDB.I18n.GetAll(lang.Code)
				if err != nil {
					ie.processedEntry(I18n, lang.Code+allKeysMarker, err)
					continue
				}
				err = ie.db.I18n.PutAll(lang.Code, trans)
				if err != nil {
					ie.processedEntry(I18n, lang.Code+allKeysMarker, err)
					continue
				}
				for k := range trans {
					ie.processedEntry(I18n, lang.Code+"_"+k, nil)
				}
				ie.processedEntry(I18n, lang.Code+allKeysMarker, nil)
			}
		}
	}
	return nil
}

func (ie *ImportExport) exportForm(id ...string) error {
	if ie.db.Form == nil {
		var err error
		ie.db.Form, err = storm.NewFormDB(ie.dbConfig)
		if err != nil {
			return err
		}
	}
	specificIds := len(id) > 0
	if len(id) == 1 {
		if ie.isProcessed(Form, id[0]) {
			return nil
		}
	}
	formComps := make(map[string]bool)

	for i := 0; true; i++ {
		items, err := ie.sysDB.Form.List(ie.auth, "", storage.IndexOptions(i).WithInclude(id))
		if err == nil && len(items) > 0 {
			for _, item := range items {
				if ie.isProcessed(Form, item.ID) {
					continue
				}
				formSrc := form.GetFormSrc(item.Data)
				err = ie.db.Form.Put(ie.auth, item)
				if err != nil {
					ie.processedEntry(Form, item.ID, err)
					continue
				}
				form.LoopComponents(formSrc, func(compId, compInstId string, compMain interface{}, comp map[string]interface{}) bool {
					if compId == "" {
						ie.processedEntry(Form, item.ID, fmt.Errorf("form contains a component with an empty id"))
					}
					formComps[compId] = true
					return true
				})
				item.Permissions.UserIdsMap(ie.neededUsers)
				ie.processedEntry(Form, item.ID, nil)
				//TODO export i18n of comps
			}
		} else {
			break
		}
		if specificIds {
			break
		}
	}
	if len(formComps) > 0 {
		for compId := range formComps {
			fbi, err := ie.sysDB.Form.GetComp(ie.auth, compId)
			if err != nil {
				ie.processedEntry(formComponent, compId, err)
				continue
			}
			err = ie.db.Form.PutComp(ie.auth, fbi)
			ie.processedEntry(formComponent, compId, err)
		}
	}

	return nil
}

func (ie *ImportExport) exportWorkflow(id ...string) error {
	if ie.db.Workflow == nil {
		var err error
		ie.db.Workflow, err = storm.NewWorkflowDB(ie.dbConfig)
		if err != nil {
			return err
		}
	}
	specificIds := len(id) > 0
	if len(id) == 1 {
		if ie.isProcessed(Workflow, id[0]) {
			return nil
		}
	}
	type NodesAfter struct {
		id    string
		store EntityType
	}
	nodes := make(map[string]*NodesAfter)
	for i := 0; true; i++ {
		items, err := ie.sysDB.Workflow.List(ie.auth, "", storage.IndexOptions(i).WithInclude(id))
		if err == nil && len(items) > 0 {
			for _, item := range items {
				if !ie.isProcessed(Workflow, item.ID) {
					err = ie.db.Workflow.Put(ie.auth, item)
					if err != nil {
						ie.processedEntry(Workflow, item.ID, err)
						continue
					}
					if item.Data != nil && item.Data.Flow != nil {
						for _, v := range item.Data.Flow.Nodes {
							nodes[v.ID] = &NodesAfter{
								id:    item.ID,
								store: StringToEntityType(v.Type),
							}
						}
					}
					item.Permissions.UserIdsMap(ie.neededUsers)
					if err != nil {
						ie.processedEntry(Workflow, item.ID, err)
						continue
					}
					ie.processedEntry(Workflow, item.ID, nil)
				}
			}
		} else {
			break
		}
		if specificIds {
			break
		}
	}
	for k, v := range nodes {
		if v.store == "" {
			continue
		}
		err := ie.Export(v.store, k)
		if err != nil {
			ie.processedEntry(Workflow, v.id, err)
		}
	}
	return nil
}

func cpProfilePhoto(from *storage.DBSet, to *storage.DBSet, item *model.User) error {
	fromPath := filepath.Join(from.User.GetBaseFilePath(), item.PhotoPath)
	toPath := filepath.Join(to.User.GetBaseFilePath(), item.PhotoPath)
	_, err := storage.CopyFileAcross(to.Files, from.Files, toPath, fromPath)
	return err
}

func (ie *ImportExport) exportUser(id ...string) error {
	if ie.db.User == nil {
		var err error
		ie.db.User, err = storm.NewUserDB(ie.dbConfig, ie.db.Files)
		if err != nil {
			return err
		}
	}
	specificIds := len(id) > 0
	if len(id) == 1 {
		if ie.isProcessed(User, id[0]) {
			return nil
		}
	}
	if !specificIds {
		ie.exportingAllUsersAnyway = true
	}
	for i := 0; true; i++ {
		items, err := ie.sysDB.User.List(ie.auth, "", storage.IndexOptions(i).WithInclude(id))
		if err == nil && len(items) > 0 {
			for _, item := range items {
				if !ie.isProcessed(User, item.ID) {
					item.Photo = ""
					err := ie.db.User.Put(ie.auth, item)
					if item.PhotoPath != "" {
						err = cpProfilePhoto(ie.sysDB, ie.db, item)
					}
					ie.processedEntry(User, item.ID, err)
				}
			}
		} else {
			break
		}
		if specificIds {
			break
		}
	}
	return nil
}

func (ie *ImportExport) exportUserData(id ...string) error {
	if ie.db.UserData == nil {
		var err error
		ie.db.UserData, err = storm.NewUserDataDB(ie.dbConfig)
		if err != nil {
			return err
		}
	}
	specificIds := len(id) > 0
	for i := 0; true; i++ {
		items, err := ie.sysDB.UserData.List(ie.auth, "", storage.IndexOptions(i).WithInclude(id), true)
		if err == nil && len(items) > 0 {
			workflows := map[string]bool{}
			wg := &sync.WaitGroup{}
			wg.Add(1)
			fileCopyErrs := map[string]error{}
			go func() {
				for _, item := range items {
					fios := ie.sysDB.UserData.GetAllFileInfosOf(item)
					for _, fio := range fios {
						path := filepath.Join(ie.db.UserData.AssetsKey(), fio.PathName())
						_, err = storage.CopyFileAcross(ie.db.Files, ie.sysDB.Files, path, fio.Path())
						if err != nil {
							fileCopyErrs[item.ID] = err
						}
					}
				}
				wg.Done()
			}()
			for _, item := range items {
				if !ie.isProcessed(UserData, item.ID) {
					err = ie.db.UserData.Put(ie.auth, item)
					if err != nil {
						ie.processedEntry(UserData, item.ID, err)
						continue
					}
					if item.WorkflowID != "" {
						workflows[item.WorkflowID] = true
					}
					item.Permissions.UserIdsMap(ie.neededUsers)
					if err != nil {
						ie.processedEntry(UserData, item.ID, err)
						continue
					}
					ie.processedEntry(UserData, item.ID, nil)
				}
			}
			if len(workflows) > 0 {
				wfIds := make([]string, len(workflows))
				i := 0
				for wfID := range workflows {
					wfIds[i] = wfID
					i++
				}
				err = ie.exportWorkflow(wfIds...)
				if err != nil {
					return err
				}
			}
			wg.Wait()
			for k, v := range fileCopyErrs {
				ie.processedEntry(UserData, k, v)
			}
		} else {
			break
		}
		if specificIds {
			break
		}
	}
	return nil
}

func (ie *ImportExport) importSettings() error {
	if ie.db.Settings == nil {
		var err error
		ie.db.Settings, err = storm.NewSettingsDB(ie.dir)
		if err != nil {
			ie.processedEntry(Settings, string(Settings), err)
			return err
		}
	}

	s, err := ie.db.Settings.Get()
	if err != nil { //does not exist
		return nil
	}
	if !ie.auth.AccessRights().IsGrantedFor(model.ROOT) {
		ie.processedEntry(Settings, string(Settings), fmt.Errorf("no authority to import"))
		//return nil to not break the export, just ignore the call
		return nil
	}
	err = ie.sysDB.Settings.Put(s)
	if err != nil {
		ie.processedEntry(Settings, string(Settings), err)
		return nil
	}
	ie.processedEntry(Settings, string(Settings), nil)
	return nil
}

func (ie *ImportExport) importI18n() error {
	if ie.db.I18n == nil {
		var err error
		ie.db.I18n, err = storm.NewI18nDB(ie.dbConfig)
		if err != nil {
			return err
		}
	}
	langs, err := ie.db.I18n.GetAllLangs()
	if err == nil && len(langs) > 0 {
		for _, lang := range langs {
			if ie.skipExistingOnImport {
				if !ie.sysDB.I18n.HasLang(lang.Code) {
					err = ie.sysDB.I18n.PutLang(lang.Code, lang.Enabled)
					ie.processedEntry(I18n, lang.Code, err)
				}
			} else {
				err = ie.sysDB.I18n.PutLang(lang.Code, lang.Enabled)
				ie.processedEntry(I18n, lang.Code, err)
			}
			var trans map[string]string
			trans, err = ie.db.I18n.GetAll(lang.Code)
			if err != nil {
				ie.processedEntry(I18n, lang.Code+allKeysMarker, err)
				continue
			}
			if ie.skipExistingOnImport {
				var transExisting map[string]string
				transExisting, err = ie.sysDB.I18n.GetAll(lang.Code)
				if err != nil {
					ie.processedEntry(I18n, lang.Code+allKeysMarker, err)
					continue
				}

				for k, v := range trans {
					if _, exists := transExisting[k]; exists {
						continue
					} else {
						err = ie.sysDB.I18n.Put(lang.Code, k, v)
						ie.processedEntry(I18n, lang.Code+"_"+k, err)
					}
				}
				ie.processedEntry(I18n, lang.Code+allKeysMarker, err)
			} else {
				err = ie.sysDB.I18n.PutAll(lang.Code, trans)
				for k, t := range trans {
					if k == "" {
						ie.processedEntry(I18n, lang.Code+"_"+k, fmt.Errorf("invalid argument: key can not be empty {key:'%s',text:'%s'}", k, t))
						continue
					}
					ie.processedEntry(I18n, lang.Code+"_"+k, err)
				}
			}
		}
	}
	return nil
}

func (ie *ImportExport) importForm() error {
	if ie.db.Form == nil {
		var err error
		ie.db.Form, err = storm.NewFormDB(ie.dbConfig)
		if err != nil {
			return err
		}
	}
	for i := 0; true; i++ {
		items, err := ie.db.Form.ListComp(ie.auth, "", storage.IndexOptions(i))
		if err == nil && len(items) > 0 {
			for _, item := range items {
				if ie.skipExistingOnImport {
					_, err = ie.sysDB.Form.GetComp(ie.auth, item.ID)
					if err == nil {
						continue
					}
				}
				err = ie.sysDB.Form.PutComp(ie.auth, item)
				if err != nil {
					ie.processedEntry(formComponent, item.ID, err)
					continue
				}
				ie.processedEntry(formComponent, item.ID, nil)
			}
		} else {
			break
		}
	}
	for i := 0; true; i++ {
		items, err := ie.db.Form.List(ie.auth, "", storage.IndexOptions(i))
		if err == nil && len(items) > 0 {
			for _, item := range items {
				if ie.skipExistingOnImport {
					_, err = ie.sysDB.Form.Get(ie.auth, item.ID)
					if err == nil {
						continue
					}
				}

				item.Permissions.UpdateUserID(ie.locatedSameUserWithDifferentID)

				err = ie.sysDB.Form.Put(ie.auth, item)
				if err != nil {
					ie.processedEntry(Form, item.ID, err)
					continue
				}
				ie.processedEntry(Form, item.ID, nil)
				//TODO export i18n of comps
			}
		} else {
			break
		}
	}
	return nil
}

func (ie *ImportExport) importWorkflow() error {
	if ie.db.Workflow == nil {
		var err error
		ie.db.Workflow, err = storm.NewWorkflowDB(ie.dbConfig)
		if err != nil {
			return err
		}
	}
	for i := 0; true; i++ {
		items, err := ie.db.Workflow.List(ie.auth, "", storage.IndexOptions(i))
		if err == nil && len(items) > 0 {
			for _, item := range items {
				if ie.skipExistingOnImport {
					_, err = ie.sysDB.Workflow.Get(ie.auth, item.ID)
					if err == nil {
						continue
					}
				}
				item.Permissions.UpdateUserID(ie.locatedSameUserWithDifferentID)

				err = ie.sysDB.Workflow.Put(ie.auth, item)
				if err != nil {
					ie.processedEntry(Workflow, item.ID, err)
					continue
				}
				ie.processedEntry(Workflow, item.ID, nil)
			}
		} else {
			break
		}
	}
	return nil
}

func (ie *ImportExport) importUser() error {
	if ie.db.User == nil {
		var err error
		ie.db.User, err = storm.NewUserDB(ie.dbConfig, ie.db.Files)
		if err != nil {
			return err
		}
	}
	for i := 0; true; i++ {
		items, err := ie.db.User.List(ie.auth, "", storage.IndexOptions(i))
		if err == nil && len(items) > 0 {
			for _, item := range items {
				var existingItem *model.User
				existingItem, err = ie.sysDB.User.Get(ie.auth, item.ID)
				if err == nil && ie.skipExistingOnImport {
					continue
				}

				if existingItem == nil {
					//treat email as an ID and update all references on this import package if user was located on the target system
					if item.Email != "" {
						existingItem, _ = ie.sysDB.User.GetByEmail(item.Email)
						if existingItem != nil {
							//provide user id correction map for entities with permissions item.ID -> existingItem.ID
							ie.locatedSameUserWithDifferentID[item.ID] = existingItem.ID
							if ie.skipExistingOnImport {
								continue
							} else {
								item.ID = existingItem.ID
								updateEmptyFields(item, existingItem)
							}
						}
					}
					//treat Ethereum address as an ID and update all references on this import package if user was located on the target system
					if item.EthereumAddr != "" {
						existingItem, _ = ie.sysDB.User.GetByBCAddress(item.EthereumAddr)
						if existingItem != nil {
							//provide user id correction map for entities with permissions item.ID -> existingItem.ID
							ie.locatedSameUserWithDifferentID[item.ID] = existingItem.ID
							if ie.skipExistingOnImport {
								continue
							} else {
								item.ID = existingItem.ID
								updateEmptyFields(item, existingItem)
							}
						}
					}
				}

				err = ie.sysDB.User.Put(ie.auth, item)
				if err != nil {
					ie.processedEntry(User, item.ID, err)
					continue
				}
				//no permission errors when writing ... we are allowed to set the photo too

				if existingItem != nil && existingItem.PhotoPath != "" && existingItem.PhotoPath != item.PhotoPath {
					//remove old photo of existingItem before the reference is lost
					ie.sysDB.Files.Delete(filepath.Join(ie.sysDB.User.GetBaseFilePath(), existingItem.PhotoPath))
				}

				if item.PhotoPath != "" {
					err = cpProfilePhoto(ie.db, ie.sysDB, item)
					if err != nil {
						continue
					}
				}
				ie.processedEntry(User, item.ID, nil)
			}
		} else {
			break
		}
	}
	return nil
}

func updateEmptyFields(of, by *model.User) {
	if of.PhotoPath == "" {
		of.PhotoPath = by.PhotoPath
	}
	if of.Name == "" {
		of.Name = by.Name
	}
	if of.Detail == "" {
		of.Detail = by.Detail
	}
	if of.Email == "" {
		of.Email = by.Email
	}
	if of.EthereumAddr == "" {
		of.EthereumAddr = by.EthereumAddr
	}
	if of.Role == 0 {
		of.Role = by.Role
	}
	if of.Data == nil {
		of.Data = by.Data
	}
}

func (ie *ImportExport) importUserData() error {
	if ie.db.UserData == nil {
		var err error
		ie.db.UserData, err = storm.NewUserDataDB(ie.dbConfig)
		if err != nil {
			return err
		}
	}
	for i := 0; true; i++ {
		items, err := ie.db.UserData.List(ie.auth, "", storage.IndexOptions(i), true)
		if err == nil && len(items) > 0 {
			for _, item := range items {
				if ie.skipExistingOnImport {
					_, err = ie.sysDB.UserData.Get(ie.auth, item.ID)
					if err == nil {
						continue
					}
				}
				item.Permissions.UpdateUserID(ie.locatedSameUserWithDifferentID)

				err = ie.sysDB.UserData.Put(ie.auth, item)
				if err != nil {
					ie.processedEntry(UserData, item.ID, err)
					continue
				}
				hadError := false
				fios := ie.db.UserData.GetAllFileInfosOf(item)
				for _, fio := range fios {
					path := filepath.Join(ie.sysDB.UserData.AssetsKey(), fio.PathName())
					_, err = storage.CopyFileAcross(ie.sysDB.Files, ie.db.Files, path, fio.Path())
					if err != nil {
						hadError = true
						ie.processedEntry(UserData, item.ID, err)
					}
				}
				if !hadError {
					ie.processedEntry(UserData, item.ID, nil)
				}
			}
		} else {
			break
		}
	}
	return nil
}
