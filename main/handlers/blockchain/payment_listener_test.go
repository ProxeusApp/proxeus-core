package blockchain

import (
	"errors"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"

	strm "github.com/asdine/storm"

	"git.proxeus.com/core/central/sys/model"

	"git.proxeus.com/core/central/sys/db/storm"
)

var errCleanupTestData = errors.New("db data has not been cleanup up after finishing tests")

func removePaymentIfExists(workflowPaymentsDB storm.WorkflowPaymentsDBInterface,
	persistedPaymentItem **model.WorkflowPaymentItem) {

	if *persistedPaymentItem != nil {
		err := workflowPaymentsDB.Remove(*persistedPaymentItem)
		if err != nil {
			panic(err.Error())
		}
	}
}

func TestPaymentEventHandling(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	workflowPaymentsDB, err := storm.NewWorkflowPaymentDB(".test_data")
	if err != nil {
		panic(err)
	}

	defer func() {
		payments, err := workflowPaymentsDB.All()
		if err != nil {
			panic(err)
		}
		if len(payments) != 0 {
			panic(errCleanupTestData)
		}

		err = os.Remove(filepath.Join(".test_data", storm.WorkflowPaymentDBDir, storm.WorkflowPaymentDB))
		if err != nil {
			panic(err.Error())
		}
		err = os.Remove(filepath.Join(".test_data", storm.WorkflowPaymentDBDir))
		if err != nil {
			panic(err.Error())
		}
		err = os.Remove(".test_data")
		if err != nil {
			panic(err.Error())
		}

	}()

	t.Run("ShouldSetPendingPaymentWithTxHashToConfirmed", func(t *testing.T) {

		var persistedPaymentItem *model.WorkflowPaymentItem
		defer removePaymentIfExists(workflowPaymentsDB, &persistedPaymentItem)

		paymentID := "1"
		paymentTxHash := "0x04f1bbf224b5876d91c74984c4d7f7768c5cc9da5b7f7afe1a31ef9115310f67"
		from := "0xe4f2604bc8300004aa4477af5f2AdBd37765F3F7"
		to := "0xe902Fb81617079236cB6eF8f34b2A1e759ef676D"
		eventTxHash := paymentTxHash
		expectedStatus := model.PaymentStatusConfirmed
		xesAmount := big.NewInt(1)

		err = runTest(mockCtrl, workflowPaymentsDB, xesAmount, eventTxHash, paymentID, paymentTxHash, from,
			to, model.PaymentStatusPending, xesAmount)
		if err != nil {
			panic(err.Error())
		}

		persistedPaymentItem, err := workflowPaymentsDB.Get(paymentID)
		if err != nil {
			panic(err)
		}

		if persistedPaymentItem.Status != expectedStatus {
			t.Errorf("Expected persistedPaymentItem to have status %s but got: %s",
				expectedStatus, persistedPaymentItem.Status)
		}

		if persistedPaymentItem.TxHash != eventTxHash {
			t.Errorf("Expected persistedPaymentItem to have TxHash %s but got: %s",
				eventTxHash, persistedPaymentItem.TxHash)
		}
	})

	t.Run("ShouldSetCreatedPaymentWithoutTxHashToConfirmed", func(t *testing.T) {

		var persistedPaymentItem *model.WorkflowPaymentItem
		defer removePaymentIfExists(workflowPaymentsDB, &persistedPaymentItem)

		paymentID := "2"
		eventTxHash := "0x94420cb493b721c627feb8f911df8546bba0f911cb4433dfabd4c9c65012593c"
		paymentTxHash := ""
		from := "0xe4f2604bc8300004aa4477af5f2AdBd37765F3F7"
		to := "0xe902Fb81617079236cB6eF8f34b2A1e759ef676D"
		expectedStatus := model.PaymentStatusConfirmed
		xesAmount := big.NewInt(100000000000000000)

		err = runTest(mockCtrl, workflowPaymentsDB, xesAmount, eventTxHash, paymentID,
			paymentTxHash, from, to, model.PaymentStatusCreated, xesAmount)
		if err != nil {
			panic(err.Error())
		}

		persistedPaymentItem, err := workflowPaymentsDB.Get(paymentID)
		if err != nil {
			panic(err)
		}

		if persistedPaymentItem.Status != expectedStatus {
			t.Errorf("Expected persistedPaymentItem to have status %s but got: %s",
				expectedStatus, persistedPaymentItem.Status)
		}

		if persistedPaymentItem.TxHash != eventTxHash {
			t.Errorf("Expected persistedPaymentItem to have TxHash %s but got: %s",
				eventTxHash, persistedPaymentItem.TxHash)
		}
	})

	t.Run("ShouldIgnoreAlreadyRedeemedPayment", func(e *testing.T) {

		var persistedPaymentItem *model.WorkflowPaymentItem
		defer removePaymentIfExists(workflowPaymentsDB, &persistedPaymentItem)

		paymentID := "3"
		eventTxHash := "0x8c962ca22918cf37e89a7bef93efe2938320c38ec113321d847d6fc48f2ba2fa"
		paymentTxHash := ""
		from := "0xe4f2604bc8300004aa4477af5f2AdBd37765F3F7"
		to := "0xe902Fb81617079236cB6eF8f34b2A1e759ef676D"
		expectedStatus := model.PaymentStatusRedeemed
		xesAmount := big.NewInt(1)

		err = runTest(mockCtrl, workflowPaymentsDB, xesAmount, eventTxHash, paymentID,
			paymentTxHash, from, to, model.PaymentStatusRedeemed, xesAmount)

		if err != strm.ErrNotFound {
			if err != nil {
				t.Errorf("Expected to have %s but got: %s", strm.ErrNotFound, err.Error())
			}
			t.Errorf("Expected to have %s but got: nil", strm.ErrNotFound)
		}

		persistedPaymentItem, err = workflowPaymentsDB.Get(paymentID)
		if err != nil {
			panic(err)
		}

		if persistedPaymentItem.Status != expectedStatus {
			t.Errorf("Expected persistedPaymentItem to have status %s but got: %s",
				expectedStatus, persistedPaymentItem.Status)
		}
	})

	t.Run("ShouldIgnorePaymentOnNotMatchingXesAmount", func(e *testing.T) {

		var persistedPaymentItem *model.WorkflowPaymentItem
		defer removePaymentIfExists(workflowPaymentsDB, &persistedPaymentItem)

		paymentID := "4"
		paymentTxHash := "0x04f1bbf224b5876d91c74984c4d7f7768c5cc9da5b7f7afe1a31ef9115310f67"
		from := "0xe4f2604bc8300004aa4477af5f2AdBd37765F3F7"
		to := "0xe902Fb81617079236cB6eF8f34b2A1e759ef676D"
		eventTxHash := paymentTxHash
		expectedStatus := model.PaymentStatusCreated
		xesAmount := big.NewInt(2)
		xesAmountEvent := big.NewInt(1)

		err = runTest(mockCtrl, workflowPaymentsDB, xesAmountEvent, eventTxHash, paymentID, paymentTxHash, from,
			to, model.PaymentStatusCreated, xesAmount)
		if err != strm.ErrNotFound {
			if err != nil {
				t.Errorf("Expected to have %s but got: %s", strm.ErrNotFound, err.Error())
			}
			t.Errorf("Expected to have %s but got: nil", strm.ErrNotFound)
		}

		persistedPaymentItem, err := workflowPaymentsDB.Get(paymentID)
		if err != nil {
			panic(err)
		}

		if persistedPaymentItem.Status != expectedStatus {
			t.Errorf("Expected persistedPaymentItem to have status %s but got: %s",
				expectedStatus, persistedPaymentItem.Status)
		}

		if persistedPaymentItem.TxHash != eventTxHash {
			t.Errorf("Expected persistedPaymentItem to have TxHash %s but got: %s",
				eventTxHash, persistedPaymentItem.TxHash)
		}
	})

	t.Run("ShouldReturnErrorOverflowOnBigXesAmount", func(e *testing.T) {

		var persistedPaymentItem *model.WorkflowPaymentItem
		defer removePaymentIfExists(workflowPaymentsDB, &persistedPaymentItem)

		paymentID := "5"
		eventTxHash := "0x8c962ca22918cf37e89a7bef93efe2938320c38ec113321d847d6fc48f2ba2fa"
		paymentTxHash := ""
		from := "0xe4f2604bc8300004aa4477af5f2AdBd37765F3F7"
		to := "0xe902Fb81617079236cB6eF8f34b2A1e759ef676D"

		xesTmp := big.NewInt(1000000000000000000)
		xesAmount := xesTmp.Mul(xesTmp, big.NewInt(1000000000000000000))

		err = runTest(mockCtrl, workflowPaymentsDB, xesAmount, eventTxHash, paymentID,
			paymentTxHash, from, to, model.PaymentStatusRedeemed, xesAmount)

		if err != xesOverflowError {
			if err != nil {
				t.Errorf("Expected to have %s but got: %s", xesOverflowError, err.Error())
			}
			t.Errorf("Expected to have %s but got: nil", xesOverflowError)
		}

		persistedPaymentItem, err = workflowPaymentsDB.Get(paymentID)
		if err != nil {
			panic(err)
		}
	})

	t.Run("ShouldReturnErrorOverflowOnNegativeXesAmount", func(e *testing.T) {

		var persistedPaymentItem *model.WorkflowPaymentItem
		defer removePaymentIfExists(workflowPaymentsDB, &persistedPaymentItem)

		paymentID := "6"
		eventTxHash := "0x8c962ca22918cf37e89a7bef93efe2938320c38ec113321d847d6fc48f2ba2fa"
		paymentTxHash := ""
		from := "0xe4f2604bc8300004aa4477af5f2AdBd37765F3F7"
		to := "0xe902Fb81617079236cB6eF8f34b2A1e759ef676D"

		xesAmount := big.NewInt(-1)

		err = runTest(mockCtrl, workflowPaymentsDB, xesAmount, eventTxHash, paymentID,
			paymentTxHash, from, to, model.PaymentStatusRedeemed, xesAmount)

		if err != xesOverflowError {
			if err != nil {
				t.Errorf("Expected to have %s but got: %s", xesOverflowError, err.Error())
			}
			t.Errorf("Expected to have %s but got: nil", xesOverflowError)
		}

		persistedPaymentItem, err = workflowPaymentsDB.Get(paymentID)
		if err != nil {
			panic(err)
		}
	})

}

func runTest(mockCtrl *gomock.Controller, workflowPaymentsDB *storm.WorkflowPaymentsDB,
	eventXesAmount *big.Int, eventTxHash, paymentID, paymentTxHash, from, to, status string, xesAmount *big.Int) error {

	adapterMock := NewMockadapter(mockCtrl)
	adapterMock.EXPECT().eventFromLog(gomock.Any(), gomock.Any(), gomock.Eq("Transfer")).Return(nil).Times(1)

	newPaymentItem := &model.WorkflowPaymentItem{
		ID:         paymentID,
		From:       from,
		To:         to,
		CreatedAt:  time.Now(),
		Status:     status,
		Xes:        xesAmount.Uint64(),
		WorkflowID: "1",
	}

	if paymentTxHash != "" {
		newPaymentItem.TxHash = paymentTxHash
	}

	err := workflowPaymentsDB.Save(newPaymentItem)
	if err != nil {
		return err
	}

	eventXes := xesAmount.Mul(eventXesAmount, big.NewInt(1000000000000000000))

	event := &XesMainTokenTransfer{
		Value:       eventXes,
		FromAddress: common.HexToAddress(from),
		ToAddress:   common.HexToAddress(to),
	}

	listener := &PaymentListener{
		xesAdapter:         adapterMock,
		workflowPaymentsDB: workflowPaymentsDB,
	}

	ethLog := &types.Log{
		TxHash: common.HexToHash(eventTxHash),
	}

	return listener.eventsHandler(ethLog, event)
}
