package test

import (
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
	logout(s)
}
