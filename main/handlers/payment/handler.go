package payment

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/ProxeusApp/proxeus-core/service"
	"github.com/ProxeusApp/proxeus-core/storage/database/db"

	cfg "github.com/ProxeusApp/proxeus-core/main/config"
	"github.com/ProxeusApp/proxeus-core/main/handlers/blockchain"
	"github.com/ProxeusApp/proxeus-core/main/www"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/labstack/echo"
)

var (
	paymentService service.PaymentService
	userService    service.UserService

	errNotAuthorized = errors.New("user not authorized")
)

type (
	createPaymentRequest struct {
		WorkflowId string `json:"workflowId"`
	}
	updatePaymentPendingRequest struct {
		TxHash string `json:"txHash"`
	}
)

func Init(paymentS service.PaymentService, userS service.UserService) {
	paymentService = paymentS
	userService = userS
}

//create a payment for a workflow
func CreateWorkflowPayment(e echo.Context) error {
	c := e.(*www.Context)

	user, err := userService.GetUser(c.Session(false))
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

	payment, err := paymentService.CreateWorkflowPayment(c.Session(false),
		createPaymentRequest.WorkflowId, user.EthereumAddr)
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

	user, err := userService.GetUser(c.Session(false))
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	payment, err := paymentService.GetWorkflowPaymentById(paymentId)
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

	user, err := userService.GetUser(c.Session(false))
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	payment, err := paymentService.GetWorkflowPayment(txHash, user.EthereumAddr, status)
	if err != nil {
		if db.NotFound(err) {
			return c.NoContent(http.StatusNotFound)
		}
		log.Println("[GetWorkflowPayment] GetByTxHashAndStatusAndFromEthAddress err: ", err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	return c.JSON(http.StatusOK, payment)
}

// Set a workflow payment from status created to status pending
func UpdateWorkflowPaymentPending(e echo.Context) error {
	c := e.(*www.Context)
	paymentId := strings.TrimSpace(c.Param("paymentId"))

	user, err := userService.GetUser(c.Session(false))
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	updatePaymentRequest := &updatePaymentPendingRequest{}
	err = c.Bind(&updatePaymentRequest)
	if err != nil {
		log.Printf("[UpdateWorkflowPayment] UpdateWorkflowPayment bind err: %s", err.Error())
		return err
	}

	err = paymentService.UpdateWorkflowPaymentPending(paymentId, updatePaymentRequest.TxHash, user.EthereumAddr)
	if err != nil {
		log.Printf("[UpdateWorkflowPayment] err: %s", err.Error())
		if err == service.ErrTxHashEmpty {
			return c.String(http.StatusBadRequest, err.Error())
		}
		return c.NoContent(http.StatusBadRequest)
	}

	return c.NoContent(http.StatusOK)
}

// Set status of workflow from created to cancelled
func CancelWorkflowPayment(e echo.Context) error {
	c := e.(*www.Context)
	paymentId := strings.TrimSpace(c.Param("paymentId"))

	user, err := userService.GetUser(c.Session(false))
	if err != nil {
		return errNotAuthorized
	}

	err = paymentService.CancelWorkflowPayment(paymentId, user.EthereumAddr)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// Set Payment for a workflow to status = Deleted. Only for superadmin
func DeleteWorkflowPayment(e echo.Context) error {
	c := e.(*www.Context)
	paymentId := strings.TrimSpace(c.Param("paymentId"))

	if paymentId == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	err := paymentService.Delete(paymentId)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// List all Payment. Only for superadmin and debugging purposes
func ListPayments(e echo.Context) error {
	c := e.(*www.Context)

	payments, err := paymentService.All()
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	return c.JSON(http.StatusOK, payments)
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

	log.Println("DEBUG PUSH LOG", l)
	blockchain.TestChannelPayment <- l
	return nil
}
