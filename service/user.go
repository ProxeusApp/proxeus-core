package service

import (
	"github.com/ProxeusApp/proxeus-core/main/www"
	"github.com/ProxeusApp/proxeus-core/sys/model"
)

type (
	UserService interface {
		GetUser(c *www.Context) (*model.User, error)
	}
	defaultUserService struct {
	}
)

func NewUserService() *defaultUserService {
	return &defaultUserService{}
}

func (me *defaultUserService) GetUser(c *www.Context) (*model.User, error) {
	sess := c.Session(false)
	return c.System().DB.User.Get(sess, sess.UserID())
}
