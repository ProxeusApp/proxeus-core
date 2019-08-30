package blockchain

import (
	"log"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"

	"git.proxeus.com/core/central/sys/db/storm"
	"git.proxeus.com/core/central/sys/model"
)

type (
	workflowPaymentItemMatcher struct {
		transactionHash string
		from            string
		to              string
		xes             uint64
	}
)

func (me *workflowPaymentItemMatcher) Matches(x interface{}) bool {
	workflowPaymentItem, ok := x.(*model.WorkflowPaymentItem)
	if !ok {
		log.Fatal("workflowPaymentItemMatcher cast error")
	}
	return workflowPaymentItem.WorkflowID == "" &&
		me.transactionHash == workflowPaymentItem.TxHash &&
		me.from == workflowPaymentItem.From &&
		me.to == workflowPaymentItem.To &&
		me.xes == workflowPaymentItem.Xes
}
func (me *workflowPaymentItemMatcher) String() string {
	return "workflowPaymentItem needs to match"
}

func TestCheckIfWorkflowNeedsPayment(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("CheckIfEventsHandlerSavesBlockchainEvent1Xes", func(t *testing.T) {
		err := runTest(mockCtrl, big.NewInt(1), true)
		if err != nil {
			log.Fatal(err.Error())
		}
	})
	t.Run("CheckIfEventsHandlerSavesBlockchainEvent130Xes", func(t *testing.T) {
		err := runTest(mockCtrl, big.NewInt(130), true)
		if err != nil {
			log.Fatal(err.Error())
		}
	})
	t.Run("CheckIfEventsHandlerSavesBlockchainEvent100000000000000000Xes", func(t *testing.T) {
		err := runTest(mockCtrl, big.NewInt(100000000000000000), true)
		if err != nil {
			log.Fatal(err.Error())
		}
	})
	t.Run("CheckIfEventsHandlerShowsOverflowError", func(t *testing.T) {
		xesTmp := big.NewInt(1000000000000000000)
		xes := xesTmp.Mul(xesTmp, big.NewInt(1000000000000000000))
		err := runTest(mockCtrl, xes, false)
		if err != xesOverflowError {
			if err != nil {
				log.Fatal("expected xesOverflowError but got: ", err.Error())
			}
			log.Fatal("expected xesOverflowError but got: nil")
		}
	})
}

func runTest(mockCtrl *gomock.Controller, xesAmount *big.Int, addPaymentExpected bool) error {
	adapterMock := NewMockadapter(mockCtrl)
	adapterMock.EXPECT().eventFromLog(gomock.Any(), gomock.Any(), gomock.Eq("Transfer")).Return(nil).Times(1)

	transactionHash := "0x04f1bbf224b5876d91c74984c4d7f7768c5cc9da5b7f7afe1a31ef9115310f67"
	from := "0xe4f2604bc8300004aa4477af5f2AdBd37765F3F7"
	to := "0xe902Fb81617079236cB6eF8f34b2A1e759ef676D"

	xesUint := xesAmount.Uint64()

	xes := xesAmount.Mul(xesAmount, big.NewInt(1000000000000000000))

	workflowPaymentItemMatcher := &workflowPaymentItemMatcher{
		transactionHash: transactionHash,
		from:            from,
		to:              to,
		xes:             xesUint,
	}
	paymentsDBMock := storm.NewMockWorkflowPaymentsDBInterface(mockCtrl)

	if addPaymentExpected {
		paymentsDBMock.EXPECT().Add(workflowPaymentItemMatcher)
	}

	event := &XesMainTokenTransfer{
		Value:       xes,
		FromAddress: common.HexToAddress(from),
		ToAddress:   common.HexToAddress(to),
	}

	listener := &Paymentlistener{
		xesAdapter:         adapterMock,
		workflowPaymentsDB: paymentsDBMock,
	}

	ethLog := &types.Log{
		TxHash: common.HexToHash(transactionHash),
	}

	return listener.eventsHandler(ethLog, event)
}
