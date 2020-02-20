package blockchain

import (
	"context"
	"errors"
	"log"
	"math/big"

	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/storage/database/db"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

type listener struct {
	logs            chan types.Log
	ethWebSocketURL string
	ethURL          string
	sub             ethereum.Subscription
}

type PaymentListener struct {
	listener
	workflowPaymentsDB storage.WorkflowPaymentsIF
	xesAdapter         Adapter
	testMode           bool
}

var TestChannelPayment chan types.Log

func NewPaymentListener(xesAdapter Adapter, ethWebSocketURL, ethURL string, workflowPaymentsDB storage.WorkflowPaymentsIF, testMode bool) *PaymentListener {
	me := &PaymentListener{}
	me.xesAdapter = xesAdapter
	me.ethWebSocketURL = ethWebSocketURL
	me.ethURL = ethURL
	me.workflowPaymentsDB = workflowPaymentsDB
	me.logs = make(chan types.Log, 200)
	TestChannelPayment = me.logs

	me.testMode = testMode

	return me
}

func (me *PaymentListener) Listen(ctx context.Context) {
	subscription := make(chan ethereum.Subscription)

	for {
		if me.testMode {
			go dummyEthConnect(subscription)
		} else {
			go ethConnectWebSocketsAsync(ctx, me.ethWebSocketURL, me.xesAdapter.GetContractAddress(), me.logs, subscription)
		}
		select {
		case sub := <-subscription:
			me.sub = sub
			log.Println("[paymentlistener] listen on contract started. contract address: ", me.xesAdapter.GetContractAddress())
			reconnect := me.listenLoop(ctx)
			if !reconnect {
				log.Printf("[paymentlistener][eventHandler] finished")
				return
			}
		case <-ctx.Done():
			log.Printf("[paymentlistener][eventHandler] done")
			return
		}
	}
	return
}

func (me *PaymentListener) listenLoop(ctx context.Context) (shouldReconnect bool) {
	for {
		select {
		case <-ctx.Done():
			return false
		case err, ok := <-me.sub.Err():
			if !ok {
				return true
			}
			log.Println("ERROR sub", err)
			return true
		case vLog, ok := <-me.logs:
			if !ok {
				return true
			}
			event := new(XesMainTokenTransfer)
			me.eventsHandler(&vLog, event)
		}
	}
}

var xesOverflowError = errors.New("overflow on xes event")

func (me *PaymentListener) eventsHandler(lg *types.Log, event *XesMainTokenTransfer) error {
	log.Printf("[PaymentListener][eventHandler] txHash: %s, value: %s, %v",
		lg.TxHash.String(), event.Value.String(), lg)
	err := me.xesAdapter.EventFromLog(event, lg, "Transfer")
	if err != nil {
		log.Println("[blockchain][listener] err: ", err.Error())
		return err
	}

	bigXes := event.Value.Div(event.Value, big.NewInt(1000000000000000000)) //to xes-ether

	if !bigXes.IsUint64() {
		log.Println("[PaymentListener][eventHandler]  error overflow on transfer event value:", event.Value)
		return xesOverflowError
	}

	err = me.workflowPaymentsDB.ConfirmPayment(lg.TxHash.String(), event.FromAddress.String(), event.ToAddress.String(), bigXes.Uint64())
	if err != nil {
		if db.NotFound(err) {
			log.Printf(" [PaymentListener][eventHandler] info: no matching payment found for txHash: %s, reason: %s", lg.TxHash.String(), err.Error())
			return err
		}

		log.Printf(" [PaymentListener][eventHandler] err: workflowPaymentsDB.ConfirmPayment for txHash: %s, err: %s", lg.TxHash.String(), err.Error())
		return err
	}
	log.Println("[PaymentListener][eventHandler] confirmed payment with hash: ", lg.TxHash.String())
	return nil
}
