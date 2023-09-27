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
	cache2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/clients/cache"
	naming_proxy2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/clients/naming_client/naming_proxy"
	constant2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/constant"
	logger2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/logger"
	model2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/model"
	util2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/util"
	"reflect"
	"strings"
)

type ConnectionEventListener struct {
	clientProxy              naming_proxy2.INamingProxy
	registeredInstanceCached cache2.ConcurrentMap
	subscribes               cache2.ConcurrentMap
}

func NewConnectionEventListener(clientProxy naming_proxy2.INamingProxy) *ConnectionEventListener {
	return &ConnectionEventListener{
		clientProxy:              clientProxy,
		registeredInstanceCached: cache2.NewConcurrentMap(),
		subscribes:               cache2.NewConcurrentMap(),
	}
}

func (c *ConnectionEventListener) OnConnected() {
	c.redoSubscribe()
	c.redoRegisterEachService()
}

func (c *ConnectionEventListener) OnDisConnect() {

}

func (c *ConnectionEventListener) redoSubscribe() {
	for _, key := range c.subscribes.Keys() {
		info := strings.Split(key, constant2.SERVICE_INFO_SPLITER)
		var err error
		var service model2.Service
		if len(info) > 2 {
			service, err = c.clientProxy.Subscribe(info[1], info[0], info[2])
		} else {
			service, err = c.clientProxy.Subscribe(info[1], info[0], "")
		}

		if err != nil {
			logger2.Warnf("redo subscribe service:%s faild:%+v", info[1], err)
			return
		}

		grpcProxy, ok := c.clientProxy.(*NamingGrpcProxy)
		if !ok {
			return
		}
		grpcProxy.serviceInfoHolder.ProcessService(&service)
	}
}

func (c *ConnectionEventListener) redoRegisterEachService() {
	for k, v := range c.registeredInstanceCached.Items() {
		info := strings.Split(k, constant2.SERVICE_INFO_SPLITER)
		serviceName := info[1]
		groupName := info[0]
		instance, ok := v.(model2.Instance)
		if !ok {
			logger2.Warnf("redo register service:%s faild,instances type not is model.instance", info[1])
			return
		}
		_, err := c.clientProxy.RegisterInstance(serviceName, groupName, instance)
		if err != nil {
			logger2.Warnf("redo register service:%s groupName:%s faild:%s", info[1], info[0], err.Error())
		}
	}
}

func (c *ConnectionEventListener) CacheInstanceForRedo(serviceName, groupName string, instance model2.Instance) {
	key := util2.GetGroupName(serviceName, groupName)
	getInstance, _ := c.registeredInstanceCached.Get(key)
	if getInstance != nil && reflect.DeepEqual(getInstance.(model2.Instance), instance) {
		return
	}
	c.registeredInstanceCached.Set(key, instance)
}

func (c *ConnectionEventListener) RemoveInstanceForRedo(serviceName, groupName string, instance model2.Instance) {
	key := util2.GetGroupName(serviceName, groupName)
	_, ok := c.registeredInstanceCached.Get(key)
	if !ok {
		return
	}
	c.registeredInstanceCached.Remove(key)
}

func (c *ConnectionEventListener) CacheSubscriberForRedo(fullServiceName, clusters string) {
	key := util2.GetServiceCacheKey(fullServiceName, clusters)
	if _, ok := c.subscribes.Get(key); !ok {
		c.subscribes.Set(key, struct{}{})
	}
}

func (c *ConnectionEventListener) RemoveSubscriberForRedo(fullServiceName, clusters string) {
	c.subscribes.Remove(util2.GetServiceCacheKey(fullServiceName, clusters))
}
