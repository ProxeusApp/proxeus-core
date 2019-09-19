package payment

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"

	"git.proxeus.com/core/central/sys/db/storm"

	strm "github.com/asdine/storm"

	"github.com/labstack/echo"

	"git.proxeus.com/core/central/main/www"
	"git.proxeus.com/core/central/sys/model"
)

var errNotAuthorized = errors.New("user not authorized")

type createPaymentRequest struct {
	WorkflowId string `json:"workflowId"`
}

//create a payment for a workflow
func CreateWorkflowPayment(e echo.Context) error {
	c := e.(*www.Context)

	user, err := getUser(c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	createPaymentRequest := &createPaymentRequest{}
	err = c.Bind(&createPaymentRequest)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	createPaymentRequest.WorkflowId = strings.TrimSpace(createPaymentRequest.WorkflowId)

	if createPaymentRequest.WorkflowId == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	workflow, err := c.System().DB.Workflow.Get(c.Session(false), createPaymentRequest.WorkflowId)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	id := uuid.NewV4().String()

	payment := &model.WorkflowPaymentItem{
		ID:         id,
		Xes:        workflow.Price,
		From:       user.EthereumAddr,
		To:         workflow.OwnerEthAddress,
		Status:     model.PaymentStatusCreated,
		CreatedAt:  time.Now(),
		WorkflowID: createPaymentRequest.WorkflowId,
	}

	err = c.System().DB.WorkflowPaymentsDB.Save(payment)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	return c.JSON(http.StatusOK, payment)
}

// Gets a payment by Id
// Payment is only returned if the from address == the user sending the request
func GetWorkflowPaymentById(e echo.Context) error {
	c := e.(*www.Context)
	paymentId := c.Param("paymentId")

	user, err := getUser(c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	payment, err := c.System().DB.WorkflowPaymentsDB.Get(paymentId)
	if err != nil {
		log.Println("[GetWorkflowPaymentById] getUserPaymentById err: ", err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	if payment.From != user.EthereumAddr {
		return c.NoContent(http.StatusBadRequest)
	}

	return c.JSON(http.StatusOK, payment)
}

// Returns payment with the given txHash and status.
// Payment is only returned if the from address == the user sending the request
func GetWorkflowPayment(e echo.Context) error {
	c := e.(*www.Context)
	txHash := c.QueryParam("txHash")
	status := c.QueryParam("status")

	user, err := getUser(c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	payment, err := getWorkflowPayment(c.System().DB.WorkflowPaymentsDB, txHash, user.EthereumAddr, status)
	if err != nil {
		if err == strm.ErrNotFound {
			return c.NoContent(http.StatusNotFound)
		}
		log.Println("[GetWorkflowPayment] GetByTxHashAndStatusAndFromEthAddress err: ", err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	return c.JSON(http.StatusOK, payment)
}

var errRequiredParamMissing = errors.New("required parameter missing")

func getWorkflowPayment(workflowPaymentsDB storm.WorkflowPaymentsDBInterface, txHash,
	ethAddr, status string) (*model.WorkflowPaymentItem, error) {

	if txHash == "" {
		log.Printf("[GetWorkflowPayment] bad request, either provide paymentId, txHash or workflowId")
		return nil, errRequiredParamMissing
	}

	payment, err := workflowPaymentsDB.GetByTxHashAndStatusAndFromEthAddress(txHash, status, ethAddr)
	if err != nil {
		log.Println("[GetWorkflowPayment] GetByTxHashAndStatusAndFromEthAddress err: ", err.Error())
		return nil, err
	}

	log.Printf("[workflowHandler][GetWorkflowPayment] ID: %s, txHash: %s", payment.ID, payment.TxHash)

	return payment, nil
}

type updatePaymentPendingRequest struct {
	TxHash string `json:"txHash"`
}

var errTxHashEmpty = errors.New("no txHash given")

// Set a workflow payment from status created to status pending
func UpdateWorkflowPaymentPending(e echo.Context) error {
	c := e.(*www.Context)
	paymentId := strings.TrimSpace(c.Param("paymentId"))

	user, err := getUser(c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	updatePaymentRequest := &updatePaymentPendingRequest{}
	err = c.Bind(&updatePaymentRequest)
	if err != nil {
		log.Printf("[UpdateWorkflowPayment] UpdateWorkflowPayment bind err: %s", err.Error())
		return err
	}

	err = updateWorkflowPaymentPending(c.System().DB.WorkflowPaymentsDB, paymentId, updatePaymentRequest.TxHash, user.EthereumAddr)
	if err != nil {
		log.Printf("[UpdateWorkflowPayment] err: %s", err.Error())
		if err == errTxHashEmpty {
			return c.String(http.StatusBadRequest, err.Error())
		}
		return c.NoContent(http.StatusBadRequest)
	}

	return c.NoContent(http.StatusOK)
}

func updateWorkflowPaymentPending(workflowPaymentsDB storm.WorkflowPaymentsDBInterface, paymentId, txHash, ethAddr string) error {
	txHash = strings.TrimSpace(txHash)
	if txHash == "" {
		return errTxHashEmpty
	}

	err := workflowPaymentsDB.Update(paymentId, model.PaymentStatusPending, txHash, ethAddr)
	if err != nil {
		log.Printf("[UpdateWorkflowPayment] WorkflowPaymentsDB.Update err: %s", err.Error())
		return err
	}

	return nil
}

// Set status of workflow from created to cancelled
func CancelWorkflowPayment(e echo.Context) error {
	c := e.(*www.Context)
	paymentId := strings.TrimSpace(c.Param("paymentId"))

	user, err := getUser(c)
	if err != nil {
		return errNotAuthorized
	}

	err = cancelWorkflowPayment(c.System().DB.WorkflowPaymentsDB, paymentId, user.EthereumAddr)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func cancelWorkflowPayment(workflowPaymentsDB storm.WorkflowPaymentsDBInterface, paymentId, ethAddr string) error {
	return workflowPaymentsDB.Cancel(paymentId, ethAddr)
}

// Set the payment status from confirmed to redeemed
func RedeemPayment(workflowPaymentsDB storm.WorkflowPaymentsDBInterface, workflowId, ethAddr string) error {
	return workflowPaymentsDB.Redeem(workflowId, ethAddr)
}

//returns a bool indicating whether a payment is required for the user for a workflow
func CheckIfWorkflowPaymentRequired(c *www.Context, workflowId string) (bool, error) {
	sess := c.Session(false)

	workflow, err := c.System().DB.Workflow.Get(sess, workflowId)
	if err != nil {
		return true, err
	}

	_, alreadyStarted, err := c.System().DB.UserData.GetByWorkflow(sess, workflow, false)
	if err != nil {
		if err != strm.ErrNotFound {
			return true, nil
		}
		//if workflow not found (strm.ErrNotFound ) still check with isPaymentRequired
	}

	return isPaymentRequired(alreadyStarted, workflow, c.Session(false).UserID()), nil
}

func isPaymentRequired(alreadyStarted bool, workflow *model.WorkflowItem, userId string) bool {
	return !alreadyStarted && workflow.Owner != userId && workflow.Price != 0
}

// Set Payment for a workflow to status = Deleted. Only for superadmin
func DeleteWorkflowPayment(e echo.Context) error {
	c := e.(*www.Context)
	paymentId := strings.TrimSpace(c.Param("paymentId"))

	if paymentId == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	err := c.System().DB.WorkflowPaymentsDB.Delete(paymentId)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// List all Payment. Only for superadmin and debugging purposes
func ListPayments(e echo.Context) error {
	c := e.(*www.Context)

	payments, err := c.System().DB.WorkflowPaymentsDB.All()
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	return c.JSON(http.StatusOK, payments)
}

func getUser(c *www.Context) (*model.User, error) {
	sess := c.Session(false)
	return c.System().DB.User.Get(sess, sess.UserID())
}
