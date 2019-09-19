package test

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

type template struct {
	permissions
	ID      string                 `json:"id" storm:"id"`
	Name    string                 `json:"name" storm:"index"`
	Detail  string                 `json:"detail"`
	Updated time.Time              `json:"updated" storm:"index"`
	Created time.Time              `json:"created" storm:"index"`
	Data    map[string]interface{} `json:"data"`
}

func TestTemplate(t *testing.T) {
	s := new(t, serverURL)
	u := registerTestUser(s)
	login(s, u)
	t1 := createSimpleTemplate(s, u, "template1-"+s.id, "test/assets/test_template.odt")
	t2 := createTemplate(s, u, "template2-"+s.id)

	deleteTemplate(s, t1.ID, false)
	deleteTemplate(s, t2.ID, true)
}

func createSimpleTemplate(s *session, u *user, name, path string) *template {

	fileContent, err := Asset(path)
	if err != nil {
		s.t.Errorf("Cannot upload asset %s", err)
	}

	t := createTemplate(s, u, name)
	uploadTemplateFile(s, t, "en", fileContent, "application/vnd.oasis.opendocument.text", filepath.Base(name))

	return t
}

func createTemplate(s *session, u *user, name string) *template {
	now := time.Now()

	t := &template{
		permissions: permissions{Owner: u.uuid},
		Name:        name,
		Created:     now,
		Updated:     now,
	}

	s.e.POST("/api/admin/template/update").WithJSON(t).Expect().Status(http.StatusOK)

	l := s.e.GET("/api/admin/template/list").Expect().Status(http.StatusOK).JSON()

	l.Path("$..name").Array().Contains(t.Name)

	for _, e := range l.Array().Iter() {
		if e.Object().Value("name").String().Raw() == t.Name {
			t.ID = e.Object().Value("id").String().Raw()
			break
		}
	}

	return t
}

func updateTemplate(s *session, t *template) *template {
	s.e.POST("/api/admin/template/update").WithQuery("id", t.ID).WithJSON(t).Expect().Status(http.StatusOK)

	expected := removeUpdatedField(toMap(t))
	s.e.GET("/api/admin/template/{id}").WithPath("id", t.ID).Expect().Status(http.StatusOK).
		JSON().Object().ContainsMap(expected)

	return t
}

func uploadTemplateFile(s *session, t *template, lang string, b []byte, contentType string, fileName string) {
	s.e.POST("/api/admin/template/upload/{id}/{lang}").WithPath("id", t.ID).WithPath("lang", lang).WithBytes(b).
		WithHeader("Content-Type", contentType).
		WithHeader("File-Name", fileName).
		WithHeader("Content-Length", strconv.Itoa(len(b))).
		Expect().Status(http.StatusOK)

	r := s.e.GET("/api/admin/template/{id}").WithPath("id", t.ID).
		Expect().Status(http.StatusOK).JSON()
	r.Path("$.data.en.name").Equal(fileName)
	r.Path("$.data.en.contentType").Equal(contentType)

}

func deleteTemplate(s *session, id string, expectEmptyList bool) {
	s.e.GET(fmt.Sprintf("/api/admin/template/%s/delete", id)).Expect().Status(http.StatusOK)
	l := s.e.GET("/api/admin/template/list").Expect()

	if expectEmptyList {
		l.Status(http.StatusNotFound)
	} else {
		l.Status(http.StatusOK).
			JSON().Path("$..name").Array().NotContains(id)
	}
}
