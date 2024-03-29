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

package cache

import (
	"encoding/json"
	"fmt"
	file2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/file"
	logger2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/logger"
	model2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/model"
	util2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/util"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/go-errors/errors"
)

func GetFileName(cacheKey string, cacheDir string) string {
	return cacheDir + string(os.PathSeparator) + cacheKey
}

func WriteServicesToFile(service model2.Service, cacheDir string) {
	file2.MkdirIfNecessary(cacheDir)
	sb, _ := json.Marshal(service)
	domFileName := GetFileName(util2.GetServiceCacheKey(service.Name, service.Clusters), cacheDir)

	err := ioutil.WriteFile(domFileName, sb, 0666)
	if err != nil {
		logger2.Errorf("failed to write name cache:%s ,value:%s ,err:%+v", domFileName, string(sb), err)
	}

}

func ReadServicesFromFile(cacheDir string) map[string]model2.Service {
	files, err := ioutil.ReadDir(cacheDir)
	if err != nil {
		logger2.Errorf("read cacheDir:%s failed!err:%+v", cacheDir, err)
		return nil
	}
	serviceMap := map[string]model2.Service{}
	for _, f := range files {
		fileName := GetFileName(f.Name(), cacheDir)
		b, err := ioutil.ReadFile(fileName)
		if err != nil {
			logger2.Errorf("failed to read name cache file:%s,err:%+v ", fileName, err)
			continue
		}

		s := string(b)
		service := util2.JsonToService(s)

		if service == nil {
			continue
		}
		cacheKey := util2.GetServiceCacheKey(util2.GetGroupName(service.Name, service.GroupName), service.Clusters)
		serviceMap[cacheKey] = *service
	}

	logger2.Infof("finish loading name cache, total: %s", strconv.Itoa(len(files)))
	return serviceMap
}

func WriteConfigToFile(cacheKey string, cacheDir string, content string) {
	file2.MkdirIfNecessary(cacheDir)
	fileName := GetFileName(cacheKey, cacheDir)
	if len(content) == 0 {
		// delete config snapshot
		if err := os.Remove(fileName); err != nil {
			logger2.Errorf("failed to delete config file,cache:%s ,value:%s ,err:%+v", fileName, content, err)
		}
		return
	}
	err := ioutil.WriteFile(fileName, []byte(content), 0666)
	if err != nil {
		logger2.Errorf("failed to write config  cache:%s ,value:%s ,err:%+v", fileName, content, err)
	}

}

func ReadConfigFromFile(cacheKey string, cacheDir string) (string, error) {
	fileName := GetFileName(cacheKey, cacheDir)
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", errors.New(fmt.Sprintf("failed to read config cache file:%s,err:%+v ", fileName, err))
	}
	return string(b), nil
}
