package test

import (
	"github.com/ProxeusApp/proxeus-core/externalnode"

	"net/http"

	uuid "github.com/satori/go.uuid"
)

const (
	fieldName_input  = "XES"
	fieldName_output = "USD_XES"
)

func testWorkflowExternalNode(s *session) {
	u := registerTestUser(s)
	login(s, u)

	w1 := createWorkflow(s, u, "workflow-"+s.id)
	f := createSimpleForm(s, u, "form-"+s.id, fieldName_input)
	externalNodeId := uuid.NewV4().String()

	w1.Data = workflowExternalNodeData(s, f.ID, externalNodeId)
	updateWorkflow(s, w1)

	configExternalNode(s, externalNodeId)

	type configData struct {
		FiatCurrency string
	}

	config := &configData{
		FiatCurrency: "USD",
	}

	node := externalnode.ExternalNodeInstance{
		ID:     externalNodeId,
		Config: config,
	}

	setExternalNodeConfig(s, externalNodeId, node)
	getExternalNodeConfig(s, externalNodeId, config)

	executeWorkflowExternalNode(s, w1)

	deleteWorkflow(s, w1.ID, true)
	deleteUser(s, u)
}

func executeWorkflowExternalNode(s *session, w *workflow) {
	expectWorkflowInCleanState(s, w)
	// filling a form
	{
		d := map[string]string{fieldName_input: "100"}
		s.e.POST("/api/document/" + w.ID + "/data").WithJSON(d).Expect().Status(http.StatusOK)

		r := s.e.POST("/api/document/" + w.ID + "/next").WithJSON(d).Expect().Status(http.StatusOK).
			JSON().Path("$.status")
		r.Path("$.steps").Array().Length().Equal(2)
		r.Path("$.userData").Object().ContainsKey(fieldName_input)
		r.Path("$.userData").Object().ContainsKey(fieldName_output)
		str := r.Path("$.userData").Object().Path("$." + fieldName_output).String()
		str.NotEmpty()

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
	s.e.GET("/api/admin/external/node-crypto-forex-rates/" + id).Expect().Status(http.StatusOK)
}

func setExternalNodeConfig(s *session, id string, config interface{}) {
	s.e.POST("/api/admin/external/config/" + id).WithJSON(config).Expect().Status(http.StatusOK)
}

func getExternalNodeConfig(s *session, id string, config interface{}) {
	s.e.GET("/api/admin/external/config/" + id).Expect().Status(http.StatusOK)
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
          "name": "node-crypto-forex-rates",
          "type": "externalNode",
          "p": {
            "x": 301,
            "y": -78
          }
        }
      }
    }
  }`
