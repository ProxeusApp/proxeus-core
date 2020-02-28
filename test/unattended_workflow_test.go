package test

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/ProxeusApp/proxeus-core/test/assets"
)

func testUnattendedWorkflow(s *session) {
	u := registerTestUser(s)

	login(s, u)
	apiKey, summary := createApiKey(s, u, "test-"+s.id)
	w := createWorkflow(s, u, "workflow-"+s.id)
	f := createSimpleForm(s, u, "form-"+s.id, "test_name")
	tpl := createSimpleTemplate(s, u, "template-"+s.id, templateOdtPath)
	w.Data = simpleWorkflowData(s.id, f.ID, tpl.ID)
	updateWorkflow(s, w)
	logout(s)

	token := getSessionToken(s, u.username, apiKey)
	id := listFirstDocument(s, token)
	schema := getDocumentSchema(s, token, id)

	data := map[string]interface{}{}
	i := 0
	for k := range schema {
		data[k] = fmt.Sprintf("value-%d", i)
		i++
	}

	r := executeAllAtOnce(s, token, id, data, false)

	expected, err := assets.Asset("test/assets/templates/test_expected.pdf")
	if err != nil {
		s.t.Errorf("Cannot upload asset %s", err)
	}

	if bytes.Compare(cleanPDF(r), cleanPDF(expected)) != 0 {
		s.t.Errorf("Wrong pdf result")
	}

	login(s, u)
	deleteWorkflow(s, w.ID, true)
	deleteApiKey(s, u, summary)
	deleteUser(s, u)
}

func listFirstDocument(s *session, token string) string {
	return s.e.GET("/api/document/list").WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Array().First().Object().Value("id").String().Raw()
}

func getDocumentSchema(s *session, token, id string) map[string]interface{} {
	schema := s.e.GET("/api/document/{id}/allAtOnce/schema").WithPath("id", id).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object().Path("$.workflow.data").Object()
	schema.ContainsKey("test_name")
	return schema.Raw()
}

func executeAllAtOnce(s *session, token, id string, data map[string]interface{}, expectZip bool) []byte {
	r := s.e.POST("/api/document/{id}/allAtOnce").WithPath("id", id).
		WithHeader("Authorization", "Bearer "+token).WithJSON(data).
		Expect().Status(http.StatusOK)
	if expectZip {
		r.ContentType("application/zip")
	} else {
		r.ContentType("application/pdf")
	}

	return []byte(r.Body().Raw())
}

func testUnattendedWorkflowAdvanced(s *session) {
	u := registerTestUser(s)
	login(s, u)
	apiKey, summary := createApiKey(s, u, "test-"+s.id)

	w1, w2 := prepareWorkflows(s, u)

	token := getSessionToken(s, u.username, apiKey)
	id := listFirstDocument(s, token)
	schema := getDocumentSchema(s, token, id)

	data := map[string]interface{}{}
	i := 0
	for k := range schema {
		data[k] = fmt.Sprintf("value-%d", i)
		i++
	}

	r := executeAllAtOnce(s, token, id, data, true)
	if len(r) < 1000 {
		s.t.Error("got zip too small")
	}

	deleteWorkflow(s, w2.ID, false)
	deleteWorkflow(s, w1.ID, true)
	deleteApiKey(s, u, summary)
	deleteUser(s, u)
}
