package ethglue

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

var ethDialler fakeETHDialler
var ethClient fakeETHClient

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
	ethDialler = fakeETHDialler{}
	ethClient = fakeETHClient{}
}

type fakeETHDialler struct {
	DialContextStub  func(ctx context.Context, rawUrl string) (ETHClientIF, error)
	DialContextCalls int

	DialStub  func(rawUrl string) (ethClient ETHClientIF, err error)
	DialCalls int
}

func (me *fakeETHDialler) Dial(rawUrl string) (ethClient ETHClientIF, err error) {
	me.DialCalls++

	if me.DialStub == nil {
		return nil, nil
	}

	return me.DialStub(rawUrl)
}

func (me *fakeETHDialler) DialContext(ctx context.Context, rawUrl string) (ETHClientIF, error) {
	me.DialContextCalls++

	if me.DialContextStub == nil {
		return nil, nil
	}

	return me.DialContextStub(ctx, rawUrl)
}

type fakeETHClient struct {
	HeaderByNumberStub  func(ctx context.Context, number *big.Int) (*types.Header, error)
	HeaderByNumberCalls int
}

func (me *fakeETHClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	me.HeaderByNumberCalls++

	if me.HeaderByNumberStub == nil {
		return nil, nil
	}

	return me.HeaderByNumberStub(ctx, number)
}

func (me *fakeETHClient) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	panic("implement me")
}
func (me *fakeETHClient) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	panic("implement me")
}
func (me *fakeETHClient) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	panic("implement me")
}
func (me *fakeETHClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	panic("implement me")
}
func (me *fakeETHClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) { panic("implement me") }
func (me *fakeETHClient) EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error) {
	panic("implement me")
}
func (me *fakeETHClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	panic("implement me")
}
func (me *fakeETHClient) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	panic("implement me")
}
func (me *fakeETHClient) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	panic("implement me")
}
