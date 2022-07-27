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

package naming_http

import (
	"errors"
	"fmt"
	naming_cache2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/clients/naming_client/naming_cache"
	constant2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/constant"
	logger2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/logger"
	nacos_server2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/nacos_server"
	model2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/model"
	util2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/util"
	"net/http"
	"strconv"
	"time"

	"github.com/buger/jsonparser"
)

// NamingHttpProxy ...
type NamingHttpProxy struct {
	clientConfig      constant2.ClientConfig
	nacosServer       *nacos_server2.NacosServer
	beatReactor       BeatReactor
	serviceInfoHolder *naming_cache2.ServiceInfoHolder
}

// NewNamingHttpProxy  create naming http proxy
func NewNamingHttpProxy(clientCfg constant2.ClientConfig, nacosServer *nacos_server2.NacosServer,
	serviceInfoHolder *naming_cache2.ServiceInfoHolder) (*NamingHttpProxy, error) {
	srvProxy := NamingHttpProxy{
		clientConfig:      clientCfg,
		nacosServer:       nacosServer,
		serviceInfoHolder: serviceInfoHolder,
	}

	srvProxy.beatReactor = NewBeatReactor(clientCfg, nacosServer)

	NewPushReceiver(serviceInfoHolder).startServer()

	return &srvProxy, nil
}

// RegisterInstance ...
func (proxy *NamingHttpProxy) RegisterInstance(serviceName string, groupName string, instance model2.Instance) (bool, error) {
	logger2.Infof("register instance namespaceId:<%s>,serviceName:<%s> with instance:<%s>",
		proxy.clientConfig.NamespaceId, serviceName, util2.ToJsonString(instance))
	serviceName = util2.GetGroupName(serviceName, groupName)
	params := map[string]string{}
	params["namespaceId"] = proxy.clientConfig.NamespaceId
	params["serviceName"] = serviceName
	params["groupName"] = groupName
	params["app"] = proxy.clientConfig.AppName
	params["clusterName"] = instance.ClusterName
	params["ip"] = instance.Ip
	params["port"] = strconv.Itoa(int(instance.Port))
	params["weight"] = strconv.FormatFloat(instance.Weight, 'f', -1, 64)
	params["enable"] = strconv.FormatBool(instance.Enable)
	params["healthy"] = strconv.FormatBool(instance.Healthy)
	params["metadata"] = util2.ToJsonString(instance.Metadata)
	params["ephemeral"] = strconv.FormatBool(instance.Ephemeral)
	_, err := proxy.nacosServer.ReqApi(constant2.SERVICE_PATH, params, http.MethodPost)
	if err != nil {
		return false, err
	}
	if instance.Ephemeral {
		beatInfo := &model2.BeatInfo{
			Ip:          instance.Ip,
			Port:        instance.Port,
			Metadata:    instance.Metadata,
			ServiceName: util2.GetGroupName(serviceName, groupName),
			Cluster:     instance.ClusterName,
			Weight:      instance.Weight,
			Period:      util2.GetDurationWithDefault(instance.Metadata, constant2.HEART_BEAT_INTERVAL, time.Second*5),
			State:       model2.StateRunning,
		}
		proxy.beatReactor.AddBeatInfo(util2.GetGroupName(serviceName, groupName), beatInfo)
	}
	return true, nil
}

// DeregisterInstance ...
func (proxy *NamingHttpProxy) DeregisterInstance(serviceName string, groupName string, instance model2.Instance) (bool, error) {
	serviceName = util2.GetGroupName(serviceName, groupName)
	logger2.Infof("deregister instance namespaceId:<%s>,serviceName:<%s> with instance:<%s:%d@%s>",
		proxy.clientConfig.NamespaceId, serviceName, instance.Ip, instance.Port, instance.ClusterName)
	proxy.beatReactor.RemoveBeatInfo(serviceName, instance.Ip, instance.Port)
	params := map[string]string{}
	params["namespaceId"] = proxy.clientConfig.NamespaceId
	params["serviceName"] = serviceName
	params["clusterName"] = instance.ClusterName
	params["ip"] = instance.Ip
	params["port"] = strconv.Itoa(int(instance.Port))
	params["ephemeral"] = strconv.FormatBool(instance.Ephemeral)
	_, err := proxy.nacosServer.ReqApi(constant2.SERVICE_PATH, params, http.MethodDelete)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetServiceList ...
func (proxy *NamingHttpProxy) GetServiceList(pageNo uint32, pageSize uint32, groupName string, selector *model2.ExpressionSelector) (model2.ServiceList, error) {
	params := map[string]string{}
	params["namespaceId"] = proxy.clientConfig.NamespaceId
	params["groupName"] = groupName
	params["pageNo"] = strconv.Itoa(int(pageNo))
	params["pageSize"] = strconv.Itoa(int(pageSize))

	if selector != nil {
		switch selector.Type {
		case "label":
			params["selector"] = util2.ToJsonString(selector)
			break
		default:
			break

		}
	}
	serviceList := model2.ServiceList{}

	api := constant2.SERVICE_BASE_PATH + "/service/list"
	result, err := proxy.nacosServer.ReqApi(api, params, http.MethodGet)
	if err != nil {
		return serviceList, err
	}
	if result == "" {
		return serviceList, errors.New("request server return empty")
	}

	count, err := jsonparser.GetInt([]byte(result), "count")
	if err != nil {
		return serviceList, errors.New(fmt.Sprintf("namespaceId:<%s> get service list pageNo:<%d> pageSize:<%d> selector:<%s> from <%s> get 'count' from <%s> error:<%+v>", proxy.clientConfig.NamespaceId, pageNo, pageSize, util2.ToJsonString(selector), groupName, result, err))
	}
	var doms []string
	_, err = jsonparser.ArrayEach([]byte(result), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		doms = append(doms, string(value))
	}, "doms")
	if err != nil {
		return serviceList, errors.New(fmt.Sprintf("namespaceId:<%s> get service list pageNo:<%d> pageSize:<%d> selector:<%s> from <%s> get 'doms' from <%s> error:<%+v> ", proxy.clientConfig.NamespaceId, pageNo, pageSize, util2.ToJsonString(selector), groupName, result, err))
	}
	serviceList.Count = count
	serviceList.Doms = doms
	return serviceList, nil
}

// ServerHealthy ...
func (proxy *NamingHttpProxy) ServerHealthy() bool {
	api := constant2.SERVICE_BASE_PATH + "/operator/metrics"
	result, err := proxy.nacosServer.ReqApi(api, map[string]string{}, http.MethodGet)
	if err != nil {
		logger2.Errorf("namespaceId:[%s] sending server healthy failed!,result:%s error:%+v", proxy.clientConfig.NamespaceId, result, err)
		return false
	}
	if result != "" {
		status, err := jsonparser.GetString([]byte(result), "status")
		if err != nil {
			logger2.Errorf("namespaceId:[%s] sending server healthy failed!,result:%s error:%+v", proxy.clientConfig.NamespaceId, result, err)
		} else {
			return status == "UP"
		}
	}
	return false
}

// QueryInstancesOfService ...
func (proxy *NamingHttpProxy) QueryInstancesOfService(serviceName, groupName, clusters string, udpPort int, healthyOnly bool) (*model2.Service, error) {
	param := make(map[string]string)
	param["namespaceId"] = proxy.clientConfig.NamespaceId
	param["serviceName"] = util2.GetGroupName(serviceName, groupName)
	param["app"] = proxy.clientConfig.AppName
	param["clusters"] = clusters
	param["udpPort"] = strconv.Itoa(udpPort)
	param["healthyOnly"] = strconv.FormatBool(healthyOnly)
	param["clientIP"] = util2.LocalIP()
	api := constant2.SERVICE_PATH + "/list"
	result, err := proxy.nacosServer.ReqApi(api, param, http.MethodGet)
	if err != nil {
		return nil, err
	}
	return util2.JsonToService(result), nil

}

// Subscribe ...
func (proxy *NamingHttpProxy) Subscribe(serviceName, groupName, clusters string) (model2.Service, error) {
	return model2.Service{}, nil
}

// Unsubscribe ...
func (proxy *NamingHttpProxy) Unsubscribe(serviceName, groupName, clusters string) error {
	return nil
}

func (proxy *NamingHttpProxy) CloseClient() {

}
