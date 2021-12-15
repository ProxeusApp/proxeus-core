package ethglue

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

type ETHDiallerIF interface {
	Dial(rawUrl string) (ethClient ETHClientIF, err error)
	DialContext(ctx context.Context, rawUrl string) (ETHClientIF, error)
}

type ETHClientIF interface {
	bind.ContractBackend
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	SuggestGasTipCap(ctx context.Context) (*big.Int, error)
}

type FakeETHDialler struct {
	DialContextStub  func(ctx context.Context, rawUrl string) (ETHClientIF, error)
	DialContextCalls int

	DialStub  func(rawUrl string) (ethClient ETHClientIF, err error)
	DialCalls int
}

func (me *FakeETHDialler) Dial(rawUrl string) (ethClient ETHClientIF, err error) {
	me.DialCalls++

	if me.DialStub == nil {
		return nil, nil
	}

	return me.DialStub(rawUrl)
}

func (me *FakeETHDialler) DialContext(ctx context.Context, rawUrl string) (ETHClientIF, error) {
	me.DialContextCalls++

	if me.DialContextStub == nil {
		return nil, nil
	}

	return me.DialContextStub(ctx, rawUrl)
}

type FakeETHClient struct {
	HeaderByNumberStub       func(ctx context.Context, number *big.Int) (*types.Header, error)
	HeaderByNumberCalls      int
	SubscribeFilterLogsStub  func(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error)
	SubscribeFilterLogsCalls int
}

func (me *FakeETHClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	me.HeaderByNumberCalls++

	if me.HeaderByNumberStub == nil {
		return nil, nil
	}

	return me.HeaderByNumberStub(ctx, number)
}

func (me *FakeETHClient) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	panic("implement me")
}
func (me *FakeETHClient) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	panic("implement me")
}
func (me *FakeETHClient) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	panic("implement me")
}
func (me *FakeETHClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	panic("implement me")
}
func (me *FakeETHClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	panic("implement me")
}
func (me *FakeETHClient) EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error) {
	panic("implement me")
}
func (me *FakeETHClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	panic("implement me")
}
func (me *FakeETHClient) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	panic("implement me")
}
func (me *FakeETHClient) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	me.SubscribeFilterLogsCalls++

	if me.SubscribeFilterLogsStub == nil {
		return nil, nil
	}

	return me.SubscribeFilterLogsStub(ctx, query, ch)
}
func (me *FakeETHClient) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	panic("implement me")
}
