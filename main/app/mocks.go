package app

import (
	"context"
	"github.com/ProxeusApp/proxeus-core/main/handlers/blockchain"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"strings"
)

type ethClientStub struct {
	EthBalance *big.Int
	XESBalance *big.Int
	erc20ABI   abi.ABI
}

func NewEthClientStub() *ethClientStub {
	erc20ABI, err := abi.JSON(strings.NewReader(blockchain.ERC20ABI))
	if err != nil {
		panic(err)
	}

	xesBalance := big.Int{}
	xesBalance.SetString("77524316000000000000000000", 10)

	return &ethClientStub{
		EthBalance: big.NewInt(12345674000000000),
		XESBalance: &xesBalance,
		erc20ABI:   erc20ABI,
	}
}

func (me ethClientStub) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {

	return &types.Header{
		Number: big.NewInt(1000),
	}, nil
}

func (me ethClientStub) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {

	return me.EthBalance, nil
}

func (me ethClientStub) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	xesSmartContractAddress := common.HexToAddress("0xA017ac5faC5941f95010b12570B812C974469c2C")
	targetAddress := common.HexToHash("0x043129ab3945D2bB75f3B5DE21487343EFBeffd2")

	// We're going to simulate three transfers of XES. 2 incoming, 1 outgoing + another non-existing
	firstTransferValue := big.Int{}
	firstTransferValue.SetString("10000000000000000000000000", 10)

	secondTransferValue := big.Int{}
	secondTransferValue.SetString("55000000000000000000000000", 10)

	thirdTransferValue := big.Int{}
	thirdTransferValue.SetString("25000000000000000000000000", 10)

	return []types.Log{
		// incoming XES transfer
		{
			Address: xesSmartContractAddress,
			Topics: []common.Hash{
				common.BytesToHash(crypto.Keccak256(me.erc20ABI.Methods["transfer"].ID())), // keccac256(transfer(to address, _value uint256))
				common.HexToHash("0xef91ecd0142ae4c5163b2cf060c0563d49188c82"),             // From
				targetAddress, // To
			},
			Data:        common.BytesToHash(firstTransferValue.Bytes()).Bytes(),
			BlockNumber: 500,
			TxHash:      common.HexToHash("0x640200e1f5b8ebd7bc3147c50437e5e74600959e58b5c4e1a1da2803e1b8663c"),
			Removed:     false,
		},
		// a second XES transfer
		{
			Address: xesSmartContractAddress,
			Topics: []common.Hash{
				common.BytesToHash(crypto.Keccak256(me.erc20ABI.Methods["transfer"].ID())),
				common.HexToHash("0xef91ecd0142ae4c5163b2cf060c0563d49188c82"),
				targetAddress,
			},
			Data:        common.BytesToHash(secondTransferValue.Bytes()).Bytes(),
			BlockNumber: 505,
			TxHash:      common.HexToHash("0x140200e1f5b8ebd7bc3147c50437e5e74600959e58b5c4e1a1da2803e1b8663c"),
			Removed:     false,
		},
		// another, NON XES transfer (non-existing ERC-20 token)
		{
			Address: common.HexToAddress("0xabcdefghiC5941f95010b12570B812C974469c2C"),
			Topics: []common.Hash{
				common.BytesToHash(crypto.Keccak256(me.erc20ABI.Methods["transfer"].ID())),
				common.HexToHash("0xef91ecd0142ae4c5163b2cf060c0563d49188c82"),
				common.HexToHash("0x043129ab3945D2bB75f3B5DE21487343EFBeffd2"),
			},
			Data:        common.BytesToHash(secondTransferValue.Bytes()).Bytes(),
			BlockNumber: 505,
			TxHash:      common.HexToHash("0x520200e1f5b8ebd7bc3147c50437e5e74600959e58b5c4e1a1da2803e1b8663c"),
			Removed:     false,
		},
		// the third XES transfer, this time it's outgoing
		{
			Address: xesSmartContractAddress,
			Topics: []common.Hash{
				common.BytesToHash(crypto.Keccak256(me.erc20ABI.Methods["transfer"].ID())),
				targetAddress,
				common.HexToHash("0xef91ecd0142ae4c5163b2cf060c0563d49188c82"),
			},
			Data:        common.BytesToHash(thirdTransferValue.Bytes()).Bytes(),
			BlockNumber: 507,
			TxHash:      common.HexToHash("0x140200e1f5b8ebd7bc3147c50437e5e74600959e58b5c4e1a1da2803e1b8663c"),
			Removed:     false,
		},
	}, nil
}
