package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

type workflow struct {
	permissions
	ID              string                 `json:"id" storm:"id"`
	Name            string                 `json:"name" storm:"index"`
	Detail          string                 `json:"detail"`
	Price           int                    `json:"price"`
	Published       bool                   `json:"published"`
	OwnerEthAddress string                 `json:"ownerEthAddress"`
	Updated         time.Time              `json:"updated" storm:"index"`
	Created         time.Time              `json:"created" storm:"index"`
	Data            map[string]interface{} `json:"data"`
}

func TestWorkflow(t *testing.T) {
	s := new(t, serverURL)
	u := registerTestUser(s)
	login(s, u)
	w1 := createSimpleWorkflow(s, u, "workflow1-"+s.id)
	w2 := createWorkflow(s, u, "workflow2-"+s.id)

	deleteWorkflow(s, w1.ID, false)
	deleteWorkflow(s, w2.ID, true)
	deleteUser(s, u)
}

func createWorkflow(s *session, u *user, name string) *workflow {
	now := time.Now()
	f := &workflow{
		permissions: permissions{Owner: u.uuid},
		Name:        name,
		Detail:      "details0",
		Created:     now,
		Updated:     now,
	}

	s.e.POST("/api/admin/workflow/update").WithJSON(f).Expect().Status(http.StatusOK)

	l := s.e.GET("/api/admin/workflow/list").Expect().Status(http.StatusOK).JSON()

	l.Path("$..name").Array().Contains(f.Name)

	for _, e := range l.Array().Iter() {
		if e.Object().Value("name").String().Raw() == f.Name {
			f.ID = e.Object().Value("id").String().Raw()
			break
		}
	}

	data := s.e.GET("/api/admin/workflow/{id}").WithPath("id", f.ID).Expect().
		Status(http.StatusOK).Body().Raw()
	err := json.Unmarshal([]byte(data), f)
	if err != nil {
		s.t.Error(err)
	}
	return f
}

func publishWorkflowWithPrice(s *session, w *workflow, xesPrice int) {
	w.Price = xesPrice
	w.Published = true
	s.e.POST("/api/admin/workflow/update").WithQueryString("publish=true&id=" + w.ID).
		WithJSON(w).Expect().Status(http.StatusOK)
}

func createSimpleWorkflow(s *session, u *user, name string) *workflow {
	w1 := createWorkflow(s, u, "workflow1-"+s.id)
	f := createSimpleForm(s, u, "form-"+s.id, fieldName)
	tpl := createSimpleTemplate(s, u, "template-"+s.id, "test/assets/test_template.odt")
	w1.Data = simpleWorkflowData(s.id, f.ID, tpl.ID)
	updateWorkflow(s, w1)
	return w1
}

func simpleWorkflowData(id string, formId, templateId string) map[string]interface{} {
	j := fmt.Sprintf(`{
    "flow": {
      "start": {
        "node": "%s",
        "p": {
          "x": -438,
          "y": -100
        }
      },
      "nodes": {
        "%s": {
          "id": "%s",
          "name": "test",
          "type": "form",
          "conns": [
            {
              "id": "%s"
            }
          ],
          "p": {
            "x": -225,
            "y": -102
          }
        },
        "%s": {
          "id": "%s",
          "name": "test",
          "type": "template",
          "p": {
            "x": -18,
            "y": -131
          }
        }
      }
    }
  }`, formId, formId, formId, templateId, templateId, templateId)

	var result map[string]interface{}

	err := json.Unmarshal([]byte(j), &result)
	if err != nil {
		return nil
	}

	return result
}

func updateWorkflow(s *session, f *workflow) *workflow {
	s.e.POST("/api/admin/workflow/update").
		WithQuery("id", f.ID).WithQuery("publish", true).WithJSON(f).
		Expect().Status(http.StatusOK)

	expected := removeUpdatedField(toMap(f))
	s.e.GET("/api/admin/workflow/{id}").WithPath("id", f.ID).Expect().Status(http.StatusOK).
		JSON().Object().ContainsMap(expected)

	return f
}

func hasUserWorkflow(s *session, containsW *workflow) {
	r := s.e.GET("/api/user/workflow/list").WithQueryString("i=0").Expect().Status(http.StatusOK)
	var found bool
	for _, val := range r.JSON().Array().Iter() {
		if val.Path("$.id").String().Raw() == containsW.ID {
			found = true
			break
		}
	}
	if !found {
		s.e.String("").Equal(containsW.ID)
	}
}

func deleteWorkflow(s *session, id string, expectEmptyList bool) {
	return
	s.e.GET(fmt.Sprintf("/api/admin/workflow/%s/delete", id)).Expect().Status(http.StatusOK)
	l := s.e.GET("/api/admin/workflow/list").Expect()

	if expectEmptyList {
		l.Status(http.StatusNotFound)
	} else {
		l.Status(http.StatusOK).
			JSON().Path("$..name").Array().NotContains(id)
	}
}
