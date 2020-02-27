package ethglue

import (
	"context"
	"math/big"

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
}
