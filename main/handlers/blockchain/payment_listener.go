package blockchain

import (
	"context"
	"errors"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"git.proxeus.com/core/central/main/ethglue"
	"git.proxeus.com/core/central/sys/db/storm"

	strm "github.com/asdine/storm"
)

type (
	listener struct {
		logs            chan types.Log
		ethWebSocketURL string
		ethURL          string
		sub             ethereum.Subscription
	}
	PaymentListener struct {
		listener
		workflowPaymentsDB storm.WorkflowPaymentsDBInterface
		xesAdapter         adapter
	}
)

func NewPaymentListener(xesAdapter adapter, ethWebSocketURL, ethURL string, workflowPaymentsDB storm.WorkflowPaymentsDBInterface) *PaymentListener {
	me := &PaymentListener{}
	me.xesAdapter = xesAdapter
	me.ethWebSocketURL = ethWebSocketURL
	me.ethURL = ethURL
	me.workflowPaymentsDB = workflowPaymentsDB
	me.logs = make(chan types.Log, 200)
	return me
}

func (me *PaymentListener) Listen(ctx context.Context) {
	var readyCh <-chan struct{}

	for {
		readyCh = me.ethConnectWebSocketsAsync(ctx)
		select {
		case <-readyCh:
			log.Println("[paymentlistener] listen on contract started. contract address: ", me.xesAdapter.getContractAddress())
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
			err := me.eventsHandler(&vLog, event)
			if err != nil {
				if err != strm.ErrNotFound { //ErrNotFound already logged in eventsHandler
					if err == xesOverflowError {
						log.Fatal("[blockchain][listener] overflow err: ", err.Error())
					}
					log.Println("[blockchain][listener] err: ", err.Error())
				}
			}
		}
	}
}

func (me *PaymentListener) ethConnectWebSocketsAsync(ctx context.Context) <-chan struct{} {

	filterAddresses := []common.Address{common.HexToAddress(me.xesAdapter.getContractAddress())}

	readyCh := make(chan struct{})
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				var err error
				ethwsconn, err := ethglue.DialContext(ctx, me.ethWebSocketURL)
				if err != nil {
					log.Printf("failed to dial for eth events, will retry (%s)\n", err)
					time.Sleep(time.Second * 4)
					continue
				}
				query := ethereum.FilterQuery{
					Addresses: filterAddresses,
				}
				ctx, cancel := context.WithTimeout(ctx, time.Duration(10*time.Second))
				me.sub, err = ethwsconn.SubscribeFilterLogs(ctx, query, me.logs)
				cancel()
				if err != nil {
					log.Printf("failed to subscribe for eth events, will retry (%s)\n", err)
					time.Sleep(time.Second * 4)
					continue
				}
				// success!
				readyCh <- struct{}{}
				return
			}
		}
	}()
	return readyCh
}

var xesOverflowError = errors.New("overflow on xes event")

func (me *PaymentListener) eventsHandler(lg *types.Log, event *XesMainTokenTransfer) error {
	log.Printf("[PaymentListener][eventHandler] txHash: %s, value: %s, %v",
		lg.TxHash.String(), event.Value.String(), lg)
	err := me.xesAdapter.eventFromLog(event, lg, "Transfer")
	if err != nil {
		return err
	}

	bigXes := event.Value.Div(event.Value, big.NewInt(1000000000000000000)) //to xes-ether

	if !bigXes.IsUint64() {
		log.Println("[PaymentListener][eventHandler]  error overflow on transfer event value:", event.Value)
		return xesOverflowError
	}

	err = me.workflowPaymentsDB.ConfirmPayment(lg.TxHash.String(), event.FromAddress.String(), event.ToAddress.String(), bigXes.Uint64())
	if err != nil {
		if err == strm.ErrNotFound {
			log.Printf(" [PaymentListener][eventHandler] info: no matching payment found for txHash: %s, reason: %s", lg.TxHash.String(), err.Error())
			return err
		}

		log.Printf(" [PaymentListener][eventHandler] err: workflowPaymentsDB.ConfirmPayment for txHash: %s, err: %s", lg.TxHash.String(), err.Error())
		return err
	}
	log.Println("[PaymentListener][eventHandler] confirmed payment with hash: ", lg.TxHash.String())

	return nil
}
