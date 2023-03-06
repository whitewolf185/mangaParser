// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/whitewolf185/mangaparser/internal/pkg/pdf_creator (interfaces: ImageGetter)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"
	sync "sync"

	gomock "github.com/golang/mock/gomock"
	err_controller "github.com/whitewolf185/mangaparser/internal/pkg/err_controller"
)

// MockImageGetter is a mock of ImageGetter interface.
type MockImageGetter struct {
	ctrl     *gomock.Controller
	recorder *MockImageGetterMockRecorder
}

// MockImageGetterMockRecorder is the mock recorder for MockImageGetter.
type MockImageGetterMockRecorder struct {
	mock *MockImageGetter
}

// NewMockImageGetter creates a new mock instance.
func NewMockImageGetter(ctrl *gomock.Controller) *MockImageGetter {
	mock := &MockImageGetter{ctrl: ctrl}
	mock.recorder = &MockImageGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockImageGetter) EXPECT() *MockImageGetterMockRecorder {
	return m.recorder
}

// GetImageAndSave mocks base method.
func (m *MockImageGetter) GetImageAndSave(arg0 context.Context, arg1 *sync.WaitGroup, arg2 *err_controller.ErrController, arg3 string, arg4 int, arg5 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetImageAndSave", arg0, arg1, arg2, arg3, arg4, arg5)
}

// GetImageAndSave indicates an expected call of GetImageAndSave.
func (mr *MockImageGetterMockRecorder) GetImageAndSave(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetImageAndSave", reflect.TypeOf((*MockImageGetter)(nil).GetImageAndSave), arg0, arg1, arg2, arg3, arg4, arg5)
}