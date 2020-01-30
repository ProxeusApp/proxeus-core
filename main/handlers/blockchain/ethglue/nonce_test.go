package ethglue

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

type pendingNonceMock struct{}

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
