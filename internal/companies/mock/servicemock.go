package mock

import (
	"companies-service/internal/models"
	"context"
	"reflect"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

// MockService is a mock of Service interface
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceRecorder
}

// MockServiceRecorder is the mock recorder for MockService
type MockServiceRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockService) EXPECT() *MockServiceRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockService) Create(
	ctx context.Context, company *models.Company,
) (*models.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, company)
	ret0, _ := ret[0].(*models.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockServiceRecorder) Create(ctx, company interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create",
		reflect.TypeOf((*MockService)(nil).Create), ctx, company)
}

// Update mocks base method
func (m *MockService) Update(
	ctx context.Context, company *models.Company,
) (*models.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, company)
	ret0, _ := ret[0].(*models.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update
func (mr *MockServiceRecorder) Update(ctx, company interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update",
		reflect.TypeOf((*MockService)(nil).Update), ctx, company)
}

// Delete mocks base method
func (m *MockService) Delete(ctx context.Context, companyID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, companyID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockServiceRecorder) Delete(ctx, companyID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete",
		reflect.TypeOf((*MockService)(nil).Delete), ctx, companyID)
}

// GetByID mocks base method
func (m *MockService) GetByID(ctx context.Context, companyID uuid.UUID) (*models.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, companyID)
	ret0, _ := ret[0].(*models.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID
func (mr *MockServiceRecorder) GetByID(ctx, companyID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID",
		reflect.TypeOf((*MockService)(nil).GetByID), ctx, companyID)
}
