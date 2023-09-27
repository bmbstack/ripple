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

package naming_client

import (
	naming_cache2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/clients/naming_client/naming_cache"
	naming_grpc2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/clients/naming_client/naming_grpc"
	naming_http2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/clients/naming_client/naming_http"
	naming_proxy2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/clients/naming_client/naming_proxy"
	constant2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/constant"
	http_agent2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/http_agent"
	nacos_server2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/nacos_server"
	model2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/model"
	util2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/util"
)

// NamingProxyDelegate ...
type NamingProxyDelegate struct {
	httpClientProxy   *naming_http2.NamingHttpProxy
	grpcClientProxy   *naming_grpc2.NamingGrpcProxy
	serviceInfoHolder *naming_cache2.ServiceInfoHolder
}

func NewNamingProxyDelegate(clientCfg constant2.ClientConfig, serverCfgs []constant2.ServerConfig,
	httpAgent http_agent2.IHttpAgent, serviceInfoHolder *naming_cache2.ServiceInfoHolder) (naming_proxy2.INamingProxy, error) {

	nacosServer, err := nacos_server2.NewNacosServer(serverCfgs, clientCfg, httpAgent, clientCfg.TimeoutMs, clientCfg.Endpoint)
	if err != nil {
		return nil, err
	}

	httpClientProxy, err := naming_http2.NewNamingHttpProxy(clientCfg, nacosServer, serviceInfoHolder)
	if err != nil {
		return nil, err
	}

	grpcClientProxy, err := naming_grpc2.NewNamingGrpcProxy(clientCfg, nacosServer, serviceInfoHolder)
	if err != nil {
		return nil, err
	}

	return &NamingProxyDelegate{
		httpClientProxy:   httpClientProxy,
		grpcClientProxy:   grpcClientProxy,
		serviceInfoHolder: serviceInfoHolder,
	}, nil
}

func (proxy *NamingProxyDelegate) getExecuteClientProxy(instance model2.Instance) (namingProxy naming_proxy2.INamingProxy) {
	if instance.Ephemeral {
		namingProxy = proxy.grpcClientProxy
	} else {
		namingProxy = proxy.httpClientProxy
	}
	return namingProxy
}

func (proxy *NamingProxyDelegate) RegisterInstance(serviceName string, groupName string, instance model2.Instance) (bool, error) {
	return proxy.getExecuteClientProxy(instance).RegisterInstance(serviceName, groupName, instance)
}

func (proxy *NamingProxyDelegate) DeregisterInstance(serviceName string, groupName string, instance model2.Instance) (bool, error) {
	return proxy.getExecuteClientProxy(instance).DeregisterInstance(serviceName, groupName, instance)
}

func (proxy *NamingProxyDelegate) GetServiceList(pageNo uint32, pageSize uint32, groupName string, selector *model2.ExpressionSelector) (model2.ServiceList, error) {
	return proxy.grpcClientProxy.GetServiceList(pageNo, pageSize, groupName, selector)
}

func (proxy *NamingProxyDelegate) ServerHealthy() bool {
	return proxy.grpcClientProxy.ServerHealthy() || proxy.httpClientProxy.ServerHealthy()
}

func (proxy *NamingProxyDelegate) QueryInstancesOfService(serviceName, groupName, clusters string, udpPort int, healthyOnly bool) (*model2.Service, error) {
	return proxy.grpcClientProxy.QueryInstancesOfService(serviceName, groupName, clusters, udpPort, healthyOnly)
}

func (proxy *NamingProxyDelegate) Subscribe(serviceName, groupName string, clusters string) (model2.Service, error) {
	serviceNameWithGroup := util2.GetServiceCacheKey(util2.GetGroupName(serviceName, groupName), clusters)
	serviceInfo, ok := proxy.serviceInfoHolder.ServiceInfoMap.Get(serviceNameWithGroup)
	if !ok {
		result, err := proxy.grpcClientProxy.Subscribe(serviceName, groupName, clusters)
		if err != nil {
			return model2.Service{}, err
		}
		serviceInfo = result
	}
	service := serviceInfo.(model2.Service)
	proxy.serviceInfoHolder.ProcessService(&service)
	return service, nil
}

func (proxy *NamingProxyDelegate) Unsubscribe(serviceName, groupName, clusters string) error {
	proxy.serviceInfoHolder.StopUpdateIfContain(util2.GetGroupName(serviceName, groupName), clusters)
	return proxy.grpcClientProxy.Unsubscribe(serviceName, groupName, clusters)
}

func (proxy *NamingProxyDelegate) CloseClient() {
	proxy.grpcClientProxy.CloseClient()
}
