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

package naming_cache

import (
	"fmt"
	cache2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/clients/cache"
	logger2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/logger"
	monitor2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/monitor"
	model2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/model"
	util2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/util"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type ServiceInfoHolder struct {
	ServiceInfoMap       cache2.ConcurrentMap
	updateCacheWhenEmpty bool
	cacheDir             string
	notLoadCacheAtStart  bool
	subCallback          *SubscribeCallback
	UpdateTimeMap        cache2.ConcurrentMap
}

func NewServiceInfoHolder(namespace, cacheDir string, updateCacheWhenEmpty, notLoadCacheAtStart bool) *ServiceInfoHolder {
	cacheDir = cacheDir + string(os.PathSeparator) + "naming" + string(os.PathSeparator) + namespace
	serviceInfoHolder := &ServiceInfoHolder{
		updateCacheWhenEmpty: updateCacheWhenEmpty,
		notLoadCacheAtStart:  notLoadCacheAtStart,
		cacheDir:             cacheDir,
		subCallback:          NewSubscribeCallback(),
		UpdateTimeMap:        cache2.NewConcurrentMap(),
		ServiceInfoMap:       cache2.NewConcurrentMap(),
	}

	if !notLoadCacheAtStart {
		serviceInfoHolder.loadCacheFromDisk()
	}
	return serviceInfoHolder
}

func (s *ServiceInfoHolder) loadCacheFromDisk() {
	serviceMap := cache2.ReadServicesFromFile(s.cacheDir)
	if serviceMap == nil || len(serviceMap) == 0 {
		return
	}
	for k, v := range serviceMap {
		s.ServiceInfoMap.Set(k, v)
	}
}

func (s *ServiceInfoHolder) ProcessServiceJson(data string) {
	s.ProcessService(util2.JsonToService(data))
}

func (s *ServiceInfoHolder) ProcessService(service *model2.Service) {
	if service == nil {
		return
	}
	if !s.updateCacheWhenEmpty {
		//if instance list is empty,not to update cache
		if service.Hosts == nil || len(service.Hosts) == 0 {
			logger2.Warnf("instance list is empty, updateCacheWhenEmpty is set to false, callback is not triggered. service name:%s", service.Name)
			return
		}
	}

	cacheKey := util2.GetServiceCacheKey(util2.GetGroupName(service.Name, service.GroupName), service.Clusters)
	oldDomain, ok := s.ServiceInfoMap.Get(cacheKey)
	if ok && oldDomain.(model2.Service).LastRefTime >= service.LastRefTime {
		logger2.Warnf("out of date data received, old-t: %d, new-t: %d", oldDomain.(model2.Service).LastRefTime, service.LastRefTime)
		return
	}

	s.UpdateTimeMap.Set(cacheKey, uint64(util2.CurrentMillis()))
	s.ServiceInfoMap.Set(cacheKey, *service)
	if !ok || checkInstanceChanged(oldDomain, service) {
		logger2.Infof("service key:%s was updated to:%s", cacheKey, util2.ToJsonString(service))
		cache2.WriteServicesToFile(*service, s.cacheDir)
		s.subCallback.ServiceChanged(cacheKey, service)
	}
	monitor2.GetServiceInfoMapSizeMonitor().Set(float64(s.ServiceInfoMap.Count()))
}

func (s *ServiceInfoHolder) GetServiceInfo(serviceName, groupName, clusters string) (model2.Service, bool) {
	cacheKey := util2.GetServiceCacheKey(util2.GetGroupName(serviceName, groupName), clusters)
	//todo FailoverReactor
	service, ok := s.ServiceInfoMap.Get(cacheKey)
	if ok {
		return service.(model2.Service), ok
	}
	return model2.Service{}, ok
}

func (s *ServiceInfoHolder) RegisterCallback(serviceName string, clusters string, callbackFunc *func(services []model2.Instance, err error)) {
	s.subCallback.AddCallbackFunc(serviceName, clusters, callbackFunc)
}

func (s *ServiceInfoHolder) DeregisterCallback(serviceName string, clusters string, callbackFunc *func(services []model2.Instance, err error)) {
	s.subCallback.RemoveCallbackFunc(serviceName, clusters, callbackFunc)
}

func (s *ServiceInfoHolder) StopUpdateIfContain(serviceName, clusters string) {
	cacheKey := util2.GetServiceCacheKey(serviceName, clusters)
	s.ServiceInfoMap.Remove(cacheKey)
}

func (s *ServiceInfoHolder) IsSubscribed(serviceName, clusters string) bool {
	return s.subCallback.IsSubscribed(serviceName, clusters)
}

func checkInstanceChanged(oldDomain interface{}, service *model2.Service) bool {
	if oldDomain == nil {
		return true
	}
	oldService := oldDomain.(model2.Service)
	return isServiceInstanceChanged(&oldService, service)
}

// return true when service instance changed ,otherwise return false.
func isServiceInstanceChanged(oldService, newService *model2.Service) bool {
	oldHostsLen := len(oldService.Hosts)
	newHostsLen := len(newService.Hosts)
	if oldHostsLen != newHostsLen {
		return true
	}
	// compare refTime
	oldRefTime := oldService.LastRefTime
	newRefTime := newService.LastRefTime
	if oldRefTime > newRefTime {
		logger2.Warn(fmt.Sprintf("out of date data received, old-t: %v , new-t:  %v", oldRefTime, newRefTime))
		return false
	}
	// sort instance list
	oldInstance := oldService.Hosts
	newInstance := make([]model2.Instance, len(newService.Hosts))
	copy(newInstance, newService.Hosts)
	sortInstance(oldInstance)
	sortInstance(newInstance)
	return !reflect.DeepEqual(oldInstance, newInstance)
}

type instanceSorter []model2.Instance

func (s instanceSorter) Len() int {
	return len(s)
}
func (s instanceSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s instanceSorter) Less(i, j int) bool {
	insI, insJ := s[i], s[j]
	// using ip and port to sort
	ipNum1, _ := strconv.Atoi(strings.ReplaceAll(insI.Ip, ".", ""))
	ipNum2, _ := strconv.Atoi(strings.ReplaceAll(insJ.Ip, ".", ""))
	if ipNum1 < ipNum2 {
		return true
	}
	if insI.Port < insJ.Port {
		return true
	}
	return false
}

// sort instances
func sortInstance(instances []model2.Instance) {
	sort.Sort(instanceSorter(instances))
}
