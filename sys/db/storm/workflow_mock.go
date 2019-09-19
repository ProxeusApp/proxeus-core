// Code generated by MockGen. DO NOT EDIT.
// Source: sys/db/storm/workflow.go

// Package storm is a generated GoMock package.
package storm

import (
	reflect "reflect"

	storm0 "github.com/asdine/storm"
	gomock "github.com/golang/mock/gomock"

	model "git.proxeus.com/core/central/sys/model"
)

// MockWorkflowDBInterface is a mock of WorkflowDBInterface interface
type MockWorkflowDBInterface struct {
	ctrl     *gomock.Controller
	recorder *MockWorkflowDBInterfaceMockRecorder
}

// MockWorkflowDBInterfaceMockRecorder is the mock recorder for MockWorkflowDBInterface
type MockWorkflowDBInterfaceMockRecorder struct {
	mock *MockWorkflowDBInterface
}

// NewMockWorkflowDBInterface creates a new mock instance
func NewMockWorkflowDBInterface(ctrl *gomock.Controller) *MockWorkflowDBInterface {
	mock := &MockWorkflowDBInterface{ctrl: ctrl}
	mock.recorder = &MockWorkflowDBInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockWorkflowDBInterface) EXPECT() *MockWorkflowDBInterfaceMockRecorder {
	return m.recorder
}

// ListPublished mocks base method
func (m *MockWorkflowDBInterface) ListPublished(auth model.Authorization, contains string, options map[string]interface{}) ([]*model.WorkflowItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPublished", auth, contains, options)
	ret0, _ := ret[0].([]*model.WorkflowItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListPublished indicates an expected call of ListPublished
func (mr *MockWorkflowDBInterfaceMockRecorder) ListPublished(auth, contains, options interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPublished", reflect.TypeOf((*MockWorkflowDBInterface)(nil).ListPublished), auth, contains, options)
}

// List mocks base method
func (m *MockWorkflowDBInterface) List(auth model.Authorization, contains string, options map[string]interface{}) ([]*model.WorkflowItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", auth, contains, options)
	ret0, _ := ret[0].([]*model.WorkflowItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockWorkflowDBInterfaceMockRecorder) List(auth, contains, options interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockWorkflowDBInterface)(nil).List), auth, contains, options)
}

// GetPublished mocks base method
func (m *MockWorkflowDBInterface) GetPublished(auth model.Authorization, id string) (*model.WorkflowItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPublished", auth, id)
	ret0, _ := ret[0].(*model.WorkflowItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPublished indicates an expected call of GetPublished
func (mr *MockWorkflowDBInterfaceMockRecorder) GetPublished(auth, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPublished", reflect.TypeOf((*MockWorkflowDBInterface)(nil).GetPublished), auth, id)
}

// Get mocks base method
func (m *MockWorkflowDBInterface) Get(auth model.Authorization, id string) (*model.WorkflowItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", auth, id)
	ret0, _ := ret[0].(*model.WorkflowItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockWorkflowDBInterfaceMockRecorder) Get(auth, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockWorkflowDBInterface)(nil).Get), auth, id)
}

// GetList mocks base method
func (m *MockWorkflowDBInterface) GetList(auth model.Authorization, ids []string) ([]*model.WorkflowItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetList", auth, ids)
	ret0, _ := ret[0].([]*model.WorkflowItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetList indicates an expected call of GetList
func (mr *MockWorkflowDBInterfaceMockRecorder) GetList(auth, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetList", reflect.TypeOf((*MockWorkflowDBInterface)(nil).GetList), auth, ids)
}

// Put mocks base method
func (m *MockWorkflowDBInterface) Put(auth model.Authorization, item *model.WorkflowItem) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", auth, item)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put
func (mr *MockWorkflowDBInterfaceMockRecorder) Put(auth, item interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockWorkflowDBInterface)(nil).Put), auth, item)
}

// put mocks base method
func (m *MockWorkflowDBInterface) put(auth model.Authorization, item *model.WorkflowItem, updated bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "put", auth, item, updated)
	ret0, _ := ret[0].(error)
	return ret0
}

// put indicates an expected call of put
func (mr *MockWorkflowDBInterfaceMockRecorder) put(auth, item, updated interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "put", reflect.TypeOf((*MockWorkflowDBInterface)(nil).put), auth, item, updated)
}

// getDB mocks base method
func (m *MockWorkflowDBInterface) getDB() *storm0.DB {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getDB")
	ret0, _ := ret[0].(*storm0.DB)
	return ret0
}

// getDB indicates an expected call of getDB
func (mr *MockWorkflowDBInterfaceMockRecorder) getDB() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getDB", reflect.TypeOf((*MockWorkflowDBInterface)(nil).getDB))
}

// updateWF mocks base method
func (m *MockWorkflowDBInterface) updateWF(auth model.Authorization, item *model.WorkflowItem, tx storm0.Node) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "updateWF", auth, item, tx)
	ret0, _ := ret[0].(error)
	return ret0
}

// updateWF indicates an expected call of updateWF
func (mr *MockWorkflowDBInterfaceMockRecorder) updateWF(auth, item, tx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "updateWF", reflect.TypeOf((*MockWorkflowDBInterface)(nil).updateWF), auth, item, tx)
}

// Delete mocks base method
func (m *MockWorkflowDBInterface) Delete(auth model.Authorization, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", auth, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockWorkflowDBInterfaceMockRecorder) Delete(auth, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockWorkflowDBInterface)(nil).Delete), auth, id)
}

// saveOnly mocks base method
func (m *MockWorkflowDBInterface) saveOnly(item *model.WorkflowItem, tx storm0.Node) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "saveOnly", item, tx)
	ret0, _ := ret[0].(error)
	return ret0
}

// saveOnly indicates an expected call of saveOnly
func (mr *MockWorkflowDBInterfaceMockRecorder) saveOnly(item, tx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "saveOnly", reflect.TypeOf((*MockWorkflowDBInterface)(nil).saveOnly), item, tx)
}

// Import mocks base method
func (m *MockWorkflowDBInterface) Import(imex *Imex) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Import", imex)
	ret0, _ := ret[0].(error)
	return ret0
}

// Import indicates an expected call of Import
func (mr *MockWorkflowDBInterfaceMockRecorder) Import(imex interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Import", reflect.TypeOf((*MockWorkflowDBInterface)(nil).Import), imex)
}

// Export mocks base method
func (m *MockWorkflowDBInterface) Export(imex *Imex, id ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{imex}
	for _, a := range id {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Export", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Export indicates an expected call of Export
func (mr *MockWorkflowDBInterfaceMockRecorder) Export(imex interface{}, id ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{imex}, id...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Export", reflect.TypeOf((*MockWorkflowDBInterface)(nil).Export), varargs...)
}

// Close mocks base method
func (m *MockWorkflowDBInterface) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockWorkflowDBInterfaceMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockWorkflowDBInterface)(nil).Close))
}
