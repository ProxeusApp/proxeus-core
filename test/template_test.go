package test

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ProxeusApp/proxeus-core/test/assets"
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

const templateOdtPath = "test/assets/templates/test_template.odt"

func testTemplate(s *session) {
	u := registerTestUser(s)
	login(s, u)
	t1 := createSimpleTemplate(s, u, "template1-"+s.id, templateOdtPath)
	t2 := createTemplate(s, u, "template2-"+s.id)

	f1 := createSimpleForm(s, u, "form-for-template", fieldName)
	formLookupForTemplate(s, f1)
	templatePreviews(s, t1, f1)

	deleteTemplate(s, t1.ID, false)
	deleteTemplate(s, t2.ID, true)
	deleteUser(s, u)
}

func createSimpleTemplate(s *session, u *user, name, path string) *template {
	odtFileBytes, err := assets.Asset(path)
	if err != nil {
		s.t.Errorf("Cannot read asset %s", err)
	}

	t := createTemplate(s, u, name)
	uploadTemplateFile(s, t, "en", odtFileBytes, "application/vnd.oasis.opendocument.text", filepath.Base(name))

	return t
}

func createTemplate(s *session, u *user, name string) *template {
	now := time.Now()

	t := &template{
		permissions: permissions{Owner: u.uuid},
		Name:        name,
		Detail:      "test",
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

func formLookupForTemplate(s *session, f *form) {
	f.Data["data"] = nil
	array := s.e.GET("/api/admin/template/ide/form").Expect().JSON().Array()
	arrayContainsMap(s.t, array, f)
	array = s.e.GET("/api/admin/template/ide/form").WithQuery("c", f.Name[2:]).
		Expect().JSON().Array()
	arrayContainsMap(s.t, array, f)

	s.e.GET("/api/admin/form/component").WithQueryString("l=1000").Expect().
		JSON().Object().ContainsKey("HC1")
}

func templatePreviews(s *session, t *template, f *form) {
	s.e.GET("/api/admin/template/ide/active/" + t.ID + "/en").Expect().Status(http.StatusOK)

	pdf1 := s.e.GET("/api/admin/template/ide/download/" + t.ID).Expect().
		Status(http.StatusOK).Body().Raw()

	s.e.POST("/api/admin/form/test/data/" + f.ID).WithJSON(map[string]string{fieldName: "value2"}).Expect().
		Status(http.StatusOK)

	pdf2 := s.e.GET("/api/admin/template/ide/download/" + t.ID).Expect().
		Status(http.StatusOK).Body().Raw()

	if len(pdf1) < 1000 || len(pdf2) < 1000 {
		s.t.Error("pdf too small to be valid", len(pdf1), len(pdf2))
	}
	if pdf1 == pdf2 {
		s.t.Error("pdf preview should have changed")
	}

	for _, format := range []string{"docx", "doc", "odt", "pdf"} {
		s.e.GET("/api/admin/template/ide/download/" + t.ID).WithQueryString("format=" + format).Expect().
			Status(http.StatusOK).Body().Length().Gt(1000)
	}

	s.e.GET("/api/admin/template/ide/tmplAssistanceDownload").Expect().
		Status(http.StatusOK).Body().Length().Gt(1000)

	// check if we can download what we've uploaded
	odtFileBytes, err := assets.Asset(templateOdtPath)
	if err != nil {
		s.t.Errorf("Cannot read asset %s", err)
	}
	s.e.GET("/api/admin/template/ide/download/" + t.ID).WithQueryString("raw").Expect().
		Status(http.StatusOK).Body().Equal(string(odtFileBytes))

	// used when template is not active (no preview)
	s.e.GET("/api/admin/template/download/" + t.ID + "/en").WithQueryString("raw").Expect().
		Status(http.StatusOK).Body().Equal(string(odtFileBytes))
	s.e.GET("/api/admin/template/download/" + t.ID + "/en").Expect().
		Status(http.StatusOK).Body().Length().Gt(1000)
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

	// TODO: fix vars on DS side
	s.e.GET("/api/admin/template/vars").WithQuery("c", t.ID).Expect().
		Status(http.StatusNotFound)
}

func deleteTemplate(s *session, id string, expectEmptyList bool) {
	// delete file
	s.e.GET(fmt.Sprintf("/api/admin/template/ide/delete/%s/en", id)).Expect().Status(http.StatusOK)
	// delete template
	s.e.GET(fmt.Sprintf("/api/admin/template/%s/delete", id)).Expect().Status(http.StatusOK)
	l := s.e.GET("/api/admin/template/list").Expect()

	if expectEmptyList {
		l.Status(http.StatusNotFound)
	} else {
		l.Status(http.StatusOK).
			JSON().Path("$..name").Array().NotContains(id)
	}
}
