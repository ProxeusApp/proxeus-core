package storm

import (
	"bytes"
	"encoding/base64"
	"regexp"

	"github.com/ProxeusApp/proxeus-core/sys/db"

	//"encoding/json"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/msgpack"
	"github.com/asdine/storm/q"
	"github.com/disintegration/imaging"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/ProxeusApp/proxeus-core/sys/model"
)

//TODO(ave) how come we have even private methods in an interface... Doesn't look like the correct approach for an interface.
type UserDBInterface interface {
	GetDB() *storm.DB
	GetBaseFilePath() string
	Login(name, pw string) (*model.User, error)
	Count() (int, error)
	List(auth model.Authorization, contains string, options map[string]interface{}) ([]*model.User, error)
	Get(auth model.Authorization, id string) (*model.User, error)
	GetByBCAddress(bcAddress string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	UpdateEmail(id, email string) error
	Put(auth model.Authorization, item *model.User) error
	PutPw(id, pass string) error
	ImportUser(auth model.Authorization, item *model.User) error
	SetTinyUserIconBase64(item *model.User) error
	TinyUserIconBase64(reader *os.File) (string, error)
	GetProfilePhoto(auth model.Authorization, id string, writer io.Writer) (n int64, err error)
	ReadPhoto(u *model.User) (*os.File, error)
	PutProfilePhoto(auth model.Authorization, id string, reader io.Reader) (written int64, err error)
	Import(imex *Imex) error
	Export(imex *Imex, id ...string) error
	CpProfilePhoto(imex *Imex, from UserDBInterface, to UserDBInterface, item *model.User) (err error)
	APIKey(key string) (*model.User, error)
	CreateApiKey(auth model.Authorization, userId, apiKeyName string) (string, error)
	DeleteApiKey(auth model.Authorization, userId, hiddenApiKey string) error
	Close() error
}

type UserDB struct {
	db           *storm.DB
	baseFilePath string
}

//userHeavyDataBucket helps us to load the data of the model.User entity when it is requested by metaOnly = false
const userHeavyDataBucket = "user_data"

const userApiKeyBucket = "user_api_key"   //api key -> user id
const userApiKeysBucket = "user_api_keys" //user id -> [api key1, api key2...]

//userVersion helps us to keep track of the structure version between persistence and memory
const userVersion = "user_version"

//passwordBucket helps us to keep it away from the actual structure
//it is only needed for login and password reset
const passwordBucket = "pw_bucket"

func NewUserDB(dir string) (*UserDB, error) {
	var err error
	var msgpackDb *storm.DB
	baseDir := filepath.Join(dir, "user")
	err = ensureDir(baseDir)
	if err != nil {
		return nil, err
	}
	assetDir := filepath.Join(baseDir, "assets")
	err = ensureDir(assetDir)
	if err != nil {
		return nil, err
	}
	msgpackDb, err = storm.Open(filepath.Join(baseDir, "users"), storm.Codec(msgpack.Codec))
	if err != nil {
		return nil, err
	}
	udb := &UserDB{db: msgpackDb}
	udb.baseFilePath = assetDir

	example := &model.User{}
	udb.db.Init(example)
	var version int

	verr := udb.db.Get(userVersion, userVersion, &version)
	if verr == nil && version != example.GetVersion() {
		log.Println("upgrade db", version, "mem", example.GetVersion())
	}
	err = udb.db.Set(userVersion, userVersion, example.GetVersion())
	if err != nil {
		return nil, err
	}
	return udb, nil
}

func (me *UserDB) GetDB() *storm.DB {
	return me.db
}

func (me *UserDB) GetBaseFilePath() string {
	return me.baseFilePath
}

//TODO refactor login into two methods and move the validation of the signature in here
func (me *UserDB) Login(name, pw string) (*model.User, error) {
	if name == "" || pw == "" {
		return nil, os.ErrInvalid
	}
	var usr model.User
	err := me.db.One("Email", name, &usr)
	if err != nil {
		return nil, err
	}
	var pass string
	err = me.db.Get(passwordBucket, usr.ID, &pass)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(pass), []byte(pw))
	if err != nil {
		return nil, os.ErrInvalid
	}
	return &usr, nil
}

func (me *UserDB) APIKey(key string) (*model.User, error) {
	if len(key) != model.ApiKeyLength {
		return nil, model.ErrAuthorityMissing
	}
	var userID string
	err := me.db.Get(userApiKeyBucket, key, &userID)
	if err != nil {
		return nil, model.ErrAuthorityMissing
	}
	var user model.User
	err = me.db.One("ID", userID, &user)
	if err != nil {
		return nil, model.ErrAuthorityMissing
	}
	return &user, nil
}

func (me *UserDB) CreateApiKey(auth model.Authorization, userId, apiKeyName string) (string, error) {
	userItem, err := me.Get(auth, userId)
	if err != nil {
		return "", err
	}
	if auth.UserID() != userItem.ID {
		return "", model.ErrAuthorityMissing
	}
	apiKey, err := userItem.NewApiKey(apiKeyName)
	if err != nil {
		return "", err
	}
	initiallyReadableApiKey := apiKey.Key

	//store the new api key
	err = me.Put(auth, userItem)
	if err != nil {
		return "", err
	}

	return initiallyReadableApiKey, nil
}

func (me *UserDB) DeleteApiKey(auth model.Authorization, userId, hiddenApiKey string) error {
	userItem, err := me.Get(auth, userId)
	if err != nil {
		return err
	}
	if auth.UserID() != userItem.ID && auth.AccessRights() != model.ROOT {
		return model.ErrAuthorityMissing
	}
	targetIndex := -1
	for i, a := range userItem.ApiKeys {
		if a.Key == hiddenApiKey {
			targetIndex = i
			break
		}
	}
	if targetIndex == -1 {
		return db.ErrNotFound
	}
	// replace the target element with the last one
	userItem.ApiKeys[targetIndex] = userItem.ApiKeys[len(userItem.ApiKeys)-1]
	// discard the last element
	userItem.ApiKeys = userItem.ApiKeys[:len(userItem.ApiKeys)-1]

	tx, err := me.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.Save(userItem)
	if err != nil {
		return err
	}

	var existingApiKeys []*model.ApiKey
	err = tx.Get(userApiKeysBucket, userId, &existingApiKeys)
	if err != nil {
		return err
	}

	if len(existingApiKeys) == 0 {
		return db.ErrNotFound
	}
	var apiKey string
	targetIndex = -1
	for i, a := range existingApiKeys {
		if model.MatchesApiKey(hiddenApiKey, a.Key) {
			targetIndex = i
			apiKey = a.Key
		}
	}
	if targetIndex > -1 {
		// replace the target element with the last one
		existingApiKeys[targetIndex] = existingApiKeys[len(existingApiKeys)-1]
		// discard the last element
		existingApiKeys = existingApiKeys[:len(existingApiKeys)-1]
		err = tx.Set(userApiKeysBucket, userId, existingApiKeys)
		if err != nil {
			return err
		}
		err = tx.Delete(userApiKeyBucket, apiKey)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (me *UserDB) Count() (int, error) {
	return me.db.Count(&model.User{})
}

func (me *UserDB) List(auth model.Authorization, contains string, options map[string]interface{}) ([]*model.User, error) {
	contains = containsCaseInsensitiveReg(contains)
	params := makeSimpleQuery(options)
	items := make([]*model.User, 0)
	tx, err := me.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	matchers := make([]q.Matcher, 0)
	if contains != "" {
		if auth.AccessRights().IsGrantedForUserModifications() {
			matchers = append(matchers,
				q.Or(
					q.Re("Email", contains),
					q.Re("Name", contains),
					q.Re("Detail", contains),
					q.Re("EthereumAddr", contains),
				),
			)
		} else {
			matchers = append(matchers,
				q.Or(
					q.And(
						q.Eq("WantToBeFound", true),
						q.Or(
							q.Re("Email", contains),
							q.Re("Name", contains),
							q.Re("Detail", contains),
							q.Re("EthereumAddr", contains),
						),
					),
					q.And(
						q.Eq("WantToBeFound", false),
						q.Re("EthereumAddr", contains),
					),
				),
			)
		}
	}
	if len(params.exclude) > 0 {
		matchers = append(matchers,
			q.And(
				q.Not(q.In("ID", params.exclude)),
			),
		)
	}
	if len(params.include) > 0 {
		matchers = append(matchers,
			q.And(
				q.In("ID", params.include),
			),
		)
	}
	err = tx.Select(matchers...).
		Limit(params.limit).
		Skip(params.index).
		OrderBy("Updated").
		Reverse().
		Find(&items)
	if err != nil {
		return nil, err
	}
	if !params.metaOnly {
		for _, item := range items {
			if item.CheckIfAuthIsAllowedToReadPersonalData(auth) {
				//error handling not needed
				_ = tx.Get(userHeavyDataBucket, item.ID, &item.Data)

				//	//error handling not needed
				_ = me.SetTinyUserIconBase64(item)
			}
		}
	}
	return items, nil
}

func (me *UserDB) Get(auth model.Authorization, id string) (*model.User, error) {
	var err error
	var user model.User
	err = me.db.One("ID", id, &user)
	if err != nil {
		return nil, err
	}
	userItem := &user
	if userItem.CheckIfAuthIsAllowedToReadPersonalData(auth) {
		_ = me.db.Get(userHeavyDataBucket, userItem.ID, &userItem.Data)
	}
	return userItem, nil
}

func (me *UserDB) GetByBCAddress(bcAddress string) (*model.User, error) {
	var user model.User
	//case insensitive as the address is stored as a string
	bcAddress = "(?i)" + regexp.QuoteMeta(bcAddress)
	err := me.db.Select(q.Re("EthereumAddr", bcAddress)).Limit(1).First(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (me *UserDB) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := me.db.One("Email", email, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (me *UserDB) UpdateEmail(id, email string) error {
	return me.db.UpdateField(&model.User{ID: id}, "Email", email)
}

func (me *UserDB) PutPw(id, pass string) error {
	pw := []byte(pass)
	cost, err := bcrypt.Cost(pw)
	if err != nil || cost <= 0 {
		pw, err = bcrypt.GenerateFromPassword(pw, bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		pass = string(pw)
	}
	return me.db.Set(passwordBucket, id, &pass)
}

func (me *UserDB) Put(auth model.Authorization, item *model.User) error {
	return me.put(auth, item, true)
}

func (me *UserDB) ImportUser(auth model.Authorization, item *model.User) error {
	return me.put(auth, item, false)
}

func (me *UserDB) put(auth model.Authorization, item *model.User, updated bool) error {
	if item == nil {
		return os.ErrInvalid
	}
	tx, err := me.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	now := time.Now()
	if item.ID == "" {
		u2 := uuid.NewV4()
		item.ID = u2.String()
		if !auth.AccessRights().IsGrantedFor(item.Role) {
			return model.ErrAuthorityMissing
		}
		item.Created = now
		item.Updated = now

		return me.save(item, tx)
	} else {
		existing, err := me.Get(auth, item.ID)
		if err == storm.ErrNotFound {
			err = nil
			if !auth.AccessRights().IsGrantedFor(item.Role) {
				return model.ErrAuthorityMissing
			}
			return me.save(item, tx)
		}
		if err != nil {
			return err
		}
		if auth.AccessRights().IsGrantedFor(existing.Role) &&
			auth.AccessRights().IsGrantedFor(item.Role) {
			if updated {
				item.Updated = now
			}
			return me.save(item, tx)
		} else {
			return model.ErrAuthorityMissing
		}
	}
}

func (me *UserDB) save(u *model.User, tx storm.Node) error {
	err := me.updateApiKeys(u, tx)
	if err != nil {
		return err
	}
	err = tx.Save(u)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (me *UserDB) updateApiKeys(u *model.User, tx storm.Node) error {
	newKeys := make([]model.ApiKey, 0)
	for _, a := range u.ApiKeys {
		if a.IsNew() {
			newKeys = append(newKeys, *a)
			err := tx.Set(userApiKeyBucket, a.Key, u.ID)
			if err != nil {
				return err
			}
			a.HideKey()
		}
	}
	if len(newKeys) > 0 {
		var existingKeys []model.ApiKey
		_ = tx.Get(userApiKeysBucket, u.ID, &existingKeys)
		if len(existingKeys) > 0 {
			newKeys = append(newKeys, existingKeys...)
		}
		err := tx.Set(userApiKeysBucket, u.ID, newKeys)
		if err != nil {
			return err
		}
	}
	return nil
}

func (me *UserDB) SetTinyUserIconBase64(item *model.User) error {
	f, err := me.ReadPhoto(item)
	if err == nil {
		item.Photo, err = me.TinyUserIconBase64(f)
	}
	return err
}

func (me *UserDB) TinyUserIconBase64(reader *os.File) (string, error) {
	//b := &bytes.Buffer{}
	//io.Copy(b, reader)
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

func (me *UserDB) GetProfilePhoto(auth model.Authorization, id string, writer io.Writer) (n int64, err error) {
	u, err := me.Get(auth, id)
	if err != nil {
		return 0, os.ErrNotExist
	}
	u.CheckIfAuthIsAllowedToReadPersonalData(auth)
	var tmplFile *os.File
	tmplFile, err = me.ReadPhoto(u)
	if err != nil {
		if tmplFile != nil {
			tmplFile.Close()
		}
		return 0, os.ErrNotExist
	}
	defer tmplFile.Close()
	return io.Copy(writer, tmplFile)
}

func (me *UserDB) ReadPhoto(u *model.User) (*os.File, error) {
	if u.PhotoPath == "" {
		return nil, os.ErrNotExist
	}
	return os.OpenFile(filepath.Join(me.GetBaseFilePath(), u.PhotoPath), os.O_RDONLY, 0600)
}

func (me *UserDB) PutProfilePhoto(auth model.Authorization, id string, reader io.Reader) (written int64, err error) {
	if id == "" {
		return 0, os.ErrInvalid
	}
	u, err := me.Get(auth, id)
	if err != nil {
		return 0, err
	}
	if u.ID != auth.UserID() && !auth.AccessRights().IsGrantedForUserModifications() {
		return 0, model.ErrAuthorityMissing
	}
	u.Updated = time.Now()
	if u.PhotoPath == "" {
		u2 := uuid.NewV4()
		u.PhotoPath = u2.String()
		err = me.db.Save(u)
		if err != nil {
			return 0, err
		}
	}
	var tmplFile *os.File
	tmplFile, err = os.OpenFile(filepath.Join(me.GetBaseFilePath(), u.PhotoPath), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
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

func (me *UserDB) Import(imex *Imex) error {
	err := me.init(imex)
	if err != nil {
		return err
	}
	for i := 0; true; i++ {
		items, err := imex.db.User.List(imex.auth, "", map[string]interface{}{"index": i, "limit": 1000, "metaOnly": false})
		if err == nil && len(items) > 0 {
			for _, item := range items {
				var existingItem *model.User
				existingItem, err = imex.sysDB.User.Get(imex.auth, item.ID)
				if err == nil && imex.skipExistingOnImport {
					continue
				}

				if existingItem == nil {
					//treat email as an ID and update all references on this import package if user was located on the target system
					if item.Email != "" {
						existingItem, _ = imex.sysDB.User.GetByEmail(item.Email)
						if existingItem != nil {
							//provide user id correction map for entities with permissions item.ID -> existingItem.ID
							imex.locatedSameUserWithDifferentID[item.ID] = existingItem.ID
							if imex.skipExistingOnImport {
								continue
							} else {
								item.ID = existingItem.ID
								updateEmptyFields(item, existingItem)
							}
						}
					}
					//treat Ethereum address as an ID and update all references on this import package if user was located on the target system
					if item.EthereumAddr != "" {
						existingItem, _ = imex.sysDB.User.GetByBCAddress(item.EthereumAddr)
						if existingItem != nil {
							//provide user id correction map for entities with permissions item.ID -> existingItem.ID
							imex.locatedSameUserWithDifferentID[item.ID] = existingItem.ID
							if imex.skipExistingOnImport {
								continue
							} else {
								item.ID = existingItem.ID
								updateEmptyFields(item, existingItem)
							}
						}
					}
				}

				if len(item.ApiKeys) > 0 {
					_ = imex.db.User.GetDB().Get(userApiKeysBucket, item.ID, &item.ApiKeys)
				}

				err = imex.sysDB.User.ImportUser(imex.auth, item)
				if err != nil {
					imex.processedEntry(imexUser, item.ID, err)
					continue
				}
				//no permission errors when writing ... we are allowed to set the photo too

				if existingItem != nil && existingItem.PhotoPath != "" && existingItem.PhotoPath != item.PhotoPath {
					//remove old photo of existingItem before the reference is lost
					_ = os.Remove(filepath.Join(imex.sysDB.User.GetBaseFilePath(), existingItem.PhotoPath))
				}

				if item.PhotoPath != "" {
					err = me.CpProfilePhoto(imex, imex.db.User, imex.sysDB.User, item)
					if err != nil {
						continue
					}
				}
				imex.processedEntry(imexUser, item.ID, nil)
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

const imexUser = "User"

func (me *UserDB) init(imex *Imex) error {
	var err error
	if imex.db.User == nil {
		imex.db.User, err = NewUserDB(imex.dir)
	}
	return err
}

func (me *UserDB) Export(imex *Imex, id ...string) error {
	err := me.init(imex)
	if err != nil {
		return err
	}
	specificIds := len(id) > 0
	if len(id) == 1 {
		if imex.isProcessed(imexUser, id[0]) {
			return nil
		}
	}
	if !specificIds {
		imex.exportingAllUsersAnyway = true
	}
	var tx storm.Node
	for i := 0; true; i++ {
		items, err := imex.sysDB.User.List(imex.auth, "", map[string]interface{}{"include": id, "index": i, "limit": 1000, "metaOnly": false})
		if err == nil && len(items) > 0 {
			tx, err = imex.db.User.GetDB().Begin(true)
			if err != nil {
				return err
			}
			for _, item := range items {
				if !imex.isProcessed(imexUser, item.ID) {
					item.Photo = ""
					if len(item.ApiKeys) > 0 {
						var apiKeys []*model.ApiKey
						err = imex.sysDB.User.GetDB().Get(userApiKeysBucket, item.ID, &apiKeys)
						if err != nil || len(apiKeys) == 0 {
							item.ApiKeys = nil
						} else {
							for _, a := range apiKeys {
								err = tx.Set(userApiKeyBucket, a.Key, item.ID)
								if err != nil {
									break
								}
							}
							if err != nil {
								imex.processedEntry(imexUser, item.ID, err)
								continue
							}
							err = tx.Set(userApiKeysBucket, item.ID, apiKeys)
							if err != nil {
								imex.processedEntry(imexUser, item.ID, err)
								continue
							}
						}
					}
					err = tx.Save(item)
					if err != nil {
						imex.processedEntry(imexUser, item.ID, err)
						continue
					}
					if item.PhotoPath != "" {
						err = me.CpProfilePhoto(imex, imex.sysDB.User, imex.db.User, item)
						if err != nil {
							imex.processedEntry(imexUser, item.ID, err)
							continue
						}
					}
					imex.processedEntry(imexUser, item.ID, nil)
				}
			}
			err = tx.Commit()
			if err != nil {
				_ = tx.Rollback()
				return err
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

func (me *UserDB) CpProfilePhoto(imex *Imex, from UserDBInterface, to UserDBInterface, item *model.User) (err error) {
	var readFile *os.File
	readFile, err = os.Open(filepath.Join(from.GetBaseFilePath(), item.PhotoPath))
	if os.IsExist(err) {
		err = nil
	}
	if err != nil {
		if readFile != nil {
			_ = readFile.Close()
		}
		imex.processedEntry(imexUser, item.ID, err)
		return
	}
	var exportFile *os.File
	exportFile, err = os.OpenFile(filepath.Join(to.GetBaseFilePath(), item.PhotoPath), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if os.IsExist(err) {
		err = nil
	}
	if err != nil {
		if exportFile != nil {
			_ = exportFile.Close()
		}
		if readFile != nil {
			_ = readFile.Close()
		}
		imex.processedEntry(imexUser, item.ID, err)
		return
	}
	_, err = io.Copy(exportFile, readFile)
	if err != nil {
		_ = readFile.Close()
		_ = exportFile.Close()
		imex.processedEntry(imexUser, item.ID, err)
		return
	}
	_ = readFile.Close()
	_ = exportFile.Close()
	return
}

func (me *UserDB) Close() error {
	return me.db.Close()
}
