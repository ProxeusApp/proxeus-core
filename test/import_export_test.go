package test

import (
	"net/http"
	"testing"
)

func TestImportExport(t *testing.T) {
	s := new(t, serverURL)
	u := registerTestUser(s)
	login(s, u)

	b1 := s.e.GET("api/export").WithQueryString("include=UserData").
		Expect().Body().Raw()
	s.e.POST("api/import").WithQueryString("skipExisting=false").WithBytes([]byte(b1)).
		Expect().Status(http.StatusOK)

	deleteUser(s, u)
}

func exportWorkflow(s *session, w *workflow) []byte {
	// export by name
	b1 := s.e.GET("api/workflow/export").WithQueryString("include=workflow&contains=" + w.Name).
		Expect().Body().Raw()
	// export by id
	b2 := s.e.GET("api/workflow/export").WithQueryString("id=" + w.ID).
		Expect().Body().Raw()

	if len(b1) < 1000 || len(b2) < 1000 {
		s.t.Error("export files too small", len(b1), len(b2))
	}

	stats := s.e.GET("api/export/results").Expect().JSON().Path("$.results").Object()
	stats.Path("$.Form").Object().Keys().Length().Equal(2)
	stats.Path("$.FormComponent").Object().Keys().Length().Equal(1)
	stats.Path("$.Template").Object().Keys().Length().Equal(1)
	stats.Path("$.User").Object().Keys().Length().Equal(1)

	return []byte(b2)
}

func importWorkflowNewUser(t *testing.T, exportedW []byte, expectedW *workflow) {
	s := new(t, serverURL)
	u := registerTestUser(s)
	login(s, u)

	s.e.POST("api/import").WithQueryString("skipExisting=false").WithBytes(exportedW).
		Expect().Status(http.StatusOK)

	// FIXME
	//stats := s.e.GET("api/import/results").Expect().JSON().Path("$.results").Object()
	//stats.Path("$.Form").Object().Keys().Length().Equal(2)
	//stats.Path("$.FormComponent").Object().Keys().Length().Equal(1)
	//stats.Path("$.Template").Object().Keys().Length().Equal(1)
	//stats.Path("$.User").Object().Keys().Length().Equal(1)
	//stats.Path("$.Workflow").Object().Keys().Length().Equal(2)
	//
	//s.e.GET("/api/admin/workflow/{id}").WithPath("id", expectedW.ID).Expect().Status(http.StatusOK).
	//	JSON().Object().ContainsMap(removeUpdatedField(toMap(expectedW)))

	deleteUser(s, u)
}
