package test

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
	"time"

	w "git.proxeus.com/core/central/sys/workflow"
)

func TestUnattendedWorkflow(t *testing.T) {
	s := new(t, serverURL)
	u := registerTestUser(s)

	login(s, u)
	apiKey, summary := createApiKey(s, u, "test-"+s.id)
	w := createWorkflow(s, u, "workflow-"+s.id)
	f := createSimpleForm(s, u, "form-"+s.id, "test_name")
	tpl := createSimpleTemplate(s, u, "template-"+s.id, "test/assets/test_template.odt")
	w.Data = simpleWorkflowData(s.id, f.ID, tpl.ID)
	updateWorkflow(s, w)
	logout(s)

	token := getSessionToken(s, u.username, apiKey)
	id := listFirstDocument(s, token)
	schema := getDocumentSchema(s, token, id)

	data := map[string]interface{}{}
	i := 0
	for k, _ := range schema {
		data[k] = fmt.Sprintf("value-%d", i)
		i++
	}

	r := executeAllAtOnce(s, token, id, data)

	expected, err := Asset("test/assets/test_expected.pdf")
	if err != nil {
		s.t.Errorf("Cannot upload asset %s", err)
	}

	if bytes.Compare(cleanPDF(r), cleanPDF(expected)) != 0 {
		t.Errorf("Wrong pdf result")
	}

	login(s, u)
	deleteWorkflow(s, w.ID, true)
	deleteApiKey(s, u, summary)
	deleteUser(s, u)
}

type workflowItem struct {
	permissions
	ID      string    `json:"id" storm:"id"`
	Name    string    `json:"name" storm:"index"`
	Detail  string    `json:"detail"`
	Updated time.Time `json:"updated" storm:"index"`
	Created time.Time `json:"created" storm:"index"`
	Price   uint64    `json:"price" storm:"index"`

	Data            *w.Workflow `json:"data"`
	OwnerEthAddress string      `json:"ownerEthAddress"` //only used in frontend
	Deactivated     bool        `json:"deactivated"`
}

func listFirstDocument(s *session, token string) string {
	return s.e.GET("/api/document/list").WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Array().First().Object().Value("id").String().Raw()
}

func getDocumentSchema(s *session, token, id string) map[string]interface{} {
	schema := s.e.GET("/api/document/{id}/allAtOnce/schema").WithPath("id", id).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object().Path("$.workflow.data").Object()
	schema.ContainsKey("test_name")
	return schema.Raw()
}

func executeAllAtOnce(s *session, token, id string, data map[string]interface{}) []byte {
	r := s.e.POST("/api/document/{id}/allAtOnce").WithPath("id", id).WithHeader("Authorization", "Bearer "+token).WithJSON(data).Expect().ContentType("application/pdf").Body().Raw()

	return []byte(r)
}
