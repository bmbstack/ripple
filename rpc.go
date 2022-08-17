package ripple

import (
	"fmt"
	. "github.com/bmbstack/ripple/helper"
	constant2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/constant"
	"github.com/bmbstack/ripple/nacos/rpcxnacos/serverplugin"
	"github.com/bmbstack/ripple/util"
	"github.com/labstack/gommon/color"
	"github.com/smallnest/rpcx/server"
	"log"
	"strings"
)

// NewRpcServerNacos create rpc server
func NewRpcServerNacos(nacos NacosConfig) (*server.Server, *serverplugin.NacosRegisterPlugin) {
	if IsEmpty(nacos.Server) {
		fmt.Println(color.Green("RPC: Just RPC service caller, not RPC service provider"))
		return nil, nil
	}
	s := server.NewServer()
	fmt.Println(color.Green("RPC: RPC service provider"))
	clientConfig := constant2.ClientConfig{
		TimeoutMs:            10 * 1000,
		BeatInterval:         5 * 1000,
		NamespaceId:          nacos.NamespaceId,
		CacheDir:             nacos.CacheDir,
		LogDir:               nacos.LogDir,
		UpdateThreadNum:      20,
		NotLoadCacheAtStart:  true,
		UpdateCacheWhenEmpty: true,
	}

	serverConfig := []constant2.ServerConfig{{
		IpAddr: nacos.Host,
		Port:   nacos.Port,
	}}

	arr := strings.Split(nacos.Server, ":")
	address := fmt.Sprintf("%s:%s", util.InternalIP(), arr[len(arr)-1:][0])
	plugin := &serverplugin.NacosRegisterPlugin{
		ServiceAddress: fmt.Sprintf("tcp@%s", address),
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
	return s, plugin
}
