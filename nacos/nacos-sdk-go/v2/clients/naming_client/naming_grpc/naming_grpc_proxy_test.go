package naming_grpc

import (
	model2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/model"
)

type MockNamingGrpc struct {
}

func (m *MockNamingGrpc) RegisterInstance(serviceName string, groupName string, instance model2.Instance) (bool, error) {
	return true, nil
}

func (m *MockNamingGrpc) DeregisterInstance(serviceName string, groupName string, instance model2.Instance) (bool, error) {
	return true, nil
}

func (m *MockNamingGrpc) GetServiceList(pageNo uint32, pageSize uint32, groupName string, selector *model2.ExpressionSelector) (model2.ServiceList, error) {
	return model2.ServiceList{Doms: []string{""}}, nil
}

func (m *MockNamingGrpc) ServerHealthy() bool {
	return true
}

func (m *MockNamingGrpc) QueryInstancesOfService(serviceName, groupName, clusters string, udpPort int, healthyOnly bool) (*model2.Service, error) {
	return &model2.Service{}, nil
}

func (m *MockNamingGrpc) Subscribe(serviceName, groupName, clusters string) (model2.Service, error) {
	return model2.Service{}, nil
}

func (m *MockNamingGrpc) Unsubscribe(serviceName, groupName, clusters string) {}
