package service

import (
	"github.com/ProxeusApp/proxeus-core/sys"
	"github.com/ProxeusApp/proxeus-core/sys/model"
)

type (
	UserService interface {
		GetUser(auth model.Auth) (*model.User, error)
		GetById(auth model.Auth, id string) (*model.User, error)
		GetUserDataById(auth model.Auth, id string) (*model.UserDataItem, error)
		CreateApiKeyHandler(auth model.Auth, userId, apiKeyName string) (string, error)
		DeleteApiKey(auth model.Auth, userId, hiddenApiKey string) error
		DeleteUserData(auth model.Auth, id string) error
	}
	defaultUserService struct {
		*baseService
	}
)

func NewUserService(system *sys.System) *defaultUserService {
	return &defaultUserService{baseService: &baseService{system: system}}
}

func (me *defaultUserService) GetUser(auth model.Auth) (*model.User, error) {
	return me.userDB().Get(auth, auth.UserID())
}
func (me *defaultUserService) GetById(auth model.Auth, id string) (*model.User, error) {
	return me.userDB().Get(auth, id)
}

func (me *defaultUserService) GetUserDataById(auth model.Auth, id string) (*model.UserDataItem, error) {
	return me.userDataDB().Get(auth, id)
}

func (me *defaultUserService) DeleteUserData(auth model.Auth, id string) error {
	return me.userDataDB().Delete(auth, me.filesDB(), id)
}

func (me *defaultUserService) CreateApiKeyHandler(auth model.Auth, userId, apiKeyName string) (string, error) {
	return me.userDB().CreateApiKey(auth, userId, apiKeyName)
}

func (me *defaultUserService) DeleteApiKey(auth model.Auth, userId, hiddenApiKey string) error {
	return me.userDB().DeleteApiKey(auth, userId, hiddenApiKey)
}
