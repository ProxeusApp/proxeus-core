package mysql

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/disintegration/imaging"
	"github.com/go-sql-driver/mysql"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"

	d "git.proxeus.com/core/central/sys/db"
	"git.proxeus.com/core/central/sys/model"
)

type MysqlUser struct {
	db           *sql.DB
	baseFilePath string
}

func NewUserStore(db *sql.DB, baseFilePath string) (res *MysqlUser, err error) {
	usr := &MysqlUser{db: db, baseFilePath: baseFilePath}
	usr.setup()
	return usr, nil
}

func (me *MysqlUser) setup() (err error) {
	_, err = me.db.Exec("Update `user` u SET u.role = (select IF(role_id=2,10,5) from user_to_role utr where utr.user_id=u.id) where u.role is null;")
	_, _ = me.db.Exec("ALTER TABLE `user` ADD COLUMN `last_updated` DATETIME NULL DEFAULT NULL AFTER `id`;")
	_, _ = me.db.Exec("ALTER TABLE `user` ADD COLUMN `ethereum_public_address` VARCHAR(255) NULL DEFAULT NULL AFTER `id`;")
	_, _ = me.db.Exec("ALTER TABLE `user` ADD COLUMN `role` BIGINT(20) NULL DEFAULT NULL AFTER `id`;")
	_, _ = me.db.Exec("ALTER TABLE `user` ADD COLUMN `detail` LONGTEXT NULL DEFAULT NULL AFTER `id`;")
	_, _ = me.db.Exec("ALTER TABLE `user` ADD COLUMN `accessible` LONGTEXT NULL DEFAULT NULL AFTER `id`;")
	_, _ = me.db.Exec("ALTER TABLE `user` ADD COLUMN `photo` LONGTEXT NULL DEFAULT NULL AFTER `id`;")
	return
}

var ErrBadParams = errors.New("bad params")
var ErrBadCredentials = errors.New("bad credentials")

//Login returns an user *Item if and only if the execution was successfully
//otherwise a nil, someErr is returned
func (me *MysqlUser) Login(ltype, name, uid string) (item *model.User, err error) {
	log.Println(ltype, name, uid)
	if uid != "" {
		var (
			id          sql.NullString
			accessible  sql.NullString
			role        sql.NullInt64
			dbemail     sql.NullString
			dbpassword  sql.NullString
			lastUpdated mysql.NullTime
			dateCreated mysql.NullTime
		)
		if ltype == "walletSign" {
			// check DB and validate for error
			err = me.db.QueryRow("SELECT id,`accessible`,role,email,password,last_updated,date_created FROM user WHERE ethereum_public_address=?", uid).Scan(&id, &accessible, &role, &dbemail, &dbpassword, &lastUpdated, &date_created)
			if err == sql.ErrNoRows {
				err = d.ErrNotFound
				return
			}
			if err != nil {
				return
			}
		} else {
			if name != "" && uid != "" {
				// check DB and validate for error
				err = me.db.QueryRow("SELECT id,`accessible`,role,email,password,last_updated,date_created FROM user WHERE email=?", name).Scan(&id, &accessible, &role, &dbemail, &dbpassword, &lastUpdated, &date_created)
				if err != nil {
					return
				}
				// check for password
				if id.Valid && dbemail.Valid && dbpassword.Valid {
					err = bcrypt.CompareHashAndPassword([]byte(dbpassword.String), []byte(uid))
					if err != nil {
						err = ErrBadCredentials
						return
					}
				} else {
					err = ErrBadCredentials
					return
				}
			}
		}
		if id.Valid && role.Valid {
			auth := &model.User{}
			auth.ID = id.String
			auth.Role = model.Role(role.Int64)
			return me.get(auth, id.String)
		}
	}
	err = ErrBadParams
	return
}

func (me *MysqlUser) GetPw(id string) string {
	var dbpassword sql.NullString
	err := me.db.QueryRow("SELECT `password` FROM `user` u WHERE u.id=?", id).Scan(&dbpassword)
	if err != nil {
		return ""
	}
	if dbpassword.Valid {
		return dbpassword.String
	}
	return ""
}

//List all users
func (me *MysqlUser) List(auth model.Authorization, contains string, options map[string]interface{}) ([]*model.User, error) {
	sQuery := makeSimpleQuery(options)
	var err error
	var r *sql.Rows
	cnts := "%"
	cntsCount := 0
	if contains != "" {
		cntsCount = 3
		cnts = "%" + contains + "%"
	}
	exlen := len(sQuery.exclude)
	inlen := len(sQuery.include)
	if exlen > 0 || inlen > 0 {
		params := make([]interface{}, inlen+exlen+cntsCount+2)
		ind := 0
		inQuery := &bytes.Buffer{}
		inQuery.WriteString("SELECT id, photo,`accessible`, email, ethereum_public_address, role, last_updated, date_created, profile_json FROM user f WHERE ")
		if inlen > 0 {
			inQuery.WriteString("f.id in(")
			i := 0
			for a := range sQuery.include {
				if i > 0 {
					inQuery.WriteString(",?")
				} else {
					inQuery.WriteString("?")
				}
				i++
				params[ind] = a
				ind++
			}
			inQuery.WriteString(")")
		}
		if exlen > 0 {
			if inlen > 0 {
				inQuery.WriteString(" AND ")
			}
			inQuery.WriteString("f.id not in(")
			i := 0
			for a := range sQuery.exclude {
				if i > 0 {
					inQuery.WriteString(",?")
				} else {
					inQuery.WriteString("?")
				}
				i++
				params[ind] = a
				ind++
			}
			inQuery.WriteString(")")
		}
		if contains != "" {
			params[ind] = cnts
			ind++
			params[ind] = cnts
			ind++
			params[ind] = cnts
			ind++
			inQuery.WriteString(" AND (f.detail LIKE ? OR f.ethereum_public_address LIKE ? OR f.email LIKE ?)")
		}
		params[ind] = sQuery.index
		ind++
		params[ind] = sQuery.limit
		ind++
		inQuery.WriteString(" ORDER BY f.last_updated DESC, f.date_created DESC, f.email ASC LIMIT ?,?")
		r, err = me.db.Query(inQuery.String(), params...)
	} else {
		r, err = me.db.Query("SELECT id, photo,`accessible`, email,ethereum_public_address, role, last_updated, date_created, profile_json FROM user f WHERE (f.detail LIKE ? OR f.ethereum_public_address LIKE ? OR f.email LIKE ?) ORDER BY f.last_updated DESC, f.date_created DESC, f.email ASC LIMIT ?,?", cnts, cnts, cnts, sQuery.index, sQuery.limit)
	}
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var (
		id             sql.NullString
		photo          sql.NullString
		accessible     sql.NullString
		email          sql.NullString
		ethereumPubKey sql.NullString
		role           sql.NullInt64
		usrForm        sql.NullString
		lastUpdated    mysql.NullTime
		dateCreated    mysql.NullTime
	)

	resusts := make([]*model.User, 0)
	for r.Next() {
		err = r.Scan(&id, &photo, &accessible, &email, &ethereumPubKey, &role, &lastUpdated, &dateCreated, &usrForm)
		if err != nil {
			return nil, err
		}
		item := &model.User{}
		if id.Valid {
			item.ID = id.String
		}
		if photo.Valid {
			item.PhotoPath = photo.String
		}
		var f *os.File
		f, err = me.readPhoto(item)
		if err == nil {
			item.Photo, err = me.tinyUserIconBase64(f)
			if err != nil {
				log.Println("read photo error2", err)
			}
		}

		if email.Valid {
			item.Name = email.String
		}
		if ethereumPubKey.Valid {
			item.EthereumAddr = ethereumPubKey.String
		}
		if role.Valid {
			item.Role = model.Role(role.Int64)
		}
		if lastUpdated.Valid {
			item.Updated = lastUpdated.Time
		}
		if dateCreated.Valid {
			item.Created = dateCreated.Time
		}
		if !sQuery.metaOnly && usrForm.Valid {
			var userForm map[string]interface{}
			err = json.Unmarshal([]byte(usrForm.String), &userForm)
			if err == nil {
				item.Data = userForm
			}
		}
		if !sQuery.metaOnly && item.PhotoPath != "" {
			log.Println("item.PhotoPath", item.PhotoPath)
			item.PhotoPath = filepath.Join(me.baseFilePath, item.PhotoPath)
		}

		resusts = append(resusts, item)
	}
	return resusts, nil
}

func (me *MysqlUser) GetProfilePhoto(auth model.Authorization, id string, writer io.Writer) (int64, error) {
	u, err := me.get(auth, id)
	if err != nil {
		return 0, os.ErrNotExist
	}
	var tmplFile *os.File
	tmplFile, err = me.readPhoto(u)
	if err != nil {
		if tmplFile != nil {
			tmplFile.Close()
		}
		return 0, os.ErrNotExist
	}
	defer tmplFile.Close()
	return io.Copy(writer, tmplFile)
}

func (me *MysqlUser) readPhoto(u *model.User) (*os.File, error) {
	if u.PhotoPath == "" {
		return nil, os.ErrNotExist
	}
	return os.OpenFile(filepath.Join(me.baseFilePath, u.PhotoPath), os.O_RDONLY, 0600)
}

func (me *MysqlUser) tinyUserIconBase64(reader *os.File) (string, error) {
	var (
		orgImg image.Image
		err    error
	)
	prefix := "data:image/jpeg;base64,"
	orgImg, err = jpeg.Decode(reader)
	if err != nil {
		reader.Seek(0, 0)
		orgImg, err = png.Decode(reader)
		prefix = "data:image/png;base64,"
	}
	if err != nil {
		return "", err
	}
	newImage := imaging.Resize(orgImg, 44, 0, imaging.Cosine)
	w := &bytes.Buffer{}
	//// Encode uses a Writer, use a Buffer if you need the raw []byte
	err = jpeg.Encode(w, newImage, nil)
	return prefix + base64.StdEncoding.EncodeToString(w.Bytes()), nil
}

func (me *MysqlUser) PutProfilePhoto(auth model.Authorization, id string, reader io.Reader) (int64, error) {
	if id == "" {
		return 0, os.ErrInvalid
	}
	u, err := me.get(auth, id)
	if err != nil {
		return 0, err
	}
	if u.PhotoPath == "" {
		u2 := uuid.NewV4()
		u.PhotoPath = u2.String()
		u.Updated = time.Now()
		_, err = me.db.Exec("UPDATE user f SET f.photo = ?, f.last_updated = ? WHERE f.id = ?",
			u.PhotoPath,
			u.Updated,
			u.ID)
		if err != nil {
			return 0, err
		}
	}
	var tmplFile *os.File
	tmplFile, err = os.OpenFile(filepath.Join(me.baseFilePath, u.PhotoPath), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if os.IsExist(err) {
		err = nil
	}
	if err != nil {
		if tmplFile != nil {
			tmplFile.Close()
		}
		return 0, err
	}
	defer tmplFile.Close()
	return io.Copy(tmplFile, reader)
}

func (me *MysqlUser) Get(auth model.Authorization, id string) (*model.User, error) {
	return me.get(auth, id)
}

func (me *MysqlUser) get(auth model.Authorization, id string) (*model.User, error) {
	var r *sql.Rows
	var r2 *sql.Rows
	var stmt *sql.Stmt
	var err error
	if id == "" {
		return nil, fmt.Errorf("id cant be empty")
	}
	stmt, err = me.db.Prepare("SELECT id, photo,`accessible`, email, detail, ethereum_public_address, role, last_updated, date_created, profile_json FROM `user` f WHERE f.id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	r, err = stmt.Query(id)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	r2, err = me.db.Query("SELECT profile_src_json FROM `system_settings` f WHERE f.id = 1")
	if err != nil {
		return nil, err
	}
	defer r2.Close()
	var (
		sid            sql.NullString
		photo          sql.NullString
		accessible     sql.NullString
		email          sql.NullString
		dbdetail       sql.NullString
		ethereumPubKey sql.NullString
		role           sql.NullInt64
		usrFormSrc     sql.NullString
		usrForm        sql.NullString
		lastUpdated    mysql.NullTime
		dateCreated    mysql.NullTime
	)
	item := &model.User{}
	m := map[string]interface{}{"userSrc": 0, "user": 0}
	item.Data = m
	for r2.Next() {
		err = r2.Scan(&usrFormSrc)
		if err != nil {
			return nil, err
		}
		if usrFormSrc.Valid {
			var userForm map[string]interface{}
			jsnBts := []byte(usrFormSrc.String)
			err = json.Unmarshal(jsnBts, &userForm)
			if err == nil {
				m["userSrc"] = userForm
			}
		}
	}
	for r.Next() {
		err = r.Scan(&sid, &photo, &accessible, &email, &dbdetail, &ethereumPubKey, &role, &lastUpdated, &dateCreated, &usrForm)
		if err != nil {
			return nil, err
		}
		if sid.Valid {
			item.ID = sid.String
		}
		if photo.Valid {
			item.PhotoPath = photo.String
		}

		var f *os.File
		f, err = me.readPhoto(item)
		if err == nil {
			item.Photo, err = me.tinyUserIconBase64(f)
		}

		if email.Valid {
			item.Name = email.String
		}
		if dbdetail.Valid {
			item.Detail = dbdetail.String
		}
		if ethereumPubKey.Valid {
			item.EthereumAddr = ethereumPubKey.String
		}
		if role.Valid {
			item.Role = model.Role(role.Int64)
		}
		if lastUpdated.Valid {
			item.Updated = lastUpdated.Time
		}
		if dateCreated.Valid {
			item.Created = dateCreated.Time
		}

		if usrForm.Valid {
			var userForm map[string]interface{}
			jsnBts := []byte(usrForm.String)
			err = json.Unmarshal(jsnBts, &userForm)
			if err == nil {
				m["user"] = userForm
				if role.Valid {
					userForm["role"] = int(role.Int64)
				} else {
					userForm["role"] = 0
				}
			} else {
				userForm = map[string]interface{}{}
				m["user"] = userForm
				if role.Valid {
					userForm["role"] = int(role.Int64)
				} else {
					userForm["role"] = 0
				}
			}
		}
		return item, nil
	}
	return nil, nil
}

//Put new user or update an existing one
func (me *MysqlUser) Put(auth model.Authorization, item *model.User) error {
	if item == nil {
		return errors.New("item cannot be nil")
	}
	var err error
	if item.ID == "" { //insert new item
		if item.Name == "" {
			return errors.New("name cannot be empty")
		}
		var ab sql.Result
		item.Created = time.Now()
		item.Updated = item.Created
		ab, err = me.db.Exec("INSERT INTO `user` (`ethereum_public_address`,`role`,`version`, `account_expired`, `account_locked`, `email`, `detail`, `enabled`, `password`, `password_expired`, `credit`, `guest`, `date_created`, `last_updated`, `photo`) VALUES (?,?,1,0,1,?,?,0,'',0,0,0,?,?,'');", item.EthereumAddr, item.Role, item.Name, item.Detail, item.Created, item.Updated)
		if err != nil {
			return err
		}
		mid, err := ab.LastInsertId()
		if err != nil {
			return err
		}
		item.ID = strconv.FormatInt(mid, 10)
	}
	var usrDataStr string
	var usrSrcStr string
	if item.Data != nil {
		if usrData, ok := item.Data.(map[string]interface{}); ok {
			if userm, ok := usrData["user"]; ok {
				if userMap, ok := userm.(map[string]interface{}); ok {
					var jsonData []byte
					jsonData, err = json.Marshal(userMap)
					if err != nil {
						usrDataStr = ""
					} else {
						usrDataStr = string(jsonData)
					}
				}
			}
			if usersrc, ok := usrData["userSrc"]; ok {
				if userSrcMap, ok := usersrc.(map[string]interface{}); ok {
					var jsonData []byte
					jsonData, err = json.Marshal(userSrcMap)
					if err != nil {
						usrSrcStr = ""
					} else {
						usrSrcStr = string(jsonData)
					}
				}
			}
		}
	}

	item.Updated = time.Now()
	result, err := me.db.Exec("UPDATE user f SET f.role = ?, f.detail = ?, f.ethereum_public_address = ?, f.enabled = ?, f.email = ?, f.profile_json = ?, f.last_updated = ? WHERE f.id = ?",
		item.Role,
		item.Detail,
		item.EthereumAddr,
		item.Active,
		item.Name,
		usrDataStr,
		item.Updated,
		item.ID)
	if err != nil {
		return err
	}

	if usrSrcStr != "" {
		_, err := me.db.Exec("UPDATE `system_settings` f SET f.profile_src_json = ? WHERE f.id = 1", usrSrcStr)
		if err != nil {
			return err
		}
	}
	_, err = result.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}

func (me *MysqlUser) Close() error {
	return nil
}
