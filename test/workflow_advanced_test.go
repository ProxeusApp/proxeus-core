package test

import (
	"bytes"
	"encoding/json"
	tpl "html/template"
	"net/http"
	"testing"
)

const fieldName = "test_name"

func TestWorkflowAdvanced(t *testing.T) {
	s := new(t, serverURL)
	u := registerTestUser(s)
	login(s, u)
	w1 := createWorkflow(s, u, "sub-workflow1-"+s.id)
	w2 := createWorkflow(s, u, "workflow2-"+s.id)

	f := createSimpleForm(s, u, "form-"+s.id, fieldName)
	tpl := createSimpleTemplate(s, u, "template-"+s.id, "test/assets/test_template.odt")
	w1.Data = advancedWorkflowData(t, workflow1Data, map[string]string{
		"formId":     f.ID,
		"templateId": tpl.ID,
	})
	w2.Data = advancedWorkflowData(t, workflow2Data, map[string]string{
		"formId":        f.ID,
		"subworkflowId": w1.ID,
	})
	updateWorkflow(s, w1)
	updateWorkflow(s, w2)

	documentId := executeWorkflow(s, w2)

	testDocumentActions(s, documentId)

	deleteWorkflow(s, w2.ID, false)
	deleteWorkflow(s, w1.ID, true)
}

func advancedWorkflowData(t *testing.T, data string, dataValues map[string]string) map[string]interface{} {
	tp, err := tpl.New("").Parse(data)
	if err != nil {
		t.Error(err)
	}

	var buf bytes.Buffer
	tp.Execute(&buf, dataValues)

	var result map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &result)
	if err != nil {
		t.Error(err)
	}
	return result
}

func executeWorkflow(s *session, w *workflow) string {
	{
		r := s.e.GET("/api/document/" + w.ID).Expect().Status(http.StatusOK).JSON().Path("$.status")
		r.Path("$.hasNext").Boolean().True()
		r.Path("$.hasPrev").Boolean().False()
		r.Path("$.steps").Array().Length().Gt(0)
		r.Path("$.data").NotNull()
	}
	d := map[string]string{fieldName: "value1"}
	// filling a form
	{
		s.e.POST("/api/document/" + w.ID + "/data").WithJSON(d).Expect().Status(http.StatusOK)

		r := s.e.POST("/api/document/" + w.ID + "/next").WithJSON(d).Expect().Status(http.StatusOK).
			JSON().Path("$.status")
		r.Path("$.steps").Array().Length().Equal(2)
		r.Path("$.userData").Object().ContainsKey(fieldName)
		r.Path("$.data").NotNull()
	}
	// sub-workflow
	{
		// step back to parent workflow
		s.e.GET("/api/document/" + w.ID + "/prev").Expect().Status(http.StatusOK)
		// go forward
		s.e.POST("/api/document/" + w.ID + "/next").WithJSON(d).Expect().Status(http.StatusOK).JSON()
		// execute some unattended nodes
		r := s.e.POST("/api/document/" + w.ID + "/next").WithJSON(d).Expect().Status(http.StatusOK).JSON()

		// final node
		r.Path("$.status.hasNext").Boolean().False()
	}

	// final confirmations
	{
		s.e.POST("/api/document/" + w.ID + "/next").WithQueryString("confirm").
			Expect().Status(http.StatusOK)
		r := s.e.POST("/api/document/" + w.ID + "/next").WithQueryString("confirm").
			Expect().Status(http.StatusOK)

		pdfID := r.JSON().Path("$.status.docs[0].id").String().Raw()
		s.e.GET("/api/document/" + w.ID + "/preview/" + pdfID + "/en/pdf").
			Expect().Body().Length().Gt(3000)

		rf := s.e.POST("/api/document/" + w.ID + "/next").WithQueryString("final").
			WithJSON(map[string]string{"id": w.ID}).Expect().Status(http.StatusOK)

		return rf.JSON().Path("$.id").String().Raw()
	}
}

const workflow1Data = `{
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
              "id": "{{.templateId}}"
            }
          ]
        },
        "{{.templateId}}": {
          "id": "{{.templateId}}",
          "name": "test",
          "type": "template",
          "p": {
            "x": 301,
            "y": -78
          }
        }
      }
    }
  }`

const workflow2Data = `{
    "flow": {
      "start": {
        "p": {
          "x": -68,
          "y": -19
        },
        "node": "{{.formId}}"
      },
      "nodes": {
        "3": {
          "id": "3",
          "name": "Price retriever",
          "detail": "Retrieves CHF/XES price",
          "type": "priceretriever",
          "p": {
            "x": 196,
            "y": -294
          },
          "conns": [
            {
              "id": "1234123-1234124"
            }
          ]
        },
        "1234123-1234124": {
          "id": "1234123-1234124",
          "name": "Mail Sender",
          "detail": "sends an email",
          "type": "mailsender",
          "p": {
            "x": 356,
            "y": -185
          }
        },
        "{{.formId}}": {
          "id": "{{.formId}}",
          "name": "test",
          "type": "form",
          "p": {
            "x": -85,
            "y": -210
          },
          "conns": [
            {
              "id": "{{.subworkflowId}}"
            }
          ]
        },
        "{{.subworkflowId}}": {
          "id": "{{.subworkflowId}}",
          "name": "sub-flow",
          "type": "workflow",
          "p": {
            "x": 22,
            "y": -377
          },
          "conns": [
            {
              "id": "3"
            }
          ]
        }
      }
    }
  }`
