// Code generated by MockGen. DO NOT EDIT.
// Source: clients/naming_client/naming_proxy/proxy_interface.go

// Package naming_proxy is a generated GoMock package.
package naming_proxy

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/model"
)

// MockINamingProxy is a mock of INamingProxy interface.
type MockINamingProxy struct {
	ctrl     *gomock.Controller
	recorder *MockINamingProxyMockRecorder
}

// MockINamingProxyMockRecorder is the mock recorder for MockINamingProxy.
type MockINamingProxyMockRecorder struct {
	mock *MockINamingProxy
}

// NewMockINamingProxy creates a new mock instance.
func NewMockINamingProxy(ctrl *gomock.Controller) *MockINamingProxy {
	mock := &MockINamingProxy{ctrl: ctrl}
	mock.recorder = &MockINamingProxyMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockINamingProxy) EXPECT() *MockINamingProxyMockRecorder {
	return m.recorder
}

// BatchRegisterInstance mocks base method.
func (m *MockINamingProxy) BatchRegisterInstance(serviceName, groupName string, instances []model.Instance) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BatchRegisterInstance", serviceName, groupName, instances)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BatchRegisterInstance indicates an expected call of BatchRegisterInstance.
func (mr *MockINamingProxyMockRecorder) BatchRegisterInstance(serviceName, groupName, instances interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BatchRegisterInstance", reflect.TypeOf((*MockINamingProxy)(nil).BatchRegisterInstance), serviceName, groupName, instances)
}

// CloseClient mocks base method.
func (m *MockINamingProxy) CloseClient() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CloseClient")
}

// CloseClient indicates an expected call of CloseClient.
func (mr *MockINamingProxyMockRecorder) CloseClient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseClient", reflect.TypeOf((*MockINamingProxy)(nil).CloseClient))
}

// DeregisterInstance mocks base method.
func (m *MockINamingProxy) DeregisterInstance(serviceName, groupName string, instance model.Instance) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeregisterInstance", serviceName, groupName, instance)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeregisterInstance indicates an expected call of DeregisterInstance.
func (mr *MockINamingProxyMockRecorder) DeregisterInstance(serviceName, groupName, instance interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeregisterInstance", reflect.TypeOf((*MockINamingProxy)(nil).DeregisterInstance), serviceName, groupName, instance)
}

// GetServiceList mocks base method.
func (m *MockINamingProxy) GetServiceList(pageNo, pageSize uint32, groupName, namespaceId string, selector *model.ExpressionSelector) (model.ServiceList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetServiceList", pageNo, pageSize, groupName, namespaceId, selector)
	ret0, _ := ret[0].(model.ServiceList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetServiceList indicates an expected call of GetServiceList.
func (mr *MockINamingProxyMockRecorder) GetServiceList(pageNo, pageSize, groupName, namespaceId, selector interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetServiceList", reflect.TypeOf((*MockINamingProxy)(nil).GetServiceList), pageNo, pageSize, groupName, namespaceId, selector)
}

// QueryInstancesOfService mocks base method.
func (m *MockINamingProxy) QueryInstancesOfService(serviceName, groupName, clusters string, udpPort int, healthyOnly bool) (*model.Service, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryInstancesOfService", serviceName, groupName, clusters, udpPort, healthyOnly)
	ret0, _ := ret[0].(*model.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryInstancesOfService indicates an expected call of QueryInstancesOfService.
func (mr *MockINamingProxyMockRecorder) QueryInstancesOfService(serviceName, groupName, clusters, udpPort, healthyOnly interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryInstancesOfService", reflect.TypeOf((*MockINamingProxy)(nil).QueryInstancesOfService), serviceName, groupName, clusters, udpPort, healthyOnly)
}

// RegisterInstance mocks base method.
func (m *MockINamingProxy) RegisterInstance(serviceName, groupName string, instance model.Instance) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterInstance", serviceName, groupName, instance)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterInstance indicates an expected call of RegisterInstance.
func (mr *MockINamingProxyMockRecorder) RegisterInstance(serviceName, groupName, instance interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterInstance", reflect.TypeOf((*MockINamingProxy)(nil).RegisterInstance), serviceName, groupName, instance)
}

// ServerHealthy mocks base method.
func (m *MockINamingProxy) ServerHealthy() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ServerHealthy")
	ret0, _ := ret[0].(bool)
	return ret0
}

// ServerHealthy indicates an expected call of ServerHealthy.
func (mr *MockINamingProxyMockRecorder) ServerHealthy() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ServerHealthy", reflect.TypeOf((*MockINamingProxy)(nil).ServerHealthy))
}

// Subscribe mocks base method.
func (m *MockINamingProxy) Subscribe(serviceName, groupName, clusters string) (model.Service, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe", serviceName, groupName, clusters)
	ret0, _ := ret[0].(model.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockINamingProxyMockRecorder) Subscribe(serviceName, groupName, clusters interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockINamingProxy)(nil).Subscribe), serviceName, groupName, clusters)
}

// Unsubscribe mocks base method.
func (m *MockINamingProxy) Unsubscribe(serviceName, groupName, clusters string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unsubscribe", serviceName, groupName, clusters)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unsubscribe indicates an expected call of Unsubscribe.
func (mr *MockINamingProxyMockRecorder) Unsubscribe(serviceName, groupName, clusters interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unsubscribe", reflect.TypeOf((*MockINamingProxy)(nil).Unsubscribe), serviceName, groupName, clusters)
}
