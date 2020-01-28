package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/ProxeusApp/proxeus-core/main/handlers/blockchain"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"log"
	"math/big"
	"strings"
	"sync"
)

type EthereumClient interface {
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
}

func NewEthClientBalanceService(ethClient EthereumClient) (*ethClientBalanceService, error) {
	erc20, err := abi.JSON(strings.NewReader(blockchain.ERC20ABI))
	if err != nil {
		return nil, err
	}

	return &ethClientBalanceService{
		ethClient: ethClient,
		smartContractTokensMap: map[string]string{
			"0xA017ac5faC5941f95010b12570B812C974469c2C": "XES",
			"0x9f8F72aA9304c8B593d555F12eF6589cC3A579A2": "MKR",
		},
		workersPoolSize: 8,
		erc20:           erc20,
	}, nil
}

type ethClientBalanceService struct {
	ethClient              EthereumClient
	smartContractTokensMap map[string]string
	erc20                  abi.ABI
	workersPoolSize        int
	balanceLock            sync.Mutex
}

type job struct {
	ctx          context.Context
	address      string
	startBlock   *big.Int
	endBlock     *big.Int
	balancesMap  *sync.Map
	jobsDoneChan chan bool
}

// Retrieves balances for an Ethereum address. Given an address in hexadecimal format, will return a *sync.Map of string->*big.Int, containing listed ERC20 tokens from "smartContractTokensMap" + "ETH".
// Balances are measured in wei and every ERC20 token might have a different "decimals" amount. Please refer to the token to get that number
// For ex. if "smartContractTokensMap" contains ETH, XES and MKR, calling this function will return you a map in the following format:
//
// *sync.Map {
//    "ETH": *big.Int(77524316000000000000000000),
//    "XES": *big.Int(902349094030234904393949),
//    "MKR": *big.Int(43535435324324242423333),
// }
//
// Be careful with addresses all together as they should be in format EIP-55, not lowercase!
func (me *ethClientBalanceService) GetBalancesForAddress(ctx context.Context, address string) (*sync.Map, error) {
	var toBlockNumber *big.Int = nil // Last block

	// Make sure block number exists, and retrieve it (in case of nil, will return the last block)
	blockHeader, err := me.ethClient.HeaderByNumber(ctx, toBlockNumber)
	if err != nil {
		return nil, fmt.Errorf("block %d not found. error: %v", toBlockNumber, err)
	}

	// Retrieve ether's balance
	ethBalance, err := me.ethClient.BalanceAt(ctx, common.HexToAddress(address), blockHeader.Number)
	if err != nil {
		return nil, fmt.Errorf("retrieving balance of %s. error: %v", address, err)
	}

	// Retrieve all ERC20 token balances (listed in smartContractTokensMap)
	balances, err := me.extractERC20Balances(ctx, blockHeader.Number, address)
	if err != nil {
		return nil, err
	}

	balances.Store("ETH", ethBalance)

	log.Println("Total balances", balances)

	return balances, nil
}

func (me *ethClientBalanceService) extractERC20Balances(ctx context.Context, toBlockNumber *big.Int, address string) (*sync.Map, error) {
	var (
		balancesMap  sync.Map
		jobsChan     = make(chan job, 1000)
		jobsDoneChan = make(chan bool, 100)
		errChan      = make(chan error, 1)
	)
	defer close(jobsChan)
	defer close(jobsDoneChan)
	defer close(errChan)

	// Split into block chunks as we don't want (can't) to process the whole blockchain at once
	startBlocks, endBlocks, err := me.getBlockChunks(big.NewInt(int64(0)), toBlockNumber, 3000)
	if err != nil {
		return nil, err
	}

	// Create a pool of workers
	for workerId := 1; workerId <= me.workersPoolSize; workerId++ {
		go me.worker(workerId, jobsChan, errChan)
	}

	go func() {
		for i, startBlock := range startBlocks {
			endBlock := endBlocks[i]

			jobsChan <- job{
				ctx:          ctx,
				address:      address,
				startBlock:   startBlock,
				endBlock:     endBlock,
				balancesMap:  &balancesMap,
				jobsDoneChan: jobsDoneChan,
			}
		}
		log.Println("All jobs sent to chan")
	}()

	// Wait until all jobs are processed. Each one could return an error
	for r := 0; r < len(startBlocks); r++ {
		select {
		case err := <-errChan:
			log.Printf("An error occurred %v", err)
			return nil, err
		case <-jobsDoneChan:
		}
	}

	return &balancesMap, nil
}

// Expensive operation of retrieving all event logs between two blocks. Whenever a "job" is sent to the worker,
// it first processes it, then waits for another.
func (me *ethClientBalanceService) worker(workerId int, jobs <-chan job, errChan chan error) {
	log.Printf("Worker %d ready", workerId)

	for job := range jobs {
		// Find events "Transfer" on all defined ERC20's smart contracts
		query := ethereum.FilterQuery{
			Addresses: me.smartContractAddresses(),
			FromBlock: job.startBlock,
			ToBlock:   job.endBlock,
			Topics: [][]common.Hash{{
				me.erc20.Events["Transfer"].ID(),
			},
			},
		}

		logs, err := me.ethClient.FilterLogs(job.ctx, query)
		if err != nil {
			errChan <- err
		}

		for _, eventLog := range logs {
			tokenCode, found := me.smartContractTokensMap[eventLog.Address.Hex()]
			if !found {
				log.Printf("Token %s not found, we don't have a mapping to smart contract. address %s", tokenCode, eventLog.Address.Hex())
				continue
			}

			me.balanceLock.Lock()

			transferEvent, err := me.parseTransferEventFromLog(eventLog)
			if err != nil {
				errChan <- err
			}

			balanceInterface, _ := job.balancesMap.LoadOrStore(tokenCode, big.NewInt(0))
			addressBalance := balanceInterface.(*big.Int)

			if transferEvent.IsReceiver(job.address) {
				newBalance := big.NewInt(0).Add(addressBalance, transferEvent.Value)
				log.Printf("Detected an incoming transfer of %d %s (txHash %s). New balance: %d", transferEvent.Value, tokenCode, eventLog.TxHash.Hex(), newBalance)
				job.balancesMap.Store(tokenCode, newBalance)
			}

			if transferEvent.IsSender(job.address) {
				newBalance := big.NewInt(0).Sub(addressBalance, transferEvent.Value)
				log.Printf("Detected an outgoing transfer of %d %s (txHash %s). New balance: %d", transferEvent.Value, tokenCode, eventLog.TxHash.Hex(), newBalance)
				job.balancesMap.Store(tokenCode, newBalance)
			}

			me.balanceLock.Unlock()
		}

		job.jobsDoneChan <- true
	}
}

func (me *ethClientBalanceService) parseTransferEventFromLog(eventLog types.Log) (blockchain.ERC20TransferEvent, error) {
	transferEvent := blockchain.ERC20TransferEvent{}
	err := me.erc20.Unpack(&transferEvent, "Transfer", eventLog.Data)
	if err != nil {
		return transferEvent, fmt.Errorf("unpacking 'Transfer' from Data from transaction %s. Error %v", eventLog.TxHash.Hex(), err)
	}

	transferEvent.From = common.BytesToAddress(eventLog.Topics[1].Bytes())
	transferEvent.To = common.BytesToAddress(eventLog.Topics[2].Bytes())
	return transferEvent, nil
}

func (me *ethClientBalanceService) smartContractAddresses() []common.Address {
	addresses := make([]common.Address, len(me.smartContractTokensMap))

	i := 0
	for key, _ := range me.smartContractTokensMap {
		addresses[i] = common.HexToAddress(key)
		i++
	}

	return addresses
}

/*
   Splits blocks in chunks of size. For example, given startBlock 10, toBlock 30 and size 3
   should return an array with  [10, 13], [14, 17], [18, 21],..
*/
func (me *ethClientBalanceService) getBlockChunks(startBlock *big.Int, toBlock *big.Int, size int) ([]*big.Int, []*big.Int, error) {
	if startBlock == nil || toBlock == nil {
		return nil, nil, errors.New("startBlock and toBlock parameters can't be nil")
	}

	var startBlocks []*big.Int
	var toBlocks []*big.Int

	counter := 0
	for i := startBlock; i.CmpAbs(toBlock) < 0; i.Add(i, big.NewInt(int64(size))) {
		startBlockTemp := new(big.Int).Set(i) // we want it by value not reference
		toBlockTemp := new(big.Int).Set(i)
		toBlockTemp.Add(toBlockTemp, big.NewInt(int64(size)))

		if toBlockTemp.Cmp(toBlock) > 0 {
			toBlockTemp.Set(toBlock)
		}

		if counter == 0 {
			startBlocks = append(startBlocks, startBlockTemp)
		} else {
			startBlockTemp.Add(startBlockTemp, big.NewInt(1))
			startBlocks = append(startBlocks, startBlockTemp)
		}

		toBlocks = append(toBlocks, toBlockTemp)
		counter++
	}

	return startBlocks, toBlocks, nil
}
