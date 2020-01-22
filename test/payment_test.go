package test

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ProxeusApp/proxeus-core/sys/model"
)

func TestPayment(t *testing.T) {
	s := new(t, serverURL)
	s2 := new(t, serverURL)

	u := registerTestUser(s)
	u2 := registerTestUser(s2)

	login(s, u)
	login(s2, u2)

	setEthKey(s, u)
	setEthKey(s2, u2)

	w := createSimpleWorkflow(s, u, "workflow1-"+s.id)
	publishWorkflowWithPrice(s, w, 1)
	// check if workflow is public
	hasUserWorkflow(s2, w)

	p := initiatePayment(s2, w)
	updateTxHash(s2, p, randomHash())
	putMockedPayment(s2, p)
	paymentInStatus(s2, p, "confirmed")
	canStartWorkflow(s2, w)

	deleteUser(s, u)
	deleteUser(s2, u2)
}

func initiatePayment(s *session, w *workflow) *model.WorkflowPaymentItem {
	s.e.GET("/api/admin/payments/check").WithQuery("workflowId", w.ID).Expect().
		Status(http.StatusNotFound)
	data := s.e.POST("/api/admin/payments").WithJSON(map[string]string{"workflowId": w.ID}).
		Expect().Status(http.StatusOK).Body().Raw()

	var p model.WorkflowPaymentItem
	json.Unmarshal([]byte(data), &p)

	s.e.String(p.ID).Length().Gt(10)
	return &p
}

func updateTxHash(s *session, p *model.WorkflowPaymentItem, txHash string) {
	s.e.PUT("/api/admin/payments/" + p.ID).WithJSON(map[string]string{"txHash": txHash}).
		Expect().Status(http.StatusOK)
	p.TxHash = txHash
}

func paymentInStatus(s *session, p *model.WorkflowPaymentItem, expectedStatus string) {
	var status string
	for i := 0; i < 10; i++ {
		status = s.e.GET("/api/admin/payments/"+p.ID).WithQuery("txHash", p.TxHash).
			Expect().Status(http.StatusOK).JSON().Path("$.status").String().Raw()
		if status == expectedStatus {
			break
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	s.e.String(status).Equal(expectedStatus)
	p.Status = "confirmed"
	expected := removeTimeFields(toMap(p))
	s.e.GET("/api/admin/payments").WithQuery("txHash", p.TxHash).
		WithQuery("status", "confirmed").Expect().Status(http.StatusOK).
		JSON().Object().ContainsMap(expected)
}

func canStartWorkflow(s *session, w *workflow) {
	s.e.GET("/api/document/" + w.ID).Expect().Status(http.StatusOK)
}

func putMockedPayment(s *session, p *model.WorkflowPaymentItem) {
	req := struct {
		TxHash string
		From   string
		To     string
	}{
		TxHash: p.TxHash,
		From:   p.From,
		To:     p.To,
	}
	s.e.PUT("/api/test/payments").WithJSON(req).Expect().Status(http.StatusOK)
}

func randomHash() string {
	var h common.Hash
	for i := 0; i < len(h); i++ {
		h[i] = byte(rand.Uint32())
	}
	return h.String()
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
