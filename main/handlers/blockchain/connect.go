package blockchain

import (
	"context"
	"log"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ProxeusApp/proxeus-core/main/handlers/blockchain/ethglue"
)

func ethConnectWebSocketsAsync(ctx context.Context, webSocketURL, contract string, logs chan<- types.Log, sub chan<- ethereum.Subscription) {

	filterAddresses := []common.Address{common.HexToAddress(contract)}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			var err error
			ethwsconn, err := ethglue.DialContext(ctx, webSocketURL)
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

func dummyEthConnect(ready chan<- ethereum.Subscription) {
	ready <- newDummySubscription()
}
