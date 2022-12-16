package service

import (
	"fmt"
	"log"
	"os"

	uuid "github.com/satori/go.uuid"

	"github.com/ProxeusApp/proxeus-core/main/handlers/blockchain"
	"github.com/ProxeusApp/proxeus-core/main/handlers/helpers"
	"github.com/ProxeusApp/proxeus-core/main/www"
	"github.com/ProxeusApp/proxeus-core/storage/database/db"
	"github.com/ProxeusApp/proxeus-core/sys/model"
)

type (

	// AuthenticationService is an interface that provides user authentication functions
	AuthenticationService interface {

		// LoginWithUsernamePassword performs a login with username and password
		LoginWithUsernamePassword(email, password string) (*model.User, error)

		// LoginWithWallet logs the user in when the login is done with a signature (e.g. Metamask)
		LoginWithWallet(challenge, signature string) (bool, *model.User, error)

		// ChangeEmail Changes the email of a user
		ChangeEmail(tokenID string) (*model.TokenRequest, error)

		// ResetPasswordRequest sends an email request to the user to reset the password
		ResetPasswordRequest(translator www.Translator, scheme, host string, m *model.TokenRequest) (model.TokenRequest, error)

		// ResetPassword resets the users password
		ResetPassword(translator www.Translator, password, tokenID string) (error, map[string]interface{})

		// RegisterRequest sends an email request to the user to register
		RegisterRequest(translator www.Translator, scheme, host string, m *model.TokenRequest) (model.TokenRequest, error)

		// Register registers the user on the platform
		Register(translator www.Translator, tokenID, password string) (error, map[string]interface{})

		// PutTokenRequest saves the tokenRequest in the sessionDB
		PutTokenRequest(token *model.TokenRequest) error
	}
	defaultAuthenticationService struct {
		userService     UserService
		workflowService WorkflowService
		emailService    EmailService
	}
)

var rootRole = &model.User{Role: model.ROOT}

func NewAuthenticationService(userS UserService, workflowS WorkflowService, emailS EmailService) *defaultAuthenticationService {
	return &defaultAuthenticationService{userService: userS, workflowService: workflowS, emailService: emailS}
}

// LoginWithUsernamePassword performs a login with username and password
func (me *defaultAuthenticationService) LoginWithUsernamePassword(email, password string) (*model.User, error) {
	return userDB().Login(email, password)
}

// LoginWithWallet logs the user in when the login is done with a signature (e.g. Metamask)
func (me *defaultAuthenticationService) LoginWithWallet(challenge, signature string) (bool, *model.User, error) {
	created := false
	var address string
	var err error
	address, err = blockchain.VerifySignInChallenge(challenge, signature)
	if err != nil {
		return false, nil, err
	}
	var usr *model.User
	usr, err = me.userService.GetByBCAddress(address)
	if db.NotFound(err) {
		stngs, err := settingsDB().Get()
		if err != nil {
			return false, nil, err
		}
		it := &model.User{
			EthereumAddr: address,
			Role:         model.StringToRole(stngs.DefaultRole),
		}
		it.Name = "created by ethereum account sign"
		err = userDB().Put(rootRole, it)
		if err != nil {
			return false, nil, err
		}
		created = true
		usr, err = me.userService.GetByBCAddress(address)
		if err == nil {
			me.workflowService.CopyWorkflows(rootRole, usr)
			if stngs.BlockchainNet == "goerli" && stngs.AirdropEnabled == "true" {
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

// ChangeEmail Changes the email of a user
func (me *defaultAuthenticationService) ChangeEmail(tokenID string) (*model.TokenRequest, error) {
	tokenRequest, err := sessionDB().GetTokenRequest(model.TokenChangeEmail, tokenID)
	if err != nil {
		return nil, err
	}
	err = userDB().UpdateEmail(tokenRequest.UserID, tokenRequest.Email)
	if err != nil {
		return tokenRequest, err
	}
	err = sessionDB().DeleteTokenRequest(tokenRequest)
	return tokenRequest, err
}

// ResetPasswordRequest sends an email request to the user to reset the password
func (me *defaultAuthenticationService) ResetPasswordRequest(translator www.Translator, scheme, host string, m *model.TokenRequest) (model.TokenRequest, error) {

	var token model.TokenRequest
	usr, err := me.userService.GetByEmail(m.Email)
	if err != nil {
		// always return ok if provided email was valid
		// otherwise public users can test what email accounts exist
		return token, nil
	}

	token.Email = m.Email
	token.Token = uuid.NewV4().String()
	token.UserID = usr.ID
	token.Type = model.TokenResetPassword
	err = sessionDB().PutTokenRequest(&token)
	if err != nil {
		return token, err
	}

	subject := translator.T("Reset Password")
	body := fmt.Sprintf(
		"Hi %s,\n\nif you requested a password reset, please go on and click on this link to reset your password\n%s\n\nIf you didn't request it, please ignore this email.\n\nProxeus",
		usr.Name,
		helpers.AbsoluteURLWithScheme(scheme, host, "/reset/password/", token.Token),
	)

	err = me.emailService.Send(m.Email, subject, body)

	return token, err
}

// ResetPassword resets the users password
func (me *defaultAuthenticationService) ResetPassword(translator www.Translator, password, tokenID string) (error, map[string]interface{}) {
	errors := map[string]interface{}{}
	if len(password) < 6 {
		passwordError := []map[string]interface{}{{"msg": translator.T("Password not strong enough")}}
		errors["password"] = passwordError
		return os.ErrInvalid, errors
	}
	r, err := sessionDB().GetTokenRequest(model.TokenResetPassword, tokenID)
	if err != nil {
		return err, errors
	}
	err = userDB().PutPw(r.UserID, password)
	if err != nil {
		return err, errors
	}
	return sessionDB().DeleteTokenRequest(r), errors
}

// RegisterRequest sends an email request to the user to register
func (me *defaultAuthenticationService) RegisterRequest(translator www.Translator, scheme, host string, m *model.TokenRequest) (model.TokenRequest, error) {
	var token model.TokenRequest

	if usr, _ := me.userService.GetByEmail(m.Email); usr != nil {
		// always return ok if provided email was valid
		// otherwise public users can test what email accounts exist
		return token, nil
	}

	settings, err := settingsDB().Get()

	token.Email = m.Email
	token.Token = uuid.NewV4().String()
	token.Role = model.StringToRole(settings.DefaultRole)
	token.Type = model.TokenRegister
	if settings.TestMode == "true" && m.Role > 0 {
		token.Role = m.Role
	}
	err = sessionDB().PutTokenRequest(&token)
	if err != nil {
		return token, err
	}

	subject := translator.T("Register")
	body := fmt.Sprintf(
		"Hi there,\n\nplease proceed with your registration by visiting this link:\n%s\n\nIf you didn't request this, please ignore this email.\n\nProxeus",
		helpers.AbsoluteURLWithScheme(scheme, host, "/register/", token.Token),
	)
	err = me.emailService.Send(m.Email, subject, body)
	return token, err
}

// Register registers the user on the platform
func (me *defaultAuthenticationService) Register(translator www.Translator, tokenID, password string) (error, map[string]interface{}) {
	errors := map[string]interface{}{}
	if len(password) < 6 {
		errors["password"] = []map[string]interface{}{{"msg": translator.T("Password not strong enough")}}
		return nil, errors
	}
	r, err := sessionDB().GetTokenRequest(model.TokenRegister, tokenID)
	if err != nil {
		return err, errors
	}
	newUser := &model.User{Email: r.Email, Role: r.Role}
	err = userDB().Put(rootRole, newUser)
	if err != nil {
		return err, errors
	}

	me.workflowService.CopyWorkflows(rootRole, newUser)

	err = userDB().PutPw(newUser.ID, password)
	if err != nil {
		return err, errors
	}
	err = sessionDB().DeleteTokenRequest(r)
	return err, errors
}

// PutTokenRequest saves the tokenRequest in the sessionDB
func (me *defaultAuthenticationService) PutTokenRequest(token *model.TokenRequest) error {
	return sessionDB().PutTokenRequest(token)
}
