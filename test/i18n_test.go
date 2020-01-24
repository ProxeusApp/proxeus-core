package test

import (
	"net/http"
	"testing"
)

func TestI18n(t *testing.T) {
	s := new(t, serverURL)
	u := registerTestUser(s)
	login(s, u)

	s.e.GET("/admin").Expect().Status(http.StatusOK)

	s.e.GET("/api/admin/i18n/meta").Expect().Status(http.StatusOK).
		JSON().Path("$.langListSize").Equal(1)

	s.e.GET("/api/admin/i18n/all").Expect().Status(http.StatusOK).
		JSON().Object().Keys().Length().Gt(100)

	s.e.GET("/api/i18n/search").WithQuery("c", "login").Expect().Status(http.StatusOK).
		JSON().Object().Keys().Length().Ge(1)

	s.e.GET("/api/i18n/search").WithJSON([]string{"API keys"}).Expect().Status(http.StatusOK).
		JSON().Object().Keys().Length().Ge(1)
	s.e.GET("/api/i18n/search").WithJSON("API keys").Expect().Status(http.StatusOK).
		JSON().Object().Keys().Length().Ge(1)

	deleteUser(s, u)
}
