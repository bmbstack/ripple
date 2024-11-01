package remotecfgx

import (
	"github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/clients"
	"github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/clients/config_client"
	"github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/vo"
	"github.com/pkg/errors"
	"github.com/sagikazarmark/crypt/backend"
)

type Client struct {
	client config_client.IConfigClient
	source vo.ConfigParam
}

func New(param vo.NacosClientParam, source vo.ConfigParam) (backend.Store, error) {
	client, err := clients.NewConfigClient(param)
	if err != nil {
		return nil, err
	}
	instance := &Client{
		client: client,
		source: source,
	}
	if err != nil {
		return nil, err
	}
	return instance, nil
}
func (c *Client) Get(path string) ([]byte, error) {
	content, err := c.client.GetConfig(c.source)
	if err != nil {
		return nil, err
	}
	return []byte(content), nil
}

func (c *Client) List(key string) (backend.KVPairs, error) {
	return nil, errors.New("Interface not implemented.")
}

func (c *Client) Set(key string, value []byte) error {
	return errors.New("Interface not implemented.")
}

func (c *Client) Watch(key string, stop chan bool) <-chan *backend.Response {
	respChan := make(chan *backend.Response)
	go func() {
		defer func() {
			close(respChan)
			_ = c.client.CancelListenConfig(c.source)
		}()

		err := c.client.ListenConfig(vo.ConfigParam{
			DataId: c.source.DataId,
			Group:  c.source.Group,
			OnChange: func(namespace, group, dataId, data string) {
				retValue, err := c.Get(key)
				if err != nil {
					respChan <- &backend.Response{Value: nil, Error: err}
					return
				}
				respChan <- &backend.Response{Value: retValue, Error: nil}
			},
		})
		if err != nil {
			return
		}

		<-stop
	}()

	return respChan
}
