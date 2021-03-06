// Code generated by MockGen. DO NOT EDIT.
// Source: ./state/factory.go

// Package mock_state is a generated GoMock package.
package mock_state

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	action "github.com/iotexproject/iotex-core/blockchain/action"
	hash "github.com/iotexproject/iotex-core/pkg/hash"
	state "github.com/iotexproject/iotex-core/state"
	big "math/big"
	reflect "reflect"
)

// MockFactory is a mock of Factory interface
type MockFactory struct {
	ctrl     *gomock.Controller
	recorder *MockFactoryMockRecorder
}

// MockFactoryMockRecorder is the mock recorder for MockFactory
type MockFactoryMockRecorder struct {
	mock *MockFactory
}

// NewMockFactory creates a new mock instance
func NewMockFactory(ctrl *gomock.Controller) *MockFactory {
	mock := &MockFactory{ctrl: ctrl}
	mock.recorder = &MockFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFactory) EXPECT() *MockFactoryMockRecorder {
	return m.recorder
}

// Start mocks base method
func (m *MockFactory) Start(arg0 context.Context) error {
	ret := m.ctrl.Call(m, "Start", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start
func (mr *MockFactoryMockRecorder) Start(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockFactory)(nil).Start), arg0)
}

// Stop mocks base method
func (m *MockFactory) Stop(arg0 context.Context) error {
	ret := m.ctrl.Call(m, "Stop", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop
func (mr *MockFactoryMockRecorder) Stop(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockFactory)(nil).Stop), arg0)
}

// LoadOrCreateState mocks base method
func (m *MockFactory) LoadOrCreateState(arg0 string, arg1 uint64) (*state.State, error) {
	ret := m.ctrl.Call(m, "LoadOrCreateState", arg0, arg1)
	ret0, _ := ret[0].(*state.State)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadOrCreateState indicates an expected call of LoadOrCreateState
func (mr *MockFactoryMockRecorder) LoadOrCreateState(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadOrCreateState", reflect.TypeOf((*MockFactory)(nil).LoadOrCreateState), arg0, arg1)
}

// Balance mocks base method
func (m *MockFactory) Balance(arg0 string) (*big.Int, error) {
	ret := m.ctrl.Call(m, "Balance", arg0)
	ret0, _ := ret[0].(*big.Int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Balance indicates an expected call of Balance
func (mr *MockFactoryMockRecorder) Balance(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Balance", reflect.TypeOf((*MockFactory)(nil).Balance), arg0)
}

// Nonce mocks base method
func (m *MockFactory) Nonce(arg0 string) (uint64, error) {
	ret := m.ctrl.Call(m, "Nonce", arg0)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Nonce indicates an expected call of Nonce
func (mr *MockFactoryMockRecorder) Nonce(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Nonce", reflect.TypeOf((*MockFactory)(nil).Nonce), arg0)
}

// State mocks base method
func (m *MockFactory) State(arg0 string) (*state.State, error) {
	ret := m.ctrl.Call(m, "State", arg0)
	ret0, _ := ret[0].(*state.State)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// State indicates an expected call of State
func (mr *MockFactoryMockRecorder) State(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "State", reflect.TypeOf((*MockFactory)(nil).State), arg0)
}

// CachedState mocks base method
func (m *MockFactory) CachedState(arg0 string) (*state.State, error) {
	ret := m.ctrl.Call(m, "CachedState", arg0)
	ret0, _ := ret[0].(*state.State)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CachedState indicates an expected call of CachedState
func (mr *MockFactoryMockRecorder) CachedState(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CachedState", reflect.TypeOf((*MockFactory)(nil).CachedState), arg0)
}

// RootHash mocks base method
func (m *MockFactory) RootHash() hash.Hash32B {
	ret := m.ctrl.Call(m, "RootHash")
	ret0, _ := ret[0].(hash.Hash32B)
	return ret0
}

// RootHash indicates an expected call of RootHash
func (mr *MockFactoryMockRecorder) RootHash() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RootHash", reflect.TypeOf((*MockFactory)(nil).RootHash))
}

// Height mocks base method
func (m *MockFactory) Height() (uint64, error) {
	ret := m.ctrl.Call(m, "Height")
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Height indicates an expected call of Height
func (mr *MockFactoryMockRecorder) Height() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Height", reflect.TypeOf((*MockFactory)(nil).Height))
}

// NewWorkingSet mocks base method
func (m *MockFactory) NewWorkingSet() (state.WorkingSet, error) {
	ret := m.ctrl.Call(m, "NewWorkingSet")
	ret0, _ := ret[0].(state.WorkingSet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewWorkingSet indicates an expected call of NewWorkingSet
func (mr *MockFactoryMockRecorder) NewWorkingSet() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewWorkingSet", reflect.TypeOf((*MockFactory)(nil).NewWorkingSet))
}

// RunActions mocks base method
func (m *MockFactory) RunActions(arg0 uint64, arg1 []*action.Transfer, arg2 []*action.Vote, arg3 []*action.Execution, arg4 []action.Action) (hash.Hash32B, error) {
	ret := m.ctrl.Call(m, "RunActions", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(hash.Hash32B)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RunActions indicates an expected call of RunActions
func (mr *MockFactoryMockRecorder) RunActions(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunActions", reflect.TypeOf((*MockFactory)(nil).RunActions), arg0, arg1, arg2, arg3, arg4)
}

// Commit mocks base method
func (m *MockFactory) Commit(arg0 state.WorkingSet) error {
	ret := m.ctrl.Call(m, "Commit", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Commit indicates an expected call of Commit
func (mr *MockFactoryMockRecorder) Commit(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*MockFactory)(nil).Commit), arg0)
}

// GetCodeHash mocks base method
func (m *MockFactory) GetCodeHash(arg0 hash.PKHash) (hash.Hash32B, error) {
	ret := m.ctrl.Call(m, "GetCodeHash", arg0)
	ret0, _ := ret[0].(hash.Hash32B)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCodeHash indicates an expected call of GetCodeHash
func (mr *MockFactoryMockRecorder) GetCodeHash(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCodeHash", reflect.TypeOf((*MockFactory)(nil).GetCodeHash), arg0)
}

// GetCode mocks base method
func (m *MockFactory) GetCode(arg0 hash.PKHash) ([]byte, error) {
	ret := m.ctrl.Call(m, "GetCode", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCode indicates an expected call of GetCode
func (mr *MockFactoryMockRecorder) GetCode(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCode", reflect.TypeOf((*MockFactory)(nil).GetCode), arg0)
}

// SetCode mocks base method
func (m *MockFactory) SetCode(arg0 hash.PKHash, arg1 []byte) error {
	ret := m.ctrl.Call(m, "SetCode", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetCode indicates an expected call of SetCode
func (mr *MockFactoryMockRecorder) SetCode(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCode", reflect.TypeOf((*MockFactory)(nil).SetCode), arg0, arg1)
}

// GetContractState mocks base method
func (m *MockFactory) GetContractState(arg0 hash.PKHash, arg1 hash.Hash32B) (hash.Hash32B, error) {
	ret := m.ctrl.Call(m, "GetContractState", arg0, arg1)
	ret0, _ := ret[0].(hash.Hash32B)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContractState indicates an expected call of GetContractState
func (mr *MockFactoryMockRecorder) GetContractState(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContractState", reflect.TypeOf((*MockFactory)(nil).GetContractState), arg0, arg1)
}

// SetContractState mocks base method
func (m *MockFactory) SetContractState(arg0 hash.PKHash, arg1, arg2 hash.Hash32B) error {
	ret := m.ctrl.Call(m, "SetContractState", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetContractState indicates an expected call of SetContractState
func (mr *MockFactoryMockRecorder) SetContractState(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetContractState", reflect.TypeOf((*MockFactory)(nil).SetContractState), arg0, arg1, arg2)
}

// candidates mocks base method
func (m *MockFactory) candidates() (uint64, []*state.Candidate) {
	ret := m.ctrl.Call(m, "candidates")
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].([]*state.Candidate)
	return ret0, ret1
}

// candidates indicates an expected call of candidates
func (mr *MockFactoryMockRecorder) candidates() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "candidates", reflect.TypeOf((*MockFactory)(nil).candidates))
}

// CandidatesByHeight mocks base method
func (m *MockFactory) CandidatesByHeight(arg0 uint64) ([]*state.Candidate, error) {
	ret := m.ctrl.Call(m, "CandidatesByHeight", arg0)
	ret0, _ := ret[0].([]*state.Candidate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CandidatesByHeight indicates an expected call of CandidatesByHeight
func (mr *MockFactoryMockRecorder) CandidatesByHeight(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CandidatesByHeight", reflect.TypeOf((*MockFactory)(nil).CandidatesByHeight), arg0)
}

// MockActionHandler is a mock of ActionHandler interface
type MockActionHandler struct {
	ctrl     *gomock.Controller
	recorder *MockActionHandlerMockRecorder
}

// MockActionHandlerMockRecorder is the mock recorder for MockActionHandler
type MockActionHandlerMockRecorder struct {
	mock *MockActionHandler
}

// NewMockActionHandler creates a new mock instance
func NewMockActionHandler(ctrl *gomock.Controller) *MockActionHandler {
	mock := &MockActionHandler{ctrl: ctrl}
	mock.recorder = &MockActionHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockActionHandler) EXPECT() *MockActionHandlerMockRecorder {
	return m.recorder
}

// handle mocks base method
func (m *MockActionHandler) handle(arg0 action.Action) error {
	ret := m.ctrl.Call(m, "handle", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// handle indicates an expected call of handle
func (mr *MockActionHandlerMockRecorder) handle(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "handle", reflect.TypeOf((*MockActionHandler)(nil).handle), arg0)
}
