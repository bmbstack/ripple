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

package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	grpc2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/api/grpc"
	rpc_request2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/remote/rpc/rpc_request"
	rpc_response2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/remote/rpc/rpc_response"
	util2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/util"
	"time"

	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/grpc"
)

type GrpcConnection struct {
	*Connection
	client         grpc2.RequestClient
	biStreamClient grpc2.BiRequestStream_RequestBiStreamClient
}

func NewGrpcConnection(serverInfo ServerInfo, connectionId string, conn *grpc.ClientConn,
	client grpc2.RequestClient, biStreamClient grpc2.BiRequestStream_RequestBiStreamClient) *GrpcConnection {
	return &GrpcConnection{
		Connection: &Connection{
			serverInfo:   serverInfo,
			connectionId: connectionId,
			abandon:      false,
			conn:         conn,
		},
		client:         client,
		biStreamClient: biStreamClient,
	}
}
func (g *GrpcConnection) request(request rpc_request2.IRequest, timeoutMills int64, client *RpcClient) (rpc_response2.IResponse, error) {
	p := convertRequest(request)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutMills)*time.Millisecond)
	defer cancel()
	responsePayload, err := g.client.Request(ctx, p)
	if err != nil {
		return nil, err
	}

	responseFunc, ok := rpc_response2.ClientResponseMapping[responsePayload.Metadata.GetType()]

	if !ok {
		return nil, errors.New(fmt.Sprintf("request:%s,unsupported response type:%s", request.GetRequestType(),
			responsePayload.Metadata.GetType()))
	}
	response := responseFunc()
	err = json.Unmarshal(responsePayload.GetBody().Value, response)
	return response, err
}

func (g *GrpcConnection) close() {
	g.Connection.close()
}

func (g *GrpcConnection) biStreamSend(payload *grpc2.Payload) error {
	return g.biStreamClient.Send(payload)
}

func convertRequest(r rpc_request2.IRequest) *grpc2.Payload {
	Metadata := grpc2.Metadata{
		Type:     r.GetRequestType(),
		Headers:  r.GetHeaders(),
		ClientIp: util2.LocalIP(),
	}
	return &grpc2.Payload{
		Metadata: &Metadata,
		Body:     &any.Any{Value: []byte(r.GetBody(r))},
	}
}

func convertResponse(r rpc_response2.IResponse) *grpc2.Payload {
	Metadata := grpc2.Metadata{
		Type:     r.GetResponseType(),
		ClientIp: util2.LocalIP(),
	}
	return &grpc2.Payload{
		Metadata: &Metadata,
		Body:     &any.Any{Value: []byte(r.GetBody())},
	}
}
