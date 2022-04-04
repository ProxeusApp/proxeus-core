package service

import "github.com/ProxeusApp/proxeus-core/sys/model"

type (

	// ApiService is an interface that provides api key functions
	ApiService interface {
		// CreateApiKey creates and returns a new api key
		CreateApiKey(auth model.Auth, userId, apiKeyName string) (string, error)

		// DeleteApiKey removes an existing API key
		DeleteApiKey(auth model.Auth, userId, hiddenApiKey string) error

		// AuthenticateWithApiKey tries to authenticate the user with the supplied API key and returns the user object or an error
		AuthenticateWithApiKey(apiKey string) (*model.User, error)

		// AuthenticateApiKeyForUser tries to authenticate with supplied API key and user ID
		AuthenticateApiKeyForUser(apiKey string, userId string) (*model.User, error)
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

// AuthenticateWithApiKey tries to authenticate the user with the supplied API key and returns the user object or an error
func (me *defaultApiService) AuthenticateWithApiKey(apiKey string) (*model.User, error) {
	return userDB().GetByApiKey(apiKey, "")
}

// AuthenticateApiKeyForUser tries to authenticate with supplied API key and user ID
func (me *defaultApiService) AuthenticateApiKeyForUser(apiKey string, userId string) (*model.User, error) {
	return userDB().GetByApiKey(apiKey, userId)
}
