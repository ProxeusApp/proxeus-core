package payment

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	"github.com/ProxeusApp/proxeus-core/main/www"
	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/storage/database"
	"github.com/ProxeusApp/proxeus-core/storage/database/db"
	sm "github.com/ProxeusApp/proxeus-core/storage/database/mock"
	"github.com/ProxeusApp/proxeus-core/sys"
	"github.com/ProxeusApp/proxeus-core/sys/model"
)

func setupPaymentRequestTest(httpMethod, targetUrl, body string) (*www.Context, *httptest.ResponseRecorder, *model.User, *model.User) {
	e := echo.New()
	req := httptest.NewRequest(httpMethod, targetUrl, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	sessionStore := sessions.NewCookieStore([]byte("secret_Dummy_1234"), []byte("12345678901234567890123456789012"))
	c.Set("_session_store", sessionStore)
	sysSession := &sys.Session{S: &model.Session{}}
	sysSession.S.UsrID = "1"

	c.Set("sys.session", sysSession)
	wwwContext := &www.Context{Context: c}
	wwwContext.SetRequest(req)

	user := &model.User{}
	user.EthereumAddr = "0x00"

	ownerUser := &model.User{}
	ownerUser.EthereumAddr = "0x3"

	return wwwContext, rec, user, ownerUser
}

type paymentResponse struct {
	Id string `json:id`
}

func TestCreateWorkflowPayment(t *testing.T) {

	mockCtrl, workflowPaymentsDB := up(t)
	defer down(mockCtrl, workflowPaymentsDB)

	body := `{"workflowId":"552a2f0e-c6c4-403b-8aaf-2d9ebf55eb8f"}`

	wwwContext, rec, user, _ := setupPaymentRequestTest(http.MethodPost, "/api/admin/payments", body)

	userDBMock := sm.NewMockUserIF(mockCtrl)
	userDBMock.EXPECT().Get(gomock.Any(), gomock.Eq("1")).Return(user, nil).Times(1)

	workflow := &model.WorkflowItem{Price: 2000000000000000000}
	workflow.Owner = user.EthereumAddr
	workflowDBMock := sm.NewMockWorkflowIF(mockCtrl)
	workflowDBMock.EXPECT().Get(gomock.Any(), gomock.Any()).Return(workflow, nil)

	system := &sys.System{}
	system.DB = &storage.DBSet{WorkflowPayments: workflowPaymentsDB, User: userDBMock, Workflow: workflowDBMock}
	www.SetSystem(system)

	t.Run("ShouldCreatePaymentItem", func(t *testing.T) {
		if assert.NoError(t, CreateWorkflowPayment(wwwContext)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response = paymentResponse{}
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			if err != nil {
				panic(err)
			}

			err = workflowPaymentsDB.Remove(&model.WorkflowPaymentItem{ID: response.Id})
			if err != nil {
				panic(err)
			}
		}
	})
}

func TestGetWorkflowPaymentById(t *testing.T) {
	mockCtrl, workflowPaymentsDB := up(t)
	defer down(mockCtrl, workflowPaymentsDB)

	t.Run("ShouldReturnPayment", func(t *testing.T) {
		paymentId := "1"

		wwwContext, rec, user, _ := setupPaymentRequestTest(http.MethodGet,
			fmt.Sprintf("/api/admin/payments/%s", paymentId), "{}")

		wwwContext.SetParamNames("paymentId")
		wwwContext.SetParamValues(paymentId)

		userDBMock := sm.NewMockUserIF(mockCtrl)
		userDBMock.EXPECT().Get(gomock.Any(), gomock.Eq("1")).Return(user, nil).Times(1)

		system := &sys.System{}
		system.DB = &storage.DBSet{WorkflowPayments: workflowPaymentsDB, User: userDBMock}
		www.SetSystem(system)

		err := workflowPaymentsDB.Save(&model.WorkflowPaymentItem{ID: paymentId, From: user.EthereumAddr})
		if err != nil {
			panic(err)
		}

		if assert.NoError(t, GetWorkflowPaymentById(wwwContext)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response = paymentResponse{}
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			if err != nil {
				panic(err)
			}

			err = workflowPaymentsDB.Remove(&model.WorkflowPaymentItem{ID: response.Id})
			if err != nil {
				panic(err)
			}
		}
	})

	t.Run("ShouldNotReturnPayment", func(t *testing.T) {
		paymentId := "2"

		wwwContext, rec, user, userOwner := setupPaymentRequestTest(http.MethodGet,
			fmt.Sprintf("/api/admin/payments/%s", paymentId), "{}")

		wwwContext.SetParamNames("paymentId")
		wwwContext.SetParamValues(paymentId)

		userDBMock := sm.NewMockUserIF(mockCtrl)
		userDBMock.EXPECT().Get(gomock.Any(), gomock.Eq("1")).Return(user, nil).Times(1)

		system := &sys.System{}
		system.DB = &storage.DBSet{WorkflowPayments: workflowPaymentsDB, User: userDBMock}
		www.SetSystem(system)

		//here pass "userOwner" instead of "user"
		err := workflowPaymentsDB.Save(&model.WorkflowPaymentItem{ID: paymentId, From: userOwner.EthereumAddr})
		if err != nil {
			panic(err)
		}

		if assert.NoError(t, GetWorkflowPaymentById(wwwContext)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)

			err = workflowPaymentsDB.Remove(&model.WorkflowPaymentItem{ID: paymentId})
			if err != nil {
				panic(err)
			}
		}
	})
}

func TestGetWorkflowPayment(t *testing.T) {
	mockCtrl, workflowPaymentsDB := up(t)
	defer down(mockCtrl, workflowPaymentsDB)

	t.Run("ShouldReturnPayment", func(t *testing.T) {
		paymentId := "3"
		txHash := "0x3"
		from := "0x4"
		status := model.PaymentStatusConfirmed

		err := workflowPaymentsDB.Save(&model.WorkflowPaymentItem{ID: paymentId, TxHash: txHash, From: from, Status: status})
		if err != nil {
			panic(err)
		}

		payment, err := getWorkflowPayment(workflowPaymentsDB, txHash, from, status)

		assert.Nil(t, err)
		assert.Equal(t, paymentId, payment.ID)

		err = workflowPaymentsDB.Remove(payment)
		if err != nil {
			panic(err)
		}
	})
	t.Run("ShouldNotReturnPaymentIfFromNotMatching", func(t *testing.T) {
		paymentId := "4"
		txHash := "0x3"
		from := "0x4"
		status := model.PaymentStatusConfirmed

		err := workflowPaymentsDB.Save(&model.WorkflowPaymentItem{ID: paymentId, TxHash: txHash, From: "0x5", Status: status})
		if err != nil {
			panic(err)
		}

		payment, err := getWorkflowPayment(workflowPaymentsDB, txHash, from, status)

		if !db.NotFound(err) {
			t.Errorf("Expected to have not found but got: %v", err)
		}
		assert.Nil(t, payment)

		err = workflowPaymentsDB.Remove(&model.WorkflowPaymentItem{ID: paymentId})
		if err != nil {
			panic(err)
		}
	})
}

func TestUpdateWorkflowPaymentPending(t *testing.T) {

	mockCtrl, workflowPaymentsDB := up(t)
	defer down(mockCtrl, workflowPaymentsDB)

	t.Run("ShouldUpdatePayment", func(t *testing.T) {
		paymentId := "3"
		txHash := "0x3"
		from := "0x4"

		err := workflowPaymentsDB.Save(&model.WorkflowPaymentItem{ID: paymentId, From: from, Status: model.PaymentStatusCreated})
		if err != nil {
			panic(err)
		}

		assert.NoError(t, updateWorkflowPaymentPending(workflowPaymentsDB, paymentId, txHash, from))

		payment, err := workflowPaymentsDB.Get(paymentId)
		if err != nil {
			panic(err)
		}

		assert.Equal(t, paymentId, payment.ID)
		assert.Equal(t, model.PaymentStatusPending, payment.Status)
		assert.Equal(t, txHash, payment.TxHash)

		err = workflowPaymentsDB.Remove(payment)
		if err != nil {
			panic(err)
		}
	})
	t.Run("ShouldReturnErrorOnUpdatePaymentIfFromNotMatching", func(t *testing.T) {
		paymentId := "4"
		txHash := "0x4"
		from := "0x5"

		err := workflowPaymentsDB.Save(&model.WorkflowPaymentItem{ID: paymentId, From: from, Status: model.PaymentStatusCreated})
		if err != nil {
			panic(err)
		}

		err = updateWorkflowPaymentPending(workflowPaymentsDB, paymentId, txHash, "0x6")
		if !db.NotFound(err) {
			t.Error("expected not found err, got", err)
		}

		payment, err := workflowPaymentsDB.Get(paymentId)
		if err != nil {
			panic(err)
		}

		assert.Equal(t, paymentId, payment.ID)
		assert.Equal(t, model.PaymentStatusCreated, payment.Status)
		assert.Equal(t, "", payment.TxHash)

		err = workflowPaymentsDB.Remove(payment)
		if err != nil {
			panic(err)
		}
	})
}

func TestCancelWorkflowPayment(t *testing.T) {
	mockCtrl, workflowPaymentsDB := up(t)
	defer down(mockCtrl, workflowPaymentsDB)

	t.Run("ShouldCancelWorkflowPayment", func(t *testing.T) {
		paymentId := "4"
		txHash := "0x4"
		from := "0x5"

		err := workflowPaymentsDB.Save(&model.WorkflowPaymentItem{ID: paymentId, TxHash: txHash, From: from,
			Status: model.PaymentStatusCreated})
		if err != nil {
			panic(err)
		}

		assert.NoError(t, cancelWorkflowPayment(workflowPaymentsDB, paymentId, from))

		payment, err := workflowPaymentsDB.Get(paymentId)
		if err != nil {
			panic(err)
		}

		assert.Equal(t, model.PaymentStatusCancelled, payment.Status)

		err = workflowPaymentsDB.Remove(payment)
		if err != nil {
			panic(err)
		}
	})

	t.Run("ShouldNotCancelWorkflowPaymentIfStatusIsNotPending", func(t *testing.T) {
		paymentId := "5"
		txHash := "0x5"
		from := "0x6"

		err := workflowPaymentsDB.Save(&model.WorkflowPaymentItem{ID: paymentId, TxHash: txHash, From: from,
			Status: model.PaymentStatusConfirmed})
		if err != nil {
			panic(err)
		}

		assert.Error(t, cancelWorkflowPayment(workflowPaymentsDB, paymentId, from))

		payment, err := workflowPaymentsDB.Get(paymentId)
		if err != nil {
			panic(err)
		}

		assert.Equal(t, model.PaymentStatusConfirmed, payment.Status)

		err = workflowPaymentsDB.Remove(payment)
		if err != nil {
			panic(err)
		}
	})

	t.Run("ShouldNotCancelWorkflowPaymentIfFromNotMatching", func(t *testing.T) {
		paymentId := "6"
		txHash := "0x6"
		from := "0x7"

		err := workflowPaymentsDB.Save(&model.WorkflowPaymentItem{ID: paymentId, TxHash: txHash, From: from,
			Status: model.PaymentStatusCreated})
		if err != nil {
			panic(err)
		}

		assert.Error(t, cancelWorkflowPayment(workflowPaymentsDB, paymentId, "0x8"))

		payment, err := workflowPaymentsDB.Get(paymentId)
		if err != nil {
			panic(err)
		}

		assert.Equal(t, model.PaymentStatusCreated, payment.Status)

		err = workflowPaymentsDB.Remove(payment)
		if err != nil {
			panic(err)
		}
	})
}

func TestRedeemPayment(t *testing.T) {
	mockCtrl, workflowPaymentsDB := up(t)
	defer down(mockCtrl, workflowPaymentsDB)

	t.Run("ShouldRedeemWorkflowPayment", func(t *testing.T) {
		paymentId := "7"
		txHash := "0x7"
		from := "0x8"
		workflowId := "01-02"

		err := workflowPaymentsDB.Save(&model.WorkflowPaymentItem{ID: paymentId, TxHash: txHash, From: from,
			Status: model.PaymentStatusConfirmed, WorkflowID: workflowId})
		if err != nil {
			panic(err)
		}

		assert.NoError(t, RedeemPayment(workflowPaymentsDB, workflowId, from))

		payment, err := workflowPaymentsDB.Get(paymentId)
		if err != nil {
			panic(err)
		}

		assert.Equal(t, model.PaymentStatusRedeemed, payment.Status)

		err = workflowPaymentsDB.Remove(payment)
		if err != nil {
			panic(err)
		}
	})

	t.Run("ShouldRedeemNewerPaymentItemIfTwoAvailable", func(t *testing.T) {
		err := workflowPaymentsDB.Save(&model.WorkflowPaymentItem{ID: "8first", TxHash: "0x8", From: "0x9",
			Status: model.PaymentStatusConfirmed, WorkflowID: "01-03"})
		if err != nil {
			panic(err)
		}

		err = workflowPaymentsDB.Save(&model.WorkflowPaymentItem{ID: "9second", TxHash: "0x9", From: "0x9",
			Status: model.PaymentStatusConfirmed, WorkflowID: "01-03"})
		if err != nil {
			panic(err)
		}

		assert.NoError(t, RedeemPayment(workflowPaymentsDB, "01-03", "0x9"))

		firstPayment, err := workflowPaymentsDB.Get("8first")
		if err != nil {
			panic(err)
		}

		assert.Equal(t, model.PaymentStatusConfirmed, firstPayment.Status)

		secondPayment, err := workflowPaymentsDB.Get("9second")
		if err != nil {
			panic(err)
		}

		assert.Equal(t, model.PaymentStatusRedeemed, secondPayment.Status)

		err = workflowPaymentsDB.Remove(firstPayment)
		if err != nil {
			panic(err)
		}
		err = workflowPaymentsDB.Remove(secondPayment)
		if err != nil {
			panic(err)
		}
	})

	t.Run("ShouldNotRedeemWorkflowPaymentIfStatusIsPending", func(t *testing.T) {
		paymentId := "9"
		txHash := "0x9"
		from := "0x10"
		workflowId := "01-03"

		err := workflowPaymentsDB.Save(&model.WorkflowPaymentItem{ID: paymentId, TxHash: txHash, From: from,
			Status: model.PaymentStatusPending, WorkflowID: workflowId})
		if err != nil {
			panic(err)
		}

		assert.Error(t, RedeemPayment(workflowPaymentsDB, workflowId, from))

		payment, err := workflowPaymentsDB.Get(paymentId)
		if err != nil {
			panic(err)
		}

		assert.Equal(t, model.PaymentStatusPending, payment.Status)

		err = workflowPaymentsDB.Remove(payment)
		if err != nil {
			panic(err)
		}
	})

	t.Run("ShouldNotRedeemWorkflowPaymentIfFromNotMatching", func(t *testing.T) {
		paymentId := "10"
		txHash := "0x10"
		from := "0x11"
		workflowId := "01-04"

		err := workflowPaymentsDB.Save(&model.WorkflowPaymentItem{ID: paymentId, TxHash: txHash, From: from,
			Status: model.PaymentStatusConfirmed, WorkflowID: workflowId})
		if err != nil {
			panic(err)
		}

		assert.Error(t, RedeemPayment(workflowPaymentsDB, workflowId, "0x12"))

		payment, err := workflowPaymentsDB.Get(paymentId)
		if err != nil {
			panic(err)
		}

		assert.Equal(t, model.PaymentStatusConfirmed, payment.Status)

		err = workflowPaymentsDB.Remove(payment)
		if err != nil {
			panic(err)
		}
	})
}

func TestCheckIfWorkflowPaymentRequired(t *testing.T) {
	mockCtrl, workflowPaymentsDB := up(t)
	defer down(mockCtrl, workflowPaymentsDB)

	t.Run("ShouldRequirePaymentIfWorkflowNotForFree", func(t *testing.T) {
		permissions := &model.Permissions{Owner: "1"}
		workflow := &model.WorkflowItem{Price: 2, Permissions: *permissions}
		assert.True(t, isPaymentRequired(false, workflow, "2"))
	})
	t.Run("ShouldNotRequirePaymentIfWorkflowNotForFreeButAlreadyStarted", func(t *testing.T) {
		permissions := &model.Permissions{Owner: "1"}
		workflow := &model.WorkflowItem{Price: 2, Permissions: *permissions}
		assert.False(t, isPaymentRequired(true, workflow, "2"))
	})
	t.Run("ShouldNotRequirePaymentIfWorkflowIsFree", func(t *testing.T) {
		permissions := &model.Permissions{Owner: "1"}
		workflow := &model.WorkflowItem{Price: 0, Permissions: *permissions}
		assert.False(t, isPaymentRequired(false, workflow, "2"))
	})
	t.Run("ShouldNotRequirePaymentForWorkflowOwner", func(t *testing.T) {
		permissions := &model.Permissions{Owner: "1"}
		workflow := &model.WorkflowItem{Price: 2, Permissions: *permissions}
		assert.False(t, isPaymentRequired(false, workflow, "1"))
	})
}

var errCleanupTestData = errors.New("db data has not been cleanup up after finishing tests")

func up(t *testing.T) (*gomock.Controller, storage.WorkflowPaymentsIF) {
	mockCtrl := gomock.NewController(t)

	workflowPaymentsDB, err := database.NewWorkflowPaymentDB(database.DBConfig{Dir: ".test_data"})
	if err != nil {
		panic(err)
	}
	return mockCtrl, workflowPaymentsDB
}

func down(mockCtrl *gomock.Controller, workflowPaymentsDB storage.WorkflowPaymentsIF) {

	mockCtrl.Finish()

	payments, err := workflowPaymentsDB.All()
	if err != nil {
		panic(err)
	}
	if len(payments) != 0 {
		panic(errCleanupTestData)
	}

	err = os.Remove(filepath.Join(".test_data", database.WorkflowPaymentDBDir, database.WorkflowPaymentDB))
	if err != nil {
		panic(err.Error())
	}
	err = os.Remove(filepath.Join(".test_data", database.WorkflowPaymentDBDir))
	if err != nil {
		panic(err.Error())
	}
	err = os.Remove(".test_data")
	if err != nil {
		panic(err.Error())
	}

}
