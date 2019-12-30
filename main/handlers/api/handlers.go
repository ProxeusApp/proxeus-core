package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ProxeusApp/proxeus-core/storage/portable"

	"github.com/ProxeusApp/proxeus-core/main/app"
	cfg "github.com/ProxeusApp/proxeus-core/main/config"
	"github.com/ProxeusApp/proxeus-core/main/handlers/blockchain"
	"github.com/ProxeusApp/proxeus-core/main/handlers/payment"
	"github.com/ProxeusApp/proxeus-core/main/helpers"
	"github.com/ProxeusApp/proxeus-core/main/www"
	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/storage/db"
	"github.com/ProxeusApp/proxeus-core/sys"
	"github.com/ProxeusApp/proxeus-core/sys/eio"
	"github.com/ProxeusApp/proxeus-core/sys/email"
	"github.com/ProxeusApp/proxeus-core/sys/file"
	"github.com/ProxeusApp/proxeus-core/sys/model"
	"github.com/ProxeusApp/proxeus-core/sys/session"
	"github.com/ProxeusApp/proxeus-core/sys/utils"
	"github.com/ProxeusApp/proxeus-core/sys/validate"
	workflow2 "github.com/ProxeusApp/proxeus-core/sys/workflow"

	strm "github.com/asdine/storm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

var filenameRegex = regexp.MustCompile(`^[^\s][\p{L}\d.,_\-&: ]{3,}[^\s]$`)

func html(c echo.Context, p string) error {
	bts, err := sys.ReadAllFile(p)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	return c.HTMLBlob(http.StatusOK, bts)
}

func SharedByLinkHTMLHandler(c echo.Context) error {
	log.Println("SharedByLinkHTMLHandler")
	//c.Param("type") TODO different html by type for user data
	return html(c, "app.html")
}

func PublicIndexHTMLHandler(c echo.Context) error {
	return html(c, "frontend.html")
}

func UserBackendHTMLHandler(c echo.Context) error {
	return html(c, "user.html")
}

func AdminIndexHandler(c echo.Context) error {
	return html(c, "app.html")
}

type ImportExportResult struct {
	Filename  string                    `json:"filename"`
	Timestamp time.Time                 `json:"timestamp"`
	Results   portable.ProcessedResults `json:"results"`
}

func GetExport(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	var exportEntities []portable.EntityType
	if strings.ToLower(c.Request().Method) == "get" {
		spl := strings.Split(c.QueryParam("include"), ",")
		for _, s := range spl {
			s = strings.TrimSpace(s)
			entity := portable.StringToEntityType(s)
			if entity != "" {
				exportEntities = append(exportEntities, entity)
			}
		}
	} else {
		_ = c.Bind(&exportEntities)
	}
	if len(exportEntities) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}
	return Export(sess, exportEntities, c)
}

func PostImport(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusBadRequest)
	}
	fileName, _ := url.QueryUnescape(c.Request().Header.Get("File-Name"))
	skipExistingStr := c.QueryParam("skipExisting")
	skipExisting, _ := strconv.ParseBool(skipExistingStr)
	results, err := c.System().Import(c.Request().Body, sess, skipExisting)
	_ = c.Request().Body.Close()
	if err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}
	sess.Put("lastImport", &ImportExportResult{Filename: fileName, Timestamp: time.Now(), Results: results})
	return c.NoContent(http.StatusOK)

}

func ExportUserData(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	var id []string
	if c.QueryParam("id") != "" {
		id = []string{c.QueryParam("id")}
	} else if c.QueryParam("contains") != "" {
		items, _ := c.System().DB.UserData.List(sess, c.QueryParam("contains"), storage.Options{Limit: 1000}, false)
		if len(items) > 0 {
			id = make([]string, len(items))
			for i, item := range items {
				id[i] = item.ID
			}
		}
	}
	return Export(sess, []portable.EntityType{portable.UserData}, c, id...)
}

func ExportSettings(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	return Export(sess, []portable.EntityType{portable.Settings}, c, "Settings")
}

func ExportUser(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	var id []string
	if c.QueryParam("id") != "" {
		id = []string{c.QueryParam("id")}
	} else if c.QueryParam("contains") != "" {
		items, _ := c.System().DB.User.List(sess, c.QueryParam("contains"), storage.Options{Limit: 1000})
		if len(items) > 0 {
			id = make([]string, len(items))
			for i, item := range items {
				id[i] = item.ID
			}
		}
	}
	return Export(sess, []portable.EntityType{portable.User}, c, id...)
}

func Export(sess *session.Session, exportEntities []portable.EntityType, e echo.Context, id ...string) error {
	c := e.(*www.Context)
	if len(exportEntities) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}
	resp := c.Response()
	resp.Header().Set("Content-Disposition", fmt.Sprintf(`%s; filename="proxeus.db"`, "attachment"))
	resp.Committed = true //prevents from-> http: multiple response.WriteHeader calls
	var (
		results portable.ProcessedResults
		err     error
	)
	if len(id) > 0 && len(exportEntities) == 1 {
		results, err = c.System().ExportSingle(resp.Writer, sess, exportEntities[0], id...)
	} else {
		results, err = c.System().Export(resp.Writer, sess, exportEntities)
	}
	if err != nil {
		return c.String(http.StatusUnprocessableEntity, err.Error())
	}
	sess.Put("lastExport", &ImportExportResult{Timestamp: time.Now(), Results: results})
	return c.NoContent(http.StatusOK)
}

func GetExportResults(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	return results("lastExport", sess, c)
}

func GetImportResults(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	return results("lastImport", sess, c)
}

func results(key string, sess *session.Session, c echo.Context) error {
	if _, exists := c.QueryParams()["delete"]; exists {
		del := c.QueryParam("delete")
		if del == "" {
			sess.Delete(key)
		} else {
			var imexResults *ImportExportResult
			_ = sess.Get(key, &imexResults)
			if imexResults != nil && imexResults.Results != nil {
				delete(imexResults.Results, portable.EntityType(del))
				sess.Put(key, imexResults)
			}
		}
	}
	var imexResults *ImportExportResult
	_ = sess.Get(key, &imexResults)
	return c.JSON(http.StatusOK, imexResults)
}

func GetInit(e echo.Context) error {
	c := e.(*www.Context)
	configured, err := c.System().Configured()
	if err != nil && !os.IsNotExist(err) {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	settings := c.System().GetSettings()
	if len(settings.PlatformDomain) == 0 {
		settings.PlatformDomain = e.Request().Host
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"settings": settings, "configured": configured})
}

var root = &model.User{Role: model.ROOT}

func PostInit(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	type usr struct {
		Email    string     `json:"email" validate:"required=false,email=true"`
		Password string     `json:"password" validate:"required=false,matches=^.{6}"`
		Role     model.Role `json:"role"`
	}
	type Init struct {
		Settings *model.Settings `json:"settings"`
		User     *usr            `json:"user"`
	}
	var err error
	yes, _ := c.System().Configured()
	d := &Init{User: &usr{}}
	_ = c.Bind(d)
	if yes {
		d.User = nil
	}
	err = validate.Struct(d)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}
	err = c.System().PutSettings(d.Settings)
	if err != nil {
		fmt.Println("Error during PostInit: ", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	if !yes {
		u := &model.User{Email: d.User.Email, Role: d.User.Role}
		uex, _ := c.System().DB.User.GetByEmail(u.Email)
		if uex == nil {
			err = c.System().DB.User.Put(root, u)
			if err != nil {
				fmt.Println("Error during PostInit: ", err)
				return c.NoContent(http.StatusInternalServerError)
			}
			err = c.System().DB.User.PutPw(u.ID, d.User.Password)
			if err != nil {
				fmt.Println("Error during PostInit: ", err)
				return c.NoContent(http.StatusInternalServerError)
			}
		}
	}
	return c.NoContent(http.StatusOK)
}

func ConfigHandler(version string) echo.HandlerFunc {
	return func(e echo.Context) error {
		c := e.(*www.Context)
		sess := c.Session(false)
		var roles []model.RoleSet
		if sess != nil {
			roles = sess.AccessRights().RolesInRange()
		}
		stngs := c.System().GetSettings()
		return c.JSON(http.StatusOK, map[string]interface{}{
			"roles":                      roles,
			"blockchainNet":              strings.Replace(stngs.BlockchainNet, "mainnet", "main", 1),
			"blockchainProxeusFSAddress": stngs.BlockchainContractAddress,
			"version":                    version,
		})
	}
}

type loginForm struct {
	Signature string
	Name      string `json:"name" form:"name"`
	Email     string `json:"email" form:"email"`
	Password  string `json:"password" form:"password"`
	Address   string `json:"address" form:"address"`
}

func UpdateAddress(e echo.Context) error {
	c := e.(*www.Context)
	loginForm := new(loginForm)
	_ = c.Bind(loginForm)
	sess := c.Session(false)
	u := getUserFromSession(c, sess)
	if sess == nil || u == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	challengeError := map[string]interface{}{"etherPK": []map[string]interface{}{{"msg": "challenge.invalid"}}}
	var challenge string
	_ = sess.Get("challenge", &challenge)
	if challenge == "" {
		return c.JSON(http.StatusUnprocessableEntity, challengeError)
	}
	address, err := blockchain.VerifySignInChallenge(challenge, loginForm.Signature)
	sess.Delete("challenge")
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, challengeError)
	}
	item, err := c.System().DB.User.GetByBCAddress(address)
	if item != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"etherPK": []map[string]interface{}{{"msg": "Please choose another account."}}})
	}
	item, err = c.System().DB.User.Get(sess, sess.UserID())
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	item.EthereumAddr = address
	err = c.System().DB.User.Put(sess, item)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	sess.Put("user", item)
	return c.NoContent(http.StatusOK)
}

//root only feature to switch user by address - useful for permission checks
func SwitchUserHandler(e echo.Context) error {
	c := e.(*www.Context)
	user, err := c.System().DB.User.GetByBCAddress(e.Param("address"))
	if err != nil || user == nil {
		return c.NoContent(http.StatusBadRequest)
	}
	//create a new session only if role, id or name has changed
	sess := c.SessionWithUser(user)
	if sess == nil {
		return c.NoContent(http.StatusBadRequest)
	}
	sess.Put("user", user)
	return c.NoContent(http.StatusOK)
}

func LoginHandler(e echo.Context) (err error) {
	c := e.(*www.Context)
	loginForm := new(loginForm)
	if err := c.Bind(loginForm); err != nil {
		return err
	}
	var user *model.User
	created := false
	// Check if posted data contains a signature.
	// If it does then it's an attempt at login by signature,
	// else basic auth
	if loginForm.Signature != "" {
		sess := c.Session(false)
		if sess == nil {
			return c.NoContent(http.StatusBadRequest)
		}
		var challenge string
		_ = sess.Get("challenge", &challenge)
		if challenge == "" {
			return errors.New("challenge.invalid")
		}
		log.Println("Session challenge", challenge)
		created, user, err = LoginWithWallet(c, challenge, loginForm.Signature)
		if user != nil {
			sess.Delete("challenge")
		}
	} else {
		user, err = c.System().DB.User.Login(loginForm.Email, loginForm.Password)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
	}

	if user != nil {
		//create a new session only if role, id or name has changed
		sess := c.SessionWithUser(user)
		if sess == nil {
			return c.NoContent(http.StatusBadRequest)
		}
		sess.Put("user", user)
	} else {
		return c.NoContent(http.StatusBadRequest)
	}
	var referer []byte
	cookie, err := c.Cookie("R")
	if err == nil {
		referer, err = base64.RawURLEncoding.DecodeString(cookie.Value)
		if err != nil || len(referer) > 0 {
			// reset redirect cookie
			c.SetCookie(&http.Cookie{
				Name:   "R",
				Path:   "/",
				MaxAge: -1,
			})
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"location": redirectAfterLogin(user.Role, string(referer)),
		"created":  created,
	})
}

func LoginWithWallet(c *www.Context, challenge, signature string) (bool, *model.User, error) {
	created := false
	var address string
	var err error
	address, err = blockchain.VerifySignInChallenge(challenge, signature)
	if err != nil {
		return false, nil, err
	}
	var usr *model.User
	usr, err = c.System().DB.User.GetByBCAddress(address)
	if err == db.ErrNotFound || (err != nil && err.Error() == strm.ErrNotFound.Error()) {
		stngs := c.System().GetSettings()
		it := &model.User{
			EthereumAddr: address,
			Role:         stngs.DefaultRole,
		}
		it.Name = "created by ethereum account sign"
		err = c.System().DB.User.Put(root, it)
		if err != nil {
			return false, nil, err
		}
		created = true
		usr, err = c.System().DB.User.GetByBCAddress(address)
		if err == nil {
			copyWorkflows(c, usr)
			if c.System().GetSettings().BlockchainNet == "ropsten" && c.System().GetSettings().AirdropEnabled == "true" {
				go func() {
					defer func() {
						if r := recover(); r != nil {
							log.Println("airdrop recover with err ", r)
						}
					}()
					blockchain.GiveTokens(address)
				}()
			}
		}

	}
	return created, usr, err
}

func GetSessionTokenHandler(e echo.Context) (err error) {
	c := e.(*www.Context)

	username, apiKey := c.BasicAuth()

	if username == "" || apiKey == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	user, err := c.System().DB.User.APIKey(apiKey)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if user == nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if user.Email != username && user.EthereumAddr != username {
		return c.NoContent(http.StatusBadRequest)
	}

	//create a new session only if role, id or name has changed
	sess := c.SessionWithUser(user)
	if sess == nil {
		return c.NoContent(http.StatusBadRequest)
	}
	sess.Put("user", user)

	c.Response().Header().Del("Set-Cookie")

	return c.JSON(http.StatusOK, map[string]string{
		"token": sess.ID(),
	})
}

func DeleteSessionTokenHandler(e echo.Context) (err error) {
	c := e.(*www.Context)
	c.EndSession()
	return c.NoContent(http.StatusOK)
}

type TokenRequest struct {
	Email  string     `json:"email" validate:"email=true,required=true"`
	Token  string     `json:"token"`
	UserID string     `json:"userID"`
	Role   model.Role `json:"role"`
}

func InviteRequest(e echo.Context) (err error) {
	c := e.(*www.Context)
	sess := c.Session(false)
	m := &TokenRequest{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	err = validate.Struct(m)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}
	if !sess.AccessRights().IsGrantedFor(m.Role) {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"role": []map[string]interface{}{{"msg": c.I18n().T("No authority to grant this role.")}}})
	}
	stngs := c.System().GetSettings()
	if m.Role == 0 {
		m.Role = stngs.DefaultRole
	}
	if usr, err := c.System().DB.User.GetByEmail(m.Email); usr == nil {
		resetKey := m.Email + "_register"
		var token *TokenRequest
		err = c.System().Cache.Get(resetKey, &token)
		if err != nil {
			token = m
			u2 := uuid.NewV4()
			token.Token = u2.String()
			err = c.System().EmailSender.Send(&email.Email{
				From:    stngs.EmailFrom,
				To:      []string{m.Email},
				Subject: c.I18n().T("Invitation"),
				Body: fmt.Sprintf(
					"Hi there,\n\nyou have been invited to join Proxeus. If you would like to benefit from the invitation, please proceed by visiting this link:\n%s\n\nProxeus",
					helpers.AbsoluteURL(c, "/register/", token.Token),
				),
			})
			if err != nil {
				return c.String(http.StatusFailedDependency, c.I18n().T("couldn't send the email"))
			}
			err = c.System().Cache.Put(resetKey, token)
			if err != nil {
				return c.NoContent(http.StatusInternalServerError)
			}
			err = c.System().Cache.Put(token.Token, token)
			if err != nil {
				return c.NoContent(http.StatusInternalServerError)
			}
		}
	} else {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"email": []map[string]interface{}{{"msg": c.I18n().T("Account with that email already exists.")}}})
	}
	return c.NoContent(http.StatusOK)
}

func RegisterRequest(e echo.Context) (err error) {
	c := e.(*www.Context)
	m := &TokenRequest{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	err = validate.Struct(m)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}
	stngs := c.System().GetSettings()
	m.Role = stngs.DefaultRole

	if usr, _ := c.System().DB.User.GetByEmail(m.Email); usr != nil {
		// always return ok if provided email was valid
		// otherwise public users can test what email accounts exist
		return c.NoContent(http.StatusOK)
	}

	resetKey := m.Email + "_register"
	var token *TokenRequest

	err = c.System().Cache.Get(resetKey, &token)
	if err == nil {
		return c.NoContent(http.StatusOK)
	}

	token = m
	u2 := uuid.NewV4()
	token.Token = u2.String()

	if c.System().TestMode {
		c.Response().Header().Set("X-Test-Token", token.Token)
	} else {
		err = c.System().EmailSender.Send(&email.Email{
			From:    stngs.EmailFrom,
			To:      []string{m.Email},
			Subject: c.I18n().T("Register"),
			Body: fmt.Sprintf(
				"Hi there,\n\nplease proceed with your registration by visiting this link:\n%s\n\nIf you didn't request this, please ignore this email.\n\nProxeus",
				helpers.AbsoluteURL(c, "/register/", token.Token),
			),
		})
		if err != nil {
			return c.NoContent(http.StatusExpectationFailed)
		}
	}

	err = c.System().Cache.Put(resetKey, token)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	err = c.System().Cache.Put(token.Token, token)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func Register(e echo.Context) error {
	c := e.(*www.Context)
	token := c.Param("token")
	var tokenRequest *TokenRequest
	err := c.System().Cache.Get(token, &tokenRequest)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	p := &struct {
		Password string `json:"password"`
	}{}
	if err := c.Bind(&p); err != nil {
		return err
	}
	if len(p.Password) < 6 {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"password": []map[string]interface{}{{"msg": c.I18n().T("Password not strong enough")}}})
	}
	newUser := &model.User{Email: tokenRequest.Email, Role: tokenRequest.Role}
	err = c.System().DB.User.Put(root, newUser)
	if err != nil {
		return c.NoContent(http.StatusExpectationFailed)
	}

	copyWorkflows(c, newUser)

	err = c.System().DB.User.PutPw(newUser.ID, p.Password)
	if err != nil {
		return c.NoContent(http.StatusExpectationFailed)
	}
	err = c.System().Cache.Remove(tokenRequest.Email + "_register")
	if err != nil {
		return c.NoContent(http.StatusExpectationFailed)
	}
	err = c.System().Cache.Remove(tokenRequest.Token)
	if err != nil {
		return c.NoContent(http.StatusExpectationFailed)
	}
	return c.NoContent(http.StatusOK)
}

func ResetPasswordRequest(e echo.Context) (err error) {
	c := e.(*www.Context)
	m := &TokenRequest{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	err = validate.Struct(m)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}
	if usr, err := c.System().DB.User.GetByEmail(m.Email); err == nil {
		resetKey := m.Email + "_reset_pw"
		var token *TokenRequest
		err = c.System().Cache.Get(resetKey, &token)
		if err != nil {
			token = m
			u2 := uuid.NewV4()
			token.Token = u2.String()
			token.UserID = usr.ID
			stngs := c.System().GetSettings()
			err = c.System().EmailSender.Send(&email.Email{
				From:    stngs.EmailFrom,
				To:      []string{m.Email},
				Subject: c.I18n().T("Reset Password"),
				Body: fmt.Sprintf(
					"Hi %s,\n\nif you requested a password reset, please go on and click on this link to reset your password\n%s\n\nIf you didn't request it, please ignore this email.\n\nProxeus",
					usr.Name,
					helpers.AbsoluteURL(c, "/reset/password/", token.Token),
				),
			})
			if err != nil {
				return c.NoContent(http.StatusExpectationFailed)
			}
			err = c.System().Cache.Put(resetKey, token)
			if err != nil {
				return c.NoContent(http.StatusInternalServerError)
			}
			err = c.System().Cache.Put(token.Token, token)
			if err != nil {
				return c.NoContent(http.StatusInternalServerError)
			}
		}
	}
	// always return ok if provided email was valid
	// otherwise public users can test what email accounts exist
	return c.NoContent(http.StatusOK)
}

func ResetPassword(e echo.Context) error {
	c := e.(*www.Context)
	token := c.Param("token")
	var resetPwByEmail *TokenRequest
	err := c.System().Cache.Get(token, &resetPwByEmail)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	p := &struct {
		Password string `json:"password"`
	}{}
	if err := c.Bind(&p); err != nil {
		return err
	}
	if len(p.Password) < 6 {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"password": []map[string]interface{}{{"msg": c.I18n().T("Password not strong enough")}}})
	}
	err = c.System().DB.User.PutPw(resetPwByEmail.UserID, p.Password)
	if err != nil {
		return c.NoContent(http.StatusExpectationFailed)
	}
	err = c.System().Cache.Remove(resetPwByEmail.Email + "_reset_pw")
	if err != nil {
		return c.NoContent(http.StatusExpectationFailed)
	}
	err = c.System().Cache.Remove(resetPwByEmail.Token)
	if err != nil {
		return c.NoContent(http.StatusExpectationFailed)
	}
	return c.NoContent(http.StatusOK)
}

func ChangeEmailRequest(e echo.Context) (err error) {
	c := e.(*www.Context)
	m := &TokenRequest{}
	_ = c.Bind(&m)
	err = validate.Struct(m)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}
	sess := c.Session(false)
	if sess == nil || sess.AccessRights() == model.PUBLIC {
		return c.NoContent(http.StatusUnauthorized)
	}
	if usr, err := c.System().DB.User.GetByEmail(m.Email); usr == nil {
		resetKey := m.Email + "_change_email"
		var token *TokenRequest
		err = c.System().Cache.Get(resetKey, &token)
		if err != nil {
			token = m
			u2 := uuid.NewV4()
			usr, _ = c.System().DB.User.Get(sess, sess.UserID())
			if usr == nil {
				return c.NoContent(http.StatusUnauthorized)
			}
			token.Token = u2.String()
			token.UserID = sess.UserID()
			stngs := c.System().GetSettings()
			err = c.System().EmailSender.Send(&email.Email{
				From:    stngs.EmailFrom,
				To:      []string{m.Email},
				Subject: c.I18n().T("Change Email"),
				Body: fmt.Sprintf(
					"Hi %s,\n\nif you have requested an email change, please go on and click on this link to validate it:\n%s\n\nIf you didn't request it, please ignore this email.\n\nProxeus",
					usr.Name,
					helpers.AbsoluteURL(c, "/change/email/", token.Token),
				),
			})
			if err != nil {
				return c.NoContent(http.StatusExpectationFailed)
			}
			err = c.System().Cache.Put(resetKey, token)
			if err != nil {
				return c.NoContent(http.StatusInternalServerError)
			}
			err = c.System().Cache.Put(token.Token, token)
			if err != nil {
				return c.NoContent(http.StatusInternalServerError)
			}
		}
	} else {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"email": []map[string]interface{}{{"msg": c.I18n().T("Please choose another one.")}}})
	}
	return c.NoContent(http.StatusOK)
}

func ChangeEmail(e echo.Context) error {
	c := e.(*www.Context)
	token := c.Param("token")
	var tokenRequest *TokenRequest
	err := c.System().Cache.Get(token, &tokenRequest)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	err = c.System().DB.User.UpdateEmail(tokenRequest.UserID, tokenRequest.Email)
	if err != nil {
		return c.NoContent(http.StatusExpectationFailed)
	}
	err = c.System().Cache.Remove(tokenRequest.Email + "_change_email")
	if err != nil {
		return c.NoContent(http.StatusExpectationFailed)
	}
	err = c.System().Cache.Remove(tokenRequest.Token)
	if err != nil {
		return c.NoContent(http.StatusExpectationFailed)
	}
	sess := c.Session(false)
	if sess != nil && sess.UserID() == tokenRequest.UserID {
		sess.Delete("user")
		getUserFromSession(c, sess)
	}
	return c.NoContent(http.StatusOK)
}

func LogoutHandler(e echo.Context) error {
	c := e.(*www.Context)
	c.EndSession()
	return c.JSON(http.StatusOK, map[string]string{"location": "/"})
}

func redirectAfterLogin(role model.Role, referer string) string {
	suggestedByRole := ""
	switch role {
	case model.GUEST:
		suggestedByRole = "/guest"
	case model.USER:
		suggestedByRole = "/user"
	case model.CREATOR, model.ADMIN, model.SUPERADMIN, model.ROOT:
		suggestedByRole = "/admin"
	default:
		suggestedByRole = "/"
	}
	if referer != "" && strings.HasPrefix(referer, suggestedByRole) {
		return referer
	}
	return suggestedByRole
}

// apiChallengeHandler replies with a message to be signed
func ChallengeHandler(e echo.Context) error {
	c := e.(*www.Context)
	langMessage := c.I18n().T("Account sign message")
	challengeHex := blockchain.CreateSignInChallenge(langMessage)
	sess := c.Session(true)
	if sess == nil {
		return c.NoContent(http.StatusBadRequest)
	}
	sess.Put("challenge", challengeHex)
	return c.HTML(http.StatusOK, challengeHex)
}

func MeHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusNotFound)
	}
	u := getUserFromSession(c, sess)
	if u != nil {
		return c.JSON(http.StatusOK, u)
	}
	return c.NoContent(http.StatusNotFound)
}

type UserWithPw struct {
	model.User
	Password string `json:"password,omitempty"`
}

func MeUpdateHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusBadRequest)
	}
	item := &UserWithPw{}
	if err := c.Bind(&item); err != nil {
		return err
	}
	if sess.UserID() == item.ID {
		u, err := c.System().DB.User.Get(sess, sess.UserID())
		if err != nil {
			return c.String(http.StatusUnprocessableEntity, err.Error())
		}
		if u != nil {
			u.Name = item.Name
			u.Detail = item.Detail
			u.WantToBeFound = item.WantToBeFound
			err = c.System().DB.User.Put(sess, u)
			if err != nil {
				return c.String(http.StatusUnprocessableEntity, err.Error())
			}
			if len(item.Password) >= 6 {
				err = c.System().DB.User.PutPw(u.ID, item.Password)
				if err != nil {
					return c.String(http.StatusUnprocessableEntity, err.Error())
				}
			}
			sess.Put("user", u)
			return c.NoContent(http.StatusOK)
		}
	}
	return c.NoContent(http.StatusBadRequest)
}

func PutProfilePhotoHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusBadRequest)
	}
	_, err := c.System().DB.User.PutProfilePhoto(sess, sess.UserID(), c.Request().Body)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	return c.NoContent(http.StatusOK)
}

func GetProfilePhotoHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusNotFound)
	}
	id := c.QueryParam("id")
	if id == "" {
		id = sess.UserID()
	}
	_, err := c.System().DB.User.GetProfilePhoto(sess, id, c.Response().Writer)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	c.Response().Committed = true
	c.Response().Header().Set("Content-Type", "image/jpeg")
	return c.NoContent(http.StatusOK)
}

// Check if a payment is required for current user for the workflow.
// Return http OK if no payment is required or if payment is required and a payment a matching payment with status = "confirmed" is found
func CheckForWorkflowPayment(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(true)
	if sess == nil {
		return c.NoContent(http.StatusNotFound)
	}
	workflowId := strings.TrimSpace(c.QueryParam("workflowId"))

	if workflowId == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	user, err := c.System().DB.User.Get(sess, sess.UserID())
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	paymentRequired, err := payment.CheckIfWorkflowPaymentRequired(c, workflowId)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	if paymentRequired {
		_, err := c.System().DB.WorkflowPayments.GetByWorkflowIdAndFromEthAddress(workflowId, user.EthereumAddr, []string{model.PaymentStatusConfirmed})
		if err != nil {
			if err == strm.ErrNotFound {
				return c.NoContent(http.StatusNotFound)
			}
			return c.NoContent(http.StatusBadRequest)
		}
	}

	return c.NoContent(http.StatusOK)
}

var errNoPaymentFound = errors.New("no payment for workflow")

func DocumentHandler(e echo.Context) error {
	c := e.(*www.Context)
	ID := c.Param("ID")
	var err error
	sess := c.Session(true)
	if sess == nil {
		return c.NoContent(http.StatusNotFound)
	}

	var st *app.Status

	var wf *model.WorkflowItem
	wf, err = c.System().DB.Workflow.GetPublished(sess, ID)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	docApp := getDocApp(c, sess, ID)
	if docApp == nil {
		paymentRequired, err := payment.CheckIfWorkflowPaymentRequired(c, ID)
		if err != nil {
			return c.String(http.StatusNotFound, err.Error())
		}

		if paymentRequired {
			sess := c.Session(false)
			user, err := c.System().DB.User.Get(sess, sess.UserID())
			if err != nil {
				return c.NoContent(http.StatusBadRequest)
			}
			err = payment.RedeemPayment(c.System().DB.WorkflowPayments, wf.ID, user.EthereumAddr)
			if err != nil {
				log.Println("[redeemPayment] ", err.Error())
				return c.String(http.StatusUnprocessableEntity, errNoPaymentFound.Error())
			}
		}

		usrDataItem, _, err := c.System().DB.UserData.GetByWorkflow(sess, wf, false)
		if err != nil {
			if err != strm.ErrNotFound {
				return c.String(http.StatusNotFound, err.Error())
			}

			usrDataItem = &model.UserDataItem{
				WorkflowID: wf.ID,
				Name:       wf.Name,
				Detail:     wf.Detail,
			}
			err := c.System().DB.UserData.Put(sess, usrDataItem)
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}
		}

		docApp, err = app.NewDocumentApp(usrDataItem, sess, c.System(), ID, sess.SessionDir())
		if err != nil {
			return c.String(http.StatusUnprocessableEntity, err.Error())
		}
		sess.Put("docApp_"+ID, docApp)
	}

	st, err = docApp.Current(nil)
	if err == nil {
		return c.JSON(http.StatusOK, map[string]interface{}{"name": docApp.WF().Name, "status": st})
	} else {
		if er, ok := err.(validate.ErrorMap); ok {
			er.Translate(func(key string, args ...string) string {
				return c.I18n().T(key, args)
			})
			return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"errors": er, "name": docApp.WF().Name, "status": st})
		}
		return c.String(http.StatusBadRequest, err.Error())
	}
}

func DocumentDeleteHandler(e echo.Context) error {
	c := e.(*www.Context)
	ID := c.Param("ID")
	sess := c.Session(false)
	if sess != nil {
		userDataItem, err := c.System().DB.UserData.Get(sess, ID)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		sess.Delete("docApp_" + userDataItem.WorkflowID)
		err = c.System().DB.UserData.Delete(sess, ID)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}
	return c.NoContent(http.StatusBadRequest)
}

func DocumentEditHandler(e echo.Context) error {
	c := e.(*www.Context)
	ID := c.Param("ID")
	sess := c.Session(false)
	if sess != nil {
		formInput, err := helpers.ParseDataFromReq(c)
		if err != nil {
			return err
		}
		if n, ok := formInput["name"]; ok {
			if fname, ok := n.(string); ok {
				if len(fname) < 80 && filenameRegex.MatchString(fname) {
					usrDataItem, err := c.System().DB.UserData.Get(sess, ID)
					if err != nil {
						return c.String(http.StatusBadRequest, err.Error())
					}
					if n, ok := formInput["detail"]; ok {
						if fdetail, ok := n.(string); ok {
							usrDataItem.Detail = fdetail
						}
					}
					usrDataItem.Name = fname
					err = c.System().DB.UserData.Put(sess, usrDataItem)
					if err != nil {
						return c.String(http.StatusBadRequest, err.Error())
					}
					return c.NoContent(http.StatusOK)
				}
			}
		}
	}
	return c.NoContent(http.StatusUnprocessableEntity)
}

func getUserFromSession(c *www.Context, s *session.Session) (user *model.User) {
	if s == nil {
		return nil
	}
	err := s.Get("user", &user)
	if err != nil {
		if s.ID() != "" {
			id := s.UserID()
			user, err = c.System().DB.User.Get(s, id)
			if err != nil {
				return nil
			}
			s.Put("user", user)
		}
	}
	return user
}

func DocumentNextHandler(e echo.Context) error {
	c := e.(*www.Context)
	ID := c.Param("ID")
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	docApp := getDocApp(c, sess, ID)
	if docApp == nil {
		return c.String(http.StatusBadRequest, "app does not exist")
	}
	var (
		status *app.Status
	)

	formInput, err := helpers.ParseDataFromReq(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	status, err = docApp.Next(formInput)
	if err == nil && !status.HasNext {
		//done
		_, _, status, err = docApp.Confirm(c.Lang(), formInput, c.System().DB.UserData)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		//after tx success
		if _, ok := c.QueryParams()["final"]; ok {
			dataID := docApp.DataID
			sess.Delete("docApp_" + ID)
			var item *model.UserDataItem
			item, err = c.System().DB.UserData.Get(sess, dataID)
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}
			item.Finished = true
			err = c.System().DB.UserData.Put(sess, item)
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}

			return c.JSON(http.StatusOK, map[string]interface{}{"id": dataID})
		}
	}

	resData := map[string]interface{}{
		"status": status,
		"id":     docApp.DataID,
	}
	if err != nil {
		if er, ok := err.(validate.ErrorMap); ok {
			er.Translate(func(key string, args ...string) string {
				return c.I18n().T(key, args)
			})
			resData["errors"] = er
			return c.JSON(http.StatusUnprocessableEntity, resData)
		}
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, resData)
}

func DocumentPrevHandler(e echo.Context) error {
	c := e.(*www.Context)
	ID := c.Param("ID")
	if ID != "" {
		sess := c.Session(false)
		if sess != nil {
			docApp := getDocApp(c, sess, ID)
			if docApp != nil {
				resData := map[string]interface{}{
					"status": docApp.Previous(),
				}
				return c.JSON(http.StatusOK, resData)
			}
		}
	}
	return c.NoContent(http.StatusBadRequest)
}

func DocumentDataHandler(e echo.Context) error {
	c := e.(*www.Context)
	ID := c.Param("ID")
	sess := c.Session(false)
	if sess != nil {
		docApp := getDocApp(c, sess, ID)
		if docApp != nil {
			input, err := helpers.ParseDataFromReq(c)
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}
			verrs, err := docApp.UpdateData(input, false)
			if len(verrs) > 0 {
				verrs.Translate(func(key string, args ...string) string {
					return c.I18n().T(key, args)
				})
				return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"errors": verrs})
			} else if err == nil {
				return c.NoContent(http.StatusOK)
			}
		}
	}
	return c.NoContent(http.StatusBadRequest)
}

func DocumentFileGetHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess != nil {
		ID := c.Param("ID")
		docApp := getDocApp(c, sess, ID)
		if docApp != nil {
			finfo, err := docApp.GetFile(c.Param("inputName"))
			if err == nil && finfo != nil {
				c.Response().Header().Set("Content-Type", finfo.ContentType())
				c.Response().Committed = true //prevents from-> http: multiple response.WriteHeader calls
				_, err = finfo.Read(c.Response().Writer)
				if err == nil {
					return c.NoContent(http.StatusOK)
				}
			}
			if docApp.DataID != "" {
				dataPath := fmt.Sprintf("input.%s", c.Param("inputName"))
				f, err := c.System().DB.UserData.GetDataFile(sess, docApp.DataID, dataPath)
				if err == nil {
					c.Response().Header().Set("Content-Type", f.ContentType())
					c.Response().Committed = true //prevents from-> http: multiple response.WriteHeader calls
					_, err = f.Read(c.Response().Writer)
					if err == nil {
						return c.NoContent(http.StatusOK)
					}
				}
			}
		}
	}
	return c.NoContent(http.StatusNotFound)
}

func DocumentFilePostHandler(e echo.Context) error {
	c := e.(*www.Context)
	fieldname := c.Param("inputName")
	fileName, _ := url.QueryUnescape(c.Request().Header.Get("File-Name"))
	if fieldname != "" {
		sess := c.Session(false)
		if sess != nil {
			ID := c.Param("ID")
			docApp := getDocApp(c, sess, ID)
			if docApp != nil {
				defer c.Request().Body.Close()
				verrs, err := docApp.UpdateFile(fieldname, file.Meta{Name: fileName, ContentType: c.Request().Header.Get("Content-Type"), Size: 0}, c.Request().Body)
				if len(verrs) > 0 {
					verrs.Translate(func(key string, args ...string) string {
						return c.I18n().T(key, args)
					})
					return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"errors": verrs})
				} else if err == nil {
					var finfo *file.IO
					finfo, err = docApp.GetFile(fieldname)
					if err == nil && finfo != nil {
						return c.JSON(http.StatusOK, finfo)
					}
					return c.NoContent(http.StatusNoContent)
				}
			}
		}
	}
	return c.NoContent(http.StatusBadRequest)
}

func DocumentPreviewHandler(e echo.Context) error {
	c := e.(*www.Context)
	ID := c.Param("ID")
	tmplID := c.Param("templateID")
	lang := c.Param("lang")
	format := c.Param("format")
	if ID != "" && tmplID != "" && lang != "" {
		sess := c.Session(false)
		if sess != nil {
			docApp := getDocApp(c, sess, ID)
			if docApp != nil {
				err := docApp.Preview(tmplID, lang, format, c.Response())
				if err == nil {
					return nil
				} else if !os.IsNotExist(err) {
					if err, ok := err.(net.Error); ok && err.Timeout() {
						return c.NoContent(http.StatusServiceUnavailable)
					}
					return c.NoContent(http.StatusBadRequest)
				}
			}
		}
	}
	return c.NoContent(http.StatusNotFound)
}

func getDocApp(c *www.Context, sess *session.Session, ID string) *app.DocumentFlowInstance {
	if sess != nil {
		var docApp *app.DocumentFlowInstance
		sessDocAppID := "docApp_" + ID
		err := sess.Get(sessDocAppID, &docApp)
		if err != nil {
			return nil
		}
		if docApp != nil && docApp.NeedToBeInitialized() {
			err = docApp.Init(sess, c.System())
			if err != nil {
				log.Println("Init err", err)
				return nil
			}
		}
		return docApp
	}
	return nil
}

func UserDocumentListHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	contains := c.QueryParam("c")
	settings := helpers.RequestOptions(c)
	items, err := c.System().DB.UserData.List(sess, contains, settings, false)
	if err == nil && items != nil {
		return c.JSON(http.StatusOK, items)
	}
	return c.NoContent(http.StatusBadRequest)
}

func UserDocumentGetHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	id := c.Param("ID")
	items, err := c.System().DB.UserData.Get(sess, id)
	if err == nil && items != nil {
		return c.JSON(http.StatusOK, items)
	}
	return c.NoContent(http.StatusNotFound)
}

func UserDocumentSignatureRequestGetCurrentUserHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	user, err := c.System().DB.User.Get(sess, sess.UserID())
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	ethAddr := user.EthereumAddr
	if len(ethAddr) != 42 {
		return c.NoContent(http.StatusNotFound)
	}
	signatureRequests, err := c.System().DB.SignatureRequests.GetBySignatory(ethAddr)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	type SignatureRequestItemComplete struct {
		ID          string  `json:"id"`
		DocID       string  `json:"docID"`
		Hash        string  `json:"hash"`
		From        string  `json:"requestorName"`
		FromAddr    string  `json:"requestorAddr"`
		RequestedAt *string `json:"requestedAt,omitempty"`
		Rejected    bool    `json:"rejected"`
		RejectedAt  *string `json:"rejectedAt,omitempty"`
		Revoked     bool    `json:"revoked"`
		RevokedAt   *string `json:"revokedAt,omitempty"`
	}

	type SignatureRequests []SignatureRequestItemComplete

	var requests = *new(SignatureRequests)
	for _, sigreq := range *signatureRequests {
		var requesterName string
		requester, err := c.System().DB.User.GetByBCAddress(sigreq.Requestor)
		if err != nil && err.Error() != "not found" {
			return c.NoContent(http.StatusInternalServerError)
		}
		if requester != nil {
			requesterName = requester.Name
		}

		var reqAt string
		reqAt = sigreq.RequestedAt.Format("2.1.2006 15:04")
		var rejAt string
		if sigreq.Rejected {
			rejAt = sigreq.RejectedAt.Format("2.1.2006 15:04")
		}
		var revAt string
		revAt = sigreq.RevokedAt.Format("2.1.2006 15:04")

		reqitem := SignatureRequestItemComplete{
			sigreq.DocId,
			sigreq.DocPath,
			sigreq.Hash,
			requesterName,
			sigreq.Requestor,
			&reqAt,
			sigreq.Rejected,
			&rejAt,
			sigreq.Revoked,
			&revAt,
		}
		if !sigreq.Revoked {
			reqitem.RevokedAt = nil
		}
		if !sigreq.Rejected {
			reqitem.RejectedAt = nil
		}

		requests = append(requests, reqitem)
	}
	return c.JSON(http.StatusOK, requests)
}

func UserDeleteHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	//remove documents / workflow instances of user
	userDataDB := c.System().DB.UserData
	workflowInstances, err := userDataDB.List(sess, "", storage.Options{}, false)
	if err != nil && err.Error() != "not found" {
		return c.NoContent(http.StatusInternalServerError)
	}
	for _, workflowInstance := range workflowInstances {
		err = userDataDB.Delete(sess, workflowInstance.ID)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	//set workflow templates to deactivated
	workflowDB := c.System().DB.Workflow
	workflows, err := workflowDB.List(sess, "", storage.Options{})
	if err != nil && err.Error() != "not found" {
		return c.NoContent(http.StatusInternalServerError)
	}
	for _, workflow := range workflows {
		if workflow.OwnedBy(sess) {
			workflow.Deactivated = true
			err = workflowDB.Put(sess, workflow)
			if err != nil {
				return c.NoContent(http.StatusInternalServerError)
			}
		}
	}

	// unset user data and set inactive
	userDB := c.System().DB.User
	user, err := userDB.Get(sess, sess.UserID())
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	user.Active = false
	user.EthereumAddr = "0x"
	user.Email = ""
	user.Name = ""
	user.Photo = ""
	user.PhotoPath = ""
	user.WantToBeFound = false

	err = userDB.Put(sess, user)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return LogoutHandler(c)
}

func UserDocumentSignatureRequestGetByDocumentIDHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	docId := c.Param("docID")
	id := c.Param("ID")

	signatureRequests, err := c.System().DB.SignatureRequests.GetByID(id, docId)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	type SignatureRequestItemMinimal struct {
		SignatoryName string  `json:"signatoryName"`
		SignatoryAddr string  `json:"signatoryAddr"`
		RequestedAt   *string `json:"requestedAt,omitempty"`
		Rejected      bool    `json:"rejected"`
		RejectedAt    *string `json:"rejectedAt,omitempty"`
		Revoked       bool    `json:"revoked"`
		RevokedAt     *string `json:"revokedAt,omitempty"`
	}

	type SignatureRequests []SignatureRequestItemMinimal

	var requests = *new(SignatureRequests)
	for _, sigreq := range *signatureRequests {
		signatoryName := *new(string)
		item, err := c.System().DB.User.GetByBCAddress(sigreq.Signatory)
		if err == nil {
			signatoryName = item.Name
		}
		var reqAt string
		reqAt = sigreq.RequestedAt.Format("2.1.2006 15:04")
		var rejAt string
		if sigreq.Rejected {
			rejAt = sigreq.RejectedAt.Format("2.1.2006 15:04")
		}
		var revAt string
		revAt = sigreq.RevokedAt.Format("2.1.2006 15:04")

		reqitem := SignatureRequestItemMinimal{
			signatoryName,
			sigreq.Signatory,
			&reqAt,
			sigreq.Rejected,
			&rejAt,
			sigreq.Revoked,
			&revAt,
		}
		if !sigreq.Rejected {
			reqitem.RejectedAt = nil
		}

		requests = append(requests, reqitem)
	}
	return c.JSON(http.StatusOK, requests)
}

func UserDocumentSignatureRequestAddHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	docId := c.Param("docID")
	id := c.Param("ID")

	signatory := c.FormValue("signatory")
	fileInfo, err := c.System().DB.UserData.GetDataFile(sess, id, docId)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	if !strings.HasPrefix(docId, "docs") {
		return c.NoContent(http.StatusNotFound)
	}

	var documentBytes []byte
	documentBytes, err = fileInfo.ReadAll()
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	docHash := crypto.Keccak256Hash(documentBytes).String()

	signatoryObj, err := c.System().DB.User.GetByBCAddress(signatory)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	fileObj, err := c.System().DB.UserData.Get(sess, id)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	if fileObj.Permissions.Grant == nil || !fileObj.Permissions.Grant[signatoryObj.UserID()].IsRead() {
		if fileObj.Permissions.Grant == nil {
			fileObj.Permissions.Grant = make(map[string]model.Permission)
		}
		fileObj.Permissions.Grant[signatoryObj.UserID()] = model.Permission{byte(1)}
		fileObj.Permissions.Change(sess, &fileObj.Permissions)

		err = c.System().DB.UserData.Put(sess, fileObj)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	fileObj, _ = c.System().DB.UserData.Get(sess, id)

	requestor, err := c.System().DB.User.Get(sess, sess.UserID())
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	requestorAddr := requestor.EthereumAddr

	requestItem := model.SignatureRequestItem{
		DocId:       id,
		DocPath:     docId,
		Hash:        docHash,
		Requestor:   requestorAddr,
		RequestedAt: time.Now(),
		Signatory:   signatory,
		Rejected:    false,
	}

	signatureRequests, err := c.System().DB.SignatureRequests.GetByID(id, docId)

	if err == nil {
		for _, sigreq := range *signatureRequests {
			if sigreq.Signatory == signatory &&
				sigreq.Hash == docHash &&
				sigreq.Rejected == false &&
				sigreq.Revoked == false {
				return c.String(http.StatusConflict, c.I18n().T("Request already exists"))
			}
		}
	}

	err = c.System().DB.SignatureRequests.Add(&requestItem)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	sig, err := c.System().DB.User.GetByBCAddress(signatory)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	if len(sig.Email) > 3 {
		/*
			Your signature was requested for a document on <platform base URL>by <Name> (<Email>)<Ethereum-Addr>

			The requestor would like you to review and sign the document on the platform.

			To check your pending signature requests, please log in <here (link to requests, if logged in>
		*/
		emailFrom := c.System().GetSettings().EmailFrom
		c.System().EmailSender.Send(&email.Email{From: emailFrom, To: []string{sig.Email}, Subject: c.I18n().T("New signature request received"), Body: "<div>Your signature was requested for a document from " + c.Request().Host + " <br />by " + requestor.Name + " (" + requestor.Email + ")<br />" + requestorAddr + "<br /><br />The requestor would like you to review and sign the document on the platform.<br /><br />To check your pending signature requests, please log in <a href='" + helpers.AbsoluteURL(c, "/user/signature-requests") + "'>here</a></div>"})
	}

	return c.NoContent(http.StatusOK)
}

func UserDocumentSignatureRequestRejectHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	docId := c.Param("docID")
	id := c.Param("ID")

	item, err := c.System().DB.User.Get(sess, sess.UserID())
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	signatoryAddr := item.EthereumAddr
	signatureRequests, err := c.System().DB.SignatureRequests.GetByID(id, docId)
	signatureRequest := (*signatureRequests)[0]
	req := signatureRequest.Requestor

	err = c.System().DB.SignatureRequests.SetRejected(id, docId, signatoryAddr)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	requestorAddr, err := c.System().DB.User.GetByBCAddress(req)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	if len(requestorAddr.Email) > 3 {
		/*
			Your signature request for a document on <platform base URL> from <timestamp> has been rejected by <Name> (<Email>)<Ethereum-Addr>

			You may send another request if you think this was by mistake.

		*/
		emailFrom := c.System().GetSettings().EmailFrom
		c.System().EmailSender.Send(&email.Email{From: emailFrom, To: []string{requestorAddr.Email}, Subject: c.I18n().T("Signature request rejected"), Body: "<div>Your signature request for a document on " + c.Request().Host + " from " + signatureRequest.RequestedAt.Format("2.1.2006 15:04") + " <br />has been rejected by  " + item.Name + " (" + item.Email + ")<br />" + item.EthereumAddr + "<br /><br />You may send another request if you think this was by mistake.</div>"})
	}

	return c.NoContent(http.StatusOK)
}

func UserDocumentSignatureRequestRevokeHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	docId := c.Param("docID")
	id := c.Param("ID")
	signatory := c.FormValue("signatory")

	sig, err := c.System().DB.User.GetByBCAddress(signatory)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	signatoryEmail := sig.Email

	err = c.System().DB.SignatureRequests.SetRevoked(id, docId, signatory)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	requestor, err := c.System().DB.User.Get(sess, sess.UserID())
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	if len(signatoryEmail) > 3 {
		/*
			Earlier you may have received a signature request from <base URL>by <Name> (<Email>)<Ethereum-Addr>

			The requestor has retracted the request. You may still log in and view the request, but can no longer sign the document.

			To check your signature requests, please log in <here (link to requests, if logged in>
		*/
		c.System().EmailSender.Send(&email.Email{To: []string{signatoryEmail}, Subject: c.I18n().T("New signature request received"), Body: "<div>Earlier you may have received a signature request from " + c.Request().Host + " by " + requestor.Name + " (" + requestor.Email + ")<br />" + requestor.EthereumAddr + "<br /><br />The requestor has retracted the request. You may still log in and view the request, but can no longer sign the document.<br /><br />To check your signature requests, please log in <a href='" + helpers.AbsoluteURL(c, "/user/signature-requests") + "'>here</a></div>"})
	}

	return c.NoContent(http.StatusOK)
}

func PutTestSignature(e echo.Context) error {
	c := e.(*www.Context)
	if !c.System().TestMode {
		return echo.ErrBadRequest
	}
	var req struct {
		TxHash     string
		FileHash   string
		SignerAddr string
	}
	c.Bind(&req)
	l := types.Log{
		Address: common.HexToAddress(cfg.Config.XESContractAddress),
		Topics: []common.Hash{
			common.HexToHash("0xe898b82efcc77a621bbc620d08e84d0b44e341fd7720cc544de745bdec8ebece"),
			common.HexToHash(req.FileHash),
			common.HexToHash("0x000000000000000000000000" + req.SignerAddr[2:]),
		},
		TxHash: common.HexToHash(req.TxHash),
	}

	blockchain.TestChannelSignature <- l
	return nil
}

func UserDocumentFileHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	format := c.QueryParam("format")

	dataPath := c.Param("dataPath")
	id := c.Param("ID")
	fileInfo, err := c.System().DB.UserData.GetDataFile(sess, id, dataPath)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	inlineOrAttachment := "attachment"
	if _, ok := c.QueryParams()["inline"]; ok {
		inlineOrAttachment = "inline"
	}

	if strings.HasPrefix(dataPath, "docs") {
		//final doc
		resp := c.Response()
		if fileInfo.ContentType() != "" {
			resp.Header().Set("Content-Type", fileInfo.ContentType())
		}

		fileName := fileInfo.NameWithExt("pdf")
		contentDisposition := fmt.Sprintf(`%s; filename="%s"`, inlineOrAttachment, fileName)
		resp.Header().Set("Content-Disposition", contentDisposition)
		resp.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
		resp.Committed = true //prevents from-> http: multiple response.WriteHeader calls
		_, err = fileInfo.Read(resp.Writer)
		if err != nil {
			return c.NoContent(http.StatusNotFound)
		}
	} else {
		//template with format
		dat, files, _ := c.System().DB.UserData.GetDataAndFiles(sess, id, "input")
		formt := eio.Format(format)
		dsResp, err := c.System().DS.Compile(eio.Template{
			Format:       formt,
			Data:         map[string]interface{}{"input": dat},
			TemplatePath: fileInfo.Path(),
			Assets:       files,
		})
		if err != nil {
			return c.String(http.StatusNotFound, err.Error())
		}
		resp := c.Response()
		resp.Header().Set("Content-Type", dsResp.Header.Get("Content-Type"))
		resp.Header().Set("Content-Length", dsResp.Header.Get("Content-Length"))
		resp.Header().Set("Content-Pages", dsResp.Header.Get("Content-Pages"))
		resp.Header().Set("Content-Disposition", fmt.Sprintf("%s;filename=\"%s\"", inlineOrAttachment, fileInfo.NameWithExt(formt.String())))
		defer dsResp.Body.Close()
		resp.Committed = true //prevents from-> http: multiple response.WriteHeader calls
		_, err = io.Copy(resp.Writer, dsResp.Body)
	}
	return c.NoContent(http.StatusOK)
}

func AdminUserUpdateHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	if strings.Contains(c.Request().Header.Get("Content-Type"), "application/json") {
		body, _ := ioutil.ReadAll(c.Request().Body)
		item := model.User{}
		item.ID = c.QueryParam("id")
		err := json.Unmarshal(body, &item)
		if err == nil {
			exItem, err := c.System().DB.User.Get(sess, item.ID)
			if err != nil {
				return c.NoContent(http.StatusNotFound)
			}
			exItem.Name = item.Name
			exItem.Detail = item.Detail
			exItem.Role = item.Role
			err = c.System().DB.User.Put(sess, exItem)
			if err != nil {
				return c.NoContent(http.StatusInternalServerError)
			}
			if sess.UserID() == exItem.ID {
				sess.Put("user", exItem)
			}
			return c.JSON(http.StatusOK, exItem)
		}
	}
	return c.NoContent(http.StatusBadRequest)
}

func AdminUserGetHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	itemID := c.Param("ID")
	item, err := c.System().DB.User.Get(sess, itemID)
	if err != nil {
		return err
	}
	if item == nil {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, item)
}

func AdminUserListHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	contains := c.QueryParam("c")
	settings := helpers.RequestOptions(c)
	dat, err := c.System().DB.User.List(sess, contains, settings)
	if err != nil || dat == nil {
		if err == model.ErrAuthorityMissing {
			return c.NoContent(http.StatusUnauthorized)
		}
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, dat)
}

func WorkflowSchema(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	id := c.Param("ID")
	wf, err := c.System().DB.Workflow.Get(sess, id)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	fieldsAndRules := utils.GetAllFormFieldsWithRulesOf(wf.Data, sess, c.System())
	wfDetails := &struct {
		*model.WorkflowItem
		Data interface{} `json:"data"`
	}{wf, fieldsAndRules}
	result := &struct {
		Workflow interface{} `json:"workflow"`
	}{Workflow: wfDetails}
	return c.JSON(http.StatusOK, result)
}

func WorkflowExecuteAtOnce(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	inputData, err := helpers.ParseDataFromReq(c)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	id := c.Param("ID")
	wItem, err := c.System().DB.Workflow.Get(sess, id)
	if err != nil || wItem.Data == nil {
		return c.NoContent(http.StatusNotFound)
	}
	err = app.ExecuteWorkflowAtOnce(c, sess, wItem, inputData)
	if err != nil {
		if er, ok := err.(validate.ErrorMap); ok {
			er.Translate(func(key string, args ...string) string {
				return c.I18n().T(key, args)
			})
			return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"errors": er})
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func CreateApiKeyHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	id := c.Param("ID")
	name := c.QueryParam("name")
	if name == "" {
		return c.String(http.StatusBadRequest, "please provide a name for your key")
	}
	apiKey, err := c.System().DB.User.CreateApiKey(sess, id, name)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if id == sess.UserID() {
		sess.Delete("user")
	}
	return c.String(http.StatusOK, apiKey)
}

func DeleteApiKeyHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	id := c.Param("ID")
	hiddenApiKey := c.QueryParam("hiddenApiKey")
	err := c.System().DB.User.DeleteApiKey(sess, id, hiddenApiKey)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if id == sess.UserID() {
		sess.Delete("user")
	}
	return c.NoContent(http.StatusOK)
}

func copyWorkflows(c *www.Context, newUser *model.User) {
	log.Println("Copy workflows to new user, if any...")
	// If some default workflows have to be assigned to the user, then clone them
	workflowIds := strings.Split(c.System().GetSettings().DefaultWorkflowIds, ",")
	workflows, err := c.System().DB.Workflow.GetList(root, workflowIds)
	if err != nil {
		log.Printf("Can't retrieve list of workflows (%v). Please check the ids exist. Error: %s", workflowIds, err.Error())
	}
	for _, workflow := range workflows {
		w := workflow.Clone()
		w.OwnerEthAddress = newUser.EthereumAddr
		w.Owner = newUser.ID
		newNodes := make(map[string]*workflow2.Node)
		oldToNewIdsMap := make(map[string]string)
		for oldId, node := range w.Data.Flow.Nodes {
			if node.Type == "form" {
				form, er := c.System().DB.Form.Get(root, node.ID)
				if er != nil {
					log.Println(err.Error())
				}
				f := form.Clone()
				er = c.System().DB.Form.Put(newUser, &f)
				if er != nil {
					log.Println("can't put form" + err.Error())
				}

				oldToNewIdsMap[node.ID] = f.ID
				node.ID = f.ID
				newNodes[node.ID] = node
				delete(w.Data.Flow.Nodes, oldId)

			} else if node.Type == "template" {
				template, er := c.System().DB.Template.Get(root, node.ID)
				if er != nil {
					log.Println(err.Error())
				}
				t := template.Clone()
				er = c.System().DB.Template.Put(newUser, &t)
				if er != nil {
					log.Println("can't put template" + err.Error())
				}
				oldToNewIdsMap[node.ID] = t.ID
				node.ID = t.ID
				newNodes[node.ID] = node
				delete(w.Data.Flow.Nodes, oldId)
			} else {
				newNodes[node.ID] = node
			}
		}
		oldStartNodeId := w.Data.Flow.Start.NodeID
		if _, ok := oldToNewIdsMap[oldStartNodeId]; ok {
			w.Data.Flow.Start.NodeID = oldToNewIdsMap[oldStartNodeId]
		}

		// Now go through all connections and map them with the new ids
		for _, node := range newNodes {
			for _, connection := range node.Connections {
				if _, ok := oldToNewIdsMap[connection.NodeID]; ok {
					connection.NodeID = oldToNewIdsMap[connection.NodeID]
				}
			}
		}
		w.Data.Flow.Nodes = newNodes
		c.System().DB.Workflow.Put(newUser, &w)
	}
}
