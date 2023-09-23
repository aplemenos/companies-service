package mock

import (
	"context"
	"reflect"

	"github.com/golang/mock/gomock"
	"github.com/segmentio/kafka-go"
)

// MockKafka is a mock of kafka interface
type MockKafka struct {
	ctrl     *gomock.Controller
	recorder *MockKafkaRecorder
}

// MockKafkaRecorder is the recorder for MockKafka
type MockKafkaRecorder struct {
	mock *MockKafka
}

// NewMockKafka creates a new mock instance
func NewMockKafka(ctrl *gomock.Controller) *MockKafka {
	mock := &MockKafka{ctrl: ctrl}
	mock.recorder = &MockKafkaRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockKafka) EXPECT() *MockKafkaRecorder {
	return m.recorder
}

// PublishMessage mocks base method
func (m *MockKafka) PublishMessage(
	ctx context.Context, msgs ...kafka.Message,
) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublishMessage", ctx, msgs)
	ret0, _ := ret[0].(error)
	return ret0
}

// PublishMessage indicates an expected call of PublishMessage
func (mr *MockKafkaRecorder) PublishMessage(ctx, msgs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishMessage",
		reflect.TypeOf((*MockKafka)(nil).PublishMessage), ctx, msgs)
}

// Close mocks base method
func (m *MockKafka) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of CLose
func (mr *MockKafkaRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close",
		reflect.TypeOf((*MockKafka)(nil).Close))
}
