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
	"context"
	"errors"
	"fmt"
	cache2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/clients/cache"
	constant2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/constant"
	logger2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/logger"
	monitor2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/monitor"
	nacos_server2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/nacos_server"
	model2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/model"
	util2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/util"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/buger/jsonparser"
	"golang.org/x/sync/semaphore"
)

type BeatReactor struct {
	beatMap             cache2.ConcurrentMap
	nacosServer         *nacos_server2.NacosServer
	beatThreadCount     int
	beatThreadSemaphore *semaphore.Weighted
	beatRecordMap       cache2.ConcurrentMap
	clientCfg           constant2.ClientConfig
	mux                 *sync.Mutex
}

const DefaultBeatThreadNum = 20

var ctx = context.Background()

func NewBeatReactor(clientCfg constant2.ClientConfig, nacosServer *nacos_server2.NacosServer) BeatReactor {
	br := BeatReactor{}
	br.beatMap = cache2.NewConcurrentMap()
	br.nacosServer = nacosServer
	br.clientCfg = clientCfg
	br.beatThreadCount = DefaultBeatThreadNum
	br.beatRecordMap = cache2.NewConcurrentMap()
	br.beatThreadSemaphore = semaphore.NewWeighted(int64(br.beatThreadCount))
	br.mux = new(sync.Mutex)
	return br
}

func buildKey(serviceName string, ip string, port uint64) string {
	return serviceName + constant2.NAMING_INSTANCE_ID_SPLITTER + ip + constant2.NAMING_INSTANCE_ID_SPLITTER + strconv.Itoa(int(port))
}

func (br *BeatReactor) AddBeatInfo(serviceName string, beatInfo *model2.BeatInfo) {
	logger2.Infof("adding beat: <%s> to beat map", util2.ToJsonString(beatInfo))
	k := buildKey(serviceName, beatInfo.Ip, beatInfo.Port)
	defer br.mux.Unlock()
	br.mux.Lock()
	if data, ok := br.beatMap.Get(k); ok {
		beatInfo = data.(*model2.BeatInfo)
		atomic.StoreInt32(&beatInfo.State, int32(model2.StateShutdown))
		br.beatMap.Remove(k)
	}
	br.beatMap.Set(k, beatInfo)
	beatInfo.Metadata = util2.DeepCopyMap(beatInfo.Metadata)
	monitor2.GetDom2BeatSizeMonitor().Set(float64(br.beatMap.Count()))
	go br.sendInstanceBeat(k, beatInfo)
}

func (br *BeatReactor) RemoveBeatInfo(serviceName string, ip string, port uint64) {
	logger2.Infof("remove beat: %s@%s:%d from beat map", serviceName, ip, port)
	k := buildKey(serviceName, ip, port)
	defer br.mux.Unlock()
	br.mux.Lock()
	data, exist := br.beatMap.Get(k)
	if exist {
		beatInfo := data.(*model2.BeatInfo)
		atomic.StoreInt32(&beatInfo.State, int32(model2.StateShutdown))
	}
	monitor2.GetDom2BeatSizeMonitor().Set(float64(br.beatMap.Count()))
	br.beatMap.Remove(k)

}

func (br *BeatReactor) sendInstanceBeat(k string, beatInfo *model2.BeatInfo) {
	for {
		br.beatThreadSemaphore.Acquire(ctx, 1)
		//如果当前实例注销，则进行停止心跳
		if atomic.LoadInt32(&beatInfo.State) == int32(model2.StateShutdown) {
			logger2.Infof("instance[%s] stop heartBeating", k)
			br.beatThreadSemaphore.Release(1)
			return
		}

		//进行心跳通信
		beatInterval, err := br.SendBeat(beatInfo)
		if err != nil {
			logger2.Errorf("beat to server return error:%+v", err)
			br.beatThreadSemaphore.Release(1)
			t := time.NewTimer(beatInfo.Period)
			<-t.C
			continue
		}
		if beatInterval > 0 {
			beatInfo.Period = time.Duration(time.Millisecond.Nanoseconds() * beatInterval)
		}

		br.beatRecordMap.Set(k, util2.CurrentMillis())
		br.beatThreadSemaphore.Release(1)

		t := time.NewTimer(beatInfo.Period)
		<-t.C
	}
}

func (br *BeatReactor) SendBeat(info *model2.BeatInfo) (int64, error) {
	logger2.Infof("namespaceId:<%s> sending beat to server:<%s>",
		br.clientCfg.NamespaceId, util2.ToJsonString(info))
	params := map[string]string{}
	params["namespaceId"] = br.clientCfg.NamespaceId
	params["serviceName"] = info.ServiceName
	params["beat"] = util2.ToJsonString(info)
	api := constant2.SERVICE_BASE_PATH + "/instance/beat"
	result, err := br.nacosServer.ReqApi(api, params, http.MethodPut)
	if err != nil {
		return 0, err
	}
	if result != "" {
		interVal, err := jsonparser.GetInt([]byte(result), "clientBeatInterval")
		if err != nil {
			return 0, errors.New(fmt.Sprintf("namespaceId:<%s> sending beat to server:<%s> get 'clientBeatInterval' from <%s> error:<%+v>", br.clientCfg.NamespaceId, util2.ToJsonString(info), result, err))
		} else {
			return interVal, nil
		}
	}
	return 0, nil
}
