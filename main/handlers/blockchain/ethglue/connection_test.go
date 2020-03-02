package ethglue

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

var ethDialler FakeETHDialler
var ethClient FakeETHClient

func TestDialler_DialContext(t *testing.T) {
	ctx := context.Background()

	dialler := NewCustomDialler(&ethDialler, 5)

	t.Run("when ethClient instantiated", func(t *testing.T) {
		resetStubs()
		ethDialler.DialContextStub = func(ctx context.Context, rawUrl string) (clientIF ETHClientIF, err error) {
			return &ethClient, nil
		}

		dialler.DialContext(ctx, "ws://any.eth")

		assert.Equal(t, 1, ethDialler.DialContextCalls, "should have called DialContext")
		assert.Equal(t, 1, ethClient.HeaderByNumberCalls, "should call HEaderByNumber")
	})

	t.Run("when ethClient instantiate successfully and ethClient fails to retrieve block number", func(t *testing.T) {
		resetStubs()
		ethDialler.DialContextStub = func(ctx context.Context, rawUrl string) (clientIF ETHClientIF, err error) {
			return &ethClient, nil
		}
		ethClient.HeaderByNumberStub = func(ctx context.Context, number *big.Int) (header *types.Header, err error) {
			return nil, errors.New("any")
		}

		client, err := dialler.DialContext(ctx, "ws://any.eth")
		assert.NotNil(t, err, "should return error")
		assert.Nil(t, client, "should not return any client")
	})

	t.Run("when ethClient instantiate successfully and blockNumber is retrieved", func(t *testing.T) {
		resetStubs()
		ethDialler.DialContextStub = func(ctx context.Context, rawUrl string) (clientIF ETHClientIF, err error) {
			return &ethClient, nil
		}
		ethClient.HeaderByNumberStub = func(ctx context.Context, number *big.Int) (header *types.Header, err error) {
			return &types.Header{
				Number: big.NewInt(500),
			}, nil
		}

		client, err := dialler.DialContext(ctx, "ws://any.eth")

		assert.Nil(t, err, "should not return error")
		assert.Equal(t, &ethClient, client, "should return the client")
	})

	t.Run("when ethDialler fails first time", func(t *testing.T) {
		resetStubs()
		ethDialler.DialContextStub = func(ctx context.Context, rawUrl string) (clientIF ETHClientIF, err error) {
			return nil, errors.New("connection failed")
		}

		client, err := dialler.DialContext(ctx, "ws://any.eth")

		assert.NotNil(t, err, "should return error")
		assert.Nil(t, client, "should not return any client")
		assert.Equal(t, 0, ethClient.HeaderByNumberCalls, "should never try to retrieve header by number")
	})
}

func resetStubs() {
	ethDialler = FakeETHDialler{}
	ethClient = FakeETHClient{}
}
