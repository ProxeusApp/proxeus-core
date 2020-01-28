package app

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestEthClientBalanceService_parseTransferEventFromLog(t *testing.T) {
	ctx := context.Background()

	balanceService, err := NewEthClientBalanceService(NewEthClientStub())
	assert.Nil(t, err)

	balances, err := balanceService.GetBalancesForAddress(ctx, "0x043129ab3945D2bB75f3B5DE21487343EFBeffd2")

	assert.Nil(t, err)

	ethBalance, ethFound := balances.Load("ETH")
	xesBalance, xesFound := balances.Load("XES")

	assert.True(t, ethFound)
	assert.Equal(t, big.NewInt(12345674000000000), ethBalance)
	assert.True(t, xesFound)

	expectedXES := big.Int{}
	expectedXES.SetString("40000000000000000000000000", 10)
	assert.Equal(t, &expectedXES, xesBalance)
}

func TestEthClientBalanceService_smartContractAddresses(t *testing.T) {
	balanceService := ethClientBalanceService{
		smartContractTokensMap: map[string]string{
			"0x1217AC5fAC5941F95010B12570b812c974469c98": "CODE",
			"0x6B175474e89094C44da98B954eedeAC495271D12": "CODE2",
		},
	}
	addresses := balanceService.smartContractAddresses()
	addressExistsMap := addressesToMap(addresses)

	assert.True(t, addressExistsMap["0x1217AC5fAC5941F95010B12570b812c974469c98"])
	assert.True(t, addressExistsMap["0x6B175474e89094C44da98B954eedeAC495271D12"])
}

func addressesToMap(addresses []common.Address) map[string]bool {
	addressExistsMap := make(map[string]bool)
	for _, address := range addresses {
		addressExistsMap[address.Hex()] = true
	}
	return addressExistsMap
}

func TestEthereumBlockchainConnector_getBlockChunks(t *testing.T) {
	balanceService := ethClientBalanceService{}

	t.Run("normal chunk", func(t *testing.T) {
		// I expect [10 21 31 41] and [20 30 40 42]
		startBlocks, endBlocks, err := balanceService.getBlockChunks(big.NewInt(10), big.NewInt(42), 10)

		assert.Nil(t, err)
		assert.Equal(t, []*big.Int{big.NewInt(10), big.NewInt(21), big.NewInt(31), big.NewInt(41)}, startBlocks)
		assert.Equal(t, []*big.Int{big.NewInt(20), big.NewInt(30), big.NewInt(40), big.NewInt(42)}, endBlocks)
	})

	t.Run("nil as toBlock", func(t *testing.T) {
		_, _, err := balanceService.getBlockChunks(big.NewInt(10), nil, 10)

		assert.NotNil(t, err, "no nil should be accepted as input block")
	})

	t.Run("nil as fromBlock", func(t *testing.T) {
		_, _, err := balanceService.getBlockChunks(nil, big.NewInt(10), 10)

		assert.NotNil(t, err, "no nil should be accepted as input block")
	})
}
