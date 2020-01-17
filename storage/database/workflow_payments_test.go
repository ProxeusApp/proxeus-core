package database

import (
	"testing"
	"time"

	"github.com/ProxeusApp/proxeus-core/sys/model"

	. "github.com/onsi/gomega"
)

func TestWorkflowPayments(t *testing.T) {
	RegisterTestingT(t)
	wp := testDBSet.WorkflowPayments

	item := &model.WorkflowPaymentItem{
		ID:         "1",
		WorkflowID: "w1",
		TxHash:     "2323232",
		From:       "0xfefe",
		To:         "0xbeef",
		Status:     model.PaymentStatusCreated,
		Xes:        123,
	}
	item2 := &model.WorkflowPaymentItem{
		ID:         "2",
		WorkflowID: "w2",
		TxHash:     "4532234",
		From:       "0xfefe",
		To:         "0xbeef",
		Status:     model.PaymentStatusCreated,
	}
	item3 := &model.WorkflowPaymentItem{
		ID:         "3",
		WorkflowID: "w3",
		TxHash:     "122112",
		From:       "0xfefe",
		To:         "0xbeef",
		Status:     model.PaymentStatusCreated,
	}

	// add
	Expect(wp.Save(item)).To(Succeed())
	Expect(wp.Save(item2)).To(Succeed())
	Expect(wp.Save(item3)).To(Succeed())
	Expect(wp.Update(item2.ID, model.PaymentStatusPending, item2.TxHash, item2.From)).To(Succeed())
	item2.Status = model.PaymentStatusPending
	Expect(wp.All()).To(equalJSON([]*model.WorkflowPaymentItem{item, item2, item3}))

	// get
	Expect(wp.Get(item.ID)).To(equalJSON(item))
	Expect(wp.GetByTxHashAndStatusAndFromEthAddress(item.TxHash, item.Status, item.From)).
		To(equalJSON(item))
	Expect(wp.GetByWorkflowIdAndFromEthAddress(item2.WorkflowID, item2.From, nil)).
		To(equalJSON(item2))
	Expect(wp.GetByWorkflowIdAndFromEthAddress(item2.WorkflowID, item2.From, []string{item2.Status})).
		To(equalJSON(item2))

	// managing payment state
	Expect(wp.ConfirmPayment(item.TxHash, item.From, item.To, item.Xes)).To(Succeed())
	Expect(wp.Redeem(item.WorkflowID, item.From)).To(Succeed())
	Expect(wp.Cancel(item3.ID, item3.From)).To(Succeed())

	// abandon
	Expect(wp.SetAbandonedToTimeoutBeforeTime(time.Now())).To(Succeed())

	// delete all
	Expect(wp.Delete(item.ID)).To(Succeed())
	Expect(wp.Remove(item2)).To(Succeed())
}
