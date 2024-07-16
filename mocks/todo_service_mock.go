// Code generated by MockGen. DO NOT EDIT.
// Source: ./services/todo_service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	Models "github.com/go-todo1/Models"
	gomock "github.com/golang/mock/gomock"
)

// MockTodoService is a mock of TodoService interface.
type MockTodoService struct {
	ctrl     *gomock.Controller
	recorder *MockTodoServiceMockRecorder
}

// MockTodoServiceMockRecorder is the mock recorder for MockTodoService.
type MockTodoServiceMockRecorder struct {
	mock *MockTodoService
}

// NewMockTodoService creates a new mock instance.
func NewMockTodoService(ctrl *gomock.Controller) *MockTodoService {
	mock := &MockTodoService{ctrl: ctrl}
	mock.recorder = &MockTodoServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTodoService) EXPECT() *MockTodoServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockTodoService) Create(todo Models.TodoModel) (Models.TodoModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", todo)
	ret0, _ := ret[0].(Models.TodoModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockTodoServiceMockRecorder) Create(todo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTodoService)(nil).Create), todo)
}

// Delete mocks base method.
func (m *MockTodoService) Delete(id uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockTodoServiceMockRecorder) Delete(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTodoService)(nil).Delete), id)
}

// Update mocks base method.
func (m *MockTodoService) Update(id uint, todo Models.TodoModel) (Models.TodoModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", id, todo)
	ret0, _ := ret[0].(Models.TodoModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockTodoServiceMockRecorder) Update(id, todo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockTodoService)(nil).Update), id, todo)
}
