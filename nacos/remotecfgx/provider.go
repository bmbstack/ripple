package remotecfgx

import (
	"bytes"
	"errors"
	"github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/vo"
	cryptconfig "github.com/sagikazarmark/crypt/config"
	"github.com/spf13/viper"
	"io"
	"os"
	"sync"
)

var (
	ErrUnsupportedProvider      = errors.New("This configuration manager is not supported")
	ErrUnsupportedEncrypt       = errors.New("The Nacos configuration manager is not support encrypted")
	ErrMissingRegistrationParam = errors.New("The Nacos configuration Missing registration parameters")
)

type remoteConfigProvider struct{}

func (rc remoteConfigProvider) Get(rp viper.RemoteProvider) (io.Reader, error) {
	cm, err := getConfigManager(rp)
	if err != nil {
		return nil, err
	}
	b, err := cm.Get(rp.Path())
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

func (rc remoteConfigProvider) Watch(rp viper.RemoteProvider) (io.Reader, error) {
	return rc.Get(rp)
}

func (rc remoteConfigProvider) WatchChannel(rp viper.RemoteProvider) (<-chan *viper.RemoteResponse, chan bool) {
	cm, err := getConfigManager(rp)
	if err != nil {
		return nil, nil
	}
	quit := make(chan bool)
	quitwc := make(chan bool)
	viperResponsCh := make(chan *viper.RemoteResponse)
	cryptoResponseCh := cm.Watch(rp.Path(), quit)
	// need this function to convert the Channel response form crypt.Response to viper.Response
	go func(cr <-chan *cryptconfig.Response, vr chan<- *viper.RemoteResponse, quitwc <-chan bool, quit chan<- bool) {
		for {
			select {
			case <-quitwc:
				quit <- true
				return
			case resp := <-cr:
				vr <- &viper.RemoteResponse{
					Error: resp.Error,
					Value: resp.Value,
				}
			}
		}
	}(cryptoResponseCh, viperResponsCh, quitwc, quit)

	return viperResponsCh, quitwc
}

func getConfigManager(rp viper.RemoteProvider) (cryptconfig.ConfigManager, error) {
	if rp.SecretKeyring() != "" {
		kr, err := os.Open(rp.SecretKeyring())
		if err != nil {
			return nil, err
		}
		defer func() {
			err = kr.Close()
			panic(err)
		}()

		switch rp.Provider() {
		case "etcd":
			return cryptconfig.NewEtcdConfigManager([]string{rp.Endpoint()}, kr)
		case "consul":
			return cryptconfig.NewConsulConfigManager([]string{rp.Endpoint()}, kr)
		case "nacos":
			return nil, ErrUnsupportedEncrypt
		default:
			return nil, ErrUnsupportedProvider
		}
	} else {
		switch rp.Provider() {
		case "etcd":
			return cryptconfig.NewStandardEtcdConfigManager([]string{rp.Endpoint()})
		case "consul":
			return cryptconfig.NewStandardConsulConfigManager([]string{rp.Endpoint()})
		case "nacos":
			return NewStandardNacosConfigManager(rp.Path())
		default:
			return nil, ErrUnsupportedProvider
		}
	}
}

var (
	configMap sync.Map
)

func RegisterConfig(path string, clientParam vo.NacosClientParam, configParam vo.ConfigParam) {
	configMap.Store(path, &remoteConfig{
		clientParam: clientParam,
		configParam: configParam,
	})
}

type remoteConfig struct {
	clientParam vo.NacosClientParam
	configParam vo.ConfigParam
}

func NewStandardNacosConfigManager(path string) (cryptconfig.ConfigManager, error) {
	config, ok := configMap.Load(path)
	if !ok {
		return nil, ErrMissingRegistrationParam
	}
	nConfig, ok := config.(*remoteConfig)
	if !ok {
		return nil, ErrMissingRegistrationParam
	}
	store, err := New(nConfig.clientParam, nConfig.configParam)
	if err != nil {
		return nil, err
	}
	return cryptconfig.NewStandardConfigManager(store)
}

func init() {
	viper.SupportedRemoteProviders = append(
		viper.SupportedRemoteProviders,
		"nacos",
	)
	viper.RemoteConfig = &remoteConfigProvider{}
}
