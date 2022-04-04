package test

import (
	"bytes"
	"net/http"
	// "io/ioutil"

	"github.com/ProxeusApp/proxeus-core/test/assets"
)

func testWorkflowWithFile(s *session) {
	u := registerTestUser(s)
	login(s, u)

	w1 := createWorkflow(s, u, "workflow-"+s.id)
	f := createFormWithData(s, u, "form-"+s.id, fieldName, fileFormData)
	tpl := createSimpleTemplate(s, u, "template-"+s.id, "test/assets/templates/test_template2.odt")

	w1.Data = simpleWorkflowData(s.id, f.ID, tpl.ID)
	updateWorkflow(s, w1)

	executeFileWorkflow(s, w1)

	deleteWorkflow(s, w1.ID, true)
	deleteUser(s, u)
}

func executeFileWorkflow(s *session, w *workflow) {
	expectWorkflowInCleanState(s, w)

	image, err := assets.Asset("test/assets/image.jpg")
	if err != nil {
		s.t.Errorf("Cannot read asset %s", err)
	}
	s.e.POST("/api/document/" + w.ID + "/file/" + fieldName).WithBytes(image).
		Expect().Status(http.StatusOK)
	s.e.GET("/api/document/" + w.ID + "/file/" + fieldName).Expect().
		Status(http.StatusOK).Body().Equal(string(image))

	r := s.e.POST("/api/document/" + w.ID + "/next").Expect().Status(http.StatusOK)

	previewID := r.JSON().Path("$.status.docs[0].id").String().Raw()
	previewPDF := s.e.GET("/api/document/" + w.ID + "/preview/" + previewID + "/en/pdf").Expect().Status(http.StatusOK).Body().Raw()

	expectedPDF, err := assets.Asset("test/assets/templates/test_expected2_B.pdf")
	if err != nil {
		s.t.Errorf("Cannot read asset %s", err)
	}

	// err1 := ioutil.WriteFile("/tmp/file_preview.pdf", []byte(previewPDF), 0644)
	// if err1 != nil {
	// 	s.t.Errorf("Could not write preview")
	// }

	if bytes.Compare(cleanPDF([]byte(previewPDF)), cleanPDF(expectedPDF)) != 0 {
		s.t.Errorf("Wrong pdf result")
		s.t.Errorf("(%d vs %d bytes)", len([]byte(previewPDF)), len(expectedPDF))
	}
}

const fileFormData = `{
    "formSrc": {
      "components": {
        "mdr29f7gt3znfvoyxvyzd": {
          "_compId": "HC12",
          "_file": true,
          "_order": 1,
          "help": "help text",
          "label": "File upload",
          "name": "%s",
          "placeholder": "Click here to select a file.",
          "validate": {
            "required": true
          }
        }
      },
      "v": 2
    }
}`
