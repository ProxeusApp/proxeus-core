package test

import (
	"fmt"
	"net/http"
	"testing"
)

func TestUser(t *testing.T) {
	s := new(t, serverURL)
	u := registerTestUser(s)
	login(s, u)
	logout(s)
	login(s, u)
	deleteUser(s, u)
}

func TestUserProfile(t *testing.T) {
	s := new(t, serverURL)
	u := registerTestUser(s)
	login(s, u)

	s.e.GET("/api/config").Expect().Status(http.StatusOK).
		JSON().Object().ContainsKey("blockchainNet")

	resetPassword(s, u)
	changeEmail(s, u)
	setProfilePhoto(s)

	deleteUser(s, u)
}

func resetPassword(s *session, u *user) {
	token := s.e.POST("/api/reset/password").WithJSON(map[string]string{"email": u.username}).
		Expect().Status(http.StatusOK).Header("X-Test-Token").NotEmpty().Raw()

	const newPass = "abcd123"
	s.e.POST("/api/reset/password/{token}").WithPath("token", token).
		WithJSON(map[string]string{"password": newPass}).Expect().Status(http.StatusOK)
	u.password = newPass

	logout(s)
	login(s, u)
}

func changeEmail(s *session, u *user) {
	newEmail := fmt.Sprintf("test%s@example2.com", s.id)
	token := s.e.POST("/api/change/email").WithJSON(map[string]string{"email": newEmail}).
		Expect().Status(http.StatusOK).Header("X-Test-Token").NotEmpty().Raw()
	u.username = newEmail

	s.e.POST("/api/change/email/{token}").WithPath("token", token).
		Expect().Status(http.StatusOK)

	logout(s)
	login(s, u)
}

func setProfilePhoto(s *session) {
	const photo = "ph-data"
	s.e.POST("/api/my/profile/photo").WithBytes([]byte(photo)).Expect().Status(http.StatusOK)
	s.e.GET("/api/my/profile/photo").Expect().Status(http.StatusOK).Body().Equal(photo)
}
