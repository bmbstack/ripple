package ripple

import (
	"fmt"
	. "github.com/bmbstack/ripple/helper"
	"github.com/labstack/gommon/color"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	nserverplugin "github.com/rpcxio/rpcx-nacos/serverplugin"
	"github.com/smallnest/rpcx/server"
	"log"
)

// NewRpcServerNacos create rpc server
func NewRpcServerNacos(nacos Nacos) *server.Server {
	if IsEmpty(nacos.Server) {
		fmt.Println(color.Green("RPC: Just RPC service caller, not RPC service provider"))
		return nil
	}
	s := server.NewServer()
	fmt.Println(color.Green("RPC: RPC service provider"))
	clientConfig := constant.ClientConfig{
		TimeoutMs:            10 * 1000,
		BeatInterval:         5 * 1000,
		NamespaceId:          nacos.NamespaceId,
		CacheDir:             nacos.CacheDir,
		LogDir:               nacos.LogDir,
		UpdateThreadNum:      20,
		NotLoadCacheAtStart:  true,
		UpdateCacheWhenEmpty: true,
	}

	serverConfig := []constant.ServerConfig{{
		IpAddr: nacos.Host,
		Port:   nacos.Port,
	}}

	plugin := &nserverplugin.NacosRegisterPlugin{
		ServiceAddress: fmt.Sprintf("tcp@%s", nacos.Server),
		ClientConfig:   clientConfig,
		ServerConfig:   serverConfig,
		Cluster:        nacos.Cluster,
		Group:          nacos.Group,
	}

	err := plugin.Start()
	if err != nil {
		log.Fatal(err)
	}
	s.Plugins.Add(plugin)
	return s
}
