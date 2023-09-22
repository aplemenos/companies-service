package mock

import (
	"companies-service/internal/models"
	"context"
	"reflect"

	"github.com/golang/mock/gomock"
)

// MockRedisRepository is a mock of RedisRepository interface
type MockRedisRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRedisRepositoryRecorder
}

// MockRedisRepositoryRecorder is the recorder for MockRedisRepository
type MockRedisRepositoryRecorder struct {
	mock *MockRedisRepository
}

// NewMockRedisRepository creates a new mock instance
func NewMockRedisRepository(ctrl *gomock.Controller) *MockRedisRepository {
	mock := &MockRedisRepository{ctrl: ctrl}
	mock.recorder = &MockRedisRepositoryRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRedisRepository) EXPECT() *MockRedisRepositoryRecorder {
	return m.recorder
}

// GetByIDCtx mocks base method
func (m *MockRedisRepository) GetByIDCtx(ctx context.Context, key string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByIDCtx", ctx, key)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByIDCtx indicates an expected call of GetByIDCtx
func (mr *MockRedisRepositoryRecorder) GetByIDCtx(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock,
		"GetByIDCtx", reflect.TypeOf((*MockRedisRepository)(nil).GetByIDCtx), ctx, key)
}

// SetUserCtx mocks base method
func (m *MockRedisRepository) SetUserCtx(
	ctx context.Context, key string, seconds int, user *models.User,
) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetUserCtx", ctx, key, seconds, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetUserCtx indicates an expected call of SetUserCtx
func (mr *MockRedisRepositoryRecorder) SetUserCtx(
	ctx, key, seconds, user interface{},
) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock,
		"SetUserCtx", reflect.TypeOf((*MockRedisRepository)(nil).SetUserCtx),
		ctx, key, seconds, user)
}

// DeleteUserCtx mocks base method
func (m *MockRedisRepository) DeleteUserCtx(ctx context.Context, key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUserCtx", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUserCtx indicates an expected call of DeleteUserCtx
func (mr *MockRedisRepositoryRecorder) DeleteUserCtx(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock,
		"DeleteUserCtx", reflect.TypeOf((*MockRedisRepository)(nil).DeleteUserCtx), ctx, key)
}
