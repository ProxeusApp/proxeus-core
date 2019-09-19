package blockchain

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strconv"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	"git.proxeus.com/core/central/main/config"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"git.proxeus.com/core/central/main/ethglue"
)

const etherUnit = 1000000000000000000.0

var conn *ethclient.Client
var nonceManager ethglue.NonceManager

var mu sync.Mutex

func GiveTokens(toWallet string) {
	var err error
	conn, err = ethglue.Dial(config.Config.EthClientURL)
	if err != nil {
		log.Panic("[airdrop] Failed to connect to the Ethereum client:", err)
	}
	nonceManager.OnDial(conn)

	type Web3Keystore struct {
		Address string `json:"address"`
	}
	var keystore Web3Keystore

	keystoreJSON, err := ioutil.ReadFile(config.Config.AirdropWalletfile)
	if err != nil {
		log.Panic("[airdrop] Failed to read keystore:", err)
	}
	err = json.Unmarshal(keystoreJSON, &keystore)
	if err != nil {
		log.Panic("[airdrop] Failed to parse keystore:", err)
	}

	address := keystore.Address
	nonceManager.OnAccountChange(address)
	mu.Lock()
	defer mu.Unlock()
	FreeXES(toWallet)
	FreeEth(toWallet)
}

// FreeXES sends ropsten XES to given wallet address
func FreeXES(walletAddress string) {
	log.Println("[airdrop] [Ropsten] Prepare XES for addr:", walletAddress)
	amount := new(big.Int)
	f, err := strconv.ParseFloat(config.Config.AirdropAmountXES, 64)

	amount.SetInt64(int64(f * etherUnit))

	// make transfer
	token, err := NewToken(common.HexToAddress(config.Config.XESContractAddress), conn)
	if err != nil {
		log.Panic("[airdrop] Failed to instantiate a Token contract:", err)
	}

	keystorereader, err := os.Open(config.Config.AirdropWalletfile)
	if err != nil {
		log.Panic("[airdrop] Failed to read keystore:", err)
	}

	keystorekey, err := ioutil.ReadFile(config.Config.AirdropWalletkey)
	if err != nil {
		log.Panic("[airdrop] Failed to read keystore key:", err)
	}

	// Create an authorized transactor and spend Amount XES
	auth, err := bind.NewTransactor(keystorereader, string(keystorekey[:len(keystorekey)-1]))
	if err != nil {
		log.Panic("[airdrop] Failed to create authorized transactor:", err)
	}

	auth.Nonce = nonceManager.NextNonce()
	tx, err := token.Transfer(auth, common.HexToAddress(walletAddress), amount)
	nonceManager.OnError(err)
	if err != nil {
		log.Panic("[airdrop] Failed to request token transfer:", err)
	}
	log.Println("[airdrop] [Ropsten] Sending XES with tx:", tx.Hash().String())
}

func FreeEth(walletAddress string) {
	log.Println("[airdrop] [Ropsten] Prepare ETH for addr:", walletAddress)
	amount := new(big.Int)
	f, err := strconv.ParseFloat(config.Config.AirdropAmountEther, 64)
	amount.SetInt64(int64(f * etherUnit))

	var gasLimit = uint64(21000)
	gasPrice, err := conn.SuggestGasPrice(context.Background())
	nonceManager.OnError(err)
	if err != nil {
		log.Panic(err)
	}

	keystoreJSON, err := ioutil.ReadFile(config.Config.AirdropWalletfile)
	if err != nil {
		log.Panic("[airdrop] Failed to read keystore:", err)
	}

	keystorekey, err := ioutil.ReadFile(config.Config.AirdropWalletkey)
	if err != nil {
		log.Panic("[airdrop] Failed to read keystore key:", err)
	}

	unlockedKey, err := keystore.DecryptKey(keystoreJSON, string(keystorekey[:len(keystorekey)-1]))
	if err != nil {
		log.Panic("[airdrop] Failed to create authorized transactor:", err)
	}

	nonce := nonceManager.NextNonce()
	tx := types.NewTransaction(nonce.Uint64(), common.HexToAddress(walletAddress), amount, gasLimit, gasPrice, nil)
	// chainid 3 = ropsten, see https://github.com/ethereum/EIPs/blob/master/EIPS/eip-155.md#list-of-chain-ids
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(3)), unlockedKey.PrivateKey)
	if err != nil {
		log.Panic("[airdrop] Failed to sign transaction:", err)
	}

	err = conn.SendTransaction(context.Background(), signedTx)
	nonceManager.OnError(err)
	if err != nil {
		log.Panic("[airdrop] Failed to send transaction:", err)
	}

	log.Println("[airdrop] [Ropsten] Sending ETH with tx:", signedTx.Hash().String())
}
