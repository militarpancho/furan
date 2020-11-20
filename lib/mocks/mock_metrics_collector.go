// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/dollarshaveclub/furan/lib/metrics (interfaces: MetricsCollector)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockMetricsCollector is a mock of MetricsCollector interface
type MockMetricsCollector struct {
	ctrl     *gomock.Controller
	recorder *MockMetricsCollectorMockRecorder
}

// MockMetricsCollectorMockRecorder is the mock recorder for MockMetricsCollector
type MockMetricsCollectorMockRecorder struct {
	mock *MockMetricsCollector
}

// NewMockMetricsCollector creates a new mock instance
func NewMockMetricsCollector(ctrl *gomock.Controller) *MockMetricsCollector {
	mock := &MockMetricsCollector{ctrl: ctrl}
	mock.recorder = &MockMetricsCollectorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMetricsCollector) EXPECT() *MockMetricsCollectorMockRecorder {
	return m.recorder
}

// BuildFailed mocks base method
func (m *MockMetricsCollector) BuildFailed(arg0, arg1 string, arg2 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuildFailed", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// BuildFailed indicates an expected call of BuildFailed
func (mr *MockMetricsCollectorMockRecorder) BuildFailed(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuildFailed", reflect.TypeOf((*MockMetricsCollector)(nil).BuildFailed), arg0, arg1, arg2)
}

// BuildStarted mocks base method
func (m *MockMetricsCollector) BuildStarted(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuildStarted", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// BuildStarted indicates an expected call of BuildStarted
func (mr *MockMetricsCollectorMockRecorder) BuildStarted(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuildStarted", reflect.TypeOf((*MockMetricsCollector)(nil).BuildStarted), arg0, arg1)
}

// BuildSucceeded mocks base method
func (m *MockMetricsCollector) BuildSucceeded(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuildSucceeded", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// BuildSucceeded indicates an expected call of BuildSucceeded
func (mr *MockMetricsCollectorMockRecorder) BuildSucceeded(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuildSucceeded", reflect.TypeOf((*MockMetricsCollector)(nil).BuildSucceeded), arg0, arg1)
}

// DiskFree mocks base method
func (m *MockMetricsCollector) DiskFree(arg0 uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DiskFree", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DiskFree indicates an expected call of DiskFree
func (mr *MockMetricsCollectorMockRecorder) DiskFree(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DiskFree", reflect.TypeOf((*MockMetricsCollector)(nil).DiskFree), arg0)
}

// Duration mocks base method
func (m *MockMetricsCollector) Duration(arg0, arg1, arg2 string, arg3 []string, arg4 float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Duration", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// Duration indicates an expected call of Duration
func (mr *MockMetricsCollectorMockRecorder) Duration(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Duration", reflect.TypeOf((*MockMetricsCollector)(nil).Duration), arg0, arg1, arg2, arg3, arg4)
}

// FileNodesFree mocks base method
func (m *MockMetricsCollector) FileNodesFree(arg0 uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FileNodesFree", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// FileNodesFree indicates an expected call of FileNodesFree
func (mr *MockMetricsCollectorMockRecorder) FileNodesFree(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FileNodesFree", reflect.TypeOf((*MockMetricsCollector)(nil).FileNodesFree), arg0)
}

// Float mocks base method
func (m *MockMetricsCollector) Float(arg0, arg1, arg2 string, arg3 []string, arg4 float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Float", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// Float indicates an expected call of Float
func (mr *MockMetricsCollectorMockRecorder) Float(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Float", reflect.TypeOf((*MockMetricsCollector)(nil).Float), arg0, arg1, arg2, arg3, arg4)
}

// GCBytesReclaimed mocks base method
func (m *MockMetricsCollector) GCBytesReclaimed(arg0 uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GCBytesReclaimed", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// GCBytesReclaimed indicates an expected call of GCBytesReclaimed
func (mr *MockMetricsCollectorMockRecorder) GCBytesReclaimed(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GCBytesReclaimed", reflect.TypeOf((*MockMetricsCollector)(nil).GCBytesReclaimed), arg0)
}

// GCFailure mocks base method
func (m *MockMetricsCollector) GCFailure() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GCFailure")
	ret0, _ := ret[0].(error)
	return ret0
}

// GCFailure indicates an expected call of GCFailure
func (mr *MockMetricsCollectorMockRecorder) GCFailure() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GCFailure", reflect.TypeOf((*MockMetricsCollector)(nil).GCFailure))
}

// GCUntaggedImageRemoved mocks base method
func (m *MockMetricsCollector) GCUntaggedImageRemoved() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GCUntaggedImageRemoved")
	ret0, _ := ret[0].(error)
	return ret0
}

// GCUntaggedImageRemoved indicates an expected call of GCUntaggedImageRemoved
func (mr *MockMetricsCollectorMockRecorder) GCUntaggedImageRemoved() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GCUntaggedImageRemoved", reflect.TypeOf((*MockMetricsCollector)(nil).GCUntaggedImageRemoved))
}

// ImageSize mocks base method
func (m *MockMetricsCollector) ImageSize(arg0, arg1 int64, arg2, arg3 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ImageSize", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// ImageSize indicates an expected call of ImageSize
func (mr *MockMetricsCollectorMockRecorder) ImageSize(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImageSize", reflect.TypeOf((*MockMetricsCollector)(nil).ImageSize), arg0, arg1, arg2, arg3)
}

// KafkaConsumerFailure mocks base method
func (m *MockMetricsCollector) KafkaConsumerFailure() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "KafkaConsumerFailure")
	ret0, _ := ret[0].(error)
	return ret0
}

// KafkaConsumerFailure indicates an expected call of KafkaConsumerFailure
func (mr *MockMetricsCollectorMockRecorder) KafkaConsumerFailure() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KafkaConsumerFailure", reflect.TypeOf((*MockMetricsCollector)(nil).KafkaConsumerFailure))
}

// KafkaProducerFailure mocks base method
func (m *MockMetricsCollector) KafkaProducerFailure() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "KafkaProducerFailure")
	ret0, _ := ret[0].(error)
	return ret0
}

// KafkaProducerFailure indicates an expected call of KafkaProducerFailure
func (mr *MockMetricsCollectorMockRecorder) KafkaProducerFailure() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KafkaProducerFailure", reflect.TypeOf((*MockMetricsCollector)(nil).KafkaProducerFailure))
}

// Size mocks base method
func (m *MockMetricsCollector) Size(arg0, arg1, arg2 string, arg3 []string, arg4 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Size", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// Size indicates an expected call of Size
func (mr *MockMetricsCollectorMockRecorder) Size(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Size", reflect.TypeOf((*MockMetricsCollector)(nil).Size), arg0, arg1, arg2, arg3, arg4)
}

// TriggerCompleted mocks base method
func (m *MockMetricsCollector) TriggerCompleted(arg0, arg1 string, arg2 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TriggerCompleted", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// TriggerCompleted indicates an expected call of TriggerCompleted
func (mr *MockMetricsCollectorMockRecorder) TriggerCompleted(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TriggerCompleted", reflect.TypeOf((*MockMetricsCollector)(nil).TriggerCompleted), arg0, arg1, arg2)
}