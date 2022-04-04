package test

import (
	"encoding/base64"
	"net/http"

	"github.com/ProxeusApp/proxeus-core/sys/model"
)

func testApiKey(s *session) {
	u := registerTestUser(s)

	login(s, u)
	apiKey, hashed := createApiKey(s, u, "test-"+s.id)
	logout(s)

	token := getSessionToken(s, u.username, apiKey)
	deleteSessionToken(s, token)

	login(s, u)
	deleteApiKey(s, u, hashed)
	deleteUser(s, u)
}

func createApiKey(s *session, u *user, name string) (string, string) {
	key := s.e.GET("/api/user/create/api/key/{id}").WithPath("id", u.uuid).WithQuery("name", name).Expect().Status(http.StatusOK).Body().Raw()

	tmpHashedKey := &model.ApiKey{
		Name: "",
		Key:  key,
	}
	// tmpHashedKey.HideKey()
	s.e.GET("/api/me").Expect().Status(http.StatusOK).JSON().Path("$..Key").Array().Contains(tmpHashedKey.Key)

	return key, tmpHashedKey.Key
}

func getSessionToken(s *session, username, apiKey string) string {
	b := base64.StdEncoding.EncodeToString([]byte(username + ":" + apiKey))
	return s.e.GET("/api/session/token").WithHeader("Authorization", "Basic "+b).Expect().Status(http.StatusOK).JSON().Object().Value("token").String().NotEmpty().Raw()
}

func deleteSessionToken(s *session, token string) {
	s.e.DELETE("/api/session/token").WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).NoContent()
}

func deleteApiKey(s *session, u *user, summary string) {
	s.e.DELETE("/api/user/create/api/key/{id}").WithPath("id", u.uuid).WithQuery("hiddenApiKey", summary).Expect().Status(http.StatusOK)
	s.e.GET("/api/me").Expect().Status(http.StatusOK).JSON().Path("$..Key").Array().NotContains(summary)
}
