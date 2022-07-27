package clients

import (
	constant2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/constant"
	vo2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/vo"
	"net"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getIntranetIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

func TestSetConfigClient(t *testing.T) {
	ip := getIntranetIP()
	sc := []constant2.ServerConfig{
		*constant2.NewServerConfig(
			ip,
			8848,
		),
	}

	cc := *constant2.NewClientConfig(
		constant2.WithNamespaceId("public"),
		constant2.WithTimeoutMs(5000),
		constant2.WithNotLoadCacheAtStart(true),
		constant2.WithLogDir("/tmp/nacos/log"),
		constant2.WithCacheDir("/tmp/nacos/cache"),
		constant2.WithLogLevel("debug"),
	)

	t.Run("setConfig_error", func(t *testing.T) {
		nacosClient, err := setConfig(vo2.NacosClientParam{})
		assert.Nil(t, nacosClient)
		assert.Equal(t, "server configs not found in properties", err.Error())
	})

	t.Run("setConfig_normal", func(t *testing.T) {
		// use map params setConfig
		param := getConfigParam(map[string]interface{}{
			"serverConfigs": sc,
			"clientConfig":  cc,
		})
		nacosClientFromMap, err := setConfig(param)
		assert.Nil(t, err)
		nacosClientFromStruct, err := setConfig(vo2.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		})
		assert.Nil(t, err)
		assert.True(t, reflect.DeepEqual(nacosClientFromMap, nacosClientFromStruct))
	})
}
