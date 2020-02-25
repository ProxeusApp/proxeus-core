package service

import (
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
	}
)

func NewUserService() *defaultUserService {
	return &defaultUserService{}
}

func (me *defaultUserService) GetUser(auth model.Auth) (*model.User, error) {
	return userDB().Get(auth, auth.UserID())
}
func (me *defaultUserService) GetById(auth model.Auth, id string) (*model.User, error) {
	return userDB().Get(auth, id)
}

func (me *defaultUserService) GetUserDataById(auth model.Auth, id string) (*model.UserDataItem, error) {
	return userDataDB().Get(auth, id)
}

func (me *defaultUserService) DeleteUserData(auth model.Auth, id string) error {
	return userDataDB().Delete(auth, filesDB(), id)
}

func (me *defaultUserService) CreateApiKeyHandler(auth model.Auth, userId, apiKeyName string) (string, error) {
	return userDB().CreateApiKey(auth, userId, apiKeyName)
}

func (me *defaultUserService) DeleteApiKey(auth model.Auth, userId, hiddenApiKey string) error {
	return userDB().DeleteApiKey(auth, userId, hiddenApiKey)
}
