/*
 * Copyright 1999-2020 Alibaba Group Holding Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package naming_grpc

import (
	naming_cache2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/clients/naming_client/naming_cache"
	constant2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/constant"
	logger2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/logger"
	monitor2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/monitor"
	nacos_server2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/nacos_server"
	rpc2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/remote/rpc"
	rpc_request2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/remote/rpc/rpc_request"
	rpc_response2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/remote/rpc/rpc_response"
	uuid2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/inner/uuid"
	model2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/model"
	util2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/util"
	"time"
)

// NamingGrpcProxy ...
type NamingGrpcProxy struct {
	clientConfig      constant2.ClientConfig
	nacosServer       *nacos_server2.NacosServer
	rpcClient         rpc2.IRpcClient
	eventListener     *ConnectionEventListener
	serviceInfoHolder *naming_cache2.ServiceInfoHolder
}

// NewNamingGrpcProxy create naming grpc proxy
func NewNamingGrpcProxy(clientCfg constant2.ClientConfig, nacosServer *nacos_server2.NacosServer,
	serviceInfoHolder *naming_cache2.ServiceInfoHolder) (*NamingGrpcProxy, error) {
	srvProxy := NamingGrpcProxy{
		clientConfig:      clientCfg,
		nacosServer:       nacosServer,
		serviceInfoHolder: serviceInfoHolder,
	}

	uid, err := uuid2.NewV4()
	if err != nil {
		return nil, err
	}

	labels := map[string]string{
		constant2.LABEL_SOURCE: constant2.LABEL_SOURCE_SDK,
		constant2.LABEL_MODULE: constant2.LABEL_MODULE_NAMING,
	}

	iRpcClient, err := rpc2.CreateClient(uid.String(), rpc2.GRPC, labels, srvProxy.nacosServer)
	if err != nil {
		return nil, err
	}

	srvProxy.rpcClient = iRpcClient

	rpcClient := srvProxy.rpcClient.GetRpcClient()
	rpcClient.Start()

	rpcClient.RegisterServerRequestHandler(func() rpc_request2.IRequest {
		return &rpc_request2.NotifySubscriberRequest{NamingRequest: &rpc_request2.NamingRequest{}}
	}, &rpc2.NamingPushRequestHandler{ServiceInfoHolder: serviceInfoHolder})

	srvProxy.eventListener = NewConnectionEventListener(&srvProxy)
	rpcClient.RegisterConnectionListener(srvProxy.eventListener)

	return &srvProxy, nil
}

func (proxy *NamingGrpcProxy) requestToServer(request rpc_request2.IRequest) (rpc_response2.IResponse, error) {
	start := time.Now()
	proxy.nacosServer.InjectSkAk(request.GetHeaders(), proxy.clientConfig)
	proxy.nacosServer.InjectSecurityInfo(request.GetHeaders())
	response, err := proxy.rpcClient.GetRpcClient().Request(request, int64(proxy.clientConfig.TimeoutMs))
	monitor2.GetConfigRequestMonitor(constant2.GRPC, request.GetRequestType(), rpc_response2.GetGrpcResponseStatusCode(response)).Observe(float64(time.Now().Nanosecond() - start.Nanosecond()))
	return response, err
}

// RegisterInstance ...
func (proxy *NamingGrpcProxy) RegisterInstance(serviceName string, groupName string, instance model2.Instance) (bool, error) {
	logger2.Infof("instance namespaceId:<%s>,serviceName:<%s> with instance:<%s>",
		proxy.clientConfig.NamespaceId, serviceName, util2.ToJsonString(instance))
	instanceRequest := rpc_request2.NewInstanceRequest(proxy.clientConfig.NamespaceId, serviceName, groupName, "registerInstance", instance)
	response, err := proxy.requestToServer(instanceRequest)
	proxy.eventListener.CacheInstanceForRedo(serviceName, groupName, instance)
	if err != nil {
		return false, err
	}
	return response.IsSuccess(), err
}

// DeregisterInstance ...
func (proxy *NamingGrpcProxy) DeregisterInstance(serviceName string, groupName string, instance model2.Instance) (bool, error) {
	logger2.Infof("deregister instance namespaceId:<%s>,serviceName:<%s> with instance:<%s:%d@%s>",
		proxy.clientConfig.NamespaceId, serviceName, instance.Ip, instance.Port, instance.ClusterName)
	instanceRequest := rpc_request2.NewInstanceRequest(proxy.clientConfig.NamespaceId, serviceName, groupName, "deregisterInstance", instance)
	response, err := proxy.requestToServer(instanceRequest)
	proxy.eventListener.RemoveInstanceForRedo(serviceName, groupName, instance)
	if err != nil {
		return false, err
	}
	return response.IsSuccess(), err
}

// GetServiceList ...
func (proxy *NamingGrpcProxy) GetServiceList(pageNo uint32, pageSize uint32, groupName string, selector *model2.ExpressionSelector) (model2.ServiceList, error) {
	var selectorStr string
	if selector != nil {
		switch selector.Type {
		case "label":
			selectorStr = util2.ToJsonString(selector)
		default:
			break
		}
	}
	response, err := proxy.requestToServer(rpc_request2.NewServiceListRequest(proxy.clientConfig.NamespaceId, "",
		groupName, int(pageNo), int(pageSize), selectorStr))
	if err != nil {
		return model2.ServiceList{}, err
	}
	serviceListResponse := response.(*rpc_response2.ServiceListResponse)
	return model2.ServiceList{
		Count: int64(serviceListResponse.Count),
		Doms:  serviceListResponse.ServiceNames,
	}, nil
}

// ServerHealthy ...
func (proxy *NamingGrpcProxy) ServerHealthy() bool {
	return proxy.rpcClient.GetRpcClient().IsRunning()
}

// QueryInstancesOfService ...
func (proxy *NamingGrpcProxy) QueryInstancesOfService(serviceName, groupName, clusters string, udpPort int, healthyOnly bool) (*model2.Service, error) {
	response, err := proxy.requestToServer(rpc_request2.NewServiceQueryRequest(proxy.clientConfig.NamespaceId, serviceName, groupName, clusters,
		healthyOnly, udpPort))
	if err != nil {
		return nil, err
	}
	queryServiceResponse := response.(*rpc_response2.QueryServiceResponse)
	return &queryServiceResponse.ServiceInfo, nil
}

// Subscribe ...
func (proxy *NamingGrpcProxy) Subscribe(serviceName, groupName string, clusters string) (model2.Service, error) {
	proxy.eventListener.CacheSubscriberForRedo(util2.GetGroupName(serviceName, groupName), clusters)
	request := rpc_request2.NewSubscribeServiceRequest(proxy.clientConfig.NamespaceId, serviceName,
		groupName, clusters, true)
	request.Headers["app"] = proxy.clientConfig.AppName
	response, err := proxy.requestToServer(request)
	if err != nil {
		return model2.Service{}, err
	}
	subscribeServiceResponse := response.(*rpc_response2.SubscribeServiceResponse)
	return subscribeServiceResponse.ServiceInfo, nil
}

// Unsubscribe ...
func (proxy *NamingGrpcProxy) Unsubscribe(serviceName, groupName, clusters string) error {
	proxy.eventListener.RemoveSubscriberForRedo(util2.GetGroupName(serviceName, groupName), clusters)
	_, err := proxy.requestToServer(rpc_request2.NewSubscribeServiceRequest(proxy.clientConfig.NamespaceId, serviceName, groupName,
		clusters, false))
	return err
}

func (proxy *NamingGrpcProxy) CloseClient() {
	proxy.rpcClient.GetRpcClient().Shutdown()
}
