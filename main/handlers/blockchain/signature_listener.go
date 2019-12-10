package blockchain

import (
	"context"
	"encoding/hex"
	"log"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ProxeusApp/proxeus-core/main/ethglue"
	"github.com/ProxeusApp/proxeus-core/sys/db/storm"
	"github.com/ProxeusApp/proxeus-core/sys/email"
)

type (
	Signaturelistener struct {
		listener
		signatureRequestsDB       storm.SignatureRequestsDB
		userDB                    storm.UserDBInterface
		emailSender               email.EmailSender
		domain                    string
		BlockchainContractAddress string
		ProxeusFSABI              abi.ABI
		emailFrom                 string
	}
)

func NewSignatureListener(ethWebSocketURL, ethURL, BlockchainContractAddress string, SignatureRequestsDB *storm.SignatureRequestsDB,
	UserDB storm.UserDBInterface, EmailSender email.EmailSender, ProxeusFSABI abi.ABI, domain string) *Signaturelistener {

	me := &Signaturelistener{}
	me.BlockchainContractAddress = BlockchainContractAddress
	me.ethWebSocketURL = ethWebSocketURL
	me.ethURL = ethURL
	me.ProxeusFSABI = ProxeusFSABI
	me.signatureRequestsDB = *SignatureRequestsDB
	me.emailSender = EmailSender
	me.userDB = UserDB
	me.domain = domain
	me.logs = make(chan types.Log, 200)
	return me
}

func (me *Signaturelistener) Listen(ctx context.Context) {

	var readyCh <-chan struct{}

	for {
		readyCh = me.ethConnectWebSocketsAsync(ctx)
		select {
		case <-readyCh:
			log.Println("[signaturelistener] listen on contract started. contract address: ", me.BlockchainContractAddress)
			reconnect := me.listenLoop(ctx)
			if !reconnect {
				log.Printf("[signaturelistener][eventHandler] finished")
				return
			}
		case <-ctx.Done():
			log.Printf("[signaturelistener][eventHandler] done")
			return
		}
	}
	return
}

func (me *Signaturelistener) listenLoop(ctx context.Context) (shouldReconnect bool) {
	for {
		select {
		case <-ctx.Done():
			return false
		case err, ok := <-me.sub.Err():
			if !ok {
				return true
			}
			log.Println("ERROR sub", err)
			return true
		case vLog, ok := <-me.logs:
			if !ok {
				return true
			}
			me.eventsHandler(&vLog)
		}
	}
}

func (me *Signaturelistener) ethConnectWebSocketsAsync(ctx context.Context) <-chan struct{} {

	filterAddresses := []common.Address{common.HexToAddress(me.BlockchainContractAddress)}

	readyCh := make(chan struct{})
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				var err error
				ethwsconn, err := ethglue.DialContext(ctx, me.ethWebSocketURL)
				if err != nil {
					log.Printf("failed to dial for eth events, will retry (%s)\n", err)
					continue
				}
				query := ethereum.FilterQuery{
					Addresses: filterAddresses,
				}
				ctx, cancel := context.WithTimeout(ctx, time.Duration(10*time.Second))
				me.sub, err = ethwsconn.SubscribeFilterLogs(ctx, query, me.logs)
				cancel()
				if err != nil {
					log.Printf("failed to subscribe for eth events, will retry (%s)\n", err)
					time.Sleep(time.Second * 4)
					continue
				}
				// success!
				readyCh <- struct{}{}
				return
			}
		}
	}()
	return readyCh
}

func (me *Signaturelistener) eventsHandler(lg *types.Log) {
	log.Printf("[signaturelistener][eventHandler] txHash: %s, %v", lg.TxHash.String(), lg)
	event := new(ProxeusFSFileSignedEvent)
	if err := me.eventFromLog(event, lg, "FileSignedEvent"); err == nil {
		//Search for signingrequest in db
		filehash := "0x" + hex.EncodeToString(event.Hash[:])
		signer := event.Signer.String()
		items, err := me.signatureRequestsDB.GetByHashAndSigner(filehash, signer)
		if err == nil && len(*items) > 0 {
			item := (*items)[0]
			requestorAddr, err := me.userDB.GetByBCAddress(item.Requestor)
			if err != nil {
				log.Println("Coudln't retrieve requestor for event: ", err)
				return
			}
			signatory, err := me.userDB.GetByBCAddress(signer)
			if err != nil {
				log.Println("Coudln't retrieve signer for event: ", err)
				return
			}
			if err == nil {
				if requestorAddr != nil && signatory != nil && len(requestorAddr.Email) > 3 {
					me.emailSender.Send(&email.Email{To: []string{requestorAddr.Email}, Subject: "Signature granted", Body: "<div>Your signature request for a document on dev.proxeus.com from " + item.RequestedAt.Format("2.1.2006 15:04") + "<br />has been accepted by " + signatory.Name + " (" + signatory.Email + ")<br />" + signatory.EthereumAddr + "<br /><br />The document is now cryptographically signed by this user on " + me.domain + ".</div>"})
				}
			}
		}
	}
}

func (me *Signaturelistener) eventFromLog(out interface{}, lg *types.Log, eventType string) error {
	pfsLogUnpacker := bind.NewBoundContract(common.HexToAddress(me.BlockchainContractAddress), me.ProxeusFSABI,
		nil, nil, nil)

	err := pfsLogUnpacker.UnpackLog(out, eventType, *lg)
	if err != nil {
		return err // not our event type
	}
	return nil
}
