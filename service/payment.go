package service

import (
	"errors"
	"github.com/ProxeusApp/proxeus-core/storage/database/db"
	"github.com/ProxeusApp/proxeus-core/sys/model"
	uuid "github.com/satori/go.uuid"
	"log"
	"strings"
	"time"
)

type (
	PaymentService interface {
		CreateWorkflowPayment(auth model.Auth, workflowId, ethAddress string) (*model.WorkflowPaymentItem, error)
		GetWorkflowPaymentById(paymentId string) (*model.WorkflowPaymentItem, error)
		GetWorkflowPayment(txHash, ethAddresses, status string) (*model.WorkflowPaymentItem, error)
		UpdateWorkflowPaymentPending(paymentId, txHash, ethAddress string) error
		CancelWorkflowPayment(paymentId, ethAddress string) error
		RedeemPayment(workflowId, ethAddr string) error
		CheckIfWorkflowPaymentRequired(auth model.Auth, workflowId string) (bool, error)
		CheckForWorkflowPayment(auth model.Auth, workflowId string) error
		Delete(paymentId string) error
		All() ([]*model.WorkflowPaymentItem, error)
	}

	DefaultPaymentService struct {
		//Important: Pass system to service (and not e.g. system.DB.WorkflowPayments because system.DB variable is replaced on calling api/handlers.PostInit()
		userService UserService
	}
)

var (
	errRequiredParamMissing = errors.New("required parameter missing")
	ErrTxHashEmpty          = errors.New("no txHash given")
)

//Important: Pass system to service (and not e.g. system.DB.WorkflowPayments because system.DB variable is replaced on calling api/handlers.PostInit()
func NewPaymentService(userService UserService) *DefaultPaymentService {
	return &DefaultPaymentService{userService: userService}
}

// CreateWorkflowPayment creates a workflow payment for a workflow by the ethAddress with the status "created".
func (me *DefaultPaymentService) CreateWorkflowPayment(auth model.Auth, workflowId, ethAddress string) (*model.WorkflowPaymentItem, error) {
	workflow, err := workflowDB().Get(auth, workflowId)
	if err != nil {
		return nil, err
	}

	payment := &model.WorkflowPaymentItem{
		ID:         uuid.NewV4().String(),
		Xes:        workflow.Price,
		From:       ethAddress,
		To:         workflow.OwnerEthAddress,
		Status:     model.PaymentStatusCreated,
		CreatedAt:  time.Now(),
		WorkflowID: workflowId,
	}

	return payment, paymentsDB().Save(payment)
}

// GetWorkflowPaymentById returns a WorkflowPaymentItem for the id
func (me *DefaultPaymentService) GetWorkflowPaymentById(paymentId string) (*model.WorkflowPaymentItem, error) {
	return paymentsDB().Get(paymentId)
}

//GetWorkflowPayment returns a WorkflowPaymentItem that matches the txHash, ethAddress and status
func (me *DefaultPaymentService) GetWorkflowPayment(txHash, ethAddresses, status string) (*model.WorkflowPaymentItem, error) {
	if txHash == "" {
		log.Printf("[GetWorkflowPayment] bad request, either provide paymentId, txHash or workflowId")
		return nil, errRequiredParamMissing
	}

	payment, err := paymentsDB().GetByTxHashAndStatusAndFromEthAddress(txHash, status, ethAddresses)
	if err != nil {
		log.Println("[GetWorkflowPayment] GetByTxHashAndStatusAndFromethAddressess err: ", err.Error())
		return nil, err
	}

	log.Printf("[workflowHandler][GetWorkflowPayment] ID: %s, txHash: %s", payment.ID, payment.TxHash)

	return payment, nil
}

// UpdateWorkflowPaymentPending updates a WorkflowPaymentItem with status "created" and sets it to status "pending"
func (me *DefaultPaymentService) UpdateWorkflowPaymentPending(paymentId, txHash, ethAddress string) error {
	txHash = strings.TrimSpace(txHash)
	if txHash == "" {
		return ErrTxHashEmpty
	}

	err := paymentsDB().Update(paymentId, model.PaymentStatusPending, txHash, ethAddress)
	if err != nil {
		log.Printf("[UpdateWorkflowPayment] WorkflowPayments.Update err: %s", err.Error())
	}
	return err
}

// CancelWorkflowPayment sets the status of a WorkflowPaymentItem to "cancelled"
func (me *DefaultPaymentService) CancelWorkflowPayment(paymentId, ethAddress string) error {
	return paymentsDB().Cancel(paymentId, ethAddress)
}

// RedeemPayment sets the payment status of WorkflowPaymentItem from "confirmed" to "redeemed"
func (me *DefaultPaymentService) RedeemPayment(workflowId, ethAddr string) error {
	return paymentsDB().Redeem(workflowId, ethAddr)
}

//CheckIfWorkflowPaymentRequired returns whether a payment is required for the user for a workflow
func (me *DefaultPaymentService) CheckIfWorkflowPaymentRequired(auth model.Auth, workflowId string) (bool, error) {
	workflow, err := workflowDB().Get(auth, workflowId)
	if err != nil {
		return true, err
	}

	_, alreadyStarted, err := userDataDB().GetByWorkflow(auth, workflow, false)
	if err != nil {
		if !db.NotFound(err) {
			return true, nil
		}
		//if workflow not found (strm.ErrNotFound ) still check with isPaymentRequired
	}

	return isPaymentRequired(alreadyStarted, workflow, auth.UserID()), nil
}

// CheckForWorkflowPayment checks whether a workflow payment is required.
// If a payment is required checks whether a payment with status "confirmed" is found.
func (me *DefaultPaymentService) CheckForWorkflowPayment(auth model.Auth, workflowId string) error {
	user, err := me.userService.GetUser(auth)
	if err != nil {
		return err
	}

	paymentRequired, err := me.CheckIfWorkflowPaymentRequired(auth, workflowId)
	if err != nil {
		return err
	}

	if paymentRequired {
		_, err = paymentsDB().GetByWorkflowIdAndFromEthAddress(workflowId, user.EthereumAddr, []string{model.PaymentStatusConfirmed})
	}

	return err
}

// Delete sets the status of a CheckForWorkflowPayment to "deleted"
func (me *DefaultPaymentService) Delete(paymentId string) error {
	return paymentsDB().Delete(paymentId)
}

// All returns a list of all WorkflowPaymentItem
func (me *DefaultPaymentService) All() ([]*model.WorkflowPaymentItem, error) {
	return paymentsDB().All()
}

func isPaymentRequired(alreadyStarted bool, workflow *model.WorkflowItem, userId string) bool {
	return !alreadyStarted && workflow.Owner != userId && workflow.Price != 0
}
