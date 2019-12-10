package test

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/ProxeusApp/proxeus-core/main/handlers/api"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/gavv/httpexpect.v2"
)

var serverURL string

type user struct {
	uuid     string
	username string
	password string
}

type session struct {
	id string
	t  *testing.T
	e  *httpexpect.Expect
}

func init() {
	serverURL = os.Getenv("PROXEUS_URL")
}

func new(t *testing.T, serverURL string) *session {
	return &session{
		id: uuid.NewV4().String(),
		t:  t,
		e:  httpexpect.New(t, serverURL),
	}
}

func TestUser(t *testing.T) {
	s := new(t, serverURL)
	u := registerTestUser(s)
	login(s, u)
	logout(s)
	login(s, u)
	deleteUser(s, u)
}

func registerTestUser(s *session) *user {
	// Register test user
	u := &user{
		username: fmt.Sprintf("test%s@example.com", s.id),
		password: s.id,
	}

	s.t.Logf("Starting test %s", s.id)
	s.t.Logf("User %s %s", u.username, u.password)

	tr := &api.TokenRequest{
		Email: u.username,
	}

	r := s.e.POST("/api/register").WithJSON(tr).Expect()

	r.Status(http.StatusOK)
	r.Header("X-Test-Token").NotEmpty() // This is only true in TESTMODE
	registrationToken := r.Header("X-Test-Token").Raw()

	p := &struct {
		Password string `json:"password"`
	}{
		Password: u.password,
	}

	r = s.e.POST("/api/register/" + registrationToken).WithJSON(p).
		Expect().
		Status(http.StatusOK)

	return u
}

func login(s *session, u *user) {
	l := &struct {
		Email    string `json:"email" form:"email"`
		Password string `json:"password" form:"password"`
	}{
		Email:    u.username,
		Password: u.password,
	}
	s.e.POST("/api/login").WithJSON(l).Expect().Status(http.StatusOK)

	me := s.e.GET("/api/me").Expect().Status(http.StatusOK).JSON().Object()
	me.ValueEqual("email", u.username)

	u.uuid = me.Value("id").String().Raw()
}

func logout(s *session) {
	s.e.POST("/api/logout").Expect().Status(http.StatusOK)
	s.e.GET("/api/me").Expect().Status(http.StatusNotFound)
}

func deleteUser(s *session, u *user) {
	s.e.POST("/api/user/delete").Expect().Status(http.StatusOK)

	l := &struct {
		Email    string `json:"email" form:"email"`
		Password string `json:"password" form:"password"`
	}{
		Email:    u.username,
		Password: u.password,
	}
	s.e.POST("/api/login").WithJSON(l).Expect().Status(http.StatusBadRequest)
}
