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

func (me *DefaultPaymentService) GetWorkflowPaymentById(paymentId string) (*model.WorkflowPaymentItem, error) {
	return paymentsDB().Get(paymentId)
}

func (me *DefaultPaymentService) GetWorkflowPayment(txHash, ethAddressess, status string) (*model.WorkflowPaymentItem, error) {
	if txHash == "" {
		log.Printf("[GetWorkflowPayment] bad request, either provide paymentId, txHash or workflowId")
		return nil, errRequiredParamMissing
	}

	payment, err := paymentsDB().GetByTxHashAndStatusAndFromEthAddress(txHash, status, ethAddressess)
	if err != nil {
		log.Println("[GetWorkflowPayment] GetByTxHashAndStatusAndFromethAddressess err: ", err.Error())
		return nil, err
	}

	log.Printf("[workflowHandler][GetWorkflowPayment] ID: %s, txHash: %s", payment.ID, payment.TxHash)

	return payment, nil
}

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

func (me *DefaultPaymentService) CancelWorkflowPayment(paymentId, ethAddress string) error {
	return paymentsDB().Cancel(paymentId, ethAddress)
}

// Set the payment status from confirmed to redeemed
func (me *DefaultPaymentService) RedeemPayment(workflowId, ethAddr string) error {
	return paymentsDB().Redeem(workflowId, ethAddr)
}

//returns a bool indicating whether a payment is required for the user for a workflow
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

func (me *DefaultPaymentService) Delete(paymentId string) error {
	return paymentsDB().Delete(paymentId)
}

func (me *DefaultPaymentService) All() ([]*model.WorkflowPaymentItem, error) {
	return paymentsDB().All()
}

func isPaymentRequired(alreadyStarted bool, workflow *model.WorkflowItem, userId string) bool {
	return !alreadyStarted && workflow.Owner != userId && workflow.Price != 0
}
