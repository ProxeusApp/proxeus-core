package service

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ProxeusApp/proxeus-core/main/handlers/helpers"
	"github.com/ProxeusApp/proxeus-core/main/www"
	"github.com/ProxeusApp/proxeus-core/sys/model"
	"github.com/ethereum/go-ethereum/crypto"
	uuid "github.com/satori/go.uuid"
	"log"
	"os"
	"strings"
	"time"
)

type (
	SignatureService interface {
		GetById(id, docId string) (SignatureRequests, error)
		AddAndNotify(auth model.Auth, translator www.Translator, id, docId, signatory, host, scheme string) error
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
	SignatureRequests []SignatureRequestItemMinimal
)

var ErrSignatureRequestAlreadyExists = errors.New("Request already exists")

func NewSignatureService(fileService FileService, userService UserService, emailService EmailService) SignatureService {
	return &DefaultSignatureService{fileService: fileService, userService: userService, emailService: emailService}
}

func (me *DefaultSignatureService) GetById(id, docId string) (SignatureRequests, error) {
	signatureRequests, err := signatureRequestDB().GetByID(id, docId)
	if err != nil {
		log.Println("[signatureService][GetById] signatureRequestDB().GetByID err: ", err.Error())
		return nil, os.ErrNotExist
	}

	var requests = *new(SignatureRequests)
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
