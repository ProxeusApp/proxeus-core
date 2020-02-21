package blockchain

import (
	"context"
	"encoding/hex"
	"log"

	"github.com/ProxeusApp/proxeus-core/storage"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ProxeusApp/proxeus-core/sys/email"
)

type Signaturelistener struct {
	listener
	signatureRequestsDB storage.SignatureRequestsIF
	userDB              storage.UserIF
	emailSender         email.EmailSender
	domain              string
	ProxeusFSABI        abi.ABI
	emailFrom           string
	logSubscriber       LogSubscriber
}

var TestChannelSignature chan types.Log

func NewSignatureListener(SignatureRequestsDB storage.SignatureRequestsIF,
	UserDB storage.UserIF, EmailSender email.EmailSender, ProxeusFSABI abi.ABI, domain string, logSubscriber LogSubscriber) *Signaturelistener {

	me := &Signaturelistener{}
	me.ProxeusFSABI = ProxeusFSABI
	me.signatureRequestsDB = SignatureRequestsDB
	me.emailSender = EmailSender
	me.userDB = UserDB
	me.domain = domain
	me.logs = make(chan types.Log, 200)
	me.logSubscriber = logSubscriber

	TestChannelSignature = me.logs

	return me
}

func (me *Signaturelistener) Listen(ctx context.Context) {

	subscription := make(chan ethereum.Subscription)

	for {
		go me.logSubscriber.Subscribe(ctx, me.logs, subscription)
		select {
		case sub := <-subscription:
			me.sub = sub
			log.Println("[signaturelistener] listen on contract started on: ", me.logSubscriber)
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
	pfsLogUnpacker := bind.NewBoundContract(common.HexToAddress(me.logSubscriber.GetContractAddress()), me.ProxeusFSABI,
		nil, nil, nil)

	err := pfsLogUnpacker.UnpackLog(out, eventType, *lg)
	if err != nil {
		return err // not our event type
	}
	return nil
}
