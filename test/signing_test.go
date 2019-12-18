package test

import (
	"net/http"
)

func testSigning(s *session, u *user, documentID string) {
	s2 := new(s.t, serverURL)
	u2 := registerTestUser(s2)
	login(s2, u2)

	setEthKey(s, u)
	setEthKey(s2, u2)

	requestSignature(s, documentID, u2)
	revokeSignature(s, documentID, u2)

	requestSignature(s, documentID, u2)
	rejectSignature(s2, documentID)

	expectSignatureCount(s, documentID, 2)

	requestSignature(s, documentID, u2)
	putMockedSignature(s2, u, lastRequestedFileHash(s2))

	expectSignatureCount(s, documentID, 3)

	deleteUser(s2, u2)
}

func requestSignature(s *session, documentID string, fromUser *user) {
	s.e.POST("/api/user/document/signingRequests/"+documentID+"/docs[0]/add").
		WithFormField("signatory", fromUser.EthereumAddr).Expect().Status(http.StatusOK)
}

func revokeSignature(s *session, documentID string, fromUser *user) {
	s.e.POST("/api/user/document/signingRequests/"+documentID+"/docs[0]/revoke").
		WithFormField("signatory", fromUser.EthereumAddr).Expect().Status(http.StatusOK)
}

func rejectSignature(s *session, documentID string) {
	s.e.POST("/api/user/document/signingRequests/" + documentID + "/docs[0]/reject").
		Expect().Status(http.StatusOK)
}

func expectSignatureCount(s *session, documentID string, expected int) {
	s.e.GET("/api/user/document/signingRequests/" + documentID + "/docs[0]").
		Expect().Status(http.StatusOK).JSON().Array().Length().Equal(expected)
}

func lastRequestedFileHash(s *session) string {
	return s.e.GET("/api/user/document/signingRequests").
		Expect().Status(http.StatusOK).JSON().Array().Element(0).
		Path("$.hash").String().Raw()
}

func putMockedSignature(s *session, u *user, fileHash string) {
	req := struct {
		TxHash     string
		FileHash   string
		SignerAddr string
	}{
		TxHash:     randomHash(),
		FileHash:   fileHash,
		SignerAddr: u.EthereumAddr,
	}
	s.e.PUT("/api/test/signatures").WithJSON(req).Expect().Status(http.StatusOK)
}
