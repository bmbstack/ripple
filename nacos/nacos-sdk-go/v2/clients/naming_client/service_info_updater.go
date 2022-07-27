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
	naming_proxy2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/clients/naming_client/naming_proxy"
	logger2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/logger"
	model2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/model"
	util2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/util"
	"time"
)

type ServiceInfoUpdater struct {
	serviceInfoHolder *naming_cache2.ServiceInfoHolder
	updateThreadNum   int
	namingProxy       naming_proxy2.INamingProxy
}

func NewServiceInfoUpdater(serviceInfoHolder *naming_cache2.ServiceInfoHolder, updateThreadNum int,
	namingProxy naming_proxy2.INamingProxy) *ServiceInfoUpdater {

	return &ServiceInfoUpdater{
		serviceInfoHolder: serviceInfoHolder,
		updateThreadNum:   updateThreadNum,
		namingProxy:       namingProxy,
	}
}

func (s *ServiceInfoUpdater) asyncUpdateService() {
	sema := util2.NewSemaphore(s.updateThreadNum)
	for {
		for _, v := range s.serviceInfoHolder.ServiceInfoMap.Items() {
			service := v.(model2.Service)
			lastRefTime, ok := s.serviceInfoHolder.UpdateTimeMap.Get(util2.GetServiceCacheKey(util2.GetGroupName(service.Name, service.GroupName),
				service.Clusters))
			if !ok {
				lastRefTime = uint64(0)
			}
			if uint64(util2.CurrentMillis())-lastRefTime.(uint64) > service.CacheMillis {
				sema.Acquire()
				go func() {
					s.updateServiceNow(service.Name, service.GroupName, service.Clusters)
					sema.Release()
				}()
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func (s *ServiceInfoUpdater) updateServiceNow(serviceName, groupName, clusters string) {
	result, err := s.namingProxy.QueryInstancesOfService(serviceName, groupName, clusters, 0, false)
	if err != nil {
		logger2.Errorf("QueryList return error!serviceName:%s cluster:%s err:%+v", serviceName, clusters, err)
		return
	}
	// TODO modify
	result.Clusters = clusters
	s.serviceInfoHolder.ProcessService(result)
}
