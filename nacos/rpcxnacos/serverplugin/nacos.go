package serverplugin

import (
	"errors"
	clients2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/clients"
	naming_client2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/clients/naming_client"
	constant2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/constant"
	vo2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/vo"
	"strings"

	"github.com/smallnest/rpcx/log"
	"github.com/smallnest/rpcx/util"
)

// NacosRegisterPlugin implements consul registry.
type NacosRegisterPlugin struct {
	// service address, for example, tcp@127.0.0.1:8972, quic@127.0.0.1:1234
	ServiceAddress string
	// nacos client config
	ClientConfig constant2.ClientConfig
	// nacos server config
	ServerConfig []constant2.ServerConfig
	Cluster      string
	Group        string
	Weight       float64

	// Registered services
	Services []string

	namingClient naming_client2.INamingClient
}

// Start starts to connect consul cluster
func (p *NacosRegisterPlugin) Start() error {
	namingClient, err := clients2.CreateNamingClient(map[string]interface{}{
		"clientConfig":  p.ClientConfig,
		"serverConfigs": p.ServerConfig,
	})
	if err != nil {
		return err
	}

	p.namingClient = namingClient

	return nil
}

// Stop unregister all services.
func (p *NacosRegisterPlugin) Stop() error {
	_, ip, port, _ := util.ParseRpcxAddress(p.ServiceAddress)

	for _, name := range p.Services {
		inst := vo2.DeregisterInstanceParam{
			Ip:          ip,
			Ephemeral:   true,
			Port:        uint64(port),
			ServiceName: name,
			Cluster:     p.Cluster,
			GroupName:   p.Group,
		}
		_, err := p.namingClient.DeregisterInstance(inst)
		if err != nil {
			log.Errorf("faield to deregister %s: %v", name, err)
		}
	}

	return nil
}

// Register handles registering event.
// this service is registered at BASE/serviceName/thisIpAddress node
func (p *NacosRegisterPlugin) Register(name string, rcvr interface{}, metadata string) (err error) {
	if strings.TrimSpace(name) == "" {
		return errors.New("Register service `name` can't be empty")
	}

	network, ip, port, err := util.ParseRpcxAddress(p.ServiceAddress)
	if err != nil {
		log.Errorf("failed to parse rpcx addr in Register: %v", err)
		return err
	}

	meta := util.ConvertMeta2Map(metadata)
	meta["network"] = network

	inst := vo2.RegisterInstanceParam{
		Ip:          ip,
		Port:        uint64(port),
		ServiceName: name,
		Metadata:    meta,
		ClusterName: p.Cluster,
		GroupName:   p.Group,
		Weight:      p.Weight,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	}

	_, err = p.namingClient.RegisterInstance(inst)
	if err != nil {
		log.Errorf("failed to register %s: %v", name, err)
		return err
	}

	p.Services = append(p.Services, name)

	return
}

func (p *NacosRegisterPlugin) RegisterFunction(serviceName, fname string, fn interface{}, metadata string) error {
	return p.Register(serviceName, fn, metadata)
}

func (p *NacosRegisterPlugin) Unregister(name string) (err error) {
	if len(p.Services) == 0 {
		return nil
	}

	if strings.TrimSpace(name) == "" {
		return errors.New("Unregister service `name` can't be empty")
	}

	_, ip, port, err := util.ParseRpcxAddress(p.ServiceAddress)
	if err != nil {
		log.Errorf("wrong address %s: %v", p.ServiceAddress, err)
		return err
	}

	inst := vo2.DeregisterInstanceParam{
		Ip:          ip,
		Ephemeral:   true,
		Port:        uint64(port),
		ServiceName: name,
		Cluster:     p.Cluster,
		GroupName:   p.Group,
	}
	_, err = p.namingClient.DeregisterInstance(inst)
	if err != nil {
		log.Errorf("failed to deregister %s: %v", name, err)
		return err
	}

	services := make([]string, 0, len(p.Services)-1)
	for _, s := range p.Services {
		if s != name {
			services = append(services, s)
		}
	}
	p.Services = services

	return nil
}