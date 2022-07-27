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

package clients

import (
	"errors"
	config_client2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/clients/config_client"
	nacos_client2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/clients/nacos_client"
	naming_client2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/clients/naming_client"
	constant2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/constant"
	http_agent2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/http_agent"
	vo2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/vo"
)

// CreateConfigClient use to create config client
func CreateConfigClient(properties map[string]interface{}) (iClient config_client2.IConfigClient, err error) {
	param := getConfigParam(properties)
	return NewConfigClient(param)
}

//CreateNamingClient use to create a nacos naming client
func CreateNamingClient(properties map[string]interface{}) (iClient naming_client2.INamingClient, err error) {
	param := getConfigParam(properties)
	return NewNamingClient(param)
}

func NewConfigClient(param vo2.NacosClientParam) (iClient config_client2.IConfigClient, err error) {
	nacosClient, err := setConfig(param)
	if err != nil {
		return
	}
	config, err := config_client2.NewConfigClient(nacosClient)
	if err != nil {
		return
	}
	iClient = config
	return
}

func NewNamingClient(param vo2.NacosClientParam) (iClient naming_client2.INamingClient, err error) {
	nacosClient, err := setConfig(param)
	if err != nil {
		return
	}
	naming, err := naming_client2.NewNamingClient(nacosClient)
	if err != nil {
		return
	}
	iClient = naming
	return
}

func getConfigParam(properties map[string]interface{}) (param vo2.NacosClientParam) {

	if clientConfigTmp, exist := properties[constant2.KEY_CLIENT_CONFIG]; exist {
		if clientConfig, ok := clientConfigTmp.(constant2.ClientConfig); ok {
			param.ClientConfig = &clientConfig
		}
	}
	if serverConfigTmp, exist := properties[constant2.KEY_SERVER_CONFIGS]; exist {
		if serverConfigs, ok := serverConfigTmp.([]constant2.ServerConfig); ok {
			param.ServerConfigs = serverConfigs
		}
	}
	return
}

func setConfig(param vo2.NacosClientParam) (iClient nacos_client2.INacosClient, err error) {
	client := &nacos_client2.NacosClient{}
	if param.ClientConfig == nil {
		// default clientConfig
		_ = client.SetClientConfig(constant2.ClientConfig{
			TimeoutMs:    10 * 1000,
			BeatInterval: 5 * 1000,
		})
	} else {
		err = client.SetClientConfig(*param.ClientConfig)
		if err != nil {
			return nil, err
		}
	}

	if len(param.ServerConfigs) == 0 {
		clientConfig, _ := client.GetClientConfig()
		if len(clientConfig.Endpoint) <= 0 {
			err = errors.New("server configs not found in properties")
			return nil, err
		}
		_ = client.SetServerConfig(nil)
	} else {
		err = client.SetServerConfig(param.ServerConfigs)
		if err != nil {
			return nil, err
		}
	}

	if _, _err := client.GetHttpAgent(); _err != nil {
		if clientCfg, err := client.GetClientConfig(); err == nil {
			_ = client.SetHttpAgent(&http_agent2.HttpAgent{TlsConfig: clientCfg.TLSCfg})
		}
	}
	iClient = client
	return
}
