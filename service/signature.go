package service

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	uuid "github.com/satori/go.uuid"

	"github.com/ProxeusApp/proxeus-core/main/handlers/helpers"
	"github.com/ProxeusApp/proxeus-core/main/www"
	"github.com/ProxeusApp/proxeus-core/storage/database/db"
	"github.com/ProxeusApp/proxeus-core/sys/model"
)

type (

	// SignatureService is an interface that provides document signature functions
	SignatureService interface {
		// GetById returns the signature request related to the provided documentId
		GetById(id, docId string) (SignatureRequestsMinimal, error)

		// GetForCurrentUser returns a list of the signature requests for the current user
		GetForCurrentUser(auth model.Auth) (SignatureRequestsComplete, error)

		// AddAndNotify adds a new signature request and if the signatory has provided an email address sends an email notification.
		AddAndNotify(auth model.Auth, translator www.Translator, id, docId, signatory, host, scheme string) error

		// RevokeAndNotify revokes an existing signature request and if the signatory has provided an email address sends an email notification.
		RevokeAndNotify(auth model.Auth, translator www.Translator, id, docId, signatory, host, scheme string) error

		// RejectAndNotify rejects an existing signature request and if the signatory has provided an email address sends an email notification.
		RejectAndNotify(auth model.Auth, translator www.Translator, id, docId, host string) error
	}
	DefaultSignatureService struct {
		fileService  FileService
		userService  UserService
		emailService EmailService
	}

	SignatureRequestItemMinimal struct {
		SignatoryName string  `json:"signatoryName"`
		SignatoryAddr string  `json:"signatoryAddr"`
		RequestedAt   *string `json:"requestedAt,omitempty"`
		Rejected      bool    `json:"rejected"`
		RejectedAt    *string `json:"rejectedAt,omitempty"`
		Revoked       bool    `json:"revoked"`
		RevokedAt     *string `json:"revokedAt,omitempty"`
	}

	SignatureRequestItemComplete struct {
		ID          string  `json:"id"`
		DocID       string  `json:"docID"`
		Hash        string  `json:"hash"`
		From        string  `json:"requestorName"`
		FromAddr    string  `json:"requestorAddr"`
		RequestedAt *string `json:"requestedAt,omitempty"`
		Rejected    bool    `json:"rejected"`
		RejectedAt  *string `json:"rejectedAt,omitempty"`
		Revoked     bool    `json:"revoked"`
		RevokedAt   *string `json:"revokedAt,omitempty"`
	}

	SignatureRequestsMinimal []SignatureRequestItemMinimal

	SignatureRequestsComplete []SignatureRequestItemComplete
)

var ErrSignatureRequestAlreadyExists = errors.New("Request already exists")

func NewSignatureService(fileService FileService, userService UserService, emailService EmailService) SignatureService {
	return &DefaultSignatureService{fileService: fileService, userService: userService, emailService: emailService}
}

// GetById returns the signature request related to the provided documentId
func (me *DefaultSignatureService) GetById(id, documentId string) (SignatureRequestsMinimal, error) {
	signatureRequests, err := signatureRequestDB().GetByID(id, documentId)
	if err != nil {
		log.Println("[signatureService][GetById] signatureRequestDB().GetByID err: ", err.Error())
		return nil, os.ErrNotExist
	}

	var requests = *new(SignatureRequestsMinimal)
	for _, sigreq := range *signatureRequests {
		signatoryName := *new(string)
		item, err := userDB().GetByBCAddress(sigreq.Signatory)
		if err == nil {
			signatoryName = item.Name
		}
		var reqAt string
		reqAt = sigreq.RequestedAt.Format("2.1.2006 15:04")
		var rejAt string
		if sigreq.Rejected {
			rejAt = sigreq.RejectedAt.Format("2.1.2006 15:04")
		}
		var revAt string
		revAt = sigreq.RevokedAt.Format("2.1.2006 15:04")

		reqitem := SignatureRequestItemMinimal{
			signatoryName,
			sigreq.Signatory,
			&reqAt,
			sigreq.Rejected,
			&rejAt,
			sigreq.Revoked,
			&revAt,
		}
		if !sigreq.Rejected {
			reqitem.RejectedAt = nil
		}

		requests = append(requests, reqitem)
	}
	return requests, nil
}

// GetForCurrentUser returns a list of the signature requests for the current user
func (me *DefaultSignatureService) GetForCurrentUser(auth model.Auth) (SignatureRequestsComplete, error) {
	user, err := me.userService.GetUser(auth)
	if err != nil {
		return nil, err
	}
	ethAddr := user.EthereumAddr
	if len(ethAddr) != 42 {
		return nil, os.ErrNotExist
	}
	signatureRequests, err := signatureRequestDB().GetBySignatory(ethAddr)
	if err != nil {
		return nil, os.ErrNotExist
	}

	var requests = *new(SignatureRequestsComplete)
	for _, sigreq := range *signatureRequests {
		var requesterName string
		requester, err := userDB().GetByBCAddress(sigreq.Requestor)
		if err != nil && !db.NotFound(err) {
			return nil, err
		}
		if requester != nil {
			requesterName = requester.Name
		}

		var reqAt string
		reqAt = sigreq.RequestedAt.Format("2.1.2006 15:04")
		var rejAt string
		if sigreq.Rejected {
			rejAt = sigreq.RejectedAt.Format("2.1.2006 15:04")
		}
		var revAt string
		revAt = sigreq.RevokedAt.Format("2.1.2006 15:04")

		reqitem := SignatureRequestItemComplete{
			sigreq.DocId,
			sigreq.DocPath,
			sigreq.Hash,
			requesterName,
			sigreq.Requestor,
			&reqAt,
			sigreq.Rejected,
			&rejAt,
			sigreq.Revoked,
			&revAt,
		}
		if !sigreq.Revoked {
			reqitem.RevokedAt = nil
		}
		if !sigreq.Rejected {
			reqitem.RejectedAt = nil
		}

		requests = append(requests, reqitem)
	}
	return requests, nil
}

// AddAndNotify adds a new signature request and if the signatory has provided an email address sends an email notification.
func (me *DefaultSignatureService) AddAndNotify(auth model.Auth, translator www.Translator, id, docId, signatory, host, scheme string) error {
	err := me.add(auth, id, docId, signatory)
	if err != nil {
		return err
	}

	requestedSigner, err := me.userService.GetByBCAddress(signatory)
	requester, err := me.userService.GetUser(auth)
	if err != nil {
		return err
	}

	if len(requestedSigner.Email) <= 3 {
		return nil
	}

	subject := translator.T("New signature request received")
	body := fmt.Sprintf("<div>Your signature was requested for a document from %s <br />by %s (%s)<br />%s<br /><br />"+
		"The requestor would like you to review and sign the document on the platform.<br /><br />"+
		"To check your pending signature requests, please log in <a href='%s'>here</a></div>",
		host, requester.Name, requester.Email, requester.EthereumAddr, helpers.AbsoluteURLWithScheme(scheme, host, "/user/signature-requests"))

	return me.emailService.Send(requestedSigner.Email, subject, body)
}

// RevokeAndNotify revokes an existing signature request and if the signatory has provided an email address sends an email notification.
func (me *DefaultSignatureService) RevokeAndNotify(auth model.Auth, translator www.Translator, id, docId, signatory, host, scheme string) error {

	sig, err := userDB().GetByBCAddress(signatory)
	if err != nil {
		return err
	}

	err = signatureRequestDB().SetRevoked(id, docId, signatory)
	if err != nil {
		return os.ErrNotExist
	}

	requestor, err := me.userService.GetUser(auth)
	if err != nil {
		return err
	}

	if len(sig.Email) <= 3 {
		return nil
	}

	// Earlier you may have received a signature request from <base URL>by <Name> (<Email>)<Ethereum-Addr>
	// The requestor has retracted the request. You may still log in and view the request, but can no longer sign the document.
	// To check your signature requests, please log in <here (link to requests, if logged in>

	subject := translator.T("signature request revoked")
	body := fmt.Sprintf("<div>Earlier you may have received a signature request from %s by %s (%s)<br />%s<br /><br />"+
		"The requestor has retracted the request. You may still log in and view the request, but can no longer sign the document."+
		"<br /><br />To check your signature requests, please log in <a href='%s'>here</a></div>",
		host, requestor.Name, requestor.Email, requestor.EthereumAddr, helpers.AbsoluteURLWithScheme(scheme, host, "/user/signature-requests"))
	err = me.emailService.Send(sig.Email, subject, body)
	if err != nil {
		log.Println("UserDocumentSignatureRequestRevokeHandler emailService.Send err: ", err.Error())
	}
	return err
}

// RejectAndNotify rejects an existing signature request and if the signatory has provided an email address sends an email notification.
func (me *DefaultSignatureService) RejectAndNotify(auth model.Auth, translator www.Translator, id, docId, host string) error {
	item, err := me.userService.GetUser(auth)
	if err != nil {
		return err
	}
	signatoryAddr := item.EthereumAddr
	signatureRequests, err := signatureRequestDB().GetByID(id, docId)
	if err != nil {
		return os.ErrNotExist
	}
	signatureRequest := (*signatureRequests)[0]

	err = signatureRequestDB().SetRejected(id, docId, signatoryAddr)
	if err != nil {
		return os.ErrNotExist
	}

	requestorAddr, err := userDB().GetByBCAddress(signatureRequest.Requestor)
	if err != nil {
		return err
	}

	if len(requestorAddr.Email) <= 3 {
		return nil
	}

	// Your signature request for a document on <platform base URL> from <timestamp> has been rejected by <Name> (<Email>)<Ethereum-Addr>
	// You may send another request if you think this was by mistake.

	subject := translator.T("Signature request rejected")
	body := fmt.Sprintf("<div>Your signature request for a document on %s from %s <br />has been rejected by %s (%s)<br />%s<br />"+
		"<br />You may send another request if you think this was by mistake.</div>",
		host, signatureRequest.RequestedAt.Format("2.1.2006 15:04"), item.Name, item.Email, item.EthereumAddr)
	err = me.emailService.Send(requestorAddr.Email, subject, body)
	if err != nil {
		log.Println("UserDocumentSignatureRequestRejectHandler emailService.Send err: ", err.Error())
	}
	return err
}

func (me *DefaultSignatureService) add(auth model.Auth, id, docId, signatory string) error {
	fileInfo, err := me.fileService.GetDataFile(auth, id, docId)
	if err != nil {
		log.Println("[signatureService][Add] fileService.GetDataFile err: ", err.Error())
		return os.ErrNotExist
	}

	if !strings.HasPrefix(docId, "docs") {
		return os.ErrNotExist
	}

	var documentBytes bytes.Buffer
	err = filesDB().Read(fileInfo.Path(), &documentBytes)
	if err != nil {
		log.Println("[signatureService][Add] filesDB().Read err: ", err.Error())
		return os.ErrNotExist
	}
	docHash := crypto.Keccak256Hash(documentBytes.Bytes()).String()

	signatoryObj, err := userDB().GetByBCAddress(signatory)
	if err != nil {
		return err
	}
	fileObj, err := userDataDB().Get(auth, id)
	if err != nil {
		log.Println("[signatureService][Add] userDataDB().Get err: ", err.Error())
		return os.ErrNotExist
	}
	if fileObj.Permissions.Grant == nil || !fileObj.Permissions.Grant[signatoryObj.UserID()].IsRead() {
		if fileObj.Permissions.Grant == nil {
			fileObj.Permissions.Grant = make(map[string]model.Permission)
		}
		fileObj.Permissions.Grant[signatoryObj.UserID()] = model.Permission{byte(1)}
		fileObj.Permissions.Change(auth, &fileObj.Permissions)

		err = userDataDB().Put(auth, fileObj)
		if err != nil {
			return err
		}
	}

	fileObj, err = userDataDB().Get(auth, id)
	if err != nil {
		return err
	}

	requestor, err := me.userService.GetUser(auth)
	if err != nil {
		return err
	}

	requestItem := model.SignatureRequestItem{
		ID:          uuid.NewV4().String(),
		DocId:       id,
		DocPath:     docId,
		Hash:        docHash,
		Requestor:   requestor.EthereumAddr,
		RequestedAt: time.Now(),
		Signatory:   signatory,
		Rejected:    false,
	}

	signatureRequests, err := signatureRequestDB().GetByID(id, docId)
	if err == nil {
		for _, sigreq := range *signatureRequests {
			if sigreq.Signatory == signatory &&
				sigreq.Hash == docHash &&
				sigreq.Rejected == false &&
				sigreq.Revoked == false {
				return ErrSignatureRequestAlreadyExists
			}
		}
	}

	err = signatureRequestDB().Add(&requestItem)
	if err != nil {
		log.Println("[signatureService][Add] signatureRequestDB().Add err: ", err.Error())
		return os.ErrNotExist
	}

	return nil
}
