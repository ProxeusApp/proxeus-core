package test

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/ProxeusApp/proxeus-core/sys/model"

	"github.com/ethereum/go-ethereum/crypto"
)

type user struct {
	uuid     string
	username string
	password string

	ethPrivateKey *ecdsa.PrivateKey
	EthereumAddr  string
}

func registerSuperAdmin(s *session) *user {
	return registerTestUserWithRole(s, model.SUPERADMIN)
}

func registerTestUser(s *session) *user {
	return registerTestUserWithRole(s, 0)
}

func registerTestUserWithRole(s *session, role model.Role) *user {
	// Register test user
	u := &user{
		username: fmt.Sprintf("test%s@example.com", s.id),
		password: s.id,
	}

	s.t.Logf("Starting test %s", s.id)
	s.t.Logf("User %s %s", u.username, u.password)

	tr := &model.TokenRequest{
		Email: u.username,
	}

	if role > 0 {
		tr.Role = role
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

func setEthKey(s *session, u *user) {
	challenge := s.e.GET("/api/challenge").Expect().Status(http.StatusOK).Body().Raw()
	var err error
	u.ethPrivateKey, err = crypto.GenerateKey()
	if err != nil {
		s.t.Error(err)
	}

	sig, err := crypto.Sign(signHash(challenge), u.ethPrivateKey)
	if err != nil {
		s.t.Error(err)
	}
	sig[64] += 27

	s.e.POST("/api/change/bcaddress").WithJSON(
		struct {
			Signature string `json:"signature"`
		}{Signature: "0x" + hex.EncodeToString(sig)}).Expect().Status(http.StatusOK)

	addr := s.e.GET("/api/me").Expect().Status(http.StatusOK).JSON().
		Path("$.etherPK").String()
	addr.Length().Gt(10)
	u.EthereumAddr = addr.Raw()

	s.e.POST("/api/me").WithJSON(
		struct {
			ID            string `json:"id"`
			Email         string `json:"email"`
			WantToBeFound bool   `json:"wantToBeFound"`
		}{
			ID:            u.uuid,
			Email:         u.username,
			WantToBeFound: true,
		}).Expect().Status(http.StatusOK)

	s.e.GET("/api/admin/user/list").WithQuery("c", u.EthereumAddr).Expect().
		Status(http.StatusOK).JSON().Array().Length().Equal(1)
}

func signHash(data string) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
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
