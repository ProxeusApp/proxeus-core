package airdrop

import (
	"context"
	"io/ioutil"
	"log"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"git.proxeus.com/core/central/lib/wallet"
	"git.proxeus.com/core/central/main/ethglue"
)

const fromAddress = "38ba9213c70bf6fe34f70cfc5c9b26707c6c1e85"
const serverKeystore = `{
	"address": "` + fromAddress + `",
	"crypto": {
		"cipher": "aes-128-ctr",
		"ciphertext": "",
		"cipherparams": {
			"iv": ""
		},
		"kdf": "scrypt",
		"kdfparams": {
			"dklen": 32,
			"n": 262144,
			"p": 1,
			"r": 8,
			"salt": "2a58cb089bffa3530eb8dbdfaf6bdcc0196b49d3cb718cc4297a08635088523c"
		},
		"mac": "417d89c4f4daa846e058bc24566df9ef67ea5ed0180d6e19af749498911de403"
	},
	"id": "732325b3-1a49-41c3-ae41-10eea619b7a8",
	"version": 3
}`
const contractAddress = "0x84e0b37e8f5b4b86d5d299b0b0e33686405a3919"

var etherUnit = big.NewInt(1000000000000000000)
var freeXes = big.NewInt(10)
var freeEth = big.NewInt(1)
var password = ""

var conn *ethclient.Client
var nonceManager ethglue.NonceManager

// FreeXES sends ropsten XES to given wallet address
func FreeXES(walletAddress string) {
	log.Println("[Ropsten] Prepare XES for addr:", walletAddress)
	amount := new(big.Int)
	amount.Mul(etherUnit, freeXes)

	// make transfer
	token, err := wallet.NewToken(common.HexToAddress(contractAddress), conn)
	if err != nil {
		log.Panic("Failed to instantiate a Token contract:", err)
	}
	// Create an authorized transactor and spend Amount XES
	auth, err := bind.NewTransactor(strings.NewReader(serverKeystore), password)
	if err != nil {
		log.Panic("Failed to create authorized transactor:", err)
	}
	auth.Nonce = nonceManager.NextNonce()
	tx, err := token.Transfer(auth, common.HexToAddress(walletAddress), amount)
	nonceManager.OnError(err)
	if err != nil {
		log.Panic("Failed to request token transfer:", err)
	}
	log.Println("[Ropsten] Sending XES with tx:", tx.Hash().String())
}

func FreeEth(walletAddress string) {
	log.Println("[Ropsten] Prepare ETH for addr:", walletAddress)
	// it wouldn't let me set the amount to 0.02, so just divide result by 50
	amount := new(big.Int)
	amount.Div(amount.Mul(etherUnit, freeEth), big.NewInt(50))

	// TODO get this from config
	var gasLimit = uint64(21000)

	keystoreJSON, err := ioutil.ReadAll(strings.NewReader(serverKeystore))
	if err != nil {
		log.Panic("Failed to read keystore:", err)
	}

	unlockedKey, err := keystore.DecryptKey(keystoreJSON, password)
	if err != nil {
		log.Panic("Failed to create authorized transactor:", err)
	}

	gasPrice, err := conn.SuggestGasPrice(context.Background())
	nonceManager.OnError(err)
	if err != nil {
		log.Panic(err)
	}

	nonce := nonceManager.NextNonce()
	tx := types.NewTransaction(nonce.Uint64(), common.HexToAddress(walletAddress), amount, gasLimit, gasPrice, nil)
	// chainid 3 = ropsten, see https://github.com/ethereum/EIPs/blob/master/EIPS/eip-155.md#list-of-chain-ids
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(3)), unlockedKey.PrivateKey)
	if err != nil {
		log.Panic("Failed to sign transaction:", err)
	}

	err = conn.SendTransaction(context.Background(), signedTx)
	nonceManager.OnError(err)
	if err != nil {
		log.Panic("Failed to send transaction:", err)
	}

	log.Println("[Ropsten] Sending ETH with tx:", signedTx.Hash().String())
}

var mu sync.Mutex

func GiveTokens(toWallet string) {
	var err error
	conn, err = ethglue.Dial("https://ropsten.infura.io/")
	if err != nil {
		log.Panic("Failed to connect to the Ethereum client:", err)
	}
	nonceManager.OnDial(conn)
	nonceManager.OnAccountChange(fromAddress)
	mu.Lock()
	defer mu.Unlock()
	FreeXES(toWallet)
	FreeEth(toWallet)
}
