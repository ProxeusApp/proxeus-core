package mysql

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/go-sql-driver/mysql"
	uuid "github.com/satori/go.uuid"

	"git.proxeus.com/core/central/sys/model"
)

type (
	MysqlTemplate struct {
		db       *sql.DB
		basePath string
	}
	mysqlLangs struct {
		id   uint
		code string
	}
)

//TODO wrap single updates in transactions to ensure consistency on ALL methods

func NewTemplateStore(db *sql.DB, baseFilePath string) (*MysqlTemplate, error) {
	return &MysqlTemplate{db: db, basePath: baseFilePath}, nil
}

func (me *MysqlTemplate) List(auth model.Authorization, contains string, options map[string]interface{}) ([]*model.TemplateItem, error) {
	sQuery := makeSimpleQuery(options)
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
		inQuery.WriteString("SELECT `id`, `name`, `description`, `last_updated`, `date_created` FROM `template` f WHERE f.id not in(")
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
		r, err = me.db.Query("SELECT `id`, `name`, `description`, `last_updated`, `date_created` FROM `template` f WHERE (f.name LIKE ? OR f.description LIKE ?) GROUP BY id ORDER BY f.last_updated DESC, f.date_created DESC, f.name ASC LIMIT ?,?", cnts, cnts, sQuery.index, sQuery.limit)
	}

	if err != nil {
		return nil, err
	}
	defer r.Close()
	var (
		id          sql.NullString
		name        sql.NullString
		desc        sql.NullString
		lastUpdated mysql.NullTime
		dateCreated mysql.NullTime
	)

	resusts := make([]*model.TemplateItem, 0)
	var langs []mysqlLangs
	for r.Next() {
		err = r.Scan(&id, &name, &desc, &lastUpdated, &dateCreated)
		if err != nil {
			return nil, err
		}
		item := &model.TemplateItem{}
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
		if !sQuery.metaOnly {
			if langs == nil {
				langs, err = me.getLangs()
			}
			me.setData(langs, item)
		}
		resusts = append(resusts, item)
	}
	return resusts, nil
}

func (me *MysqlTemplate) Get(auth model.Authorization, id string) (item *model.TemplateItem, err error) {
	if id == "" {
		err = fmt.Errorf("id cant be empty")
		return
	}
	var r *sql.Rows
	var stmt *sql.Stmt
	stmt, err = me.db.Prepare("SELECT id, name, description, last_updated, date_created FROM `template` f WHERE f.id = ?")
	if err != nil {
		return
	}
	defer stmt.Close()
	r, err = stmt.Query(id)
	if err != nil {
		return
	}
	defer r.Close()
	var (
		sid         sql.NullString
		name        sql.NullString
		desc        sql.NullString
		lastUpdated mysql.NullTime
		dateCreated mysql.NullTime
	)
	for r.Next() {
		item = &model.TemplateItem{}
		err = r.Scan(&sid, &name, &desc, &lastUpdated, &dateCreated)
		if err != nil {
			return
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
		var langs []mysqlLangs

		langs, err = me.getLangs()
		if err == nil {
			me.setData(langs, item)
		}
	}
	return
}
func (me *MysqlTemplate) setData(langs []mysqlLangs, item *model.TemplateItem) (err error) {
	if len(langs) > 0 {
		langTmplMap := model.TemplateLangMap{}
		var tmpl *model.Template
		for _, lang := range langs {
			_, tmpl, err = me.getLangTemplate(lang.id, item.ID)
			if tmpl != nil {
				tmpl.Lang = lang.code
				if tmpl.Path != "" {
					tmpl.Path = filepath.Join(me.basePath, tmpl.Path)
					if me.doesItExistOnDisk(tmpl.Path) {
						langTmplMap[lang.code] = tmpl
					}
				}
			}
		}
		item.Data = langTmplMap
	}
	return
}

func (me *MysqlTemplate) doesItExistOnDisk(p string) bool {
	_, err := os.Stat(p)
	if !os.IsNotExist(err) {
		return true
	}
	return false
}

func (me *MysqlTemplate) Remove(auth model.Authorization, id string) error {
	item, err := me.Get(auth, id)
	if err != nil {
		return err
	}
	if item.Data != nil {
		for _, v := range item.Data {
			err = me.deleteFileInfoAndLangTemplate(id, v.Lang)
			if err != nil {
				return err
			}
			os.Remove(v.Path)
		}
		_, err = me.db.Exec("DELETE FROM `template` WHERE `id` = ?", id)
		return err
	}
	return os.ErrInvalid
}

func (me *MysqlTemplate) RemoveTemplate(auth model.Authorization, id, lang string) error {
	_, tmpl, err := me.GetTemplate(auth, id, lang, nil)
	if err != nil {
		return err
	}
	err = me.deleteFileInfoAndLangTemplate(id, lang)
	if err != nil {
		return err
	}
	me.updateLastUpdated(id)
	return os.Remove(tmpl.Path)
}

// On return, n == len(buf) if and only if err == nil.
func (me *MysqlTemplate) GetTemplate(auth model.Authorization, id, lang string, writer io.Writer) (n int64, tmpl *model.Template, err error) {
	var langID uint
	if id == "" {
		err = ErrorTmplIDCannotBeEmpty
		return
	}
	if lang == "" {
		err = errors.New("lang cannot be empty")
		return
	}
	tmpl.Lang = lang
	langID, err = me.getLangID(tmpl.Lang)
	if err != nil {
		return
	}
	if langID > 0 {
		_, tmpl, err = me.getLangTemplate(langID, id)
		if err != nil {
			return
		}
		if tmpl.Path == "" {
			err = os.ErrNotExist
			return
		}
		filename := tmpl.Path
		tmpl.Path = filepath.Join(me.basePath, filename)
		if writer != nil {
			n, err = me.readFileFromDisk(&filename, writer)
		}
		return
	}
	err = os.ErrNotExist
	return
}

//Put creates an ID if the ID is "" but it's important to handle the permissions of this behaviour outside Put
func (me *MysqlTemplate) Put(auth model.Authorization, item *model.TemplateItem) (err error) {
	if item == nil {
		return errors.New("item cannot be nil")
	}

	err = me.insertNewTmplIfIDEmpty(item)

	if err != nil {
		return
	}

	item.Updated = time.Now()
	result, err := me.db.Exec("UPDATE `template` f SET f.name = ?, f.description = ?, f.last_updated = ? WHERE f.id = ?",
		item.Name,
		item.Detail,
		item.Updated,
		item.ID)
	if err != nil {
		return
	}
	_, err = result.RowsAffected()
	return
}

var ErrLangNotFound = errors.New("lang not found")

// On return, written == len(buf) if and only if err == nil.
func (me *MysqlTemplate) PutTemplate(auth model.Authorization, id string, tmpl *model.Template, reader io.Reader) (written int64, err error) {
	ti := &model.TemplateItem{}
	ti.ID = id
	err = me.insertNewTmplIfIDEmpty(ti)
	if err != nil {
		return
	}
	var langID uint
	langID, err = me.getLangID(tmpl.Lang)
	if err != nil {
		return
	}
	if langID > 0 {
		var (
			rtmpl      *model.Template
			fileInfoId uint
		)
		fileInfoId, rtmpl, err = me.putLangTemplate(langID, id, tmpl)
		if err != nil {
			return
		}
		rtmpl.Size, err = me.writeFileToDisk(&rtmpl.Path, reader)
		if err == nil {
			if fileInfoId > 0 {
				me.putSizeToFileInfo(&fileInfoId, rtmpl)
			}
		}
		err = me.updateLastUpdated(id)
	} else {
		err = ErrLangNotFound
		return
	}
	return
}

func (me *MysqlTemplate) updateLastUpdated(id string) error {
	_, err := me.db.Exec("UPDATE `template` f SET f.last_updated = ? WHERE f.id = ?", time.Now(), id)
	return err
}

var (
	ErrorLangIDCannotBeEmpty = errors.New("langID cannot be empty")
	ErrorTmplIDCannotBeEmpty = errors.New("langID cannot be empty")
)

func (me *MysqlTemplate) putLangTemplate(langID uint, tmplID string, tmpl *model.Template) (fileInfoId uint, rtmpl *model.Template, err error) {
	fileInfoId, rtmpl, err = me.getLangTemplate(langID, tmplID)
	if err != nil {
		return
	}
	if fileInfoId != 0 {
		_, err = me.db.Exec("UPDATE `file_info` f SET f.`file_size`= ?, f.`file_type`= ?, f.`org_filename`= ? WHERE f.id = ?", tmpl.Size, tmpl.ContentType, tmpl.Filename, fileInfoId)
		return
	}
	//entries missing in the db
	rtmpl, err = me.putFileInfo(&fileInfoId, tmpl)
	if err != nil {
		return
	}
	//var rrr sql.Result
	_, err = me.db.Exec("INSERT INTO `lang_template` (`language_id`, `version`, `template_id`, `document_id`) VALUES (?, 1, ?, ?)", langID, tmplID, fileInfoId)
	return
}

func (me *MysqlTemplate) getLangTemplate(langID uint, tmplID string) (fileInfoId uint, tmpl *model.Template, err error) {
	var r *sql.Rows
	var stmt *sql.Stmt
	if langID == 0 {
		err = ErrorLangIDCannotBeEmpty
		return
	}
	if tmplID == "" {
		err = ErrorTmplIDCannotBeEmpty
		return
	}
	stmt, err = me.db.Prepare("SELECT `document_id` FROM `lang_template` f WHERE f.language_id = ? AND f.template_id = ? LIMIT ?,?;")
	if err != nil {
		return
	}
	defer stmt.Close()
	r, err = stmt.Query(langID, tmplID, 0, 1)
	if err != nil {
		return
	}
	defer r.Close()
	var (
		id sql.NullInt64
	)
	for r.Next() {
		err = r.Scan(&id)
		if err != nil {
			return
		}
		if id.Valid {
			fileInfoId = uint(id.Int64)
			tmpl, err = me.getFileInfo(&fileInfoId)
			return
		}
	}
	return
}

func (me *MysqlTemplate) writeFileToDisk(filename *string, reader io.Reader) (written int64, err error) {
	var tmplFile *os.File
	tmplFile, err = me.openFile(filename)
	if os.IsExist(err) {
		err = nil
	}
	defer tmplFile.Close()
	written, err = io.Copy(tmplFile, reader)
	return
}

func (me *MysqlTemplate) readFileFromDisk(filename *string, writer io.Writer) (n int64, err error) {
	if filename != nil && writer != nil {
		var tmplFile *os.File
		tmplFile, err = os.OpenFile(filepath.Join(me.basePath, *filename), os.O_RDONLY, 0600)
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
	err = os.ErrNotExist
	return
}

func (me *MysqlTemplate) openFile(p *string) (*os.File, error) {
	return os.OpenFile(filepath.Join(me.basePath, *p), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
}

func (me *MysqlTemplate) getLangID(lang string) (langID uint, err error) {
	var r *sql.Rows
	var stmt *sql.Stmt
	if lang == "" {
		err = errors.New("lang cannot be empty")
		return
	}
	stmt, err = me.db.Prepare("SELECT `id` FROM `language` f WHERE f.code = ? LIMIT ?,?;")
	if err != nil {
		return
	}
	defer stmt.Close()
	r, err = stmt.Query(lang, 0, 1)
	if err != nil {
		return
	}
	defer r.Close()
	var (
		id sql.NullInt64
	)
	for r.Next() {
		err = r.Scan(&id)
		if err != nil {
			return
		}

		if id.Valid {
			langID = uint(id.Int64)
			return
		}
	}
	return
}

func (me *MysqlTemplate) getLangs() (langs []mysqlLangs, err error) {
	var r *sql.Rows
	r, err = me.db.Query("SELECT `id`, `code` FROM `language` f WHERE f.enabled = 1;")
	if err != nil {
		return
	}
	defer r.Close()
	var (
		id   sql.NullInt64
		code sql.NullString
	)
	langs = make([]mysqlLangs, 0)
	for r.Next() {
		err = r.Scan(&id, &code)
		if err != nil {
			return
		}
		lang := mysqlLangs{}
		if id.Valid {
			lang.id = uint(id.Int64)
		}
		if code.Valid {
			lang.code = code.String
			if lang.code == "*" {
				lang.code = "en"
			}
		}
		langs = append(langs, lang)
	}
	return
}

func (me *MysqlTemplate) insertNewTmplIfIDEmpty(item *model.TemplateItem) (err error) {
	if item.ID == "" { //insert new item
		item.Created = time.Now()
		item.Updated = item.Created
		for i := 0; i < 5; i++ {
			var u2 uuid.UUID
			u2 = uuid.NewV4()
			item.ID = u2.String()
			_, err = me.db.Exec("INSERT INTO `template` (`id`, `version`, `name`, `description`, `sort_index`, `contract_id`, `date_created`, `last_updated`) VALUES (?, 1, '', '', 0, 1, ?, ?)", item.ID, item.Created, item.Updated)
			if err == nil {
				break
			}
		}
	}
	return
}

func (me *MysqlTemplate) putFileInfo(fileInfoId *uint, tmpl *model.Template) (rtmpl *model.Template, err error) {
	if fileInfoId == nil || *fileInfoId == 0 { //insert new item
		var u2 uuid.UUID
		u2 = uuid.NewV4()
		tmpl.Path = u2.String()
		var r sql.Result
		r, err := me.db.Exec("INSERT INTO `file_info` (`version`, `file_size`, `file_type`, `filename`, `org_filename`, `temporary`) VALUES (1, ?, ?, ?, ?, 0)", tmpl.Size, tmpl.ContentType, tmpl.Path, tmpl.Filename)
		if err == nil {
			var id int64
			id, err = r.LastInsertId()
			*fileInfoId = uint(id)
		}
		rtmpl = tmpl
	} else {
		_, err = me.db.Exec("UPDATE `file_info` f SET f.`file_size`= ?, f.`file_type`= ?, f.`org_filename`= ? WHERE f.id = ?", tmpl.Size, tmpl.ContentType, tmpl.Filename, *fileInfoId)
		tmpl, err = me.getFileInfo(fileInfoId)
	}
	return
}

func (me *MysqlTemplate) deleteFileInfoAndLangTemplate(id, lang string) error {
	langID, err := me.getLangID(lang)
	if err != nil {
		return err
	}
	_, err = me.db.Exec("DELETE `lang_template`, `file_info` FROM `lang_template` JOIN `file_info` ON `document_id`=`id` WHERE `language_id` = ? AND `template_id` = ?", langID, id)
	return err
}

func (me *MysqlTemplate) putSizeToFileInfo(fileInfoId *uint, tmpl *model.Template) (err error) {
	_, err = me.db.Exec("UPDATE `file_info` f SET f.`file_size`= ? WHERE f.id = ?", tmpl.Size, *fileInfoId)
	return
}

func (me *MysqlTemplate) getFileInfo(fileInfoId *uint) (tmpl *model.Template, err error) {
	var r *sql.Rows
	var stmt *sql.Stmt
	stmt, err = me.db.Prepare("SELECT `filename`,`file_size`, `file_type`, `org_filename` FROM `file_info` f WHERE f.id=?;")
	if err != nil {
		return
	}
	defer stmt.Close()
	r, err = stmt.Query(*fileInfoId)
	if err != nil {
		return
	}
	defer r.Close()
	var (
		size        sql.NullInt64
		contentType sql.NullString
		filename    sql.NullString
		orgFilename sql.NullString
	)
	tmpl = &model.Template{}
	for r.Next() {
		err = r.Scan(&filename, &size, &contentType, &orgFilename)
		if err != nil {
			return
		}
		if filename.Valid {
			tmpl.Path = filename.String
		}
		if size.Valid {
			tmpl.Size = size.Int64
		}
		if contentType.Valid {
			tmpl.ContentType = contentType.String
		}
		if orgFilename.Valid {
			tmpl.Filename = orgFilename.String
		}
	}
	return
}

func (me *MysqlTemplate) Vars(auth model.Authorization, contains string) ([]string, error) {
	var r *sql.Rows
	var stmt *sql.Stmt
	var err error
	stmt, err = me.db.Prepare("SELECT varname FROM _variable_ f WHERE f.varname LIKE ? GROUP by varname ORDER BY varname ASC LIMIT ?,?;")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	cnts := "%" + contains + "%"
	r, err = stmt.Query(cnts, 0, 15)
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
			vars = append(vars, "input."+name.String)
		}
	}
	return vars, nil
}

func (me *MysqlTemplate) vars(varMap map[string]string) error {
	return nil
}

func (mt *MysqlTemplate) Close() error {
	return nil
}
