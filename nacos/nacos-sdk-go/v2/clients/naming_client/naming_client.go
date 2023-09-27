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
	nacos_client2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/clients/nacos_client"
	naming_cache2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/clients/naming_client/naming_cache"
	naming_proxy2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/clients/naming_client/naming_proxy"
	constant2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/constant"
	logger2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/logger"
	model2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/model"
	util2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/util"
	vo2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/vo"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// NamingClient ...
type NamingClient struct {
	nacos_client2.INacosClient
	serviceProxy      naming_proxy2.INamingProxy
	serviceInfoHolder *naming_cache2.ServiceInfoHolder
}

// NewNamingClient ...
func NewNamingClient(nc nacos_client2.INacosClient) (*NamingClient, error) {
	rand.Seed(time.Now().UnixNano())
	naming := &NamingClient{INacosClient: nc}
	clientConfig, err := nc.GetClientConfig()
	if err != nil {
		return naming, err
	}

	serverConfig, err := nc.GetServerConfig()
	if err != nil {
		return naming, err
	}

	httpAgent, err := nc.GetHttpAgent()
	if err != nil {
		return naming, err
	}

	if err = initLogger(clientConfig); err != nil {
		return naming, err
	}
	if clientConfig.NamespaceId == "" {
		clientConfig.NamespaceId = constant2.DEFAULT_NAMESPACE_ID
	}
	naming.serviceInfoHolder = naming_cache2.NewServiceInfoHolder(clientConfig.NamespaceId, clientConfig.CacheDir,
		clientConfig.UpdateCacheWhenEmpty, clientConfig.NotLoadCacheAtStart)

	naming.serviceProxy, err = NewNamingProxyDelegate(clientConfig, serverConfig, httpAgent, naming.serviceInfoHolder)

	go NewServiceInfoUpdater(naming.serviceInfoHolder, clientConfig.UpdateThreadNum, naming.serviceProxy).asyncUpdateService()
	if err != nil {
		return naming, err
	}

	return naming, nil
}

func initLogger(clientConfig constant2.ClientConfig) error {
	return logger2.InitLogger(logger2.BuildLoggerConfig(clientConfig))
}

// RegisterInstance ...
func (sc *NamingClient) RegisterInstance(param vo2.RegisterInstanceParam) (bool, error) {
	if param.ServiceName == "" {
		return false, errors.New("serviceName cannot be empty!")
	}
	if len(param.GroupName) == 0 {
		param.GroupName = constant2.DEFAULT_GROUP
	}
	if param.Metadata == nil {
		param.Metadata = make(map[string]string)
	}
	instance := model2.Instance{
		Ip:          param.Ip,
		Port:        param.Port,
		Metadata:    param.Metadata,
		ClusterName: param.ClusterName,
		Healthy:     param.Healthy,
		Enable:      param.Enable,
		Weight:      param.Weight,
		Ephemeral:   param.Ephemeral,
	}

	return sc.serviceProxy.RegisterInstance(param.ServiceName, param.GroupName, instance)

}

// DeregisterInstance ...
func (sc *NamingClient) DeregisterInstance(param vo2.DeregisterInstanceParam) (bool, error) {
	if len(param.GroupName) == 0 {
		param.GroupName = constant2.DEFAULT_GROUP
	}
	instance := model2.Instance{
		Ip:          param.Ip,
		Port:        param.Port,
		ClusterName: param.Cluster,
		Ephemeral:   param.Ephemeral,
	}
	return sc.serviceProxy.DeregisterInstance(param.ServiceName, param.GroupName, instance)
}

// UpdateInstance ...
func (sc *NamingClient) UpdateInstance(param vo2.UpdateInstanceParam) (bool, error) {
	if param.ServiceName == "" {
		return false, errors.New("serviceName cannot be empty!")
	}
	if len(param.GroupName) == 0 {
		param.GroupName = constant2.DEFAULT_GROUP
	}
	if param.Metadata == nil {
		param.Metadata = make(map[string]string)
	}
	instance := model2.Instance{
		Ip:          param.Ip,
		Port:        param.Port,
		Metadata:    param.Metadata,
		ClusterName: param.ClusterName,
		Healthy:     param.Healthy,
		Enable:      param.Enable,
		Weight:      param.Weight,
		Ephemeral:   param.Ephemeral,
	}

	return sc.serviceProxy.RegisterInstance(param.ServiceName, param.GroupName, instance)

}

// GetService Get service info by Group and DataId, clusters was optional
func (sc *NamingClient) GetService(param vo2.GetServiceParam) (service model2.Service, err error) {
	if len(param.GroupName) == 0 {
		param.GroupName = constant2.DEFAULT_GROUP
	}
	var ok bool
	clusters := strings.Join(param.Clusters, ",")
	service, ok = sc.serviceInfoHolder.GetServiceInfo(param.ServiceName, param.GroupName, clusters)
	if !ok {
		service, err = sc.serviceProxy.Subscribe(param.ServiceName, param.GroupName, clusters)
	}
	return service, err
}

// GetAllServicesInfo Get all instance by Namespace and Group with page
func (sc *NamingClient) GetAllServicesInfo(param vo2.GetAllServiceInfoParam) (model2.ServiceList, error) {
	if len(param.GroupName) == 0 {
		param.GroupName = constant2.DEFAULT_GROUP
	}
	clientConfig, _ := sc.GetClientConfig()
	if len(param.NameSpace) == 0 {
		if len(clientConfig.NamespaceId) == 0 {
			param.NameSpace = constant2.DEFAULT_NAMESPACE_ID
		} else {
			param.NameSpace = clientConfig.NamespaceId
		}
	}
	services, err := sc.serviceProxy.GetServiceList(param.PageNo, param.PageSize, param.GroupName, &model2.ExpressionSelector{})
	return services, err
}

// SelectAllInstances Get all instance by DataId 和 Group
func (sc *NamingClient) SelectAllInstances(param vo2.SelectAllInstancesParam) ([]model2.Instance, error) {
	if len(param.GroupName) == 0 {
		param.GroupName = constant2.DEFAULT_GROUP
	}
	clusters := strings.Join(param.Clusters, ",")
	var (
		service model2.Service
		ok      bool
		err     error
	)

	service, ok = sc.serviceInfoHolder.GetServiceInfo(param.ServiceName, param.GroupName, clusters)
	if !ok {
		service, err = sc.serviceProxy.Subscribe(param.ServiceName, param.GroupName, clusters)
	}
	if err != nil || service.Hosts == nil || len(service.Hosts) == 0 {
		return []model2.Instance{}, err
	}
	return service.Hosts, err
}

// SelectInstances Get all instance by DataId, Group and Health
func (sc *NamingClient) SelectInstances(param vo2.SelectInstancesParam) ([]model2.Instance, error) {
	if len(param.GroupName) == 0 {
		param.GroupName = constant2.DEFAULT_GROUP
	}
	var (
		service model2.Service
		ok      bool
		err     error
	)
	clusters := strings.Join(param.Clusters, ",")
	service, ok = sc.serviceInfoHolder.GetServiceInfo(param.ServiceName, param.GroupName, clusters)
	if !ok {
		service, err = sc.serviceProxy.Subscribe(param.ServiceName, param.GroupName, clusters)
		if err != nil {
			return nil, err
		}
	}
	return sc.selectInstances(service, param.HealthyOnly)
}

func (sc *NamingClient) selectInstances(service model2.Service, healthy bool) ([]model2.Instance, error) {
	if service.Hosts == nil || len(service.Hosts) == 0 {
		return []model2.Instance{}, errors.New("instance list is empty!")
	}
	hosts := service.Hosts
	var result []model2.Instance
	for _, host := range hosts {
		if host.Healthy == healthy && host.Enable && host.Weight > 0 {
			result = append(result, host)
		}
	}
	return result, nil
}

// SelectOneHealthyInstance Get one healthy instance by DataId and Group
func (sc *NamingClient) SelectOneHealthyInstance(param vo2.SelectOneHealthInstanceParam) (*model2.Instance, error) {
	if len(param.GroupName) == 0 {
		param.GroupName = constant2.DEFAULT_GROUP
	}
	var (
		service model2.Service
		ok      bool
		err     error
	)
	clusters := strings.Join(param.Clusters, ",")
	service, ok = sc.serviceInfoHolder.GetServiceInfo(param.ServiceName, param.GroupName, clusters)
	if !ok {
		service, err = sc.serviceProxy.Subscribe(param.ServiceName, param.GroupName, clusters)
		if err != nil {
			return nil, err
		}
	}

	return sc.selectOneHealthyInstances(service)
}

func (sc *NamingClient) selectOneHealthyInstances(service model2.Service) (*model2.Instance, error) {
	if service.Hosts == nil || len(service.Hosts) == 0 {
		return nil, errors.New("instance list is empty!")
	}
	hosts := service.Hosts
	var result []model2.Instance
	mw := 0
	for _, host := range hosts {
		if host.Healthy && host.Enable && host.Weight > 0 {
			cw := int(math.Ceil(host.Weight))
			if cw > mw {
				mw = cw
			}
			result = append(result, host)
		}
	}
	if len(result) == 0 {
		return nil, errors.New("healthy instance list is empty!")
	}

	instance := newChooser(result).pick()
	return &instance, nil
}

// Subscribe ...
func (sc *NamingClient) Subscribe(param *vo2.SubscribeParam) error {
	if len(param.GroupName) == 0 {
		param.GroupName = constant2.DEFAULT_GROUP
	}
	clusters := strings.Join(param.Clusters, ",")
	sc.serviceInfoHolder.RegisterCallback(util2.GetGroupName(param.ServiceName, param.GroupName), clusters, &param.SubscribeCallback)
	_, err := sc.serviceProxy.Subscribe(param.ServiceName, param.GroupName, clusters)
	return err
}

// Unsubscribe ...
func (sc *NamingClient) Unsubscribe(param *vo2.SubscribeParam) (err error) {
	clusters := strings.Join(param.Clusters, ",")
	serviceFullName := util2.GetGroupName(param.ServiceName, param.GroupName)
	sc.serviceInfoHolder.DeregisterCallback(serviceFullName, clusters, &param.SubscribeCallback)
	if sc.serviceInfoHolder.IsSubscribed(serviceFullName, clusters) {
		err = sc.serviceProxy.Unsubscribe(param.ServiceName, param.GroupName, clusters)
	}

	return err
}

// CloseClient ...
func (sc *NamingClient) CloseClient() {
	sc.serviceProxy.CloseClient()
}
