package test

import (
	"net/http"
	"testing"
)

const exampleKey = "API keys"

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

	s.e.GET("/api/i18n/search").WithJSON([]string{exampleKey}).Expect().Status(http.StatusOK).
		JSON().Object().Keys().Length().Ge(1)
	s.e.GET("/api/i18n/search").WithJSON(exampleKey).Expect().Status(http.StatusOK).
		JSON().Object().Keys().Length().Ge(1)

	s.e.POST("/api/admin/i18n/translate").WithJSON([]string{exampleKey}).Expect().Status(http.StatusOK).
		JSON().Array().First().String().Equal(exampleKey)
	s.e.POST("/api/admin/i18n/translate").WithJSON(exampleKey).Expect().Status(http.StatusOK).
		JSON().String().Equal(exampleKey)

	deleteUser(s, u)
}

func TestI18nAdmin(t *testing.T) {
	s := new(t, serverURL)
	u := registerSuperAdmin(s)
	login(s, u)

	{
		r := s.e.GET("/api/i18n/export").Expect().Status(http.StatusOK)
		r.Body().Length().Gt(100)
		r.Headers().ContainsKey("Content-Disposition")
	}

	r := s.e.GET("/api/admin/i18n/find").WithQuery("k", exampleKey).
		Expect().Status(http.StatusOK).JSON().Object()
	r.Keys().Length().Equal(1)

	lang := r.Value(exampleKey).Object().Keys().First().String().Raw()
	value := r.Value(exampleKey).Object().Value(lang).String().Raw()

	s.e.POST("/api/admin/i18n/update").
		WithJSON(map[string]map[string]string{exampleKey: {lang: value}}).
		Expect().Status(http.StatusOK)

	s.e.POST("/api/admin/i18n/lang").WithJSON(map[string]bool{lang: true}).
		Expect().Status(http.StatusOK)

	fallback := s.e.GET("/api/admin/i18n/meta").Expect().Status(http.StatusOK).
		JSON().Path("$.langFallback").String().Raw()

	s.e.POST("/api/admin/i18n/fallback").WithQuery("lang", fallback).
		Expect().Status(http.StatusOK)

	// TODO: should be possible to delete admin too
	//deleteUser(s, u)
}
