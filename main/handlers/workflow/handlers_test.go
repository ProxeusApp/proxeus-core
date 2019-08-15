package workflow

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"

	"git.proxeus.com/core/central/sys"
	"git.proxeus.com/core/central/sys/db/storm"
	"git.proxeus.com/core/central/sys/model"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	"git.proxeus.com/core/central/main/www"
	sysSess "git.proxeus.com/core/central/sys/session"
)

func setupPaymentTest(httpMethod, targetUrl string) (*www.Context, *httptest.ResponseRecorder, *model.User, *model.User) {
	e := echo.New()
	req := httptest.NewRequest(httpMethod, targetUrl, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	sessionStore := sessions.NewCookieStore([]byte("secret_Dummy_1234"), []byte("12345678901234567890123456789012"))
	c.Set("_session_store", sessionStore)
	sysSession := &sysSess.Session{}
	sysSession.SetUserID("1")

	c.Set("sys.session", sysSession)
	wwwContext := &www.Context{Context: c}
	wwwContext.SetRequest(req)

	user := &model.User{}
	user.EthereumAddr = "0x00"

	ownerUser := &model.User{}
	ownerUser.EthereumAddr = "0x3"

	return wwwContext, rec, user, ownerUser
}

func TestAddWorkflowPayment(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("AddWorkflowPaymentShouldSucceed", func(t *testing.T) {
		wwwContext, rec, user, ownerUser := setupPaymentTest(http.MethodPost, "/api/admin/workflow/1/payment/0x2222")

		userDBMock := storm.NewMockUserDBInterface(mockCtrl)
		userDBMock.EXPECT().Get(gomock.Any(), gomock.Eq("1")).Return(user, nil).Times(1)
		userDBMock.EXPECT().Get(gomock.Any(), gomock.Eq(ownerUser.EthereumAddr)).Return(ownerUser, nil).Times(1)

		workflowPaymentItem := &model.WorkflowPaymentItem{Xes: 2000000000000000000, From: user.EthereumAddr, To: "0x3", Hash: "0x5"}
		workflowPaymentsDBMock := storm.NewMockWorkflowPaymentsDBInterface(mockCtrl)
		workflowPaymentsDBMock.EXPECT().GetByTxHash(gomock.Any()).Return(workflowPaymentItem, nil).Times(1)
		workflowPaymentsDBMock.EXPECT().Add(gomock.Any()).Return(nil).Times(1)

		workflow := &model.WorkflowItem{Price: 2000000000000000000}
		workflow.Owner = ownerUser.EthereumAddr
		workflowDBMock := storm.NewMockWorkflowDBInterface(mockCtrl)
		workflowDBMock.EXPECT().Get(gomock.Any(), gomock.Any()).Return(workflow, nil)

		system := &sys.System{}
		system.DB = &storm.DBSet{WorkflowPaymentsDB: workflowPaymentsDBMock, User: userDBMock, Workflow: workflowDBMock}
		www.SetSystem(system)

		if assert.NoError(t, AddWorkflowPayment(wwwContext)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("AddWorkflowPaymentShouldFailIncorrectPayer", func(t *testing.T) {
		wwwContext, rec, user, ownerUser := setupPaymentTest(http.MethodPost, "/api/admin/workflow/1/payment/0x2222")

		userDBMock := storm.NewMockUserDBInterface(mockCtrl)
		userDBMock.EXPECT().Get(gomock.Any(), gomock.Eq("1")).Return(user, nil).Times(1)

		workflowPaymentItem := &model.WorkflowPaymentItem{Xes: 2000000000000000000, From: "0xWrong", To: "0x3", Hash: "0x5"}
		workflowPaymentsDBMock := storm.NewMockWorkflowPaymentsDBInterface(mockCtrl)
		workflowPaymentsDBMock.EXPECT().GetByTxHash(gomock.Any()).Return(workflowPaymentItem, nil).Times(1)

		workflow := &model.WorkflowItem{Price: 2000000000000000000}
		workflow.Owner = ownerUser.EthereumAddr
		workflowDBMock := storm.NewMockWorkflowDBInterface(mockCtrl)
		workflowDBMock.EXPECT().Get(gomock.Any(), gomock.Any()).Return(workflow, nil)

		system := &sys.System{}
		system.DB = &storm.DBSet{WorkflowPaymentsDB: workflowPaymentsDBMock, User: userDBMock, Workflow: workflowDBMock}
		www.SetSystem(system)

		if assert.NoError(t, AddWorkflowPayment(wwwContext)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			responseBody := rec.Body.String()
			assert.Equal(t, errPaymentFailed.Error(), responseBody)
		}
	})

	t.Run("AddWorkflowPaymentShouldFailPaymentItemWorkflowIDAlreadySet", func(t *testing.T) {
		wwwContext, rec, user, _ := setupPaymentTest(http.MethodPost, "/api/admin/workflow/1/payment/0x2222")

		workflowPaymentItem := &model.WorkflowPaymentItem{Xes: 2000000000000000000, From: user.EthereumAddr, To: "0x3", WorkflowID: "3", Hash: "0x5"}
		workflowPaymentsDBMock := storm.NewMockWorkflowPaymentsDBInterface(mockCtrl)
		workflowPaymentsDBMock.EXPECT().GetByTxHash(gomock.Any()).Return(workflowPaymentItem, nil).Times(1)

		system := &sys.System{}
		system.DB = &storm.DBSet{WorkflowPaymentsDB: workflowPaymentsDBMock}
		www.SetSystem(system)

		if assert.NoError(t, AddWorkflowPayment(wwwContext)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})
}

func TestGetWorkflowPayment(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("GetWorkflowPaymentByWorkflowIdAndFromEthAddressShouldSucceed", func(t *testing.T) {

		wwwContext, rec, user, ownerUser := setupPaymentTest(http.MethodGet, "/api/admin/workflow/1/payment")

		userDBMock := storm.NewMockUserDBInterface(mockCtrl)
		userDBMock.EXPECT().Get(gomock.Any(), gomock.Eq("1")).Return(user, nil).Times(2)
		userDBMock.EXPECT().Get(gomock.Any(), gomock.Eq(ownerUser.EthereumAddr)).Return(ownerUser, nil).Times(1)

		workflowPaymentItem := &model.WorkflowPaymentItem{Xes: 2000000000000000000, From: user.EthereumAddr, To: "0x3", WorkflowID: "3", Hash: "0x5"}
		workflowPaymentsDBMock := storm.NewMockWorkflowPaymentsDBInterface(mockCtrl)
		workflowPaymentsDBMock.EXPECT().GetByWorkflowIdAndFromEthAddress(gomock.Eq(""), gomock.Eq(user.EthereumAddr)).Return(workflowPaymentItem, nil).Times(1)

		workflow := &model.WorkflowItem{Price: 2000000000000000000}
		workflow.Owner = ownerUser.EthereumAddr
		workflowDBMock := storm.NewMockWorkflowDBInterface(mockCtrl)
		workflowDBMock.EXPECT().Get(gomock.Any(), gomock.Any()).Return(workflow, nil)

		system := &sys.System{}
		system.DB = &storm.DBSet{WorkflowPaymentsDB: workflowPaymentsDBMock, User: userDBMock, Workflow: workflowDBMock}
		www.SetSystem(system)

		if assert.NoError(t, GetWorkflowPayment(wwwContext)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			responseBody := rec.Body.String()
			successResponseJSON := `{"hash":"0x5","workflowID":"3","From":"0x00","To":"0x3","xes":2000000000000000000}`
			assert.Equal(t, successResponseJSON, strings.Trim(responseBody, "\n"))
		}
	})

	t.Run("GetWorkflowPaymentByTxHashShouldSucceed", func(t *testing.T) {

		wwwContext, rec, user, _ := setupPaymentTest(http.MethodGet, "/api/admin/workflow/1/payment?txHash=0x2222")

		workflowPaymentItem := &model.WorkflowPaymentItem{Xes: 2000000000000000000, From: user.EthereumAddr, To: "0x3", WorkflowID: "3", Hash: "0x5"}
		workflowPaymentsDBMock := storm.NewMockWorkflowPaymentsDBInterface(mockCtrl)
		workflowPaymentsDBMock.EXPECT().GetByTxHash(gomock.Eq("0x2222")).Return(workflowPaymentItem, nil).Times(1)

		system := &sys.System{}
		system.DB = &storm.DBSet{WorkflowPaymentsDB: workflowPaymentsDBMock}
		www.SetSystem(system)

		if assert.NoError(t, GetWorkflowPayment(wwwContext)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			responseBody := rec.Body.String()
			successResponseJSON := `{"hash":"0x5","workflowID":"3","From":"0x00","To":"0x3","xes":2000000000000000000}`
			assert.Equal(t, successResponseJSON, strings.Trim(responseBody, "\n"))
		}
	})

	t.Run("GetWorkflowPaymentShouldFailIncorrectPaymentAmount", func(t *testing.T) {

		wwwContext, rec, user, ownerUser := setupPaymentTest(http.MethodGet, "/api/admin/workflow/1/payment")

		userDBMock := storm.NewMockUserDBInterface(mockCtrl)
		userDBMock.EXPECT().Get(gomock.Any(), gomock.Eq("1")).Return(user, nil).Times(1)

		workflowPaymentItem := &model.WorkflowPaymentItem{Xes: 2100000000000000000, From: user.EthereumAddr, To: "0x3", WorkflowID: "3", Hash: "0x5"}
		workflowPaymentsDBMock := storm.NewMockWorkflowPaymentsDBInterface(mockCtrl)
		workflowPaymentsDBMock.EXPECT().GetByWorkflowIdAndFromEthAddress(gomock.Eq(""), gomock.Eq(user.EthereumAddr)).Return(workflowPaymentItem, nil).Times(1)

		workflow := &model.WorkflowItem{Price: 1900000000000000000}
		workflow.Owner = ownerUser.EthereumAddr
		workflowDBMock := storm.NewMockWorkflowDBInterface(mockCtrl)
		workflowDBMock.EXPECT().Get(gomock.Any(), gomock.Any()).Return(workflow, nil)

		system := &sys.System{}
		system.DB = &storm.DBSet{WorkflowPaymentsDB: workflowPaymentsDBMock, User: userDBMock, Workflow: workflowDBMock}
		www.SetSystem(system)

		if assert.NoError(t, GetWorkflowPayment(wwwContext)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			responseBody := rec.Body.String()
			assert.Equal(t, errPaymentFailed.Error(), responseBody)
		}
	})

	t.Run("GetWorkflowPaymentShouldFailIncorrectPayee", func(t *testing.T) {

		wwwContext, rec, user, ownerUser := setupPaymentTest(http.MethodGet, "/api/admin/workflow/1/payment")

		userDBMock := storm.NewMockUserDBInterface(mockCtrl)
		userDBMock.EXPECT().Get(gomock.Any(), gomock.Eq("1")).Return(user, nil).Times(2)
		userDBMock.EXPECT().Get(gomock.Any(), gomock.Eq(ownerUser.EthereumAddr)).Return(ownerUser, nil).Times(1)

		workflowPaymentItem := &model.WorkflowPaymentItem{Xes: 2100000000000000000, From: user.EthereumAddr, To: "0xWrong", WorkflowID: "3", Hash: "0x5"}
		workflowPaymentsDBMock := storm.NewMockWorkflowPaymentsDBInterface(mockCtrl)
		workflowPaymentsDBMock.EXPECT().GetByWorkflowIdAndFromEthAddress(gomock.Eq(""), gomock.Eq(user.EthereumAddr)).Return(workflowPaymentItem, nil).Times(1)

		workflow := &model.WorkflowItem{Price: 2100000000000000000}
		workflow.Owner = ownerUser.EthereumAddr
		workflowDBMock := storm.NewMockWorkflowDBInterface(mockCtrl)
		workflowDBMock.EXPECT().Get(gomock.Any(), gomock.Any()).Return(workflow, nil)

		system := &sys.System{}
		system.DB = &storm.DBSet{WorkflowPaymentsDB: workflowPaymentsDBMock, User: userDBMock, Workflow: workflowDBMock}
		www.SetSystem(system)

		if assert.NoError(t, GetWorkflowPayment(wwwContext)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			responseBody := rec.Body.String()
			assert.Equal(t, errPaymentFailed.Error(), responseBody)
		}
	})
}
