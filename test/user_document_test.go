package test

import "net/http"

func testGetDocument(s *session, u *user, documentId string) {
	getUserDocumentFile(s, documentId)
}

func getUserDocumentFile(s *session, documentId string) {
	resp := s.e.GET("/api/user/document/file/" + documentId + "/docs[0]").Expect().Status(http.StatusOK)
	resp.Headers().ContainsKey("Content-Disposition")
	resp.Headers().ContainsKey("Content-Length")
	resp.Headers().ContainsKey("Content-Disposition")
	resp.Headers().Value("Content-Type").Equal([]string{"application/pdf"})
	resp.Header("Content-Length").NotEqual("0")
}
