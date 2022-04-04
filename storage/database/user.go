package database

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/ProxeusApp/proxeus-core/storage"

	"github.com/ProxeusApp/proxeus-core/storage/database/db"
	"github.com/ProxeusApp/proxeus-core/sys/model"

	"github.com/asdine/storm/q"
	"github.com/disintegration/imaging"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserDB struct {
	db           db.DB
	baseFilePath string
	fileDB       storage.FilesIF
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

// NewUserDB returns a handle to the user database, containing the user object, incl. their credentials for login and API access
func NewUserDB(c DBConfig, fileDB storage.FilesIF) (*UserDB, error) {
	baseDir := filepath.Join(c.Dir, "user")
	db, err := db.OpenDatabase(c.Engine, c.URI, filepath.Join(baseDir, "users"))
	if err != nil {
		return nil, err
	}
	udb := &UserDB{db: db, baseFilePath: filepath.Join(baseDir, "assets")}
	udb.fileDB = fileDB

	example := &model.User{}
	udb.db.Init(example)

	udb.db.Init(userHeavyDataBucket)
	udb.db.Init(userApiKeyBucket)
	udb.db.Init(userApiKeysBucket)
	udb.db.Init(passwordBucket)

	err = udb.db.Set(userVersion, userVersion, example.GetVersion())
	if err != nil {
		return nil, err
	}
	return udb, nil
}

// GetBaseFilePath returns the base file path
func (me *UserDB) GetBaseFilePath() string {
	return me.baseFilePath
}

// Login tries to authenticate a user with the supplied credentials and returns the user object or an error
func (me *UserDB) Login(name, pw string) (*model.User, error) {
	//TODO refactor login into two methods and move the validation of the signature in here
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

// APIKey tries to authenticate the user with the supplied API key and returns the user object or an error
func (me *UserDB) GetByApiKey(key string, userID string) (*model.User, error) {
	if len(key) != model.ApiKeyLength {
		return nil, model.ErrAuthorityInvalid
	}

	// Look up in the key database
	if userID != "" {

		// Look up user by provided userid
		var user model.User
		err := me.db.One("ID", userID, &user)
		if err != nil {
			return nil, model.ErrAuthorityMissing
		}
		// Check all users keys
		for _, a := range user.ApiKeys {
			if model.MatchesApiKey(a.Key, key) {
				return &user, nil
			}
		}

		return nil, model.ErrAuthorityNotFound
	}

	// Look up across all users and keys
	var allUsers []*model.User
	err := me.db.All(&allUsers)
	if err != nil {
		return nil, err
	}
	for _, user := range allUsers {
		for _, a := range user.ApiKeys {
			if model.MatchesApiKey(a.Key, key) {
				return user, nil
			}
		}
	}

	return nil, model.ErrAuthorityNotFound
}

// CreateApiKey saves and returns a newly created random api key for a user
func (me *UserDB) CreateApiKey(auth model.Auth, userId, apiKeyName string) (string, error) {
	userItem, err := me.Get(auth, userId)
	if err != nil {
		return "", err
	}
	if auth.UserID() != userItem.ID {
		return "", model.ErrAuthorityMissing
	}

	// generate a key and store in profile
	apiKey, err := userItem.SetApiKey(apiKeyName)
	if err != nil {
		return "", err
	}

	// encrypt the key
	readableKey := apiKey.Key
	apiKey.HideKey()

	// store the updated user profile
	err = me.Put(auth, userItem)
	if err != nil {
		return "", err
	}

	// return initially readable key
	return readableKey, nil
}

// DeleteApiKey removes an existing API key by name or value
func (me *UserDB) DeleteApiKey(auth model.Auth, userId, hiddenApiKey string) error {
	userItem, err := me.Get(auth, userId)
	if err != nil {
		return err
	}
	if auth.UserID() != userItem.ID && auth.AccessRights() != model.ROOT {
		return model.ErrAuthorityMissing
	}
	targetIndex := -1
	var anApiKeyValue string = ""
	for i, a := range userItem.ApiKeys {
		if a.Key == hiddenApiKey || a.Name == hiddenApiKey {
			targetIndex = i
			anApiKeyValue = a.Key
			break
		}
	}
	if targetIndex == -1 {
		return errors.New("API key not found")
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
		return errors.New("existing api key not found")
	}
	var apiKey string
	targetIndex = -1
	for i, a := range existingApiKeys {
		if model.MatchesApiKey(anApiKeyValue, a.Key) {
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

// Count returns the user count
func (me *UserDB) Count() (int, error) {
	return me.db.Count(&model.User{})
}

// List returns references to all the user object matching the supplied filter criteria
func (me *UserDB) List(auth model.Auth, contains string, options storage.Options) ([]*model.User, error) {
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
				q.And(
					q.Or(
						q.Eq("WantToBeFound", true),
						q.Eq("ID", auth.UserID()),
					),
					q.Or(
						q.Re("Email", contains),
						q.Re("Name", contains),
						q.Re("Detail", contains),
						q.Re("EthereumAddr", contains),
					),
				),
			)
		}
	} else {
		if !auth.AccessRights().IsGrantedForUserModifications() {
			matchers = append(matchers,
				q.Or(
					q.Eq("WantToBeFound", true),
					q.Eq("ID", auth.UserID()),
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

				//error handling not needed
				_ = me.setTinyUserIconBase64(item)
			}
		}
	}
	return items, nil
}

// Get return a specific user object by machting its id
func (me *UserDB) Get(auth model.Auth, id string) (*model.User, error) {
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

// GetByBCAddress return a specific user object by matching the ethereum address
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

// GetByEmail return a specific user object by matching the email address
func (me *UserDB) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := me.db.One("Email", email, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateEmail sets a new email address for a specific user id
func (me *UserDB) UpdateEmail(id, email string) error {
	tx, err := me.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var u model.User
	err = tx.Select(q.Eq("ID", id)).First(&u)
	if err != nil {
		return err
	}
	u.Email = email
	err = tx.Update(&u)
	if err != nil {
		return err
	}
	return tx.Commit()
}

// PutPw sets a new password for a specific user id
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

// Put saves a user object into the database
func (me *UserDB) Put(auth model.Auth, item *model.User) error {
	return me.put(auth, item, true)
}

func (me *UserDB) put(auth model.Auth, item *model.User, updated bool) error {
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
		item.ID = uuid.NewV4().String()
		if !auth.AccessRights().IsGrantedFor(item.Role) {
			return model.ErrAuthorityMissing
		}
		item.Created = now
		item.Updated = now

		return me.save(item, tx)
	} else {
		existing, err := me.Get(auth, item.ID)
		if db.NotFound(err) {
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

func (me *UserDB) save(u *model.User, tx db.DB) error {
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

func (me *UserDB) updateApiKeys(u *model.User, tx db.DB) error {
	newKeys := make([]model.ApiKey, 0)
	var existingKeys []model.ApiKey
	_ = tx.Get(userApiKeysBucket, u.ID, &existingKeys)

	for _, a := range u.ApiKeys {
		if me.keyIsNew(existingKeys, a) {
			newKeys = append(newKeys, *a)
			err := tx.Set(userApiKeyBucket, a.Key, u.ID)
			if err != nil {
				return err
			}
		}
	}

	if len(newKeys) > 0 {
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

func (me *UserDB) keyIsNew(existingKeys []model.ApiKey, apiKey *model.ApiKey) bool {
	exists := false
	var tmp model.ApiKey
	tmp.Key = apiKey.Key
	for _, k := range existingKeys {
		if k.Key == tmp.Key {
			exists = true
			break
		}
	}
	return !exists
}

func (me *UserDB) setTinyUserIconBase64(item *model.User) error {
	if item.PhotoPath == "" {
		return nil
	}
	var buf bytes.Buffer
	err := me.fileDB.Read(me.fullPhotoPath(item), &buf)
	if err != nil {
		return err
	}
	b := buf.Bytes()

	prefix := "data:image/jpeg;base64,"
	orgImg, err := jpeg.Decode(&buf)
	if err != nil {
		orgImg, err = png.Decode(bytes.NewBuffer(b))
		prefix = "data:image/png;base64,"
	}
	if err != nil {
		return err
	}
	newImage := imaging.Resize(orgImg, 44, 0, imaging.Cosine)
	w := &bytes.Buffer{}
	err = jpeg.Encode(w, newImage, nil)
	data := prefix + base64.StdEncoding.EncodeToString(w.Bytes())
	item.Photo = data
	return err
}

// GetProfilePhoto returns a users photo
func (me *UserDB) GetProfilePhoto(auth model.Auth, id string, writer io.Writer) error {
	u, err := me.Get(auth, id)
	if err != nil {
		return os.ErrNotExist
	}
	u.CheckIfAuthIsAllowedToReadPersonalData(auth)
	return me.fileDB.Read(me.fullPhotoPath(u), writer)
}

func (me *UserDB) fullPhotoPath(u *model.User) string {
	return filepath.Join(me.GetBaseFilePath(), u.PhotoPath)
}

// PutProfilePhoto sets a new photo for a specific user
func (me *UserDB) PutProfilePhoto(auth model.Auth, id string, reader io.Reader) error {
	if id == "" {
		return os.ErrInvalid
	}
	u, err := me.Get(auth, id)
	if err != nil {
		return err
	}
	if u.ID != auth.UserID() && !auth.AccessRights().IsGrantedForUserModifications() {
		return model.ErrAuthorityMissing
	}
	u.Updated = time.Now()
	if u.PhotoPath == "" {
		u.PhotoPath = uuid.NewV4().String()
		err = me.db.Save(u)
		if err != nil {
			return err
		}
	}
	return me.fileDB.Write(me.fullPhotoPath(u), reader)
}

// Close closes the database
func (me *UserDB) Close() error {
	return me.db.Close()
}
