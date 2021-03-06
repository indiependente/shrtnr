// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package service is a generated GoMock package.
package service

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	models "github.com/indiependente/shrtnr/models"
	reflect "reflect"
)

// MockService is a mock of Service interface
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// Add mocks base method
func (m *MockService) Add(ctx context.Context, shortURL models.URLShortened) (models.URLShortened, error) {
	ret := m.ctrl.Call(m, "Add", ctx, shortURL)
	ret0, _ := ret[0].(models.URLShortened)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add
func (mr *MockServiceMockRecorder) Add(ctx, shortURL interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockService)(nil).Add), ctx, shortURL)
}

// Get mocks base method
func (m *MockService) Get(ctx context.Context, slug string) (models.URLShortened, error) {
	ret := m.ctrl.Call(m, "Get", ctx, slug)
	ret0, _ := ret[0].(models.URLShortened)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockServiceMockRecorder) Get(ctx, slug interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockService)(nil).Get), ctx, slug)
}

// Shorten mocks base method
func (m *MockService) Shorten(ctx context.Context, url string) (models.URLShortened, error) {
	ret := m.ctrl.Call(m, "Shorten", ctx, url)
	ret0, _ := ret[0].(models.URLShortened)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Shorten indicates an expected call of Shorten
func (mr *MockServiceMockRecorder) Shorten(ctx, url interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shorten", reflect.TypeOf((*MockService)(nil).Shorten), ctx, url)
}

// Delete mocks base method
func (m *MockService) Delete(ctx context.Context, slug string) error {
	ret := m.ctrl.Call(m, "Delete", ctx, slug)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockServiceMockRecorder) Delete(ctx, slug interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockService)(nil).Delete), ctx, slug)
}
