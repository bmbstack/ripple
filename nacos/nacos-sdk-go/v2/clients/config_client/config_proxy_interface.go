package config_client

import (
	rpc2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/remote/rpc"
	rpc_request2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/remote/rpc/rpc_request"
	rpc_response2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/remote/rpc/rpc_response"
	model2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/model"
	vo2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/vo"
)

type IConfigProxy interface {
	queryConfig(dataId, group, tenant string, timeout uint64, notify bool, client *ConfigClient) (*rpc_response2.ConfigQueryResponse, error)
	searchConfigProxy(param vo2.SearchConfigParm, tenant, accessKey, secretKey string) (*model2.ConfigPage, error)
	requestProxy(rpcClient *rpc2.RpcClient, request rpc_request2.IRequest, timeoutMills uint64) (rpc_response2.IResponse, error)
	createRpcClient(taskId string, client *ConfigClient) *rpc2.RpcClient
	getRpcClient(client *ConfigClient) *rpc2.RpcClient
}
