package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

type form struct {
	permissions
	ID      string                 `json:"id" storm:"id"`
	Name    string                 `json:"name" storm:"index"`
	Detail  string                 `json:"detail"`
	Updated time.Time              `json:"updated" storm:"index"`
	Created time.Time              `json:"created" storm:"index"`
	Data    map[string]interface{} `json:"data"`
}

func TestForm(t *testing.T) {
	s := new(t, serverURL)
	u := registerTestUser(s)
	login(s, u)
	f1 := createSimpleForm(s, u, "form1-"+s.id, "test-"+s.id)
	f2 := createForm(s, u, "form2-"+s.id)

	deleteForm(s, f1.ID, false)
	deleteForm(s, f2.ID, true)
}

func createForm(s *session, u *user, name string) *form {
	now := time.Now()
	f := &form{
		permissions: permissions{Owner: u.uuid},
		Name:        name,
		Created:     now,
		Updated:     now,
	}

	s.e.POST("/api/admin/form/update").WithJSON(f).Expect().Status(http.StatusOK)

	l := s.e.GET("/api/admin/form/list").Expect().Status(http.StatusOK).JSON()

	l.Path("$..name").Array().Contains(f.Name)

	for _, e := range l.Array().Iter() {
		if e.Object().Value("name").String().Raw() == f.Name {
			f.ID = e.Object().Value("id").String().Raw()
			break
		}
	}

	return f
}

func createSimpleForm(s *session, u *user, name, fieldName string) *form {
	f := createForm(s, u, name)
	f.Data = simpleFormData(fieldName)
	return updateForm(s, f)
}

func simpleFormData(fieldName string) map[string]interface{} {
	j := fmt.Sprintf(`{
    "formSrc": {
      "components": {
        "5zvr98w21yynozx60nhmc": {
          "_compId": "HC2",
          "_order": 0,
          "autocomplete": "on",
          "help": "test-help",
          "label": "test-label",
          "name": "%s",
          "placeholder": "test-placeholder",
          "validate": {
            "required": true
          }
        }
      },
      "v": 2
    }
  }`, fieldName)

	var result map[string]interface{}

	err := json.Unmarshal([]byte(j), &result)
	if err != nil {
		return nil
	}

	return result
}

func updateForm(s *session, f *form) *form {
	s.e.POST("/api/admin/form/update").WithQuery("id", f.ID).WithJSON(f).Expect().Status(http.StatusOK)

	expected := removeUpdatedField(toMap(f))
	updated := s.e.GET("/api/admin/form/{id}").WithPath("id", f.ID).Expect().Status(http.StatusOK).
		JSON().Object().ContainsMap(expected).Path("$.updated").String().Raw()

	ti, err := time.Parse(time.RFC3339Nano, updated)
	if err != nil {
		s.t.Error(err)
	}
	f.Updated = ti
	return f
}

func deleteForm(s *session, id string, expectEmptyList bool) {

	s.e.GET(fmt.Sprintf("/api/admin/form/%s/delete", id)).Expect().Status(http.StatusOK)
	l := s.e.GET("/api/admin/form/list").Expect()

	if expectEmptyList {
		l.Status(http.StatusNotFound)
	} else {
		l.Status(http.StatusOK).
			JSON().Path("$..name").Array().NotContains(id)
	}
}
