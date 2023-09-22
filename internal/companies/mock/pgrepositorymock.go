package mock

import (
	"companies-service/internal/models"
	"context"
	"reflect"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

// MockRepository is a mock of Repository interface
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryRecorder
}

// MockRepositoryRecorder is the recorder for MockRepository
type MockRepositoryRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepository) EXPECT() *MockRepositoryRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockRepository) Create(
	ctx context.Context, company *models.Company,
) (*models.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, company)
	ret0, _ := ret[0].(*models.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockRepositoryRecorder) Create(ctx, company interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create",
		reflect.TypeOf((*MockRepository)(nil).Create), ctx, company)
}

// Update mocks base method
func (m *MockRepository) Update(
	ctx context.Context, company *models.Company,
) (*models.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, company)
	ret0, _ := ret[0].(*models.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update
func (mr *MockRepositoryRecorder) Update(ctx, company interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update",
		reflect.TypeOf((*MockRepository)(nil).Update), ctx, company)
}

// Delete mocks base method
func (m *MockRepository) Delete(ctx context.Context, companyID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, companyID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockRepositoryRecorder) Delete(ctx, companyID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete",
		reflect.TypeOf((*MockRepository)(nil).Delete), ctx, companyID)
}

// GetByID mocks base method
func (m *MockRepository) GetByID(
	ctx context.Context, companyID uuid.UUID,
) (*models.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, companyID)
	ret0, _ := ret[0].(*models.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID
func (mr *MockRepositoryRecorder) GetByID(ctx, companyID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID",
		reflect.TypeOf((*MockRepository)(nil).GetByID), ctx, companyID)
}
