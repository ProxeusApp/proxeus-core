package test

import (
	"net/http"

	uuid "github.com/satori/go.uuid"
)

func testWorkflowExternalNode(s *session) {
	u := registerTestUser(s)
	login(s, u)

	w1 := createWorkflow(s, u, "workflow-"+s.id)
	f := createSimpleForm(s, u, "form-"+s.id, fieldName)
	externalNodeId := uuid.NewV4().String()

	w1.Data = workflowExternalNodeData(s, f.ID, externalNodeId)
	updateWorkflow(s, w1)

	configExternalNode(s, externalNodeId)
	executeWorkflowExternalNode(s, w1)

	deleteWorkflow(s, w1.ID, true)
	deleteUser(s, u)
}

func executeWorkflowExternalNode(s *session, w *workflow) {
	expectWorkflowInCleanState(s, w)
	// filling a form
	{
		d := map[string]string{fieldName: "test 100 CHF"}
		s.e.POST("/api/document/" + w.ID + "/data").WithJSON(d).Expect().Status(http.StatusOK)

		r := s.e.POST("/api/document/" + w.ID + "/next").WithJSON(d).Expect().Status(http.StatusOK).
			JSON().Path("$.status")
		r.Path("$.steps").Array().Length().Equal(2)
		r.Path("$.userData").Object().ContainsKey(fieldName)
		str := r.Path("$.userData").Object().Path("$." + fieldName).String()
		str.Contains("test")
		str.Contains("XES")
		r.Path("$.data").NotNull()
	}
}

func workflowExternalNodeData(s *session, formID string, externalNodeId string) map[string]interface{} {
	d1, err := advancedWorkflowData(workflowXData, map[string]string{
		"formId":         formID,
		"externalNodeId": externalNodeId,
	})
	if err != nil {
		s.t.Fatal(err)
	}
	return d1
}

func configExternalNode(s *session, id string) {
	s.e.GET("/api/admin/external/priceGetter/" + id).Expect().Status(http.StatusOK)
}

const workflowXData = `{
    "flow": {
      "start": {
        "p": {
          "x": -21,
          "y": 72
        },
        "node": "{{.formId}}"
      },
      "nodes": {
        "{{.formId}}": {
          "id": "{{.formId}}",
          "name": "test",
          "type": "form",
          "p": {
            "x": 97,
            "y": -89
          },
          "conns": [
            {
              "id": "{{.externalNodeId}}"
            }
          ]
        },
        "{{.externalNodeId}}": {
          "id": "{{.externalNodeId}}",
          "name": "priceGetter",
          "type": "externalNode",
          "p": {
            "x": 301,
            "y": -78
          }
        }
      }
    }
  }`
