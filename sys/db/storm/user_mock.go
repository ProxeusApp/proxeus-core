// +build !coverage
// Code generated by MockGen. DO NOT EDIT.
// Source: sys/db/storm/user.go

// Package storm is a generated GoMock package.
package storm

import (
	io "io"
	os "os"
	reflect "reflect"

	model "github.com/ProxeusApp/proxeus-core/sys/model"
	storm0 "github.com/asdine/storm"
	gomock "github.com/golang/mock/gomock"
)

// MockUserDBInterface is a mock of UserDBInterface interface
type MockUserDBInterface struct {
	ctrl     *gomock.Controller
	recorder *MockUserDBInterfaceMockRecorder
}

// MockUserDBInterfaceMockRecorder is the mock recorder for MockUserDBInterface
type MockUserDBInterfaceMockRecorder struct {
	mock *MockUserDBInterface
}

// NewMockUserDBInterface creates a new mock instance
func NewMockUserDBInterface(ctrl *gomock.Controller) *MockUserDBInterface {
	mock := &MockUserDBInterface{ctrl: ctrl}
	mock.recorder = &MockUserDBInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUserDBInterface) EXPECT() *MockUserDBInterfaceMockRecorder {
	return m.recorder
}

// GetDB mocks base method
func (m *MockUserDBInterface) GetDB() *storm0.DB {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDB")
	ret0, _ := ret[0].(*storm0.DB)
	return ret0
}

// GetDB indicates an expected call of GetDB
func (mr *MockUserDBInterfaceMockRecorder) GetDB() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDB", reflect.TypeOf((*MockUserDBInterface)(nil).GetDB))
}

// GetBaseFilePath mocks base method
func (m *MockUserDBInterface) GetBaseFilePath() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBaseFilePath")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetBaseFilePath indicates an expected call of GetBaseFilePath
func (mr *MockUserDBInterfaceMockRecorder) GetBaseFilePath() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBaseFilePath", reflect.TypeOf((*MockUserDBInterface)(nil).GetBaseFilePath))
}

// Login mocks base method
func (m *MockUserDBInterface) Login(name, pw string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", name, pw)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login
func (mr *MockUserDBInterfaceMockRecorder) Login(name, pw interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUserDBInterface)(nil).Login), name, pw)
}

// Count mocks base method
func (m *MockUserDBInterface) Count() (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count
func (mr *MockUserDBInterfaceMockRecorder) Count() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockUserDBInterface)(nil).Count))
}

// List mocks base method
func (m *MockUserDBInterface) List(auth model.Authorization, contains string, options map[string]interface{}) ([]*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", auth, contains, options)
	ret0, _ := ret[0].([]*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockUserDBInterfaceMockRecorder) List(auth, contains, options interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockUserDBInterface)(nil).List), auth, contains, options)
}

// Get mocks base method
func (m *MockUserDBInterface) Get(auth model.Authorization, id string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", auth, id)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockUserDBInterfaceMockRecorder) Get(auth, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockUserDBInterface)(nil).Get), auth, id)
}

// GetByBCAddress mocks base method
func (m *MockUserDBInterface) GetByBCAddress(bcAddress string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByBCAddress", bcAddress)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByBCAddress indicates an expected call of GetByBCAddress
func (mr *MockUserDBInterfaceMockRecorder) GetByBCAddress(bcAddress interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByBCAddress", reflect.TypeOf((*MockUserDBInterface)(nil).GetByBCAddress), bcAddress)
}

// GetByEmail mocks base method
func (m *MockUserDBInterface) GetByEmail(email string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByEmail", email)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByEmail indicates an expected call of GetByEmail
func (mr *MockUserDBInterfaceMockRecorder) GetByEmail(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByEmail", reflect.TypeOf((*MockUserDBInterface)(nil).GetByEmail), email)
}

// UpdateEmail mocks base method
func (m *MockUserDBInterface) UpdateEmail(id, email string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEmail", id, email)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateEmail indicates an expected call of UpdateEmail
func (mr *MockUserDBInterfaceMockRecorder) UpdateEmail(id, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEmail", reflect.TypeOf((*MockUserDBInterface)(nil).UpdateEmail), id, email)
}

// Put mocks base method
func (m *MockUserDBInterface) Put(auth model.Authorization, item *model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", auth, item)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put
func (mr *MockUserDBInterfaceMockRecorder) Put(auth, item interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockUserDBInterface)(nil).Put), auth, item)
}

// PutPw mocks base method
func (m *MockUserDBInterface) PutPw(id, pass string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutPw", id, pass)
	ret0, _ := ret[0].(error)
	return ret0
}

// PutPw indicates an expected call of PutPw
func (mr *MockUserDBInterfaceMockRecorder) PutPw(id, pass interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutPw", reflect.TypeOf((*MockUserDBInterface)(nil).PutPw), id, pass)
}

// put mocks base method
func (m *MockUserDBInterface) put(auth model.Authorization, item *model.User, updated bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "put", auth, item, updated)
	ret0, _ := ret[0].(error)
	return ret0
}

// put indicates an expected call of put
func (mr *MockUserDBInterfaceMockRecorder) put(auth, item, updated interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "put", reflect.TypeOf((*MockUserDBInterface)(nil).put), auth, item, updated)
}

// setTinyUserIconBase64 mocks base method
func (m *MockUserDBInterface) setTinyUserIconBase64(item *model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "setTinyUserIconBase64", item)
	ret0, _ := ret[0].(error)
	return ret0
}

// setTinyUserIconBase64 indicates an expected call of setTinyUserIconBase64
func (mr *MockUserDBInterfaceMockRecorder) setTinyUserIconBase64(item interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "setTinyUserIconBase64", reflect.TypeOf((*MockUserDBInterface)(nil).setTinyUserIconBase64), item)
}

// tinyUserIconBase64 mocks base method
func (m *MockUserDBInterface) tinyUserIconBase64(reader *os.File) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "tinyUserIconBase64", reader)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// tinyUserIconBase64 indicates an expected call of tinyUserIconBase64
func (mr *MockUserDBInterfaceMockRecorder) tinyUserIconBase64(reader interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "tinyUserIconBase64", reflect.TypeOf((*MockUserDBInterface)(nil).tinyUserIconBase64), reader)
}

// GetProfilePhoto mocks base method
func (m *MockUserDBInterface) GetProfilePhoto(auth model.Authorization, id string, writer io.Writer) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProfilePhoto", auth, id, writer)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProfilePhoto indicates an expected call of GetProfilePhoto
func (mr *MockUserDBInterfaceMockRecorder) GetProfilePhoto(auth, id, writer interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProfilePhoto", reflect.TypeOf((*MockUserDBInterface)(nil).GetProfilePhoto), auth, id, writer)
}

// readPhoto mocks base method
func (m *MockUserDBInterface) readPhoto(u *model.User) (*os.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "readPhoto", u)
	ret0, _ := ret[0].(*os.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// readPhoto indicates an expected call of readPhoto
func (mr *MockUserDBInterfaceMockRecorder) readPhoto(u interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "readPhoto", reflect.TypeOf((*MockUserDBInterface)(nil).readPhoto), u)
}

// PutProfilePhoto mocks base method
func (m *MockUserDBInterface) PutProfilePhoto(auth model.Authorization, id string, reader io.Reader) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutProfilePhoto", auth, id, reader)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PutProfilePhoto indicates an expected call of PutProfilePhoto
func (mr *MockUserDBInterfaceMockRecorder) PutProfilePhoto(auth, id, reader interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutProfilePhoto", reflect.TypeOf((*MockUserDBInterface)(nil).PutProfilePhoto), auth, id, reader)
}

// Import mocks base method
func (m *MockUserDBInterface) Import(imex *Imex) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Import", imex)
	ret0, _ := ret[0].(error)
	return ret0
}

// Import indicates an expected call of Import
func (mr *MockUserDBInterfaceMockRecorder) Import(imex interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Import", reflect.TypeOf((*MockUserDBInterface)(nil).Import), imex)
}

// Export mocks base method
func (m *MockUserDBInterface) Export(imex *Imex, id ...string) error {
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
func (mr *MockUserDBInterfaceMockRecorder) Export(imex interface{}, id ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{imex}, id...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Export", reflect.TypeOf((*MockUserDBInterface)(nil).Export), varargs...)
}

// cpProfilePhoto mocks base method
func (m *MockUserDBInterface) cpProfilePhoto(imex *Imex, from, to UserDBInterface, item *model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "cpProfilePhoto", imex, from, to, item)
	ret0, _ := ret[0].(error)
	return ret0
}

// cpProfilePhoto indicates an expected call of cpProfilePhoto
func (mr *MockUserDBInterfaceMockRecorder) cpProfilePhoto(imex, from, to, item interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "cpProfilePhoto", reflect.TypeOf((*MockUserDBInterface)(nil).cpProfilePhoto), imex, from, to, item)
}

// APIKey mocks base method
func (m *MockUserDBInterface) APIKey(key string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "APIKey", key)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// APIKey indicates an expected call of APIKey
func (mr *MockUserDBInterfaceMockRecorder) APIKey(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "APIKey", reflect.TypeOf((*MockUserDBInterface)(nil).APIKey), key)
}

// CreateApiKey mocks base method
func (m *MockUserDBInterface) CreateApiKey(auth model.Authorization, userId, apiKeyName string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateApiKey", auth, userId, apiKeyName)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateApiKey indicates an expected call of CreateApiKey
func (mr *MockUserDBInterfaceMockRecorder) CreateApiKey(auth, userId, apiKeyName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateApiKey", reflect.TypeOf((*MockUserDBInterface)(nil).CreateApiKey), auth, userId, apiKeyName)
}

// DeleteApiKey mocks base method
func (m *MockUserDBInterface) DeleteApiKey(auth model.Authorization, userId, hiddenApiKey string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteApiKey", auth, userId, hiddenApiKey)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteApiKey indicates an expected call of DeleteApiKey
func (mr *MockUserDBInterfaceMockRecorder) DeleteApiKey(auth, userId, hiddenApiKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteApiKey", reflect.TypeOf((*MockUserDBInterface)(nil).DeleteApiKey), auth, userId, hiddenApiKey)
}

// Close mocks base method
func (m *MockUserDBInterface) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockUserDBInterfaceMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockUserDBInterface)(nil).Close))
}
