package ethglue

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common"
)

type pendingNonceMock struct{}

func (_ pendingNonceMock) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	panic("implement me")
}

func (_ pendingNonceMock) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	panic("implement me")
}

func (_ pendingNonceMock) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	panic("implement me")
}

func (_ pendingNonceMock) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	panic("implement me")
}

func (_ pendingNonceMock) EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error) {
	panic("implement me")
}

func (_ pendingNonceMock) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	panic("implement me")
}

func (_ pendingNonceMock) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	panic("implement me")
}

func (_ pendingNonceMock) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	panic("implement me")
}

func (_ pendingNonceMock) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	panic("implement me")
}

func (_ pendingNonceMock) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return 1, nil
}

func TestNonceManagerSequence(t *testing.T) {
	var m NonceManager
	m.OnDial(pendingNonceMock{})
	m.OnAccountChange("0xdd")

	n := m.NextNonce()
	exp := big.NewInt(1)
	if n.Cmp(exp) != 0 {
		t.Errorf("got nonce %v expected %v", n, exp)
	}

	n = m.NextNonce()
	exp = big.NewInt(2)
	if n.Cmp(exp) != 0 {
		t.Errorf("got nonce %v expected %v", n, exp)
	}

	m.OnAccountChange("0xff")

	n = m.NextNonce()
	exp = big.NewInt(1)
	if n.Cmp(exp) != 0 {
		t.Errorf("got nonce %v expected %v", n, exp)
	}

	m.OnError(errors.New("nonce too low"))
	n = m.NextNonce()
	exp = big.NewInt(1)
	if n.Cmp(exp) != 0 {
		t.Errorf("got nonce %v expected %v", n, exp)
	}

	m.OnError(errors.New("gas required exceeds allowance or always failing transaction"))
	n = m.NextNonce()
	exp = big.NewInt(1)
	if n.Cmp(exp) != 0 {
		t.Errorf("got nonce %v expected %v", n, exp)
	}

	m.OnError(errors.New("some err"))
	n = m.NextNonce()
	exp = big.NewInt(1)
	if n.Cmp(exp) != 0 {
		t.Errorf("got nonce %v expected %v", n, exp)
	}
	m.OnError(nil)
}
