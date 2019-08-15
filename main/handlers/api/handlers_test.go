package api

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"git.proxeus.com/core/central/main/www"
	"git.proxeus.com/core/central/sys"
	"git.proxeus.com/core/central/sys/db/storm"
	"git.proxeus.com/core/central/sys/model"
)

func TestCheckIfWorkflowNeedsPayment(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("CheckIfWorkflowNeedsPaymentShouldSucceedIfPaymentFound", func(t *testing.T) {

		workflowPaymentItem := &model.WorkflowPaymentItem{Xes: 2000000000000000000, From: "0x1", To: "0x3", Hash: "0x5", WorkflowID: "3"}
		workflowPaymentsDBMock := storm.NewMockWorkflowPaymentsDBInterface(mockCtrl)
		workflowPaymentsDBMock.EXPECT().GetByWorkflowId("3").Return(workflowPaymentItem, nil).Times(1)

		workflow := &model.WorkflowItem{ID: "3", Price: 2000000000000000000}
		workflow.Owner = "33"

		result := checkIfWorkflowNeedsPayment(workflowPaymentsDBMock, workflow, "44")

		assert.NoError(t, result)
	})

	t.Run("CheckIfWorkflowNeedsPaymentShouldSucceedIfFreeWorkflow", func(t *testing.T) {

		workflowPaymentsDBMock := storm.NewMockWorkflowPaymentsDBInterface(mockCtrl)

		workflow := &model.WorkflowItem{ID: "3", Price: 0}
		workflow.Owner = "33"

		result := checkIfWorkflowNeedsPayment(workflowPaymentsDBMock, workflow, "44")

		assert.NoError(t, result)
	})

	t.Run("CheckIfWorkflowNeedsPaymentShouldSucceedIfOwnerStartsWorkflow", func(t *testing.T) {

		workflowPaymentsDBMock := storm.NewMockWorkflowPaymentsDBInterface(mockCtrl)

		workflow := &model.WorkflowItem{ID: "3", Price: 2000000000000000000}
		workflow.Owner = "33"

		result := checkIfWorkflowNeedsPayment(workflowPaymentsDBMock, workflow, "33")

		assert.NoError(t, result)
	})

	t.Run("CheckIfWorkflowNeedsPaymentShouldFailIfNoPaymentFound", func(t *testing.T) {

		workflowPaymentsDBMock := storm.NewMockWorkflowPaymentsDBInterface(mockCtrl)
		workflowPaymentsDBMock.EXPECT().GetByWorkflowId("3").Return(nil, nil).Times(1)

		workflow := &model.WorkflowItem{ID: "3", Price: 2000000000000000000}
		workflow.Owner = "33"

		result := checkIfWorkflowNeedsPayment(workflowPaymentsDBMock, workflow, "44")

		assert.EqualError(t, result, errNoPaymentFound.Error())
	})
}

func TestDeletePaymentIfExists(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("DeletePaymentIfExistsShouldSucceedIfPaymentDeleted", func(t *testing.T) {
		wwwContext := &www.Context{}

		workflowPaymentItem := &model.WorkflowPaymentItem{Xes: 2000000000000000000, From: "0x1", To: "0x3", Hash: "0x5"}
		workflowPaymentsDBMock := storm.NewMockWorkflowPaymentsDBInterface(mockCtrl)
		workflowPaymentsDBMock.EXPECT().GetByWorkflowIdAndFromEthAddress(gomock.Eq("1"), gomock.Eq("0x1")).Return(workflowPaymentItem, nil).Times(1)
		workflowPaymentsDBMock.EXPECT().Delete(gomock.Eq(workflowPaymentItem.Hash)).Return(nil).Times(1)

		system := &sys.System{}
		system.DB = &storm.DBSet{WorkflowPaymentsDB: workflowPaymentsDBMock}
		www.SetSystem(system)

		result := DeletePaymentIfExists(wwwContext, "1", "0x1")
		assert.NoError(t, result)
	})

	t.Run("DeletePaymentIfExistsShouldSucceedIfNoPaymentToDelete", func(t *testing.T) {
		wwwContext := &www.Context{}

		workflowPaymentsDBMock := storm.NewMockWorkflowPaymentsDBInterface(mockCtrl)
		err := errors.New("not found")
		workflowPaymentsDBMock.EXPECT().GetByWorkflowIdAndFromEthAddress(gomock.Eq("1"), gomock.Eq("0x1")).Return(nil, err).Times(1)

		system := &sys.System{}
		system.DB = &storm.DBSet{WorkflowPaymentsDB: workflowPaymentsDBMock}
		www.SetSystem(system)

		result := DeletePaymentIfExists(wwwContext, "1", "0x1")
		assert.NoError(t, result)
	})

	t.Run("DeletePaymentIfExistsShouldFailOnGetError", func(t *testing.T) {
		wwwContext := &www.Context{}

		workflowPaymentsDBMock := storm.NewMockWorkflowPaymentsDBInterface(mockCtrl)
		err := errors.New("some error")
		workflowPaymentsDBMock.EXPECT().GetByWorkflowIdAndFromEthAddress(gomock.Eq("1"), gomock.Eq("0x1")).Return(nil, err).Times(1)

		system := &sys.System{}
		system.DB = &storm.DBSet{WorkflowPaymentsDB: workflowPaymentsDBMock}
		www.SetSystem(system)

		result := DeletePaymentIfExists(wwwContext, "1", "0x1")
		assert.Error(t, result)
	})

	t.Run("DeletePaymentIfExistsShouldFailOnDeleteError", func(t *testing.T) {
		wwwContext := &www.Context{}

		workflowPaymentItem := &model.WorkflowPaymentItem{Xes: 2000000000000000000, From: "0x1", To: "0x3", Hash: "0x5"}
		workflowPaymentsDBMock := storm.NewMockWorkflowPaymentsDBInterface(mockCtrl)
		workflowPaymentsDBMock.EXPECT().GetByWorkflowIdAndFromEthAddress(gomock.Eq("1"), gomock.Eq("0x1")).Return(workflowPaymentItem, nil).Times(1)

		err := errors.New("some error")

		workflowPaymentsDBMock.EXPECT().Delete(gomock.Eq(workflowPaymentItem.Hash)).Return(err).Times(1)

		system := &sys.System{}
		system.DB = &storm.DBSet{WorkflowPaymentsDB: workflowPaymentsDBMock}
		www.SetSystem(system)

		result := DeletePaymentIfExists(wwwContext, "1", "0x1")
		assert.Error(t, result)
	})

}
