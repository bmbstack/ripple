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

package main

import (
	"fmt"
	clients2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/clients"
	constant2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/constant"
	model2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/model"
	util2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/util"
	vo2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/vo"
	"time"
)

func main() {
	//create ServerConfig
	sc := []constant2.ServerConfig{
		*constant2.NewServerConfig("127.0.0.1", 8848, constant2.WithContextPath("/nacos")),
	}

	//create ClientConfig
	cc := *constant2.NewClientConfig(
		constant2.WithNamespaceId(""),
		constant2.WithTimeoutMs(5000),
		constant2.WithNotLoadCacheAtStart(true),
		constant2.WithLogDir("/tmp/nacos/log"),
		constant2.WithCacheDir("/tmp/nacos/cache"),
		constant2.WithLogLevel("debug"),
	)

	// create naming client
	client, err := clients2.NewNamingClient(
		vo2.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		panic(err)
	}

	//Register
	ExampleServiceClient_RegisterServiceInstance(client, vo2.RegisterInstanceParam{
		Ip:          "10.0.0.10",
		Port:        8848,
		ServiceName: "demo.go",
		GroupName:   "group-a",
		ClusterName: "cluster-a",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai"},
	})

	//DeRegister
	ExampleServiceClient_DeRegisterServiceInstance(client, vo2.DeregisterInstanceParam{
		Ip:          "10.0.0.10",
		Port:        8848,
		ServiceName: "demo.go",
		GroupName:   "group-a",
		Cluster:     "cluster-a",
		Ephemeral:   true, //it must be true
	})

	//Register
	ExampleServiceClient_RegisterServiceInstance(client, vo2.RegisterInstanceParam{
		Ip:          "10.0.0.10",
		Port:        8848,
		ServiceName: "demo.go",
		GroupName:   "group-a",
		ClusterName: "cluster-a",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai"},
	})

	time.Sleep(1 * time.Second)

	//Get service with serviceName, groupName , clusters
	ExampleServiceClient_GetService(client, vo2.GetServiceParam{
		ServiceName: "demo.go",
		GroupName:   "group-a",
		Clusters:    []string{"cluster-a"},
	})

	//SelectAllInstance
	//GroupName=DEFAULT_GROUP
	ExampleServiceClient_SelectAllInstances(client, vo2.SelectAllInstancesParam{
		ServiceName: "demo.go",
		GroupName:   "group-a",
		Clusters:    []string{"cluster-a"},
	})

	//SelectInstances only return the instances of healthy=${HealthyOnly},enable=true and weight>0
	//ClusterName=DEFAULT,GroupName=DEFAULT_GROUP
	ExampleServiceClient_SelectInstances(client, vo2.SelectInstancesParam{
		ServiceName: "demo.go",
		GroupName:   "group-a",
		Clusters:    []string{"cluster-a"},
	})

	//SelectOneHealthyInstance return one instance by WRR strategy for load balance
	//And the instance should be health=true,enable=true and weight>0
	//ClusterName=DEFAULT,GroupName=DEFAULT_GROUP
	ExampleServiceClient_SelectOneHealthyInstance(client, vo2.SelectOneHealthInstanceParam{
		ServiceName: "demo.go",
		GroupName:   "group-a",
		Clusters:    []string{"cluster-a"},
	})

	//Subscribe key=serviceName+groupName+cluster
	//Note:We call add multiple SubscribeCallback with the same key.
	param := &vo2.SubscribeParam{
		ServiceName: "demo.go",
		GroupName:   "group-a",
		SubscribeCallback: func(services []model2.Instance, err error) {
			fmt.Printf("callback return services:%s \n\n", util2.ToJsonString(services))
		},
	}
	ExampleServiceClient_Subscribe(client, param)

	ExampleServiceClient_RegisterServiceInstance(client, vo2.RegisterInstanceParam{
		Ip:          "10.0.0.10",
		Port:        8848,
		ServiceName: "demo.go",
		GroupName:   "group-a",
		ClusterName: "cluster-a",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "beijing"},
	})
	//wait for client pull change from server
	time.Sleep(3 * time.Second)

	ExampleServiceClient_UpdateServiceInstance(client, vo2.UpdateInstanceParam{
		Ip:          "10.0.0.11", //update ip
		Port:        8848,
		ServiceName: "demo.go",
		GroupName:   "group-a",
		ClusterName: "cluster-a",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "beijing1"}, //update metadata
	})

	//wait for client pull change from server
	time.Sleep(3 * time.Second)

	//Now we just unsubscribe callback1, and callback2 will still receive change event
	ExampleServiceClient_UnSubscribe(client, param)
	ExampleServiceClient_DeRegisterServiceInstance(client, vo2.DeregisterInstanceParam{
		Ip:          "10.0.0.112",
		Ephemeral:   true,
		Port:        8848,
		ServiceName: "demo.go",
		Cluster:     "cluster-b",
	})
	//wait for client pull change from server
	time.Sleep(3 * time.Second)

	//GeAllService will get the list of service name
	//NameSpace default value is public.If the client set the namespaceId, NameSpace will use it.
	//GroupName default value is DEFAULT_GROUP
	ExampleServiceClient_GetAllService(client, vo2.GetAllServiceInfoParam{
		PageNo:   1,
		PageSize: 10,
	})
}
