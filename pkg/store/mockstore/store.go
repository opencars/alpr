// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/opencars/alpr/pkg/domain (interfaces: Store,RecognitionRepository)

// Package mockstore is a generated GoMock package.
package mockstore

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/opencars/alpr/pkg/domain"
	model "github.com/opencars/alpr/pkg/domain/model"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// Recognition mocks base method.
func (m *MockStore) Recognition() domain.RecognitionRepository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recognition")
	ret0, _ := ret[0].(domain.RecognitionRepository)
	return ret0
}

// Recognition indicates an expected call of Recognition.
func (mr *MockStoreMockRecorder) Recognition() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recognition", reflect.TypeOf((*MockStore)(nil).Recognition))
}

// MockRecognitionRepository is a mock of RecognitionRepository interface.
type MockRecognitionRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRecognitionRepositoryMockRecorder
}

// MockRecognitionRepositoryMockRecorder is the mock recorder for MockRecognitionRepository.
type MockRecognitionRepositoryMockRecorder struct {
	mock *MockRecognitionRepository
}

// NewMockRecognitionRepository creates a new mock instance.
func NewMockRecognitionRepository(ctrl *gomock.Controller) *MockRecognitionRepository {
	mock := &MockRecognitionRepository{ctrl: ctrl}
	mock.recorder = &MockRecognitionRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRecognitionRepository) EXPECT() *MockRecognitionRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockRecognitionRepository) Create(arg0 context.Context, arg1 *model.Recognition) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockRecognitionRepositoryMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRecognitionRepository)(nil).Create), arg0, arg1)
}

// FindByPlate mocks base method.
func (m *MockRecognitionRepository) FindByPlate(arg0 context.Context, arg1 string) ([]model.Recognition, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByPlate", arg0, arg1)
	ret0, _ := ret[0].([]model.Recognition)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByPlate indicates an expected call of FindByPlate.
func (mr *MockRecognitionRepositoryMockRecorder) FindByPlate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByPlate", reflect.TypeOf((*MockRecognitionRepository)(nil).FindByPlate), arg0, arg1)
}
