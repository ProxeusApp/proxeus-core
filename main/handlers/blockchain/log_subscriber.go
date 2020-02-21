package blockchain

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ProxeusApp/proxeus-core/main/handlers/blockchain/ethglue"
)

type LogSubscriber interface {
	Subscribe(ctx context.Context, logs chan<- types.Log, sub chan<- ethereum.Subscription)
	GetContractAddress() string
	String() string
}

type webSocketLogSubscriber struct {
	webSocketURL string
	contract     string
	logs         chan<- types.Log
}

func NewWebSocketLogSubscriber(webSocketURL, contract string) *webSocketLogSubscriber {
	return &webSocketLogSubscriber{
		webSocketURL: webSocketURL,
		contract:     contract,
	}
}

func (c *webSocketLogSubscriber) Subscribe(ctx context.Context, logs chan<- types.Log, sub chan<- ethereum.Subscription) {

	filterAddresses := []common.Address{common.HexToAddress(c.contract)}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			var err error
			ethwsconn, err := ethglue.DialContext(ctx, c.webSocketURL)
			if err != nil {
				log.Printf("failed to dial for eth events, will retry (%s)\n", err)
				continue
			}
			query := ethereum.FilterQuery{
				Addresses: filterAddresses,
			}
			ctx, cancel := context.WithTimeout(ctx, time.Duration(10*time.Second))
			s, err := ethwsconn.SubscribeFilterLogs(ctx, query, logs)
			cancel()
			if err != nil {
				log.Printf("failed to subscribe for eth events, will retry (%s)\n", err)
				time.Sleep(time.Second * 4)
				continue
			}
			// success!
			sub <- s
			return
		}
	}
}

func (c *webSocketLogSubscriber) GetContractAddress() string {
	return c.contract
}

func (c *webSocketLogSubscriber) String() string {
	return fmt.Sprintf("webSocketLogSubscriber(%s, %s)", c.webSocketURL, c.contract)
}

// Dummy ethereum connection for test mode
type dummySubscription struct {
	err chan error
}

func newDummySubscription() *dummySubscription {
	return &dummySubscription{
		err: make(chan error),
	}
}

func (s *dummySubscription) Unsubscribe() {}
func (s *dummySubscription) Err() <-chan error {
	return s.err
}

type dummyLogSubscriber struct{}

func NewDummyLogSubscriber() *dummyLogSubscriber {
	return &dummyLogSubscriber{}
}

func (c *dummyLogSubscriber) Subscribe(_ context.Context, _ chan<- types.Log, sub chan<- ethereum.Subscription) {
	sub <- newDummySubscription()
}

func (c *dummyLogSubscriber) GetContractAddress() string {
	return "dummyContract"
}

func (c *dummyLogSubscriber) String() string {
	return "dummyLogSubscriber"
}
