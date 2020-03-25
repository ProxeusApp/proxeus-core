package test

import (
	"bytes"
	"fmt"
	"github.com/ProxeusApp/proxeus-core/test/assets"
	"net/http"
)

func testImportExport(s *session) {
	u := registerTestUser(s)
	login(s, u)
	b1 := s.e.GET("/api/export").WithQueryString("include=UserData").
		Expect().Status(http.StatusOK).Body().Raw()
	s.e.POST("/api/import").WithQueryString("skipExisting=false").WithBytes([]byte(b1)).
		Expect().Status(http.StatusOK)

	deleteUser(s, u)
}

func exportWorkflow(s *session, w *workflow) []byte {
	// Set WantToBeFound to true in order to have the User included in the Export
	id := s.e.GET("/api/me").Expect().Status(http.StatusOK).JSON().
		Path("$.id").String().Raw()
	s.e.POST("/api/me").WithJSON(
		struct {
			ID            string `json:"id"`
			WantToBeFound bool   `json:"wantToBeFound"`
		}{
			ID:            id,
			WantToBeFound: true,
		}).Expect().Status(http.StatusOK)

	// export by name
	b1 := s.e.GET("/api/workflow/export").WithQueryString("include=workflow&contains=" + w.Name).
		Expect().Status(http.StatusOK).Body().Raw()
	// export by id
	b2 := s.e.GET("/api/workflow/export").WithQueryString("id=" + w.ID).
		Expect().Status(http.StatusOK).Body().Raw()

	if len(b1) < 1000 || len(b2) < 1000 {
		s.t.Error("export files too small", len(b1), len(b2))
	}

	stats := s.e.GET("/api/export/results").Expect().JSON().Path("$.results").Object()
	stats.Path("$.Form").Object().Keys().Length().Equal(2)
	stats.Path("$.FormComponent").Object().Keys().Length().Equal(1)
	stats.Path("$.Template").Object().Keys().Length().Equal(2)
	stats.Path("$.User").Object().Keys().Length().Equal(1)

	return []byte(b2)
}

func importWorkflow(s *session, exportedW []byte, expectedW *workflow) {
	s.e.POST("/api/import").WithQueryString("skipExisting=false").WithBytes(exportedW).
		Expect().Status(http.StatusOK)

	stats := s.e.GET("/api/import/results").Expect().JSON().Path("$.results").Object()
	stats.Path("$.Form").Object().Keys().Length().Equal(2)
	stats.Path("$.FormComponent").Object().Keys().Length().Equal(1)
	stats.Path("$.User").Object().Keys().Length().Equal(1)
	stats.Path("$.Workflow").Object().Keys().Length().Equal(2)

	s.e.GET("/api/admin/workflow/{id}").WithPath("id", expectedW.ID).Expect().Status(http.StatusOK).
		JSON().Object().ContainsMap(removeTimeFields(toMap(expectedW)))
}

func exportImportEntity(s *session, entity string) {
	b := s.e.GET("/api/{entity}/export").WithPath("entity", entity).
		WithQuery("contains", "test").
		Expect().Status(http.StatusOK).Body().Raw()

	s.e.POST("/api/import").WithQueryString("skipExisting=false").
		WithBytes([]byte(b)).Expect().Status(http.StatusOK)
}

func testImportExportAdmin(s *session) {
	u := registerSuperAdmin(s)
	login(s, u)

	w1, w2 := prepareWorkflows(s, u)

	exportEntities := []string{
		"userdata",
		"user",
		"i18n",
		"workflow",
		"form",
		"template",
	}
	for _, entity := range exportEntities {
		exportImportEntity(s, entity)
	}
	deleteWorkflow(s, w2.ID, false)
	deleteWorkflow(s, w1.ID, true)
	//deleteUser(s, u)
}

func testImportExportRoot(s *session) {
	login(s, s.root)
	exportImportEntity(s, "settings")
	testExportImportTemplate(s)
	logout(s)
}

func testExportImportTemplate(s *session) {
	t1 := createSimpleTemplate(s, s.root, "template1-"+s.id, templateOdtPath)

	// check template does exist
	s.e.GET(fmt.Sprintf("/api/admin/template/%s", t1.ID)).
		Expect().Status(http.StatusOK)

	// export new template
	b1 := s.e.GET("/api/template/export").WithQueryString("id=" + t1.ID).
		Expect().Status(http.StatusOK).Body().Raw()

	// delete template
	s.e.GET(fmt.Sprintf("/api/admin/template/%s/delete", t1.ID)).
		Expect().Status(http.StatusOK).Body().Raw()

	// check template does not exist
	s.e.GET(fmt.Sprintf("/api/admin/template/%s", t1.ID)).
		Expect().Status(http.StatusNotFound)

	// import template
	s.e.POST("/api/import").WithQueryString("skipExisting=false").WithBytes([]byte(b1)).
		Expect().Status(http.StatusOK)

	// get template file
	b2 := s.e.GET(fmt.Sprintf("/api/admin/template/download/%s/en", t1.ID)).WithQueryString("raw").
		Expect().Status(http.StatusOK).Body().Raw()

	odtFileBytes, err := assets.Asset(templateOdtPath)
	if err != nil {
		s.t.Errorf("Cannot read asset %s", err)
	}
	if !bytes.Equal(odtFileBytes, []byte(b2)) {
		s.t.Error("Export or import failed. Imported a template but did not find it afterwards")
	}
}
