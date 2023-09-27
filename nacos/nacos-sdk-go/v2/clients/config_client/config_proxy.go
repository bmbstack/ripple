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

package config_client

import (
	"encoding/json"
	"errors"
	cache2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/clients/cache"
	constant2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/constant"
	http_agent2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/http_agent"
	logger2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/logger"
	monitor2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/monitor"
	nacos_server2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/nacos_server"
	rpc2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/remote/rpc"
	rpc_request2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/remote/rpc/rpc_request"
	rpc_response2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/remote/rpc/rpc_response"
	model2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/model"
	util2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/util"
	vo2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/vo"
	"net/http"
	"strconv"
	"time"
)

type ConfigProxy struct {
	nacosServer  *nacos_server2.NacosServer
	clientConfig constant2.ClientConfig
}

func NewConfigProxy(serverConfig []constant2.ServerConfig, clientConfig constant2.ClientConfig, httpAgent http_agent2.IHttpAgent) (IConfigProxy, error) {
	proxy := ConfigProxy{}
	var err error
	proxy.nacosServer, err = nacos_server2.NewNacosServer(serverConfig, clientConfig, httpAgent, clientConfig.TimeoutMs, clientConfig.Endpoint)
	proxy.clientConfig = clientConfig
	return &proxy, err
}

func (cp *ConfigProxy) requestProxy(rpcClient *rpc2.RpcClient, request rpc_request2.IRequest, timeoutMills uint64) (rpc_response2.IResponse, error) {
	start := time.Now()
	cp.nacosServer.InjectSecurityInfo(request.GetHeaders())
	cp.injectCommHeader(request.GetHeaders())
	cp.nacosServer.InjectSkAk(request.GetHeaders(), cp.clientConfig)
	signHeaders := nacos_server2.GetSignHeaders(request.GetHeaders(), cp.clientConfig.SecretKey)
	request.PutAllHeaders(signHeaders)
	//todo Config Limiter
	response, err := rpcClient.Request(request, int64(timeoutMills))
	monitor2.GetConfigRequestMonitor(constant2.GRPC, request.GetRequestType(), rpc_response2.GetGrpcResponseStatusCode(response)).Observe(float64(time.Now().Nanosecond() - start.Nanosecond()))
	return response, err
}

func (cp *ConfigProxy) injectCommHeader(param map[string]string) {
	now := strconv.FormatInt(util2.CurrentMillis(), 10)
	param[constant2.CLIENT_APPNAME_HEADER] = cp.clientConfig.AppName
	param[constant2.CLIENT_REQUEST_TS_HEADER] = now
	param[constant2.CLIENT_REQUEST_TOKEN_HEADER] = util2.Md5(now + cp.clientConfig.AppKey)
	param[constant2.EX_CONFIG_INFO] = "true"
	param[constant2.CHARSET_KEY] = "utf-8"
}

func (cp *ConfigProxy) searchConfigProxy(param vo2.SearchConfigParm, tenant, accessKey, secretKey string) (*model2.ConfigPage, error) {
	params := util2.TransformObject2Param(param)
	if len(tenant) > 0 {
		params["tenant"] = tenant
	}
	if _, ok := params["group"]; !ok {
		params["group"] = ""
	}
	if _, ok := params["dataId"]; !ok {
		params["dataId"] = ""
	}
	var headers = map[string]string{}
	headers["accessKey"] = accessKey
	headers["secretKey"] = secretKey
	result, err := cp.nacosServer.ReqConfigApi(constant2.CONFIG_PATH, params, headers, http.MethodGet, cp.clientConfig.TimeoutMs)
	if err != nil {
		return nil, err
	}
	var configPage model2.ConfigPage
	err = json.Unmarshal([]byte(result), &configPage)
	if err != nil {
		return nil, err
	}
	return &configPage, nil
}

func (cp *ConfigProxy) queryConfig(dataId, group, tenant string, timeout uint64, notify bool, client *ConfigClient) (*rpc_response2.ConfigQueryResponse, error) {
	if group == "" {
		group = constant2.DEFAULT_GROUP
	}
	configQueryRequest := rpc_request2.NewConfigQueryRequest(group, dataId, tenant)
	configQueryRequest.Headers["notify"] = strconv.FormatBool(notify)
	iResponse, err := cp.requestProxy(cp.getRpcClient(client), configQueryRequest, timeout)
	if err != nil {
		return nil, err
	}
	response, ok := iResponse.(*rpc_response2.ConfigQueryResponse)
	if !ok {
		return nil, errors.New("ConfigQueryRequest returns type error")
	}
	if response.IsSuccess() {
		//todo LocalConfigInfoProcessor.saveSnapshot
		cacheKey := util2.GetConfigCacheKey(dataId, group, tenant)
		cache2.WriteConfigToFile(cacheKey, cp.clientConfig.CacheDir, response.Content)
		//todo LocalConfigInfoProcessor.saveEncryptDataKeySnapshot
		if response.ContentType == "" {
			response.ContentType = "text"
		}
		return response, nil
	}

	if response.GetErrorCode() == 300 {
		//todo LocalConfigInfoProcessor.saveSnapshot
		cacheKey := util2.GetConfigCacheKey(dataId, group, tenant)
		cache2.WriteConfigToFile(cacheKey, cp.clientConfig.CacheDir, "")
		//todo LocalConfigInfoProcessor.saveEncryptDataKeySnapshot
		return response, nil
	}

	if response.GetErrorCode() == 400 {
		logger2.Errorf(
			"[config_rpc_client] [sub-server-error] get server config being modified concurrently, dataId=%s, group=%s, "+
				"tenant=%s", dataId, group, tenant)
		return nil, errors.New("data being modified, dataId=" + dataId + ",group=" + group + ",tenant=" + tenant)
	}

	if response.GetErrorCode() > 0 {
		logger2.Errorf("[config_rpc_client] [sub-server-error]  dataId=%s, group=%s, tenant=%s, code=%v", dataId, group,
			tenant, response)
	}
	return response, nil
}

func (cp *ConfigProxy) createRpcClient(taskId string, client *ConfigClient) *rpc2.RpcClient {
	labels := map[string]string{
		constant2.LABEL_SOURCE: constant2.LABEL_SOURCE_SDK,
		constant2.LABEL_MODULE: constant2.LABEL_MODULE_CONFIG,
		"taskId":               taskId,
	}

	iRpcClient, _ := rpc2.CreateClient("config-"+taskId+"-"+client.uid, rpc2.GRPC, labels, cp.nacosServer)
	rpcClient := iRpcClient.GetRpcClient()
	if rpcClient.IsInitialized() {
		rpcClient.RegisterServerRequestHandler(func() rpc_request2.IRequest {
			return &rpc_request2.ConfigChangeNotifyRequest{ConfigRequest: rpc_request2.NewConfigRequest()}
		}, &ConfigChangeNotifyRequestHandler{client: client})
		rpcClient.Tenant = cp.clientConfig.NamespaceId
		rpcClient.Start()
	}
	return rpcClient
}

func (cp *ConfigProxy) getRpcClient(client *ConfigClient) *rpc2.RpcClient {
	return cp.createRpcClient("0", client)
}

type ConfigChangeNotifyRequestHandler struct {
	client *ConfigClient
}

func (c *ConfigChangeNotifyRequestHandler) Name() string {
	return "ConfigChangeNotifyRequestHandler"
}

func (c *ConfigChangeNotifyRequestHandler) RequestReply(request rpc_request2.IRequest, rpcClient *rpc2.RpcClient) rpc_response2.IResponse {
	configChangeNotifyRequest, ok := request.(*rpc_request2.ConfigChangeNotifyRequest)
	if ok {
		logger2.Infof("%s [server-push] config changed. dataId=%s, group=%s,tenant=%s", rpcClient.Name,
			configChangeNotifyRequest.DataId, configChangeNotifyRequest.Group, configChangeNotifyRequest.Tenant)

		cacheKey := util2.GetConfigCacheKey(configChangeNotifyRequest.DataId, configChangeNotifyRequest.Group,
			configChangeNotifyRequest.Tenant)
		data, ok := c.client.cacheMap.Get(cacheKey)
		if !ok {
			return nil
		}
		cData := data.(*cacheData)
		cData.isSyncWithServer = false
		c.client.notifyListenConfig()
		return &rpc_response2.NotifySubscriberResponse{
			Response: &rpc_response2.Response{ResultCode: constant2.RESPONSE_CODE_SUCCESS},
		}
	}
	return nil
}
