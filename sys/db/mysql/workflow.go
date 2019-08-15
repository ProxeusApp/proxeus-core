package mysql

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/juju/errors"
	uuid "github.com/satori/go.uuid"

	"git.proxeus.com/core/central/sys/model"
	"git.proxeus.com/core/central/sys/workflow"
)

type MysqlWorkflow struct {
	db *sql.DB
}

func NewWorkflowStore(db *sql.DB) (*MysqlWorkflow, error) {
	db.SetMaxIdleConns(6)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Second * 10)
	mwf := &MysqlWorkflow{db: db}
	return mwf, mwf.runOnInit()
}

func (me *MysqlWorkflow) runOnInit() error {
	updateWFDStructure := `
		UPDATE workflow_data wd SET 
		wd.json_data = JSON_OBJECT("input", JSON_EXTRACT(wd.json_data, "$")) 
		WHERE wd.json_data is not null AND wd.json_data != "" AND JSON_EXTRACT(wd.json_data, "$.input") is null;
`
	_, err := me.db.Exec(updateWFDStructure) ///

	_, _ = me.db.Exec("INSERT INTO `contract` (`id`,`version`,`created_by_id`,`name`) VALUES (1,1,1,'n');")
	_, _ = me.db.Exec("INSERT INTO `workflow_progress`(`id`,`version`,`json_progress`,`user_id`)VALUES(1,0,'',1);")
	_, _ = me.db.Exec("UPDATE workflow a SET a.`accessible` = '{\\\"owner\\\":\\\"1\\\"}' WHERE a.`accessible` = '{owner:\\\"1\\\"}';")
	return err
}

func (d *MysqlWorkflow) List(auth model.Authorization, contains string, more map[string]interface{}) ([]*model.WorkflowItem, error) {
	sQuery := makeSimpleQuery(more)
	var err error
	var r *sql.Rows
	cnts := "%"
	if contains != "" {
		cnts = "%" + contains + "%"
	}
	if exlen := len(sQuery.exclude); exlen > 0 {
		params := make([]interface{}, exlen+4)
		ind := 0
		inQuery := &bytes.Buffer{}
		inQuery.WriteString("SELECT id, `accessible`, `name`, description, last_updated, date_created, comma_separated_form_ids, comma_separated_template_ids FROM workflow f WHERE f.id not in(")
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
		r, err = d.db.Query(inQuery.String(), params...)
	} else {
		r, err = d.db.Query("SELECT id, `accessible`, `name`, description, last_updated, dateCreated, comma_separated_form_ids, comma_separated_template_ids FROM workflow f WHERE (f.name LIKE ? OR f.description LIKE ?) GROUP BY id ORDER BY f.last_updated DESC, f.dateCreated DESC, f.name ASC LIMIT ?,?", cnts, cnts, sQuery.index, sQuery.limit)
	}
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var (
		id          sql.NullString
		name        sql.NullString
		accessible  sql.NullString
		desc        sql.NullString
		wfData      sql.NullString
		tmplId      sql.NullString
		lastUpdated mysql.NullTime
		dateCreated mysql.NullTime
	)

	resusts := make([]*model.WorkflowItem, 0)
	for r.Next() {
		err = r.Scan(&id, &accessible, &name, &desc, &lastUpdated, &dateCreated, &wfData, &tmplId)
		if err != nil {
			return nil, err
		}
		item := &model.WorkflowItem{}
		if id.Valid {
			item.ID = id.String
		}

		if accessible.Valid {
			//err handling is not important here
			//just keep in mind to convert data before reading them again if the Permissions model changes!
			err = json.Unmarshal([]byte(accessible.String), &item.Permissions)
			if !item.Permissions.IsReadGrantedFor(auth) {
				continue
			}
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
			wf := &workflow.Workflow{}
			jsnBts := []byte(wfData.String)
			err = json.Unmarshal(jsnBts, wf)
			if err == nil {
				item.Data = wf
			} else {
				d.compatibilityConv(auth, jsnBts, item, tmplId)
			}
		}
		resusts = append(resusts, item)
	}
	return resusts, nil
}

func (mw *MysqlWorkflow) Get(auth model.Authorization, id string) (*model.WorkflowItem, error) {
	var r *sql.Rows
	var err error
	if id == "" {
		return nil, fmt.Errorf("id cant be empty")
	}
	r, err = mw.db.Query("SELECT id, `accessible`, name, description, last_updated, date_created, comma_separated_form_ids, comma_separated_template_ids FROM workflow f WHERE f.id = ?", id)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	if err != nil {
		return nil, err
	}
	var (
		sid         sql.NullString
		accessible  sql.NullString
		name        sql.NullString
		desc        sql.NullString
		wfData      sql.NullString
		tmplId      sql.NullString
		lastUpdated mysql.NullTime
		dateCreated mysql.NullTime
	)
	var item *model.WorkflowItem
	for r.Next() {
		item = &model.WorkflowItem{}
		err = r.Scan(&sid, &accessible, &name, &desc, &lastUpdated, &dateCreated, &wfData, &tmplId)
		if err != nil {
			return nil, err
		}

		if accessible.Valid {
			//err handling is not important here
			//just keep in mind to convert data before reading them again if the Permissions model changes!
			err = json.Unmarshal([]byte(accessible.String), &item.Permissions)
			if !item.Permissions.IsReadGrantedFor(auth) {
				return nil, model.ErrAuthorityMissing
			}
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
			wf := &workflow.Workflow{}
			jsnBts := []byte(wfData.String)
			err = json.Unmarshal(jsnBts, wf)
			if err == nil {
				item.Data = wf
			} else {
				mw.compatibilityConv(auth, jsnBts, item, tmplId)
			}
		}
	}
	return item, nil
}

func (me *MysqlWorkflow) compatibilityConv(auth model.Authorization, jsnBts []byte, item *model.WorkflowItem, tmplId sql.NullString) {
	//compatibility case
	var d []map[string]interface{}
	err := json.Unmarshal(jsnBts, &d)
	if err == nil {
		for _, node := range d {
			if t, ok := node["type"].(string); ok && t == "condition" {
				if cid, ok := node["id"].(string); ok {
					cases, jsCode := me.getConditionJS(cid)
					if cases != "" {
						var cs interface{}
						err = json.Unmarshal([]byte(cases), &cs)
						node["cases"] = cs
					}
					if jsCode != "" {
						node["data"] = map[string]interface{}{"js": jsCode}
					}
				}

			} else if t, ok := node["type"].(string); ok && t == "step" {
				node["type"] = "form"
			}
		}
		if tmplId.Valid {
			tmplStore, err := NewTemplateStore(me.db, "")
			if err == nil {
				tmplItem, err := tmplStore.Get(auth, tmplId.String)
				if err == nil && tmplItem != nil {
					var targetId string
					if len(d) > 0 {
						node := d[0]
						if cid, ok := node["id"].(string); ok {
							targetId = fmt.Sprintf(`{"targetId":"%s"}`, cid)
						}
						tmplJson := fmt.Sprintf(`{"id":"%s","name":"%s","description":"%s","type":"template","p":{"x":-30,"y":103,"start":{"x":-235,"y":71},"seed":742587},"connectedTo":[%s]}`, tmplItem.ID, tmplItem.Name, tmplItem.Detail, targetId)
						var ddd map[string]interface{}
						err = json.Unmarshal([]byte(tmplJson), &ddd)
						if err == nil {
							d = append([]map[string]interface{}{ddd}, d...)
						}
					} else {
						tmplJson := fmt.Sprintf(`{"id":"%s","name":"%s","description":"%s","type":"template","p":{"x":-30,"y":103,"start":{"x":-235,"y":71},"seed":742587},"connectedTo":[%s]}`, tmplItem.ID, tmplItem.Name, tmplItem.Detail, targetId)
						var ddd map[string]interface{}
						err = json.Unmarshal([]byte(tmplJson), &ddd)
						if err == nil {
							d = []map[string]interface{}{ddd}
						}
					}
				}
			}
		}
		if d != nil {
			njsbts, er := json.Marshal(d)
			if er == nil {
				newBts := CompatibilityConvert(njsbts)
				if len(newBts) > 0 {
					wf := &workflow.Workflow{}
					flow := &workflow.Flow{}
					wf.Flow = flow
					err = json.Unmarshal(newBts, flow)
					if err == nil {
						item.Data = wf
						return
					}
				}
			}
		}
		bts, err := json.Marshal(map[string]interface{}{"flow": d, "progressFlow": []interface{}{}})
		if err == nil {
			json.Unmarshal(bts, item.Data)
		}
	}
}

func (me *MysqlWorkflow) getPermissions(id, tbl string) (*model.Permissions, error) {
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
	var item *model.WorkflowItem
	for r.Next() {
		item = &model.WorkflowItem{}
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

func (mw *MysqlWorkflow) Put(auth model.Authorization, item *model.WorkflowItem) error {
	if item == nil {
		return errors.New("item cannot be nil")
	}
	var err error

	var currentPerm *model.Permissions
	currentPerm, err = mw.getPermissions(item.ID, "workflow")
	if err != nil {
		return err
	}
	if currentPerm == nil || !currentPerm.IsWriteGrantedFor(auth) {
		return model.ErrAuthorityMissing
	}
	if item.ID == "" { //insert new item
		item.Created = time.Now()
		item.Updated = item.Created
		for i := 0; i < 5; i++ {
			u2 := uuid.NewV4()
			item.ID = u2.String()
			_, err = mw.db.Exec("INSERT INTO `workflow` (`id`, `version`, `comma_separated_form_ids`, `comma_separated_template_ids`, `contract_id`, `name`, `description`, `owner_id`, `price`, `published`, `sort_index`, `uri_name`, `workflow_order`, `date_created`, `last_updated`) VALUES (?,1,'','',1,?,?,1,0,0,0,?,0,?,?)", item.ID, item.Name, item.Detail, item.ID, item.Created, item.Updated)
			if err == nil {
				break
			}
		}
		if err != nil {
			return err
		}
	}

	var permBytes []byte
	currentPerm.Change(auth, &item.Permissions)
	permBytes, err = json.Marshal(currentPerm)
	if err != nil {
		return err
	}

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
	item.Updated = time.Now()
	result, err := mw.db.Exec("UPDATE `workflow` w SET w.accessible = ?, w.name = ?, w.description = ?, w.comma_separated_form_ids = ?, w.last_updated = ? WHERE w.id = ?",
		string(permBytes),
		item.Name,
		item.Detail,
		jsonDataStr,
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

func (mw *MysqlWorkflow) getConditionJS(id string) (string, string) {
	var r *sql.Rows
	var stmt *sql.Stmt
	var err error
	stmt, err = mw.db.Prepare("SELECT `cases`, `js_code` FROM `flow_condition` f WHERE f.id = ?")
	if err != nil {
		return "", ""
	}
	defer stmt.Close()
	r, err = stmt.Query(id)
	if err != nil {
		return "", ""
	}
	defer r.Close()
	var (
		sCases  sql.NullString
		sJsCode sql.NullString
	)
	cases := ""
	jsCode := ""
	for r.Next() {
		err = r.Scan(&sCases, &sJsCode)
		if err != nil {
			return "", ""
		}

		if sCases.Valid {
			cases = sCases.String
		}
		if sJsCode.Valid {
			jsCode = sJsCode.String
		}
	}
	return cases, jsCode
}

func (mw *MysqlWorkflow) Close() error {
	if mw.db != nil {
		mw.db.Close()
	}
	return nil
}
