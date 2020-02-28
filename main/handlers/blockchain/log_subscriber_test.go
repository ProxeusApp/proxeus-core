package blockchain

import (
	"errors"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	"github.com/ProxeusApp/proxeus-core/main/handlers/blockchain/ethglue"

	"testing"
	"time"
)

var ethDialler *ethglue.FakeETHDialler
var ethClient *ethglue.FakeETHClient
var ctx context.Context
var cancel context.CancelFunc

func TestWebSocketLogSubscriber_Subscribe(t *testing.T) {
	reset()

	smartContractAddress := "0xef91ecd0142ae4c5163b2cf050c0563d49188c82"

	logs := make(chan<- types.Log, 200)

	fakeSubscription := fakeSubscription{}

	t.Run("when connection works", func(t *testing.T) {
		reset()
		defer cancel()

		sub := make(chan ethereum.Subscription)
		ethDialler.DialContextStub = func(ctx context.Context, rawUrl string) (clientIF ethglue.ETHClientIF, err error) {
			return ethClient, nil
		}
		ethClient.SubscribeFilterLogsStub = func(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (subscription ethereum.Subscription, err error) {
			return fakeSubscription, nil
		}

		socketLogSubscriber := NewWebSocketLogSubscriber(ethDialler, "ws://any.eth", smartContractAddress)
		go socketLogSubscriber.Subscribe(ctx, logs, sub)
		response := <-sub

		assert.Equal(t, fakeSubscription, response, "should send the Subscription to the channel")
		assert.Equal(t, 1, ethDialler.DialContextCalls)
		assert.Equal(t, 1, ethClient.SubscribeFilterLogsCalls, "should call SubscribeFilterLogs")
	})

	t.Run("when first connection fails", func(t *testing.T) {
		reset()
		defer cancel()

		sub := make(chan ethereum.Subscription)
		ethDialler.DialContextStub = func(ctx context.Context, rawUrl string) (clientIF ethglue.ETHClientIF, err error) {
			return nil, errors.New("connection problem")
		}
		ethClient.SubscribeFilterLogsStub = func(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (subscription ethereum.Subscription, err error) {
			return fakeSubscription, nil
		}

		socketLogSubscriber := NewWebSocketLogSubscriber(ethDialler, "ws://any.eth", smartContractAddress)
		go socketLogSubscriber.Subscribe(ctx, logs, sub)

		time.Sleep(500 * time.Millisecond) // After half a second, connection works again
		ethDialler.DialContextStub = func(ctx context.Context, rawUrl string) (clientIF ethglue.ETHClientIF, err error) {
			return ethClient, nil
		}

		response := <-sub

		assert.Equal(t, fakeSubscription, response, "should retry connection after 4 seconds and return the Subscription to the channel")
		assert.Equal(t, 2, ethDialler.DialContextCalls, "should have tried twice to connect to Dialler")
		assert.Equal(t, 1, ethClient.SubscribeFilterLogsCalls, "should have tried only once to subscribe to FilterLogs")
	})

	t.Run("when first subscription fails", func(t *testing.T) {
		reset()
		defer cancel()

		sub := make(chan ethereum.Subscription)
		ethDialler.DialContextStub = func(ctx context.Context, rawUrl string) (clientIF ethglue.ETHClientIF, err error) {
			return ethClient, nil
		}
		ethClient.SubscribeFilterLogsStub = func(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (subscription ethereum.Subscription, err error) {
			return nil, errors.New("first error")
		}

		socketLogSubscriber := NewWebSocketLogSubscriber(ethDialler, "ws://any.eth", smartContractAddress)
		go socketLogSubscriber.Subscribe(ctx, logs, sub)

		time.Sleep(1 * time.Second) // After a second becomes available again
		ethClient.SubscribeFilterLogsStub = func(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (subscription ethereum.Subscription, err error) {
			return fakeSubscription, nil
		}
		response := <-sub

		assert.Equal(t, fakeSubscription, response, "should retry after 4 seconds and return the Subscription to the channel")
		assert.Equal(t, 2, ethDialler.DialContextCalls, "should have tried twice to connect to Dialler")
		assert.Equal(t, 2, ethClient.SubscribeFilterLogsCalls, "should have tried twice to subscribe to FilterLogs")
	})

}

func reset() {
	ctx, cancel = context.WithTimeout(context.Background(), 8*time.Second)
	ethDialler = &ethglue.FakeETHDialler{}
	ethClient = &ethglue.FakeETHClient{}
}

type fakeSubscription struct {
}

func (f fakeSubscription) Unsubscribe() {
	panic("implement me")
}

func (f fakeSubscription) Err() <-chan error {
	// DO NOTHING
	return make(chan error)
}
