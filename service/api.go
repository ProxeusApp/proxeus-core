package service

import "github.com/ProxeusApp/proxeus-core/sys/model"

type (
	ApiService interface {
		CreateApiKey(auth model.Auth, userId, apiKeyName string) (string, error)
		DeleteApiKey(auth model.Auth, userId, hiddenApiKey string) error
	}
	defaultApiService struct {
	}
)

func NewApiService() *defaultApiService {
	return &defaultApiService{}
}

// CreateApiKey creates and returns a new api key
func (me *defaultApiService) CreateApiKey(auth model.Auth, userId, apiKeyName string) (string, error) {
	return userDB().CreateApiKey(auth, userId, apiKeyName)
}

// DeleteApiKey removes an existing API key
func (me *defaultApiService) DeleteApiKey(auth model.Auth, userId, hiddenApiKey string) error {
	return userDB().DeleteApiKey(auth, userId, hiddenApiKey)
}
