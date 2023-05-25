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
	"strconv"
	"strings"
	"time"

	"github.com/ProxeusApp/proxeus-core/externalnode"

	"github.com/ProxeusApp/proxeus-core/service"

	"github.com/ProxeusApp/proxeus-core/main/app"
	cfg "github.com/ProxeusApp/proxeus-core/main/config"
	"github.com/ProxeusApp/proxeus-core/main/handlers/blockchain"
	"github.com/ProxeusApp/proxeus-core/main/handlers/helpers"
	"github.com/ProxeusApp/proxeus-core/main/www"
	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/storage/database/db"
	"github.com/ProxeusApp/proxeus-core/storage/portable"
	"github.com/ProxeusApp/proxeus-core/sys"
	"github.com/ProxeusApp/proxeus-core/sys/model"
	"github.com/ProxeusApp/proxeus-core/sys/validate"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

var (
	paymentService          service.PaymentService
	userService             service.UserService
	workflowService         service.WorkflowService
	documentService         service.DocumentService
	userDocumentService     service.UserDocumentService
	fileService             service.FileService
	templateDocumentService service.TemplateDocumentService
	signatureService        service.SignatureService
	emailService            service.EmailService
	formService             service.FormService
	formComponentService    service.FormComponentService
	apiService              service.ApiService
	authService             service.AuthenticationService
	nodeService             service.NodeService
)

func Init(paymentS service.PaymentService, userS service.UserService, workflowS service.WorkflowService,
	documentS service.DocumentService, userDocumentS service.UserDocumentService, fileS service.FileService,
	templateDocumentS service.TemplateDocumentService, signatureS service.SignatureService, emailS service.EmailService,
	formS service.FormService, formCompS service.FormComponentService, apiS service.ApiService, authS service.AuthenticationService,
	nodeS service.NodeService) {

	paymentService = paymentS
	userService = userS
	workflowService = workflowS
	documentService = documentS
	userDocumentService = userDocumentS
	fileService = fileS
	templateDocumentService = templateDocumentS
	signatureService = signatureS
	emailService = emailS
	formService = formS
	formComponentService = formCompS
	apiService = apiS
	authService = authS
	nodeService = nodeS
}

func html(c echo.Context, p string) error {
	bts, err := sys.ReadAllFile(p)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	return c.HTMLBlob(http.StatusOK, bts)
}

// Returns the app.html file
func SharedByLinkHTMLHandler(c echo.Context) error {
	log.Println("SharedByLinkHTMLHandler")
	//c.Param("type") TODO different html by type for user data
	return html(c, "app.html")
}

// Returns the frontend.html file
func PublicIndexHTMLHandler(c echo.Context) error {
	return html(c, "frontend.html")
}

// Returns the user.html file
func UserBackendHTMLHandler(c echo.Context) error {
	return html(c, "user.html")
}

// Returns the app.html file
func AdminIndexHandler(c echo.Context) error {
	return html(c, "app.html")
}

type ImportExportResult struct {
	Filename  string                    `json:"filename"`
	Timestamp time.Time                 `json:"timestamp"`
	Results   portable.ProcessedResults `json:"results"`
}

// Returns exported entities
//
// @params
//
//	include string
//	EntityType  []string
//
// @returns
//
//		200 => File
//		400 => Bad Request
//	 422 => Input error
//	}
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

// Imports a file containing exported entities
//
// @params
//
//	skipExisting string
//	File File (in body)
//
// @returns
//
//		200 => OK
//		400 => Bad Request
//	 422 => Input error
//	}
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

// Exports user data
//
// @params
//
//	contains string
//	id 		 string
//
// @returns
//
//		200 => File
//		400 => Bad Request
//	 401 => StatusUnauthorized
//	 422 => Input error
//	}
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
		items, _ := userDocumentService.List(sess, c.QueryParam("contains"), storage.Options{Limit: 1000})
		if len(items) > 0 {
			id = make([]string, len(items))
			for i, item := range items {
				id[i] = item.ID
			}
		}
	}
	return Export(sess, []portable.EntityType{portable.UserData}, c, id...)
}

// Export platform settings
//
// @returns
//
//		200 => File
//		400 => Bad Request
//	 422 => Input error
//	}
func ExportSettings(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	return Export(sess, []portable.EntityType{portable.Settings}, c, "Settings")
}

// Export a user record
//
// @params
//
//	contains string
//	id 		 string
//
// @returns
//
//		200 => File
//		400 => Bad Request
//	 422 => Input error
//	}
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
		items, _ := userService.List(sess, c.QueryParam("contains"), storage.Options{Limit: 1000})
		if len(items) > 0 {
			id = make([]string, len(items))
			for i, item := range items {
				id[i] = item.ID
			}
		}
	}
	return Export(sess, []portable.EntityType{portable.User}, c, id...)
}

// Helper function for all exports
//
// @returns
//
//		200 => File
//		400 => Bad Request
//	 422 => Input error
//	}
func Export(sess *sys.Session, exportEntities []portable.EntityType, e echo.Context, id ...string) error {
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

// Returns the last exported records stored in the user session
//
// @params
// @returns
//
//		200 => imexResults
//		401 => StatusUnauthorized
//	}
func GetExportResults(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	return results("lastExport", sess, c)
}

// Returns the last imported records from the user session
//
// @params
//
//	delete boolean
//
// @returns
//
//		200 => imex results
//		401 => Unauthorized
//	}
func GetImportResults(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	return results("lastImport", sess, c)
}

func results(key string, sess *sys.Session, c echo.Context) error {
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

// Returns the platform settings and the "configured" flag that defines whether the user is configured or not
//
// @returns
//
//		200 => map[string]interface{}{"settings": settings, "configured": configured}
//		500 => Server error
//	}
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

// Initialize the system
//
//	@params initStruct => struct {
//			Settings *model.Settings `json:"settings"`
//			User     *usr            `json:"user"`
//		}
//
// @returns
//
//		200 => OK
//	 422 => unprocessable entity
//		500 => Server error
//	}
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
	type InitStruct struct {
		Settings *model.Settings `json:"settings"`
		User     *usr            `json:"user"`
	}
	var err error
	yes, _ := c.System().Configured()
	d := &InitStruct{User: &usr{}}
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
		fmt.Println("Error during PostInit settings: ", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	if !yes {
		u := &model.User{Email: d.User.Email, Role: d.User.Role}
		uex, _ := userService.GetByEmail(u.Email)
		if uex == nil {
			err = userService.Put(root, u)
			if err != nil {
				fmt.Println("Error during PostInit user: ", err)
				return c.NoContent(http.StatusInternalServerError)
			}
			err = userService.PutPassword(u.ID, d.User.Password)
			if err != nil {
				fmt.Println("Error during PostInit password: ", err)
				return c.NoContent(http.StatusInternalServerError)
			}
		}
	}
	if err := formComponentService.EnsureDefaultFormComponents(root); err != nil {
		fmt.Println("Error during form component initialisation: ", err)
	}

	return c.NoContent(http.StatusOK)
}

// Returns an object containing the following config parameters
//
//	{
//	  roles => string[] => Possible User Roles
//	  blockchainNet => string => Settings.BlockchainNet
//	  blockchainProxeusFSAddress => string => Settings.BlockchainContractAddress
//	  version => string => Proxeus version
//	}
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

// Update a users' blockchain address
//
//	@param loginForm => {
//		Signature string
//		Name      string `json:"name" form:"name"`
//		Email     string `json:"email" form:"email"`
//		Password  string `json:"password" form:"password"`
//		Address   string `json:"address" form:"address"`
//	}
//
// @returns
//
//		200 => OK
//		401 => Unauthorized
//	 422 => Challenge error/Signature error
//	 500 => User not found
//	}
func UpdateAddress(e echo.Context) error {
	c := e.(*www.Context)
	loginForm := new(loginForm)
	_ = c.Bind(loginForm)
	sess := c.Session(false)
	u := getUserFromSession(sess)
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
	item, err := userService.GetByBCAddress(address)
	if item != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"etherPK": []map[string]interface{}{{"msg": "Please choose another account."}}})
	}
	item, err = userService.GetUser(sess)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	item.EthereumAddr = address
	err = userService.Put(sess, item)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	sess.Put("user", item)
	return c.NoContent(http.StatusOK)
}

// root only feature to switch user by address - useful for permission checks
func SwitchUserHandler(e echo.Context) error {
	c := e.(*www.Context)
	user, err := userService.GetByBCAddress(e.Param("address"))
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

// Create an auth session
//
//	@params => {
//		Signature string
//		Name      string `json:"name" form:"name"`
//		Email     string `json:"email" form:"email"`
//		Password  string `json:"password" form:"password"`
//		Address   string `json:"address" form:"address"`
//	}
//
// @returns
//
//		200 => OK => {
//			"location": redirectAfterLogin(user.Role, string(referer)),
//			"created":  created,
//		}
//	 400 => Auth error
//	}
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
		created, user, err = authService.LoginWithWallet(challenge, loginForm.Signature)
		if user != nil {
			sess.Delete("challenge")
		}
	} else {
		user, err = authService.LoginWithUsernamePassword(loginForm.Email, loginForm.Password)
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

// Returns an object containing
//
//	{
//	  token => string => Session ID
//	}
func GetSessionTokenHandler(e echo.Context) (err error) {
	c := e.(*www.Context)

	username, apiKey := c.BasicAuth()

	if username == "" || apiKey == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	user, err := apiService.AuthenticateWithApiKey(apiKey)
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
		"token": sess.S.ID,
	})
}

// Ends the context's session and returns 200 => OK in any case
func DeleteSessionTokenHandler(e echo.Context) (err error) {
	c := e.(*www.Context)
	c.EndSession()
	return c.NoContent(http.StatusOK)
}

// Invite a user
//
//	@params TokenRequest struct {
//			Email  string    `json:"email" validate:"email=true,required=true"`
//			Token  string    `json:"token"`
//			UserID string    `json:"userID"`
//			Role   Role      `json:"role"`
//			Type   TokenType `json:"type"`
//		}
//
// @returns
//
//		200 => OK
//	 422 => unprocessable entity
//		500 => Server error
//	}
func InviteRequest(e echo.Context) (err error) {
	c := e.(*www.Context)
	sess := c.Session(false)
	m := &model.TokenRequest{}
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
		m.Role = model.StringToRole(stngs.DefaultRole)
	}
	if usr, err := userService.GetByEmail(m.Email); usr == nil {
		var token model.TokenRequest
		token.Email = m.Email
		token.Token = uuid.NewV4().String()
		token.Role = m.Role
		token.Type = model.TokenRegister
		err = authService.PutTokenRequest(&token)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		subject := c.I18n().T("Invitation")
		body := fmt.Sprintf(
			"Hi there,\n\nyou have been invited to join Proxeus. If you would like to benefit from the invitation, please proceed by visiting this link:\n%s\n\nProxeus",
			helpers.AbsoluteURL(c, "/register/", token.Token),
		)

		err = emailService.Send(m.Email, subject, body)
		if err != nil {
			return c.String(http.StatusFailedDependency, c.I18n().T("couldn't send the email"))
		}
	}
	return c.NoContent(http.StatusOK)
}

// Handles a registration request
//
//	@params => {
//			Email  string    `json:"email" validate:"email=true,required=true"`
//			Token  string    `json:"token"`
//			UserID string    `json:"userID"`
//			Role   Role      `json:"role"`
//			Type   TokenType `json:"type"`
//		}
//
// @returns
//
//		200 => OK
//	 417 => E-Mail error
//	 422 => Input validation error
//	 500 => Data layer error
//	}
func RegisterRequest(e echo.Context) (err error) {
	c := e.(*www.Context)
	m := &model.TokenRequest{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	err = validate.Struct(m)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	token, err := authService.RegisterRequest(c.I18n(), c.Scheme(), c.Request().Host, m)
	if err != nil {
		log.Println("[RegisterRequest] authService.RegisterRequest err: ", err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	if c.System().TestMode {
		c.Response().Header().Set("X-Test-Token", token.Token)
	}

	return c.NoContent(http.StatusOK)
}

// Persists a registration request
//
//	@params => {
//		 string => token
//	  string => password
//	}
//
// @returns
//
//		200 => OK
//	 417 => Token not found/User not found/Data Layer error
//	 422 => Input validation error
//	}
func Register(e echo.Context) error {
	c := e.(*www.Context)
	tokenID := c.Param("token")
	p := &struct {
		Password string `json:"password"`
	}{}
	if err := c.Bind(&p); err != nil {
		return err
	}

	err, errors := authService.Register(c.I18n(), tokenID, p.Password)
	if len(errors) != 0 {
		return c.JSON(http.StatusUnprocessableEntity, errors)
	}
	if err != nil {
		log.Println("[Register] authService.Register err: ", err.Error())
		return c.NoContent(http.StatusExpectationFailed)
	}

	return c.NoContent(http.StatusOK)
}

// Start a user password change request
//
//	@params => {
//			Email  string    `json:"email" validate:"email=true,required=true"`
//			Token  string    `json:"token"`
//			UserID string    `json:"userID"`
//			Role   Role      `json:"role"`
//			Type   TokenType `json:"type"`
//		}
//
// @returns
//
//		200 => OK
//		400 => Token request not found
//	 417 => E-Mail error
//	 422 => Input validation error
//	 500 => Token error
//	}
func ResetPasswordRequest(e echo.Context) (err error) {
	c := e.(*www.Context)
	m := &model.TokenRequest{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	err = validate.Struct(m)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	token, err := authService.ResetPasswordRequest(c.I18n(), c.Scheme(), c.Request().Host, m)
	if err != nil {
		log.Println("[ResetPasswordRequest] authService.ResetPasswordRequest err: ", err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	if c.System().TestMode {
		c.Response().Header().Set("X-Test-Token", token.Token)
	}

	return c.NoContent(http.StatusOK)
}

// Reset a users' password
//
// @params => string => tokenID
// @returns
//
//		200 => OK
//	 417 => Data layer error
//	 422 => Input validation error
//	}
func ResetPassword(e echo.Context) error {
	c := e.(*www.Context)
	tokenID := c.Param("token")
	p := &struct {
		Password string `json:"password"`
	}{}
	if err := c.Bind(&p); err != nil {
		return err
	}

	err, errors := authService.ResetPassword(c.I18n(), p.Password, tokenID)
	if len(errors) != 0 {
		return c.JSON(http.StatusUnprocessableEntity, errors)
	}
	if err != nil {
		log.Println("[ResetPassword] authService.ResetPassword err: ", err.Error())
		return c.NoContent(http.StatusExpectationFailed)
	}

	return c.NoContent(http.StatusOK)
}

// Start a users email change request
//
//	@param tokenRequest => {
//			Email  string    `json:"email" validate:"email=true,required=true"`
//			Token  string    `json:"token"`
//			UserID string    `json:"userID"`
//			Role   Role      `json:"role"`
//			Type   TokenType `json:"type"`
//		}
//
// @returns
//
//		200 => OK
//		401 => Unauthorized
//	 417 => E-Mail error
//	 422 => Input error
//	 500 => Token error
//	}
func ChangeEmailRequest(e echo.Context) (err error) {
	c := e.(*www.Context)
	m := &model.TokenRequest{}
	_ = c.Bind(&m)
	err = validate.Struct(m)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}
	sess := c.Session(false)
	if sess == nil || sess.AccessRights() == model.PUBLIC {
		return c.NoContent(http.StatusUnauthorized)
	}
	if usr, err := userService.GetByEmail(m.Email); usr == nil {
		var token model.TokenRequest
		usr, _ = userService.GetUser(sess)
		if usr == nil {
			return c.NoContent(http.StatusUnauthorized)
		}
		token.Email = m.Email
		token.Token = uuid.NewV4().String()
		token.UserID = sess.UserID()
		token.Type = model.TokenChangeEmail
		err = authService.PutTokenRequest(&token)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		if c.System().TestMode {
			c.Response().Header().Set("X-Test-Token", token.Token)
		}
		subject := c.I18n().T("Change Email")
		body := fmt.Sprintf(
			"Hi %s,\n\nif you have requested an email change, please go on and click on this link to validate it:\n%s\n\nIf you didn't request it, please ignore this email.\n\nProxeus",
			usr.Name,
			helpers.AbsoluteURL(c, "/change/email/", token.Token),
		)
		err = emailService.Send(m.Email, subject, body)
		if err != nil {
			return c.NoContent(http.StatusExpectationFailed)
		}
	}
	return c.NoContent(http.StatusOK)
}

// Update a users' e-mail address
//
// @param tokenID => string => Request token
// @returns
//
//		200 => OK
//		400 => Token request not found
//	 417 => Data layer error
//	}
func ChangeEmail(e echo.Context) error {
	c := e.(*www.Context)
	tokenID := c.Param("token")

	tokenRequest, err := authService.ChangeEmail(tokenID)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	sess := c.Session(false)
	if sess != nil && sess.UserID() == tokenRequest.UserID {
		sess.Delete("user")
		getUserFromSession(sess)
	}
	return c.NoContent(http.StatusOK)
}

// Remove the users' auth session
//
// @params => nil
// @returns
//
//	200 => OK => {
//		"location": "/"
//	}
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

// Returns a string in plaintext, containing the message to be signed.
// Returns 400 => Bad Request if no session is found
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

// Returns the user object from the session
//
// @params => -
// @returns
//
//		200 => User => {}
//	 404 => Not found
//	}
func MeHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusNotFound)
	}
	u := getUserFromSession(sess)
	if u != nil {
		return c.JSON(http.StatusOK, u)
	}
	return c.NoContent(http.StatusNotFound)
}

type UserWithPw struct {
	model.User
	Password string `json:"password,omitempty"`
}

// Returns the user object from the session
//
// @params => -
// @returns
//
//		200 => User => {}
//	 404 => Not found
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
		u, err := userService.GetUser(sess)
		if err != nil {
			return c.String(http.StatusUnprocessableEntity, err.Error())
		}
		if u != nil {
			u.Name = item.Name
			u.Detail = item.Detail
			u.WantToBeFound = item.WantToBeFound
			err = userService.Put(sess, u)
			if err != nil {
				return c.String(http.StatusUnprocessableEntity, err.Error())
			}
			if len(item.Password) >= 6 {
				err = userService.PutPassword(u.ID, item.Password)
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

// Returns the user object from the session
//
// @params => -
// @returns
//
//		200 => User => {}
//	 404 => Not found
func PutProfilePhotoHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusBadRequest)
	}

	err := userService.PutProfilePhoto(sess, sess.UserID(), c.Request().Body)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	return c.NoContent(http.StatusOK)
}

// Returns a profile photo
//
// @params id => string
// @returns
//
//		200 => File
//	 404 => Not found
//	}
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

	err := userService.GetProfilePhoto(sess, id, c.Response().Writer)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	c.Response().Committed = true
	c.Response().Header().Set("Content-Type", "image/jpeg")
	return c.NoContent(http.StatusOK)
}

// Check if a payment is required for current user for the workflow.
// Return http OK if no payment is required or if payment is required and a payment a matching payment
// with status = "confirmed" is found
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

	err := paymentService.CheckForWorkflowPayment(sess, workflowId)
	if err != nil {
		if db.NotFound(err) {
			return c.NoContent(http.StatusNotFound)
		}
		return c.NoContent(http.StatusBadRequest)
	}

	return c.NoContent(http.StatusOK)
}

var errNoPaymentFound = errors.New("no payment for workflow")

// Process a document
//
// @params ID => string
// @returns
//
//		200 => map[string]interface{}{"name" => string, "status" => string}
//	 400 => Bad request
//	 404 => Not found
//	 422 => Unprocessable entity
//	 500 => Server error
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
	wf, err = workflowService.GetPublished(sess, ID)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	docApp := documentService.GetDocApp(sess, ID)
	if docApp == nil {
		paymentRequired, err := paymentService.CheckIfWorkflowPaymentRequired(c.Session(false), ID)
		if err != nil {
			return c.String(http.StatusNotFound, err.Error())
		}

		if paymentRequired {
			sess := c.Session(false)
			user, err := userService.GetUser(sess)
			if err != nil {
				return c.NoContent(http.StatusBadRequest)
			}
			err = paymentService.RedeemPayment(wf.ID, user.EthereumAddr)
			if err != nil {
				log.Println("[redeemPayment] ", err.Error())
				return c.String(http.StatusUnprocessableEntity, errNoPaymentFound.Error())
			}
		}

		usrDataItem, _, err := userDocumentService.GetByWorkflow(sess, wf, false)
		if err != nil {
			if !db.NotFound(err) {
				return c.String(http.StatusNotFound, err.Error())
			}

			usrDataItem = &model.UserDataItem{
				WorkflowID: wf.ID,
				Name:       wf.Name,
				Detail:     wf.Detail,
			}
			err := userDocumentService.Put(sess, usrDataItem)
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}
		}

		docApp, err = app.NewDocumentApp(usrDataItem, sess, c.System(), ID, sess.GetSessionDir())
		if err != nil {
			return c.String(http.StatusUnprocessableEntity, err.Error())
		}
		sess.PutMemory("docApp_"+ID, docApp)
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

// Delete a document
//
// @params ID => string
// @returns
//
//		200 => OK
//	 400 => Bad request
func DocumentDeleteHandler(e echo.Context) error {
	c := e.(*www.Context)
	ID := c.Param("ID")
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusBadRequest)
	}

	err := documentService.Delete(sess, ID)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// Update a document
//
// @params ID => string, FormInput => interface{}
// @returns
//
//		200 => OK
//	 400 => Bad request
//	 422 => Unprocessable entity
func DocumentEditHandler(e echo.Context) error {
	c := e.(*www.Context)
	ID := c.Param("ID")
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnprocessableEntity)
	}
	formInput, err := helpers.ParseDataFromReq(c)
	if err != nil {
		return err
	}
	err = documentService.Edit(sess, ID, formInput)
	if err != nil {
		log.Println("[api][handlers] DocumentEditHandler Edit err: ", err.Error())
		if err == service.ErrUnableToEdit {
			return c.NoContent(http.StatusUnprocessableEntity)
		}
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.NoContent(http.StatusOK)
}

func getUserFromSession(s *sys.Session) (user *model.User) {
	if s == nil {
		return nil
	}
	err := s.Get("user", &user)
	if err != nil {
		if s.S.ID != "" {
			id := s.UserID()
			user, err = userService.GetById(s, id)
			if err != nil {
				return nil
			}
			s.Put("user", user)
		}
	}
	return user
}

// Move one step further in the document process
//
// @params 	ID => string
//
//	final => boolean,
//	FormInput => map[string]interface{}
//
// @returns
//
//		200 => map[string]interface{}
//	 400 => Bad request
//	 401 => Unauthorized
//	 422 => Unprocessable entity
func DocumentNextHandler(e echo.Context) error {
	c := e.(*www.Context)
	ID := c.Param("ID")
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	formInput, err := helpers.ParseDataFromReq(c)
	if err != nil {
		fmt.Println(err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	_, isFinal := c.QueryParams()["final"]
	docApp, status, err := documentService.Next(sess, ID, c.Lang(), formInput, isFinal)
	resData := map[string]interface{}{
		"status": status,
	}
	if docApp != nil {
		resData["id"] = docApp.DataID
	}
	if err != nil {
		if er, ok := err.(validate.ErrorMap); ok {
			er.Translate(func(key string, args ...string) string {
				return c.I18n().T(key, args)
			})
			resData["errors"] = er
			return c.JSON(http.StatusUnprocessableEntity, resData)
		}
		fmt.Println(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if isFinal {
		return c.JSON(http.StatusOK, map[string]interface{}{"id": docApp.DataID})
	}

	return c.JSON(http.StatusOK, resData)
}

// Move one step back in the document process
//
// @params 	ID => string
// @returns
//
//		200 => status => string
//	 400 => Bad request
func DocumentPrevHandler(e echo.Context) error {
	c := e.(*www.Context)
	ID := c.Param("ID")
	if ID == "" {
		return c.NoContent(http.StatusBadRequest)
	}
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusBadRequest)
	}

	status, err := documentService.Prev(sess, ID)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	resData := map[string]interface{}{
		"status": status,
	}

	return c.JSON(http.StatusOK, resData)
}

// Update a document with data
//
// @params 	ID => string
//
//	FormInput => map[string]interface{}
//
// @returns
//
//		200 => OK
//	 400 => Bad request
//	 422 => Unprocessable entity
//	 500 => Server error
func DocumentDataHandler(e echo.Context) error {
	c := e.(*www.Context)
	ID := c.Param("ID")
	sess := c.Session(false)
	inputData, err := helpers.ParseDataFromReq(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if sess == nil {
		return c.NoContent(http.StatusBadRequest)
	}

	verrs, err := documentService.Update(sess, ID, inputData)
	if err != nil {
		log.Println("[api][handler] DocumentDataHandler err: ", err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	if len(verrs) > 0 {
		verrs.Translate(func(key string, args ...string) string {
			return c.I18n().T(key, args)
		})
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"errors": verrs})
	}

	return c.NoContent(http.StatusOK)
}

// Returns a document's file output
//
// @params 	ID => string
//
//	inputName => string
//
// @returns
//
//		200 => map[string]interface{}
//	 404 => Not found
func DocumentFileGetHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)

	if sess == nil {
		return c.NoContent(http.StatusNotFound)
	}

	ID := c.Param("ID")
	finfo, err := documentService.GetFile(sess, ID, c.Param("inputName"))
	if err != nil || finfo == nil {
		return c.NoContent(http.StatusNotFound)
	}

	c.Response().Header().Set("Content-Type", finfo.ContentType())
	c.Response().Committed = true //prevents from-> http: multiple response.WriteHeader calls

	err = fileService.Read(finfo.Path(), c.Response().Writer)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.NoContent(http.StatusOK)
}

// Attach a file to a document
//
// @params 	@Fileinput
//
//	ID => string
//	inputName => string
func DocumentFilePostHandler(e echo.Context) error {
	c := e.(*www.Context)
	fieldname := c.Param("inputName")
	fileName, _ := url.QueryUnescape(c.Request().Header.Get("File-Name"))
	if fieldname == "" {
		return c.NoContent(http.StatusBadRequest)
	}
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusBadRequest)
	}
	ID := c.Param("ID")

	defer c.Request().Body.Close()
	finfo, verrs, err := documentService.UpdateFile(sess, ID, fieldname, fileName, c.Request().Header.Get("Content-Type"), c.Request().Body)
	if len(verrs) > 0 {
		verrs.Translate(func(key string, args ...string) string {
			return c.I18n().T(key, args)
		})
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"errors": verrs})
	}
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	if finfo == nil {
		return c.NoContent(http.StatusNoContent)
	}

	return c.JSON(http.StatusOK, finfo)
}

// Returns a preview of the document in PDF form
func DocumentPreviewHandler(e echo.Context) error {
	c := e.(*www.Context)
	ID := c.Param("ID")
	tmplID := c.Param("templateID")
	lang := c.Param("lang")
	format := c.Param("format")

	previewResponse, err := documentService.Preview(c.Session(false), ID, tmplID, lang, format)
	if err != nil {
		if os.IsNotExist(err) {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				return c.NoContent(http.StatusServiceUnavailable)
			}
			return c.NoContent(http.StatusBadRequest)
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	resp := c.Response()
	extension := previewResponse.File.NameWithExt(previewResponse.Format.String())
	contentDisposition := fmt.Sprintf("%s;filename=\"%s\"", "attachment", extension)
	resp.Header().Set("Content-Disposition", contentDisposition)
	resp.Header().Set("Content-Type", previewResponse.ContentType)
	resp.Header().Set("Content-Length", previewResponse.ContentLength)
	resp.Committed = true //prevents from-> http: multiple response.WriteHeader calls
	_, err = io.Copy(resp.Writer, previewResponse.Data)
	if err != nil {
		c.NoContent(http.StatusBadRequest)
	}
	defer previewResponse.Data.Close()

	return c.NoContent(http.StatusOK)
}

// Returns a list of paginated user documents
func UserDocumentListHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	contains := c.QueryParam("c")
	settings := helpers.RequestOptions(c)
	items, err := userDocumentService.List(sess, contains, settings)
	if err != nil || items == nil {
		return c.NoContent(http.StatusBadRequest)
	}
	return c.JSON(http.StatusOK, items)
}

// Returns a documents metadata
func UserDocumentGetHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	id := c.Param("ID")
	item, err := userDocumentService.Get(sess, id)
	if err != nil || item == nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, item)
}

// Returns a signature request for the current user
func UserDocumentSignatureRequestGetCurrentUserHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	requests, err := signatureService.GetForCurrentUser(sess)
	if err != nil {
		if os.IsNotExist(err) {
			return c.NoContent(http.StatusNotFound)
		}
		log.Println("[UserDocumentSignatureRequestAddHandler] signatureService.GetForCurrentUser err: ", err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, requests)
}

// Remove a user from the database
func UserDeleteHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	err := userService.DeleteUser(sess)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return LogoutHandler(c)
}

// Returns a signature request based on a provided document ID
func UserDocumentSignatureRequestGetByDocumentIDHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	docId := c.Param("docID")
	id := c.Param("ID")

	requests, err := signatureService.GetById(id, docId)
	if err != nil {
		if os.IsNotExist(err) {
			return c.NoContent(http.StatusNotFound)
		}
		log.Println("UserDocumentSignatureRequestGetByDocumentIDHandler signatureService.GetById err: ", err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, requests)
}

// Add a signature request for a document
func UserDocumentSignatureRequestAddHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	docId := c.Param("docID")
	id := c.Param("ID")

	signatory := c.FormValue("signatory")

	err := signatureService.AddAndNotify(sess, c.I18n(), id, docId, signatory, c.Request().Host, c.Scheme())
	if err != nil {
		if os.IsNotExist(err) {
			return c.NoContent(http.StatusNotFound)
		}
		if errors.Is(err, service.ErrSignatureRequestAlreadyExists) {
			return c.String(http.StatusConflict, c.I18n().T(err.Error()))
		}
		log.Println("[UserDocumentSignatureRequestAddHandler] signatureService.AddAndNotify err: ", err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// Reject a signature request for a document
func UserDocumentSignatureRequestRejectHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	docId := c.Param("docID")
	id := c.Param("ID")

	err := signatureService.RejectAndNotify(sess, c.I18n(), id, docId, c.Request().Host)
	if err != nil {
		if os.IsNotExist(err) {
			return c.NoContent(http.StatusNotFound)
		}
		log.Println("UserDocumentSignatureRequestGetByDocumentIDHandler signatureService.RejectAndNotify err: ", err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// Revoke a documents signature
func UserDocumentSignatureRequestRevokeHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	docId := c.Param("docID")
	id := c.Param("ID")
	signatory := c.FormValue("signatory")

	err := signatureService.RevokeAndNotify(sess, c.I18n(), id, docId, signatory, c.Request().Host, c.Scheme())
	if err != nil {
		if os.IsNotExist(err) {
			return c.NoContent(http.StatusNotFound)
		}
		log.Println("[UserDocumentSignatureRequestAddHandler] signatureService.RevokeAndNotify err: ", err.Error())
		return c.String(http.StatusBadRequest, err.Error())
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

// Returns a file belonging to a users document
func UserDocumentFileHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	format := c.QueryParam("format")

	dataPath := c.Param("dataPath")
	id := c.Param("ID")

	inlineOrAttachment := "attachment"
	if _, ok := c.QueryParams()["inline"]; ok {
		inlineOrAttachment = "inline"
	}

	resp := c.Response()

	if strings.HasPrefix(dataPath, "docs") {
		fileHeaderResponse, filePath, err := userDocumentService.GetDocFile(sess, id, dataPath, inlineOrAttachment)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		setResponseHeader(resp, fileHeaderResponse)
		err = fileService.Read(filePath, resp.Writer)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
	} else {
		fileHeaderResponse, dsRespBody, err := userDocumentService.GetTemplateWithFormatFile(sess, id, dataPath, format, inlineOrAttachment)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		setResponseHeader(resp, fileHeaderResponse)
		defer dsRespBody.Close()
		_, err = io.Copy(resp.Writer, dsRespBody)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
	}

	return c.NoContent(http.StatusOK)
}

func setResponseHeader(resp *echo.Response, fileHeaderResponse *service.FileHeaderResponse) {
	resp.Header().Set("Content-Type", fileHeaderResponse.ContentType)
	resp.Header().Set("Content-Length", fileHeaderResponse.ContentLength)
	resp.Header().Set("Content-Pages", fileHeaderResponse.ContentPages)
	resp.Header().Set("Content-Disposition", fileHeaderResponse.ContentDisposition)
	resp.Committed = true //prevents from-> http: multiple response.WriteHeader calls
}

// Update a user's metadata
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
			exItem, err := userService.GetById(sess, item.ID)
			if err != nil {
				return c.NoContent(http.StatusNotFound)
			}
			exItem.Name = item.Name
			exItem.Detail = item.Detail
			exItem.Role = item.Role
			err = userService.Put(sess, exItem)
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

// Returns a user by ID
func AdminUserGetHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	itemID := c.Param("ID")
	item, err := userService.GetById(sess, itemID)
	if err != nil {
		return err
	}
	if item == nil {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, item)
}

// Lists all users
func AdminUserListHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	contains := c.QueryParam("c")
	settings := helpers.RequestOptions(c)
	dat, err := userService.List(sess, contains, settings)
	if err != nil || dat == nil {
		if err == model.ErrAuthorityMissing {
			return c.NoContent(http.StatusUnauthorized)
		}
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, dat)
}

// Returns a workflows schema
func WorkflowSchema(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	id := c.Param("ID")
	workflowItem, fieldsAndRules, err := documentService.GetWorkflowSchema(sess, id)
	if err != nil {
		log.Println("[apiHandler][WorkflowSchema] GetWorkflowSchema err: ", err.Error())
		return c.NoContent(http.StatusNotFound)
	}
	wfDetails := &struct {
		*model.WorkflowItem
		Data interface{} `json:"data"`
	}{workflowItem, fieldsAndRules}
	result := &struct {
		Workflow interface{} `json:"workflow"`
	}{Workflow: wfDetails}

	return c.JSON(http.StatusOK, result)
}

// Execute a workflow at once
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
	wItem, err := workflowService.Get(sess, id)
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

// Remove the users' auth session
//
//	@params => {
//	  name => string,
//	  ID => string,
//	}
//
// @returns
//
//		200 => apiKey string
//	 400 => Data layer error/Validation error
//	 401 => Unauthorized
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
	apiKey, err := apiService.CreateApiKey(sess, id, name)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if id == sess.UserID() {
		sess.Delete("user")
	}
	return c.String(http.StatusOK, apiKey)
}

// Remove the users' auth session
//
//	@params => {
//	  ID => string,
//	}
//
// @returns
//
//		200 => apiKey string
//	 400 => Data layer error
//	 401 => Unauthorized
func DeleteApiKeyHandler(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	id := c.Param("ID")
	hiddenApiKey := c.QueryParam("hiddenApiKey")
	err := apiService.DeleteApiKey(sess, id, hiddenApiKey)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if id == sess.UserID() {
		sess.Delete("user")
	}
	return c.NoContent(http.StatusOK)
}

// Instantiate an external node and return the config URL
func ExternalConfigurationPage(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	id := c.Param("id")
	name := c.Param("name")

	externalNodeQuery, err := workflowService.InstantiateExternalNode(sess, id, name)
	if err != nil {
		return err
	}
	if externalNodeQuery == nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.Redirect(http.StatusFound, externalNodeQuery.ConfigUrl())
}

// Register an external node
//
//	@param node => {
//			ID     string `json:"id" storm:"id"`
//			Name   string `json:"name"`
//			Detail string `json:"detail"`
//			Url    string `json:"url"`
//			Secret string `json:"secret"`
//		}
func ExternalRegister(e echo.Context) error {
	c := e.(*www.Context)
	var node externalnode.ExternalNode
	err := c.Bind(&node)
	if err != nil {
		return err
	}
	return nodeService.RegisterExternalNode(new(model.User), &node)
}

// List external nodes
func ExternalList(e echo.Context) error {
	c := e.(*www.Context)
	nodes := nodeService.ListExternalNodes()

	return c.JSON(http.StatusOK, nodes)
}

// Update the config of an external store
func ExternalConfigStore(e echo.Context) error {
	c := e.(*www.Context)

	var node externalnode.ExternalNodeInstance
	err := c.Bind(&node)
	if err != nil {
		return err
	}

	//QueryFromInstanceID -> instance
	q, err := nodeService.QueryFromInstanceID(new(model.User), node.ID)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	//Add config to instance
	q.Config = node.Config

	//PutExternalNodeInstance
	err = nodeService.PutExternalNodeInstance(new(model.User), &q)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.NoContent(http.StatusOK)
}

// Returns the config of an external node
func ExternalConfigRetrieve(e echo.Context) error {
	c := e.(*www.Context)
	id := c.Param("id")
	q, err := nodeService.QueryFromInstanceID(new(model.User), id)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if q.Config == nil {
		c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, q.Config)
}
