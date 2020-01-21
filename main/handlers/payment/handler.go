package payment

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ProxeusApp/proxeus-core/storage/database/db"

	cfg "github.com/ProxeusApp/proxeus-core/main/config"
	"github.com/ProxeusApp/proxeus-core/main/handlers/blockchain"
	"github.com/ProxeusApp/proxeus-core/main/www"
	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/sys/model"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
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

	payment := &model.WorkflowPaymentItem{
		ID:         uuid.NewV4().String(),
		Xes:        workflow.Price,
		From:       user.EthereumAddr,
		To:         workflow.OwnerEthAddress,
		Status:     model.PaymentStatusCreated,
		CreatedAt:  time.Now(),
		WorkflowID: createPaymentRequest.WorkflowId,
	}

	err = c.System().DB.WorkflowPayments.Save(payment)
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

	payment, err := c.System().DB.WorkflowPayments.Get(paymentId)
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

	payment, err := getWorkflowPayment(c.System().DB.WorkflowPayments, txHash, user.EthereumAddr, status)
	if err != nil {
		if db.NotFound(err) {
			return c.NoContent(http.StatusNotFound)
		}
		log.Println("[GetWorkflowPayment] GetByTxHashAndStatusAndFromEthAddress err: ", err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	return c.JSON(http.StatusOK, payment)
}

var errRequiredParamMissing = errors.New("required parameter missing")

func getWorkflowPayment(workflowPaymentsDB storage.WorkflowPaymentsIF, txHash,
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

	err = updateWorkflowPaymentPending(c.System().DB.WorkflowPayments, paymentId, updatePaymentRequest.TxHash, user.EthereumAddr)
	if err != nil {
		log.Printf("[UpdateWorkflowPayment] err: %s", err.Error())
		if err == errTxHashEmpty {
			return c.String(http.StatusBadRequest, err.Error())
		}
		return c.NoContent(http.StatusBadRequest)
	}

	return c.NoContent(http.StatusOK)
}

func updateWorkflowPaymentPending(workflowPaymentsDB storage.WorkflowPaymentsIF, paymentId, txHash, ethAddr string) error {
	txHash = strings.TrimSpace(txHash)
	if txHash == "" {
		return errTxHashEmpty
	}

	err := workflowPaymentsDB.Update(paymentId, model.PaymentStatusPending, txHash, ethAddr)
	if err != nil {
		log.Printf("[UpdateWorkflowPayment] WorkflowPayments.Update err: %s", err.Error())
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

	err = cancelWorkflowPayment(c.System().DB.WorkflowPayments, paymentId, user.EthereumAddr)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func cancelWorkflowPayment(workflowPaymentsDB storage.WorkflowPaymentsIF, paymentId, ethAddr string) error {
	return workflowPaymentsDB.Cancel(paymentId, ethAddr)
}

// Set the payment status from confirmed to redeemed
func RedeemPayment(workflowPaymentsDB storage.WorkflowPaymentsIF, workflowId, ethAddr string) error {
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
		if !db.NotFound(err) {
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

	err := c.System().DB.WorkflowPayments.Delete(paymentId)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// List all Payment. Only for superadmin and debugging purposes
func ListPayments(e echo.Context) error {
	c := e.(*www.Context)

	payments, err := c.System().DB.WorkflowPayments.All()
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	return c.JSON(http.StatusOK, payments)
}

func getUser(c *www.Context) (*model.User, error) {
	sess := c.Session(false)
	return c.System().DB.User.Get(sess, sess.UserID())
}

func PutTestPayment(e echo.Context) error {
	c := e.(*www.Context)
	if !c.System().TestMode {
		return echo.ErrBadRequest
	}
	var req struct {
		TxHash string
		From   string
		To     string
	}
	c.Bind(&req)
	l := types.Log{
		Address: common.HexToAddress(cfg.Config.XESContractAddress),
		Topics: []common.Hash{
			common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
			common.HexToHash("0x000000000000000000000000" + req.From[2:]),
			common.HexToHash("0x000000000000000000000000" + req.To[2:]),
		},
		TxHash: common.HexToHash(req.TxHash),
	}
	l.Data = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0d, 0xe0, 0xb6, 0xb3, 0xa7, 0x64, 0x00, 0x00}

	blockchain.TestChannelPayment <- l
	return nil
}
