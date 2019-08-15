package mysql

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	uuid "github.com/satori/go.uuid"

	"git.proxeus.com/core/central/sys/db"
	"git.proxeus.com/core/central/sys/form"
	"git.proxeus.com/core/central/sys/model"
)

type (
	MysqlForm struct {
		db         *sql.DB
		varService *varService
	}
)

func NewFormStore(db *sql.DB) (*MysqlForm, error) {
	mf := &MysqlForm{db: db}
	vs := newVarService(db, mf)
	mf.varService = vs
	return mf, mf.createTables()
}

func (me *MysqlForm) createTables() error {
	varsTable := `
CREATE TABLE IF NOT EXISTS _variable_ (
  id bigint(20) NOT NULL AUTO_INCREMENT,
  varname longtext NULL,
  refs longtext NULL,
  PRIMARY KEY (id)
) ENGINE=MYISAM AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;
`
	_, err := me.db.Exec(varsTable)
	_, _ = me.db.Exec("UPDATE `html_component` hc SET hc.template = REPLACE(hc.template, 'glyphicon', 'fa') WHERE hc.template LIKE '%glyphicon%'")
	_, _ = me.db.Exec("ALTER TABLE `form` ADD COLUMN `accessible` LONGTEXT NULL DEFAULT NULL AFTER `id`;")
	_, _ = me.db.Exec("ALTER TABLE `html_component` ADD COLUMN `accessible` LONGTEXT NULL DEFAULT NULL AFTER `id`;")
	_, _ = me.db.Exec("ALTER TABLE `template` ADD COLUMN `accessible` LONGTEXT NULL DEFAULT NULL AFTER `id`;")
	_, _ = me.db.Exec("ALTER TABLE `workflow_data` ADD COLUMN `accessible` LONGTEXT NULL DEFAULT NULL AFTER `id`;")
	_, _ = me.db.Exec("ALTER TABLE `workflow` ADD COLUMN `accessible` LONGTEXT NULL DEFAULT NULL AFTER `id`;")
	return err
}

func (me *MysqlForm) List(auth model.Authorization, contains string, options map[string]interface{}) ([]*model.FormItem, error) {
	sQuery := makeSimpleQuery(options)
	var err error
	var r *sql.Rows
	cnts := "%"
	if contains != "" {
		cnts = "%" + contains + "%"
	}
	selectFields := "`id`, `name`, `description`, `last_updated`, `date_created`, `accessible`"
	if !sQuery.metaOnly {
		selectFields += ", form_src"
	}
	if exlen := len(sQuery.exclude); exlen > 0 {
		params := make([]interface{}, exlen+4)
		ind := 0
		inQuery := &bytes.Buffer{}
		inQuery.WriteString(fmt.Sprintf("SELECT %s FROM form f WHERE f.id not in(", selectFields))
		for a := range sQuery.exclude {
			if ind > 0 {
				inQuery.WriteString(",?")
			} else {
				inQuery.WriteString("?")
			}
			params[ind] = a
			ind++
		} //cnts, cnts, sQuery.index, sQuery.limit
		params[ind] = cnts
		ind++
		params[ind] = cnts
		ind++
		params[ind] = sQuery.index
		ind++
		params[ind] = sQuery.limit
		ind++
		inQuery.WriteString(") AND (f.name LIKE ? OR f.description LIKE ?) GROUP BY id ORDER BY f.last_updated DESC, f.date_created DESC, f.name ASC LIMIT ?,?")
		r, err = me.db.Query(inQuery.String(), params...)
	} else {
		r, err = me.db.Query(fmt.Sprintf("SELECT %s FROM form f WHERE (f.name LIKE ? OR f.description LIKE ?) GROUP BY id ORDER BY f.last_updated DESC, f.date_created DESC, f.name ASC LIMIT ?,?", selectFields), cnts, cnts, sQuery.index, sQuery.limit)
	}
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var (
		id          sql.NullString
		name        sql.NullString
		desc        sql.NullString
		accessible  sql.NullString
		wfData      sql.NullString
		lastUpdated mysql.NullTime
		dateCreated mysql.NullTime
	)

	resusts := make([]*model.FormItem, 0)
	for r.Next() {
		if sQuery.metaOnly {
			err = r.Scan(&id, &name, &desc, &lastUpdated, &dateCreated, &accessible)
		} else {
			err = r.Scan(&id, &name, &desc, &lastUpdated, &dateCreated, &accessible, &wfData)
		}
		if err != nil {
			return nil, err
		}
		item := &model.FormItem{}
		if id.Valid {
			item.ID = id.String
		}
		if name.Valid {
			item.Name = name.String
		}
		if desc.Valid {
			item.Detail = desc.String
		}
		if lastUpdated.Valid {
			item.Updated = lastUpdated.Time
		}
		if dateCreated.Valid {
			item.Created = dateCreated.Time
		}
		if !sQuery.metaOnly && wfData.Valid {
			var formSrc map[string]interface{}
			err = json.Unmarshal([]byte(me.i18nCompatibilityConverter(wfData.String)), &formSrc)
			if err == nil && formSrc != nil {
				if _, ok := formSrc["formSrc"]; !ok {
					item.Data = map[string]interface{}{"formSrc": formSrcCompatibilityRead(formSrc)}
				} else {
					item.Data = formSrc
				}
			}
		}
		resusts = append(resusts, item)
	}
	if len(resusts) == 0 {
		return nil, db.ErrNotFound
	}
	return resusts, nil
}

func formSrcCompatibilityRead(formSrc map[string]interface{}) map[string]interface{} {
	if formSrc != nil {
		if vs, ok := formSrc["v"]; ok {
			vn, ok := vs.(float64)
			if ok && vn > 1 {
				if _, ok := formSrc["components"]; ok {
					return formSrc
				}
			} else if _, ok := formSrc["components"]; ok {
				return formSrc
			}
		} else if _, ok := formSrc["components"]; ok {
			return map[string]interface{}{"v": 2, "components": formSrc["components"]}
		}
		return map[string]interface{}{"v": 2, "components": formSrc}
	}
	return formSrc
}

func (me *MysqlForm) Get(auth model.Authorization, id string) (*model.FormItem, error) {
	var err error
	if id == "" {
		err = fmt.Errorf("id cant be empty")
		return nil, err
	}
	var r *sql.Rows
	var stmt *sql.Stmt
	stmt, err = me.db.Prepare("SELECT id, `name`, description, last_updated, date_created, `accessible`, form_src FROM form f WHERE f.id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	r, err = stmt.Query(id)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var (
		sid         sql.NullString
		name        sql.NullString
		accessible  sql.NullString
		desc        sql.NullString
		wfData      sql.NullString
		lastUpdated mysql.NullTime
		dateCreated mysql.NullTime
	)
	var item *model.FormItem
	for r.Next() {
		item = &model.FormItem{}
		err = r.Scan(&sid, &name, &desc, &lastUpdated, &dateCreated, &accessible, &wfData)
		if err != nil {
			return nil, err
		}

		if sid.Valid {
			item.ID = sid.String
		}
		if name.Valid {
			item.Name = name.String
		}
		if desc.Valid {
			item.Detail = desc.String
		}
		if lastUpdated.Valid {
			item.Updated = lastUpdated.Time
		}
		if dateCreated.Valid {
			item.Created = dateCreated.Time
		}
		if wfData.Valid {
			var formSrc map[string]interface{}
			jsnBts := []byte(me.i18nCompatibilityConverter(wfData.String))
			err = json.Unmarshal(jsnBts, &formSrc)
			if err == nil && formSrc != nil {
				if _, ok := formSrc["formSrc"]; !ok {
					item.Data = map[string]interface{}{"formSrc": formSrcCompatibilityRead(formSrc)}
				} else {
					item.Data = formSrc
				}
			}
		}
	}
	return item, nil
}

func (me *MysqlForm) getPermissions(id, tbl string) (*model.Permissions, error) {
	var err error
	if id == "" {
		err = fmt.Errorf("id cant be empty")
		return nil, err
	}
	var r *sql.Rows
	var stmt *sql.Stmt
	stmt, err = me.db.Prepare(fmt.Sprintf("SELECT `accessible` FROM %s f WHERE f.id = ?", tbl))
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	r, err = stmt.Query(id)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var (
		accessible sql.NullString
	)
	var item *model.FormItem
	for r.Next() {
		item = &model.FormItem{}
		err = r.Scan(&accessible)
		if err != nil {
			return nil, err
		}

		if accessible.Valid {
			//err handling is not important here
			//just keep in mind to convert data before reading them again if the Permissions model changes!
			json.Unmarshal([]byte(accessible.String), &item.Permissions)
		}
	}
	return &item.Permissions, nil
}

//Put creates an ID if the ID is "" but it's important to handle the permissions of this behaviour outside Put
func (me *MysqlForm) Put(auth model.Authorization, item *model.FormItem) error {
	if item == nil {
		return errors.New("item cannot be nil")
	}
	var err error
	var currentPerm *model.Permissions
	currentPerm, err = me.getPermissions(item.ID, "form")
	if err != nil {
		return err
	}
	if item.ID == "" { //insert new item
		item.Created = time.Now()
		item.Updated = item.Created
		for i := 0; i < 5; i++ {
			u2 := uuid.NewV4()
			item.ID = u2.String()
			_, err = me.db.Exec("INSERT INTO `form` (`id`, `version`, `form_compiled`, `form_src`, `name`, `description`, `sort_index`, `uri_name`, `contract_id`, `date_created`, `last_updated`, `accessible`) VALUES (?,1,'','',?,?,0,?,1,?,?,?);", item.ID, item.Name, item.Detail, item.ID, item.Created, item.Updated, "")
			if err == nil {
				break
			}
		}
		if err != nil {
			return err
		}
	}

	item.Updated = time.Now()
	var jsonDataStr string
	if item.Data != nil {
		var jsonData []byte
		jsonData, err = json.Marshal(item.Data)
		if err != nil {
			jsonDataStr = ""
		} else {
			jsonDataStr = string(jsonData)
		}
	}
	var permBytes []byte
	currentPerm.Change(auth, &item.Permissions)
	permBytes, err = json.Marshal(currentPerm)
	if err != nil {
		return err
	}
	var result sql.Result
	result, err = me.db.Exec("UPDATE form f SET f.name = ?, f.description = ?, f.form_src = ?, f.accessible = ?, f.last_updated = ? WHERE f.id = ?",
		item.Name,
		item.Detail,
		jsonDataStr,
		string(permBytes),
		item.Updated,
		item.ID)
	if err != nil {
		return err
	}
	_, err = result.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}

func (me *MysqlForm) DelComp(auth model.Authorization, id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}
	var err error
	var currentPerm *model.Permissions
	currentPerm, err = me.getPermissions(id, "html_component")
	if err != nil {
		return err
	}
	if currentPerm == nil || !currentPerm.IsWriteGrantedFor(auth) {
		return model.ErrAuthorityMissing
	}
	_, err = me.db.Exec("DELETE from html_component WHERE id = ?;", id)
	return err
}

func (me *MysqlForm) PutComp(auth model.Authorization, comp *model.FormComponentItem) error {
	if comp == nil {
		return errors.New("item cannot be nil")
	}
	var err error
	var currentPerm *model.Permissions
	currentPerm, err = me.getPermissions(comp.ID, "html_component")
	if err != nil {
		return err
	}
	if currentPerm == nil || !currentPerm.IsWriteGrantedFor(auth) {
		return model.ErrAuthorityMissing
	}
	if comp.ID == "" { //insert new item
		for i := 0; i < 5; i++ {
			u2 := uuid.NewV4()
			comp.ID = u2.String()
			_, err = me.db.Exec("INSERT INTO html_component (`id`, `version`, `settings`, `sort_index`, `template`, `accessible`) VALUES (?,1,'', 0,'','')", comp.ID)
			if err == nil {
				break
			}
		}
		if err != nil {
			return err
		}
	}

	var jsonDataStr string
	if comp.Settings != nil {
		var jsonData []byte
		jsonData, err = json.Marshal(comp.Settings)
		if err != nil {
			jsonDataStr = ""
		} else {
			jsonDataStr = string(jsonData)
		}
	}
	var permBytes []byte
	currentPerm.Change(auth, &comp.Permissions)
	permBytes, err = json.Marshal(currentPerm)
	if err != nil {
		return err
	}
	result, err := me.db.Exec("UPDATE html_component f SET f.settings = ?, f.template = ?, f.accessible = ? WHERE f.id = ?",
		jsonDataStr,
		comp.Template,
		string(permBytes),
		comp.ID)
	if err != nil {
		return err
	}
	_, err = result.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}

func (me *MysqlForm) GetComp(auth model.Authorization, id string) (*model.FormComponentItem, error) {
	if id == "" {
		return nil, fmt.Errorf("id can't be empty")
	}
	var err error
	var currentPerm *model.Permissions
	currentPerm, err = me.getPermissions(id, "html_component")
	if err != nil {
		return nil, err
	}
	if currentPerm == nil || !currentPerm.IsReadGrantedFor(auth) {
		return nil, model.ErrAuthorityMissing
	}
	var (
		sid        sql.NullString
		accessible sql.NullString
		settings   sql.NullString
		template   sql.NullString
		stmt       *sql.Stmt
		r          *sql.Rows
	)
	stmt, err = me.db.Prepare("SELECT `id`, `settings`, `template`, `accessible` FROM `html_component` f WHERE f.id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	r, err = stmt.Query(id)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	item := &model.FormComponentItem{}
	hasItem := false
	for r.Next() {
		err = r.Scan(&sid, &settings, &template, &accessible)
		if err != nil {
			return nil, err
		}
		hasItem = true
		if sid.Valid {
			item.ID = sid.String
		}
		if accessible.Valid {
			//err handling is not important here
			//just keep in mind to convert data before reading them again if the Permissions model changes!
			json.Unmarshal([]byte(accessible.String), &item.Permissions)
		}
		if settings.Valid {
			var set interface{}
			jsnBts := []byte(me.i18nCompatibilityConverter(settings.String))
			err = json.Unmarshal(jsnBts, &set)
			if err == nil {
				item.Settings = set
			}
		}
		if template.Valid {
			item.Template = template.String
		}

	}
	if hasItem {
		return item, nil
	}
	return nil, db.ErrNotFound
}

func (me *MysqlForm) ListComp(auth model.Authorization, contains string, options map[string]interface{}) (map[string]*model.FormComponentItem, error) {
	sQuery := makeSimpleQuery(options)
	var (
		sid        sql.NullString
		accessible sql.NullString
		settings   sql.NullString
		template   sql.NullString
		err        error
		r          *sql.Rows
	)
	if contains != "" {
		cnts := "%" + contains + "%"
		r, err = me.db.Query("SELECT `id`, `settings`, `template`, `accessible` FROM `html_component` f WHERE f.settings LIKE ? OR f.template LIKE ? LIMIT ?,?", cnts, cnts, sQuery.index, sQuery.limit)
	} else {
		r, err = me.db.Query("SELECT `id`, `settings`, `template`, `accessible` FROM `html_component` f LIMIT ?,?", sQuery.index, sQuery.limit)
	}
	if err != nil {
		return nil, err
	}
	defer r.Close()
	items := map[string]*model.FormComponentItem{}
	for r.Next() {
		err = r.Scan(&sid, &settings, &template, &accessible)
		if err != nil {
			return nil, err
		}
		item := &model.FormComponentItem{}
		if sid.Valid {
			item.ID = sid.String
		}
		if accessible.Valid {
			//err handling is not important here
			//just keep in mind to convert data before reading them again if the Permissions model changes!
			json.Unmarshal([]byte(accessible.String), &item.Permissions)
			if !item.Permissions.IsReadGrantedFor(auth) {
				continue
			}
		}
		if settings.Valid {
			var set interface{}
			jsnBts := []byte(me.i18nCompatibilityConverter(settings.String))
			err = json.Unmarshal(jsnBts, &set)
			if err == nil {
				item.Settings = set
			}
		}
		if template.Valid {
			item.Template = template.String
		}
		items[sid.String] = item
	}
	return items, nil
}

func (me *MysqlForm) Vars(auth model.Authorization, contains string, options map[string]interface{}) ([]string, error) {
	sQuery := makeSimpleQuery(options)
	var r *sql.Rows
	var stmt *sql.Stmt
	var err error
	stmt, err = me.db.Prepare("SELECT varname FROM _variable_ f WHERE f.varname LIKE ? GROUP by varname ORDER BY varname ASC LIMIT ?,?;")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	cnts := "%" + contains + "%"
	r, err = stmt.Query(cnts, sQuery.index, sQuery.limit)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var (
		name sql.NullString
	)
	vars := make([]string, 0)
	for r.Next() {
		err = r.Scan(&name)
		if err != nil {
			return nil, err
		}

		if name.Valid {
			vars = append(vars, name.String)
		}
	}
	return vars, nil
}

func (me *MysqlForm) syncVarsTask(gitem interface{}) {
	if me.varsCount() == 0 {
		for i := 0; i < 100000; i++ {
			items, err := me.List(nil, "", map[string]interface{}{"metaOnly": false, "index": i, "limit": 30})
			if err == nil && len(items) > 0 {
				for _, item := range items {
					me.storeVars(item)
				}
			} else {
				break
			}
		}

	}
	me.storeVars(gitem)
}

func (me *MysqlForm) storeVars(gitem interface{}) {
	if gitem != nil {
		item, ok := gitem.(*model.FormItem)
		if ok && item.Data != nil {
			vars := make([]interface{}, 0)
			form.LoopComponents(item.Data, func(compId, compInstId string, compMain interface{}, comp map[string]interface{}) bool {
				name, ok := form.CompName(comp)
				if ok {
					vars = append(vars, name)
				}
				return true
			})
			if len(vars) > 0 {
				me.insertVars(item.ID, vars)
			}
		}
	}
}

func (me *MysqlForm) varsCount() (count int) {
	rows, err := me.db.Query("SELECT COUNT(*) as count FROM `_variable_`;")
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&count)
	}
	return
}

func (me *MysqlForm) insertVars(id string, vars []interface{}) {
	sqlStr := fmt.Sprintf("INSERT INTO `_variable_` (`varname`, `refs`) VALUES ")
	insVal := fmt.Sprintf("(?, '%s'),", id)
	for range vars {
		sqlStr += insVal
	}
	//trim the last ,
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	sqlStr += ";"
	tx, err := me.db.Begin()
	//prepare the statement
	_, err = me.db.Exec(fmt.Sprintf("DELETE from _variable_ WHERE refs = '%s';", id))
	if err == nil {
		stmt, err := me.db.Prepare(sqlStr)
		if err == nil {
			defer stmt.Close()
			_, err := stmt.Exec(vars...)
			if err == nil {
				tx.Commit()
				return
			}
		}
	}
	tx.Rollback()
}

func (me *MysqlForm) i18nCompatibilityConverter(formSrc string) string {
	if formSrc != "" {
		var i18nCompReg = regexp.MustCompile(`\"\$\{message\(code:'(.*?)'\)\}\"`)
		for _, match := range i18nCompReg.FindAllStringSubmatch(formSrc, -1) {
			if len(match) == 2 {
				formSrc = strings.Replace(formSrc, match[0], fmt.Sprintf(`{"i18n": "%s"}`, match[1]), 1)
			}
		}
	}
	return formSrc
}

func (mf *MysqlForm) Close() error {
	mf.varService.Close()
	return nil
}
