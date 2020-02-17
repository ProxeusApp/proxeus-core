package service

import (
	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/sys/model"
)

type (
	UserService interface {
		GetUser(auth model.Auth) (*model.User, error)
	}
	defaultUserService struct {
		userDB storage.UserIF
	}
)

func NewUserService(userDB storage.UserIF) *defaultUserService {
	return &defaultUserService{userDB: userDB}
}

func (me *defaultUserService) GetUser(auth model.Auth) (*model.User, error) {
	return me.userDB.Get(auth, auth.UserID())
}
