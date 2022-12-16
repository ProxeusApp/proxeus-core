package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/ProxeusApp/proxeus-core/main/config"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ProxeusApp/proxeus-core/main/handlers/blockchain/ethglue"
)

const etherUnit = 1000000000000000000

var conn ethglue.ETHClientIF
var nonceManager ethglue.NonceManager

var mu sync.Mutex

func GiveTokens(toWallet string) {
	var err error

	dialler := ethglue.NewDefaultDialler()
	conn, err = dialler.Dial(config.Config.EthClientURL)
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

// FreeXES sends testnet XES to given wallet address
func FreeXES(walletAddress string) {
	log.Println("[airdrop] [testnet] Prepare XES for addr:", walletAddress)

	amountFloat := new(big.Float)
	_, parsedF := amountFloat.SetString(config.Config.AirdropAmountXES)
	if !parsedF {
		log.Panic("[airdrop] Couldn't parse AirdropAmountXES to Float:", config.Config.AirdropAmountXES)
	}

	amount := new(big.Float)
	amount.Mul(amountFloat, big.NewFloat(etherUnit))

	amountInt := new(big.Int)
	amountDec := fmt.Sprintf("%f", amount)
	amountInt.SetString(amountDec, 10)

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
	tx, err := token.Transfer(auth, common.HexToAddress(walletAddress), amountInt)
	nonceManager.OnError(err)
	if err != nil {
		log.Panic("[airdrop] Failed to request token transfer:", err)
	}
	log.Println("[airdrop] [testnet] Sending XES with tx:", tx.Hash().String())
}

func FreeEth(walletAddress string) {
	log.Println("[airdrop] [testnet] Prepare ETH for addr:", walletAddress)

	amountFloat := new(big.Float)
	_, parsedF := amountFloat.SetString(config.Config.AirdropAmountEther)
	if !parsedF {
		log.Panic("[airdrop] Couldn't parse AirdropAmountEther to Float:", config.Config.AirdropAmountEther)
	}

	amount := new(big.Float)
	amount.Mul(amountFloat, big.NewFloat(etherUnit))

	amountInt := new(big.Int)

	amountDec := fmt.Sprintf("%f", amount)
	amountInt.SetString(amountDec, 10)

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
	tx := types.NewTransaction(nonce.Uint64(), common.HexToAddress(walletAddress), amountInt, gasLimit, gasPrice, nil)
	// chainid 5 = goerli, see https://github.com/ethereum/EIPs/blob/master/EIPS/eip-155.md#list-of-chain-ids
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(3)), unlockedKey.PrivateKey)
	if err != nil {
		log.Panic("[airdrop] Failed to sign transaction:", err)
	}

	err = conn.SendTransaction(context.Background(), signedTx)
	nonceManager.OnError(err)
	if err != nil {
		log.Panic("[airdrop] Failed to send transaction:", err)
	}

	log.Println("[airdrop] [testnet] Sending ETH with tx:", signedTx.Hash().String())
}
