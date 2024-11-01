package remotecfgx

import (
	"fmt"
	"github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/constant"
	"github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_RemoteConfig(t *testing.T) {
	clientParam := vo.NacosClientParam{
		ClientConfig: &constant.ClientConfig{
			Endpoint:    "xxxx",
			AccessKey:   "xxxx",
			SecretKey:   "xxxx",
			NamespaceId: "xxxx",
		},
	}
	var vipers []*viper.Viper
	sources := []vo.ConfigParam{
		{Group: "xxxx1", DataId: "xxxx1"},
		{Group: "xxxx2", DataId: "xxxx2"},
	}
	for _, source := range sources {
		vp := viper.New()
		vp.SetConfigType("yaml")
		path := buildKey(source.Group, source.DataId)
		RegisterConfig(path, clientParam, source)
		err := vp.AddRemoteProvider("nacos", "xxxx", path)
		if !assert.Nil(t, err) {
			panic(err)
		}
		err = vp.ReadRemoteConfig()
		if !assert.Nil(t, err) {
			panic(err)
		}
		err = vp.WatchRemoteConfigOnChannel()
		if !assert.Nil(t, err) {
			panic(err)
		}
		vipers = append(vipers, vp)
	}

	tk := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-tk.C:
			vpMerge := viper.New()
			vpMerge.SetConfigType("yaml")
			for _, vp := range vipers {
				vpMerge.MergeConfigMap(vp.AllSettings())
			}
			fmt.Println(vpMerge)
		}
	}
}

func buildKey(group, dataId string) string {
	return fmt.Sprintf("%s/%s", group, dataId)
}
