package ethglue

import (
	"context"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/labstack/gommon/log"

	"github.com/ethereum/go-ethereum/common"
)

type NonceManager struct {
	nextNonce       *big.Int
	lastNonceChange time.Time

	ethClient  pendingNonceGetter
	ethAccount common.Address

	errStreakCount int64
	mu             sync.Mutex // protects next nonce updates
}

const maxErrStreak = 3
const idleSyncTimeInMinutes = 15

type pendingNonceGetter interface {
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
}

func (n *NonceManager) NextNonce() *big.Int {
	n.mu.Lock()
	defer n.mu.Unlock()

	if time.Now().Sub(n.lastNonceChange).Minutes() >= idleSyncTimeInMinutes {
		// this is the only way to fix the gap (when our nonce is higher than it should be)
		n.syncNonce()
		return n.nextNonce
	}

	no := new(big.Int).Set(n.nextNonce)
	n.nextNonce.Add(n.nextNonce, big.NewInt(1))
	n.lastNonceChange = time.Now()
	return no
}

func (n *NonceManager) OnAccountChange(addr string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.ethAccount = common.HexToAddress(addr)
	// nonce is per eth account
	n.syncNonce()
}

func (n *NonceManager) OnDial(c pendingNonceGetter) {
	n.ethClient = c
}

func (n *NonceManager) OnError(err error) {
	if err == nil {
		atomic.StoreInt64(&n.errStreakCount, 0)
		return
	}
	n.mu.Lock()
	defer n.mu.Unlock()

	if err.Error() == "nonce too low" {
		// either transaction(s) from different machine or
		// somehow transaction not failing instantly but not increasing nonce as well
		n.nextNonce = n.pendingNonceFromNode()
		log.Errorf("[NonceManager]: synced (forward jump) pendingNonce with network: %v", n.nextNonce.Int64())
		return
	}

	if err.Error() == "gas required exceeds allowance or always failing transaction" {
		// network nonce was for sure not increased - we reuse previous one
		n.nextNonce.Add(n.nextNonce, big.NewInt(-1))
		log.Errorf("[NonceManager]: decreasing nonce due to err: %v", err)
		return
	}

	// generic case unsure what to do
	n.errStreakCount++
	if n.errStreakCount >= maxErrStreak {
		n.errStreakCount = 0
		n.nextNonce = n.pendingNonceFromNode()
		log.Errorf("[NonceManager]: synced pendingNonce with network: %v", n.nextNonce.Int64())
		return
	}
	// maybe just reuse nonce
	n.nextNonce.Add(n.nextNonce, big.NewInt(-1))
	log.Errorf("[NonceManager]: decreasing nonce due to err: %v", err)
}

func (n *NonceManager) pendingNonceFromNode() *big.Int {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(20*time.Second))
	defer cancel()
	no, e := n.ethClient.PendingNonceAt(ctx, n.ethAccount)
	if e != nil {
		log.Errorf("[NonceManager]: pendingNonceAt err %v", e)
	}
	return new(big.Int).SetUint64(no)
}

func (n *NonceManager) syncNonce() {
	n.nextNonce = n.pendingNonceFromNode()
	n.lastNonceChange = time.Now()
}
