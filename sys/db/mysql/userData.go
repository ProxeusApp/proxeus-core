package mysql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	uuid "github.com/satori/go.uuid"

	"git.proxeus.com/core/central/sys/file"
	"git.proxeus.com/core/central/sys/model"
)

type MysqlUserData struct {
	userId       string
	baseFilePath string
	db           *sql.DB
}

func NewUserDataStore(db *sql.DB, userId, baseFilePath string) (res *MysqlUserData, err error) {
	return &MysqlUserData{db: db, userId: userId, baseFilePath: baseFilePath}, nil
}

//ListData can list multiple data sets of a specific user
// dataPath: &"userid.contingent.datavar"
// contains: &"some content to search for under the given dataPath
// options: limit, index, sort, metaOnly ...
func (me *MysqlUserData) ListData(auth model.Authorization, dataPath, contains, baseUri string, options map[string]interface{}) (res []*model.UserData, err error) {
	var dpath string
	dpath, _, err = checkAndGetDataPath(&dataPath)
	if err != nil {
		return nil, os.ErrInvalid
	}
	sQuery := makeSimpleQuery(options)
	var r *sql.Rows
	if contains != "" {
		cnts := "%" + contains + "%"
		r, err = me.db.Query("SELECT wd.id, wd.user_name, wd.last_updated, wd.date_created, wd.workflow_id, wd.root_lang, wd.form_lang, wd.tmpl_lang, wd.finished, JSON_EXTRACT(wd.json_data, '"+dpath+"') as myJson FROM workflow_data wd WHERE wd.user_id = ? AND wd.json_data is not null AND wd.json_data !='' AND wd.user_name LIKE ? ORDER BY wd.last_updated DESC, wd.date_created DESC, wd.user_name ASC LIMIT ?,?;",
			me.userId, cnts, sQuery.index, sQuery.limit)
	} else {
		r, err = me.db.Query("SELECT wd.id, wd.user_name, wd.last_updated, wd.date_created, wd.workflow_id, wd.root_lang, wd.form_lang, wd.tmpl_lang, wd.finished, JSON_EXTRACT(wd.json_data, '"+dpath+"') as myJson FROM workflow_data wd WHERE wd.user_id = ? AND wd.json_data is not null AND wd.json_data !='' ORDER BY wd.last_updated DESC, wd.date_created DESC, wd.user_name ASC LIMIT ?,?;",
			me.userId, sQuery.index, sQuery.limit)
	}

	if err != nil {
		return nil, err
	}
	defer r.Close()
	var (
		id          sql.NullString
		finished    sql.NullString
		wfid        sql.NullString
		wLang       sql.NullString
		wLangForm   sql.NullString
		wLangTmpl   sql.NullString
		name        sql.NullString
		desc        sql.NullString
		wfData      sql.NullString
		lastUpdated mysql.NullTime
		dateCreated mysql.NullTime
	)

	res = make([]*model.UserData, 0)
	for r.Next() {
		err = r.Scan(&id, &name, &lastUpdated, &dateCreated, &wfid, &wLang, &wLangForm, &wLangTmpl, &finished, &wfData)
		if err != nil {
			return nil, err
		}
		usrData := &model.UserData{}
		item := &model.UserData{Data: usrData}
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
		if wfid.Valid {
			usrData.WorkflowID = wfid.String
		}
		if wLang.Valid {
			usrData.Lang = wLang.String
		}
		if wLangForm.Valid {
			usrData.LangForm = wLangForm.String
		}
		if wLangTmpl.Valid {
			usrData.LangTmpl = wLangTmpl.String
		}
		if finished.Valid {
			usrData.Finished = finished.String == "\x01"
		}
		if !sQuery.metaOnly && wfData.Valid {
			var d interface{}
			err = json.Unmarshal([]byte(wfData.String), &d)
			if err == nil && d != nil {
				usrData.Data = me.makeFileInfosForFiles(baseUri, d)
			}
		}
		res = append(res, item)
	}
	return res, nil
}

//GetDataFile returns a file proxy interface
func (me *MysqlUserData) GetDataFile(auth model.Authorization, id, dataPath, contains, baseUri string) (f file.Info, err error) {
	//var item *model.UserData
	item, err := me.getData(auth, id, dataPath, contains, baseUri, nil, true)
	if item != nil {
		//if item.Data != nil {
		//	if ud, ok := item.Data.(*model.UserData); ok {
		//		if ud.Data != nil {
		//			f, _ = ud.Data.(file.Info)
		//		}
		//	}
		//}
	}
	if err == nil && f == nil {
		err = os.ErrNotExist
	}
	return
}

//GetDataFileWriter returns a file proxy interface
func (me *MysqlUserData) GetDataFileWriter(auth model.Authorization, id, dataPath, contains, baseUri string, writer io.Writer) (n int64, f file.Info, err error) {
	f, err = me.GetDataFile(auth, id, dataPath, contains, baseUri)
	if err == nil {
		if writer != nil {
			n, err = me.readFileFromDisk(f, writer)
		}
	}
	return
}

func (me *MysqlUserData) writeFileToDisk(f file.Info, reader io.Reader) (written int64, err error) {
	if reader == nil || f == nil {
		err = os.ErrInvalid
		return
	}
	var tmplFile *os.File
	tmplFile, err = os.OpenFile(f.Path(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if os.IsExist(err) {
		err = nil
	}
	defer tmplFile.Close()
	written, err = io.Copy(tmplFile, reader)
	return
}

func (me *MysqlUserData) readFileFromDisk(f file.Info, writer io.Writer) (n int64, err error) {
	var tmplFile *os.File
	tmplFile, err = os.OpenFile(f.Path(), os.O_RDONLY, 0600)
	if err != nil {
		return
	}
	defer tmplFile.Close()
	var fstat os.FileInfo
	fstat, err = tmplFile.Stat()
	if err != nil {
		return
	}
	n, err = io.CopyN(writer, tmplFile, fstat.Size())
	return
}

var selectQuery = "SELECT wd.id, wd.user_name, wd.last_updated, wd.date_created, wd.workflow_id, wd.root_lang, wd.form_lang, wd.tmpl_lang, wd.finished, JSON_EXTRACT(wd.json_data, '%s') as myJson FROM workflow_data wd WHERE"

//GetData delivers any kind of data under the given path
// if no data exists under this path "nil" is returned
func (me *MysqlUserData) GetData(auth model.Authorization, id, dataPath, contains, baseUri string, example *model.UserData) (item interface{}, err error) {
	item, err = me.getData(auth, id, dataPath, contains, baseUri, example, true)
	return
}

func (me *MysqlUserData) getData(auth model.Authorization, id, dataPath, contains, baseUri string, example *model.UserData, makeFileInfos bool) (item interface{}, err error) {
	var dpath string
	dpath, _, err = checkAndGetDataPath(&dataPath)
	if err != nil {
		err = os.ErrInvalid
		return
	}
	var r *sql.Rows
	if id != "" {
		r, err = me.db.Query(fmt.Sprintf("%s wd.id = ? AND wd.json_data is not null AND wd.json_data !='';", fmt.Sprintf(selectQuery, dpath)), id)
	} else {
		if contains != "" {
			cnts := "%" + contains + "%"
			if example == nil {
				r, err = me.db.Query(fmt.Sprintf("%s wd.user_id = ? AND wd.json_data is not null AND wd.json_data !='' AND wd.user_name LIKE ?;", fmt.Sprintf(selectQuery, dpath)), me.userId, cnts)
			} else {
				if example.WorkflowID != "" && example.LangTmpl != "" && example.LangForm != "" && example.Lang != "" {
					r, err = me.db.Query(fmt.Sprintf("%s wd.user_id = ? AND wd.json_data is not null AND wd.json_data !='' AND wd.user_name LIKE ? AND wd.workflow_id = ? AND wd.root_lang = ? AND wd.form_lang = ? AND wd.tmpl_lang = ?;", fmt.Sprintf(selectQuery, dpath)), me.userId, cnts, example.WorkflowID, example.Lang, example.LangForm, example.LangTmpl)
				} else if example.WorkflowID != "" {
					r, err = me.db.Query(fmt.Sprintf("%s wd.user_id = ? AND wd.json_data is not null AND wd.json_data !='' AND wd.user_name LIKE ? AND wd.workflow_id = ?;", fmt.Sprintf(selectQuery, dpath)), me.userId, cnts, example.WorkflowID)
				} else if example.LangTmpl != "" && example.LangForm != "" && example.Lang != "" {
					r, err = me.db.Query(fmt.Sprintf("%s wd.user_id = ? AND wd.json_data is not null AND wd.json_data !='' AND wd.user_name LIKE ? AND wd.root_lang = ? AND wd.form_lang = ? AND wd.tmpl_lang = ?;", fmt.Sprintf(selectQuery, dpath)), me.userId, cnts, example.Lang, example.LangForm, example.LangTmpl)
				} else {
					r, err = me.db.Query(fmt.Sprintf("%s wd.user_id = ? AND wd.json_data is not null AND wd.json_data !='' AND wd.user_name LIKE ?;", fmt.Sprintf(selectQuery, dpath)), me.userId, cnts)
				}
			}
		} else {
			if example == nil {
				r, err = me.db.Query(fmt.Sprintf("%s wd.user_id = ? AND wd.json_data is not null AND wd.json_data !='';", fmt.Sprintf(selectQuery, dpath)), me.userId)
			} else {
				if example.WorkflowID != "" && example.LangTmpl != "" && example.LangForm != "" && example.Lang != "" {
					r, err = me.db.Query(fmt.Sprintf("%s wd.user_id = ? AND wd.json_data is not null AND wd.json_data !='' AND wd.workflow_id = ? AND wd.root_lang = ? AND wd.form_lang = ? AND wd.tmpl_lang = ? AND wd.finished=?;", fmt.Sprintf(selectQuery, dpath)), me.userId, example.WorkflowID, example.Lang, example.LangForm, example.LangTmpl, example.Finished)
				} else if example.WorkflowID != "" {
					r, err = me.db.Query(fmt.Sprintf("%s wd.user_id = ? AND wd.json_data is not null AND wd.json_data !='' AND wd.workflow_id = ? AND wd.finished=?;", fmt.Sprintf(selectQuery, dpath)), me.userId, example.WorkflowID, example.Finished)
				} else if example.LangTmpl != "" && example.LangForm != "" && example.Lang != "" {
					r, err = me.db.Query(fmt.Sprintf("%s wd.user_id = ? AND wd.json_data is not null AND wd.json_data !='' AND wd.root_lang = ? AND wd.form_lang = ? AND wd.tmpl_lang = ? AND wd.finished=?;", fmt.Sprintf(selectQuery, dpath)), me.userId, example.Lang, example.LangForm, example.LangTmpl, example.Finished)
				} else {
					r, err = me.db.Query(fmt.Sprintf("%s wd.user_id = ? AND wd.json_data is not null AND wd.json_data !='' AND wd.finished=?;", fmt.Sprintf(selectQuery, dpath)), me.userId, example.Finished)
				}
			}
		}
	}

	if err != nil {
		return
	}
	defer r.Close()
	var (
		mid       sql.NullString
		finished  sql.NullString
		wfid      sql.NullString
		wLang     sql.NullString
		wLangForm sql.NullString
		wLangTmpl sql.NullString
		name      sql.NullString
		//desc         sql.NullString
		wfData      sql.NullString
		lastUpdated mysql.NullTime
		dateCreated mysql.NullTime
	)

	for r.Next() {
		err = r.Scan(&mid, &name, &lastUpdated, &dateCreated, &wfid, &wLang, &wLangForm, &wLangTmpl, &finished, &wfData)
		if err != nil {
			return
		}
		usrData := &model.UserData{}
		item = &model.UserData{Data: usrData}
		//if mid.Valid {
		//	item.ID = mid.String
		//}
		//if name.Valid {
		//	item.Name = name.String
		//}
		if wfid.Valid {
			usrData.WorkflowID = wfid.String
		}
		if wLang.Valid {
			usrData.Lang = wLang.String
		}
		if wLangForm.Valid {
			usrData.LangForm = wLangForm.String
		}
		if wLangTmpl.Valid {
			usrData.LangTmpl = wLangTmpl.String
		}
		if finished.Valid {
			usrData.Finished = finished.String == "\x01"
		}
		if wfData.Valid {
			var d interface{}
			err = json.Unmarshal([]byte(wfData.String), &d)
			if err == nil && d != nil {
				if makeFileInfos {
					usrData.Data = me.makeFileInfosForFiles(baseUri, d)
				} else {
					usrData.Data = d
				}
			}
		}
		return
	}
	return
}

func (me *MysqlUserData) makeFileInfosForFiles(baseUri string, d interface{}) interface{} {
	if d != nil {
		if m, ok := d.(map[string]interface{}); ok {
			if file.IsFileInfo(m) {
				return file.FromMap(me.baseFilePath, baseUri, m)
			} else {
				me.makeFileInfos(baseUri, m, "")
			}
		}
	}
	return d
}

func (me *MysqlUserData) makeFileInfos(baseUri string, m map[string]interface{}, path string) {
	if len(m) > 0 {
		for k, item := range m {
			if item != nil {
				if itemMap, ok := item.(map[string]interface{}); ok {
					if file.IsFileInfo(itemMap) {
						var p string
						if path == "" {
							p = path + k
						} else {
							p = path + "." + k
						}
						bu := filepath.Join(baseUri, p)
						m[k] = file.FromMap(me.baseFilePath, bu, itemMap)
					} else {
						var p string
						if path == "" {
							p = path + k
						} else {
							p = path + "." + k
						}
						me.makeFileInfos(baseUri, itemMap, p)
					}
				}
			}
		}
	}
}

//PutData makes it possible to persist data under the given path
func (me *MysqlUserData) PutData(auth model.Authorization, id, dataPath string, it interface{}) (err error) {
	if it == nil {
		err = errors.New("item cannot be nil")
		return
	}

	//to prevent from compile errors only!
	var item *model.UserData

	var dp string
	_, dp, err = checkAndGetDataPath(&dataPath)
	if err != nil {
		err = os.ErrInvalid
		return
	}
	var ud *model.UserData
	var workflowID string
	if item.Data != nil {
		ud, _ = item.Data.(*model.UserData)
		if ud != nil {
			workflowID = ud.WorkflowID
		}
	}
	if id == "" { //insert new item
		n := time.Now()
		item.Created = n
		var res sql.Result
		res, err = me.db.Exec("INSERT INTO `workflow_data` (`version`, `confirmed`, `draft`, `encrypted`, `finished`, `json_data`, `json_progress_tracker`, `skip_encryption`, `user_id`, `verifiable`, `workflow_id`, `workflow_progress_id`, `last_updated`, `date_created`) VALUES (1,0,0,0,0,'','',0,?,0,?,1,?,?);", me.userId, workflowID, n, n)
		if err != nil {
			return
		}
		var iid int64
		iid, err = res.LastInsertId()
		if err != nil {
			return
		}
		item.ID = strconv.FormatInt(iid, 10)
		id = item.ID
	}

	if item.ID == "" {
		item.ID = id
	}

	root := ""
	//var rootItem *model.UserData
	var rootItem interface{}
	var mergedUsrData interface{} //TODO sync in tx
	if item.Data != nil {
		rootItem, err = me.getData(auth, id, root, "", "", nil, false)
		if err != nil {
			return
		}
		log.Println(rootItem, dp)
	}
	var jsonDataStr string
	if mergedUsrData != nil {
		var jsonData []byte
		jsonData, err = json.Marshal(mergedUsrData)
		if err != nil {
			jsonDataStr = ""
		} else {
			jsonDataStr = string(jsonData)
		}
	}
	var result sql.Result
	un := time.Now()
	item.Updated = un
	if item.Name == "" && item.Detail == "" {
		if mergedUsrData != nil {
			if ud != nil {
				result, err = me.db.Exec("UPDATE workflow_data f SET f.json_data = ?, f.finished = ?, f.last_updated = ? WHERE f.id = ?",
					jsonDataStr,
					ud.Finished,
					un,
					item.ID)
			} else {
				result, err = me.db.Exec("UPDATE workflow_data f SET f.json_data = ?, f.last_updated = ? WHERE f.id = ?",
					jsonDataStr,
					un,
					item.ID)
			}
		} else {
			return
		}
	} else {
		if mergedUsrData != nil {
			if ud != nil {
				result, err = me.db.Exec("UPDATE workflow_data f SET f.user_name = ?, f.json_data = ?, f.finished = ?, f.last_updated = ? WHERE f.id = ?",
					item.Name,
					//item.Detail,
					jsonDataStr,
					ud.Finished,
					un,
					item.ID)
			} else {
				result, err = me.db.Exec("UPDATE workflow_data f SET f.user_name = ?, f.json_data = ?, f.last_updated = ? WHERE f.id = ?",
					item.Name,
					//item.Detail,
					jsonDataStr,
					un,
					item.ID)
			}
		} else {
			result, err = me.db.Exec("UPDATE workflow_data f SET f.user_name = ?, f.last_updated = ? WHERE f.id = ?",
				item.Name,
				un,
				item.ID)
		}
	}
	if err != nil {
		return err
	}
	_, err = result.RowsAffected()
	if err != nil {
		return
	}
	return
}

func (me *MysqlUserData) PutDataFile(auth model.Authorization, id, dataPath string, f file.Info, reader io.Reader) (written int64, err error) {
	var existingFileInfo file.Info
	existingFileInfo, err = me.GetDataFile(auth, id, dataPath, "", "")
	if err == nil && existingFileInfo != nil {
		if err == nil {
			written, err = me.writeFileToDisk(existingFileInfo, reader)
			if err == nil {
				item := &model.UserData{ID: id, Data: map[string]interface{}{"name": f.Name(), "path": filepath.Base(existingFileInfo.Path()), "contentType": f.ContentType(), "size": written}}
				err = me.PutData(auth, id, dataPath, item)
			}
		}
	} else if os.IsNotExist(err) {
		// new file
		var u2 uuid.UUID
		u2 = uuid.NewV4()
		filename := u2.String()
		if err == nil {
			newFileInfo := file.Output(me.baseFilePath, "", filename, f.Meta())
			written, err = me.writeFileToDisk(newFileInfo, reader)
			if err == nil {
				item := &model.UserData{ID: id, Data: map[string]interface{}{"name": f.Name(), "path": filename, "contentType": f.ContentType(), "size": written}}
				err = me.PutData(auth, id, dataPath, item)
			}
		}
	}
	return
}

func (me *MysqlUserData) NewFileInfo(auth model.Authorization, name, contentType string, size int64) file.Info {
	return file.New(me.baseFilePath, file.Meta{Name: name, ContentType: contentType, Size: size})
}

func (me *MysqlUserData) Close() error {
	return nil
}

func checkAndGetDataPath(dataPath *string) (string, string, error) {
	if dataPath == nil || *dataPath == "" {
		return "$", "", nil
	}
	if !IsValidDataPath(*dataPath) {
		return "", "", os.ErrInvalid
	}
	pcs := dataPathUnmarshal(*dataPath)
	dp := "$"
	for _, p := range pcs {
		dp += fmt.Sprintf(`."%s"`, p)
	}
	return dp, *dataPath, nil
}

var isValidDataPathRegex = regexp.MustCompile(`^[a-zA-Z]?[a-zA-Z0-9\.\-]+[^\.]$`)

func IsValidDataPath(dataPath string) bool {
	return isValidDataPathRegex.MatchString(dataPath)
}

func dataPathUnmarshal(dataPath string) []string {
	return strings.Split(dataPath, ".")
}

func mergeData(all interface{}, path string, toMerge interface{}) (merged interface{}) {
	if toMerge == nil {
		return nil
	}
	if strings.HasPrefix(path, "$") {
		path = path[1:]
	}
	if all != nil {
		if allMap, ok := all.(map[string]interface{}); ok {
			mergeInner(allMap, dataPathUnmarshal(path), 0, toMerge)
			return allMap
		}
	}
	newMap := map[string]interface{}{}
	mergeInner(newMap, dataPathUnmarshal(path), 0, toMerge)
	return newMap
}

func mergeInner(target map[string]interface{}, path []string, i int, toMerge interface{}) {
	var t interface{}
	if path[i] == "" {
		t = target
	} else {
		t = target[path[i]]
	}
	if t != nil {
		if len(path)-1 == i {
			if tm, ok := t.(map[string]interface{}); ok {
				if tmm, ok := toMerge.(map[string]interface{}); ok {
					for k, v := range tmm {
						tm[k] = v
					}
				} else {
					target[path[i]] = toMerge
				}
			} else {
				target[path[i]] = toMerge
			}
		} else {
			if tmm, ok := t.(map[string]interface{}); ok {
				mergeInner(tmm, path, i+1, toMerge)
			} else {
				newMap := map[string]interface{}{}
				target[path[i]] = newMap
				mergeInner(newMap, path, i+1, toMerge)
			}
		}
	} else {
		if len(path)-1 == i {
			target[path[i]] = toMerge
		} else {
			newMap := map[string]interface{}{}
			target[path[i]] = newMap
			mergeInner(newMap, path, i+1, toMerge)
		}
	}
}
