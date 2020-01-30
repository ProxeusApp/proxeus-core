package test

import (
	"encoding/base64"
	"net/http"
)

func testApiKey(s *session) {
	u := registerTestUser(s)

	login(s, u)
	apiKey, summary := createApiKey(s, u, "test-"+s.id)
	logout(s)

	token := getSessionToken(s, u.username, apiKey)
	deleteSessionToken(s, token)

	login(s, u)
	deleteApiKey(s, u, summary)
	deleteUser(s, u)
}

func createApiKey(s *session, u *user, name string) (string, string) {
	key := s.e.GET("/api/user/create/api/key/{id}").WithPath("id", u.uuid).WithQuery("name", name).Expect().Status(http.StatusOK).Body().Raw()

	summary := key[:4] + "..." + key[len(key)-4:]
	s.e.GET("/api/me").Expect().Status(http.StatusOK).JSON().Path("$..Key").Array().Contains(summary)

	return key, summary
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
