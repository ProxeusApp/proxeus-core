package test

func testDocumentActions(s *session, u *user, documentID string) {
	testGetDocument(s, u, documentID)
	testSigning(s, u, documentID)
}
