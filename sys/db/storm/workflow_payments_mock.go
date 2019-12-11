// +build !coverage
// Code generated by MockGen. DO NOT EDIT.
// Source: sys/db/storm/workflow_payments.go

// Package storm is a generated GoMock package.
package storm

import (
	reflect "reflect"
	time "time"

	model "github.com/ProxeusApp/proxeus-core/sys/model"
	gomock "github.com/golang/mock/gomock"
)

// MockWorkflowPaymentsDBInterface is a mock of WorkflowPaymentsDBInterface interface
type MockWorkflowPaymentsDBInterface struct {
	ctrl     *gomock.Controller
	recorder *MockWorkflowPaymentsDBInterfaceMockRecorder
}

// MockWorkflowPaymentsDBInterfaceMockRecorder is the mock recorder for MockWorkflowPaymentsDBInterface
type MockWorkflowPaymentsDBInterfaceMockRecorder struct {
	mock *MockWorkflowPaymentsDBInterface
}

// NewMockWorkflowPaymentsDBInterface creates a new mock instance
func NewMockWorkflowPaymentsDBInterface(ctrl *gomock.Controller) *MockWorkflowPaymentsDBInterface {
	mock := &MockWorkflowPaymentsDBInterface{ctrl: ctrl}
	mock.recorder = &MockWorkflowPaymentsDBInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockWorkflowPaymentsDBInterface) EXPECT() *MockWorkflowPaymentsDBInterfaceMockRecorder {
	return m.recorder
}

// GetByTxHashAndStatusAndFromEthAddress mocks base method
func (m *MockWorkflowPaymentsDBInterface) GetByTxHashAndStatusAndFromEthAddress(txHash, status, from string) (*model.WorkflowPaymentItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByTxHashAndStatusAndFromEthAddress", txHash, status, from)
	ret0, _ := ret[0].(*model.WorkflowPaymentItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByTxHashAndStatusAndFromEthAddress indicates an expected call of GetByTxHashAndStatusAndFromEthAddress
func (mr *MockWorkflowPaymentsDBInterfaceMockRecorder) GetByTxHashAndStatusAndFromEthAddress(txHash, status, from interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByTxHashAndStatusAndFromEthAddress", reflect.TypeOf((*MockWorkflowPaymentsDBInterface)(nil).GetByTxHashAndStatusAndFromEthAddress), txHash, status, from)
}

// Get mocks base method
func (m *MockWorkflowPaymentsDBInterface) Get(paymentId string) (*model.WorkflowPaymentItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", paymentId)
	ret0, _ := ret[0].(*model.WorkflowPaymentItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockWorkflowPaymentsDBInterfaceMockRecorder) Get(paymentId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockWorkflowPaymentsDBInterface)(nil).Get), paymentId)
}

// ConfirmPayment mocks base method
func (m *MockWorkflowPaymentsDBInterface) ConfirmPayment(txHash, from, to string, xes uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConfirmPayment", txHash, from, to, xes)
	ret0, _ := ret[0].(error)
	return ret0
}

// ConfirmPayment indicates an expected call of ConfirmPayment
func (mr *MockWorkflowPaymentsDBInterfaceMockRecorder) ConfirmPayment(txHash, from, to, xes interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConfirmPayment", reflect.TypeOf((*MockWorkflowPaymentsDBInterface)(nil).ConfirmPayment), txHash, from, to, xes)
}

// GetByWorkflowIdAndFromEthAddress mocks base method
func (m *MockWorkflowPaymentsDBInterface) GetByWorkflowIdAndFromEthAddress(workflowID, fromEthAddr string, statuses []string) (*model.WorkflowPaymentItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByWorkflowIdAndFromEthAddress", workflowID, fromEthAddr, statuses)
	ret0, _ := ret[0].(*model.WorkflowPaymentItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByWorkflowIdAndFromEthAddress indicates an expected call of GetByWorkflowIdAndFromEthAddress
func (mr *MockWorkflowPaymentsDBInterfaceMockRecorder) GetByWorkflowIdAndFromEthAddress(workflowID, fromEthAddr, statuses interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByWorkflowIdAndFromEthAddress", reflect.TypeOf((*MockWorkflowPaymentsDBInterface)(nil).GetByWorkflowIdAndFromEthAddress), workflowID, fromEthAddr, statuses)
}

// SetAbandonedToTimeoutBeforeTime mocks base method
func (m *MockWorkflowPaymentsDBInterface) SetAbandonedToTimeoutBeforeTime(beforeTime time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetAbandonedToTimeoutBeforeTime", beforeTime)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetAbandonedToTimeoutBeforeTime indicates an expected call of SetAbandonedToTimeoutBeforeTime
func (mr *MockWorkflowPaymentsDBInterfaceMockRecorder) SetAbandonedToTimeoutBeforeTime(beforeTime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetAbandonedToTimeoutBeforeTime", reflect.TypeOf((*MockWorkflowPaymentsDBInterface)(nil).SetAbandonedToTimeoutBeforeTime), beforeTime)
}

// Save mocks base method
func (m *MockWorkflowPaymentsDBInterface) Save(item *model.WorkflowPaymentItem) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", item)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save
func (mr *MockWorkflowPaymentsDBInterfaceMockRecorder) Save(item interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockWorkflowPaymentsDBInterface)(nil).Save), item)
}

// Update mocks base method
func (m *MockWorkflowPaymentsDBInterface) Update(paymentId, status, txHash, from string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", paymentId, status, txHash, from)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockWorkflowPaymentsDBInterfaceMockRecorder) Update(paymentId, status, txHash, from interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockWorkflowPaymentsDBInterface)(nil).Update), paymentId, status, txHash, from)
}

// Cancel mocks base method
func (m *MockWorkflowPaymentsDBInterface) Cancel(paymentId, from string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Cancel", paymentId, from)
	ret0, _ := ret[0].(error)
	return ret0
}

// Cancel indicates an expected call of Cancel
func (mr *MockWorkflowPaymentsDBInterfaceMockRecorder) Cancel(paymentId, from interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cancel", reflect.TypeOf((*MockWorkflowPaymentsDBInterface)(nil).Cancel), paymentId, from)
}

// Redeem mocks base method
func (m *MockWorkflowPaymentsDBInterface) Redeem(workflowId, from string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Redeem", workflowId, from)
	ret0, _ := ret[0].(error)
	return ret0
}

// Redeem indicates an expected call of Redeem
func (mr *MockWorkflowPaymentsDBInterfaceMockRecorder) Redeem(workflowId, from interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Redeem", reflect.TypeOf((*MockWorkflowPaymentsDBInterface)(nil).Redeem), workflowId, from)
}

// Delete mocks base method
func (m *MockWorkflowPaymentsDBInterface) Delete(paymentId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", paymentId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockWorkflowPaymentsDBInterfaceMockRecorder) Delete(paymentId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockWorkflowPaymentsDBInterface)(nil).Delete), paymentId)
}

// Remove mocks base method
func (m *MockWorkflowPaymentsDBInterface) Remove(payment *model.WorkflowPaymentItem) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", payment)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove
func (mr *MockWorkflowPaymentsDBInterfaceMockRecorder) Remove(payment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockWorkflowPaymentsDBInterface)(nil).Remove), payment)
}

// All mocks base method
func (m *MockWorkflowPaymentsDBInterface) All() ([]*model.WorkflowPaymentItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "All")
	ret0, _ := ret[0].([]*model.WorkflowPaymentItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// All indicates an expected call of All
func (mr *MockWorkflowPaymentsDBInterfaceMockRecorder) All() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "All", reflect.TypeOf((*MockWorkflowPaymentsDBInterface)(nil).All))
}

// Close mocks base method
func (m *MockWorkflowPaymentsDBInterface) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockWorkflowPaymentsDBInterfaceMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockWorkflowPaymentsDBInterface)(nil).Close))
}
