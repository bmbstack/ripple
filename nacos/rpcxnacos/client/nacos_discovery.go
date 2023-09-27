package client

import (
	"fmt"
	clients2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/clients"
	naming_client2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/clients/naming_client"
	constant2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/constant"
	model2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/model"
	vo2 "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/vo"
	"sync"
	"time"

	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/log"
	"github.com/smallnest/rpcx/util"
)

// NacosDiscovery is a nacos service discovery.
// It always returns the registered servers in nacos.
type NacosDiscovery struct {
	servicePath string
	// nacos client config
	ClientConfig constant2.ClientConfig
	// nacos server config
	ServerConfig []constant2.ServerConfig
	Cluster      string
	Group        string

	namingClient naming_client2.INamingClient

	pairsMu sync.RWMutex
	pairs   []*client.KVPair
	chans   []chan []*client.KVPair
	mu      sync.Mutex

	filter                  client.ServiceDiscoveryFilter
	RetriesAfterWatchFailed int

	stopCh chan struct{}
}

// NewNacosDiscovery returns a new NacosDiscovery.
func NewNacosDiscovery(servicePath, cluster, group string, clientConfig constant2.ClientConfig, serverConfig []constant2.ServerConfig) (client.ServiceDiscovery, error) {
	d := &NacosDiscovery{
		servicePath:  servicePath,
		Cluster:      cluster,
		Group:        group,
		ClientConfig: clientConfig,
		ServerConfig: serverConfig,
		stopCh:       make(chan struct{}),
	}

	namingClient, err := clients2.CreateNamingClient(map[string]interface{}{
		"clientConfig":  clientConfig,
		"serverConfigs": serverConfig,
	})
	if err != nil {
		log.Errorf("failed to create NacosDiscovery: %v", err)
		return nil, err
	}

	d.namingClient = namingClient

	d.fetch()
	go d.watch()
	return d, nil
}

func (d *NacosDiscovery) fetch() {
	service, err := d.namingClient.GetService(vo2.GetServiceParam{
		ServiceName: d.servicePath,
		Clusters:    []string{d.Cluster},
		GroupName:   d.Group,
	})
	if err != nil {
		log.Errorf("failed to get service %s: %v", d.servicePath, err)
		return
	}
	pairs := make([]*client.KVPair, 0, len(service.Hosts))
	for _, inst := range service.Hosts {
		network := inst.Metadata["network"]
		ip := inst.Ip
		port := inst.Port
		key := fmt.Sprintf("%s@%s:%d", network, ip, port)
		pair := &client.KVPair{Key: key, Value: util.ConvertMap2String(inst.Metadata)}
		if d.filter != nil && !d.filter(pair) {
			continue
		}
		pairs = append(pairs, pair)
	}

	d.pairsMu.Lock()
	d.pairs = pairs
	d.pairsMu.Unlock()
}

// Clone clones this ServiceDiscovery with new servicePath.
func (d *NacosDiscovery) Clone(servicePath string) (client.ServiceDiscovery, error) {
	return NewNacosDiscovery(servicePath, d.Cluster, d.Group, d.ClientConfig, d.ServerConfig)
}

// SetFilter sets the filer.
func (d *NacosDiscovery) SetFilter(filter client.ServiceDiscoveryFilter) {
	d.filter = filter
}

// GetServices returns the servers
func (d *NacosDiscovery) GetServices() []*client.KVPair {
	d.pairsMu.RLock()
	defer d.pairsMu.RUnlock()

	return d.pairs
}

// WatchService returns a nil chan.
func (d *NacosDiscovery) WatchService() chan []*client.KVPair {
	d.mu.Lock()
	defer d.mu.Unlock()

	ch := make(chan []*client.KVPair, 10)
	d.chans = append(d.chans, ch)
	return ch
}

func (d *NacosDiscovery) RemoveWatcher(ch chan []*client.KVPair) {
	d.mu.Lock()
	defer d.mu.Unlock()

	var chans []chan []*client.KVPair
	for _, c := range d.chans {
		if c == ch {
			continue
		}

		chans = append(chans, c)
	}

	d.chans = chans
}

func (d *NacosDiscovery) watch() {
	param := &vo2.SubscribeParam{
		ServiceName: d.servicePath,
		Clusters:    []string{d.Cluster},
		GroupName:   d.Group,
		SubscribeCallback: func(services []model2.Instance, err error) {
			pairs := make([]*client.KVPair, 0, len(services))
			for _, inst := range services {
				network := inst.Metadata["network"]
				ip := inst.Ip
				port := inst.Port
				key := fmt.Sprintf("%s@%s:%d", network, ip, port)
				pair := &client.KVPair{Key: key, Value: util.ConvertMap2String(inst.Metadata)}
				if d.filter != nil && !d.filter(pair) {
					continue
				}
				pairs = append(pairs, pair)
			}
			d.pairsMu.Lock()
			d.pairs = pairs
			d.pairsMu.Unlock()

			d.mu.Lock()
			for _, ch := range d.chans {
				ch := ch
				go func() {
					defer func() {
						recover()
					}()
					select {
					case ch <- d.pairs:
					case <-time.After(time.Minute):
						log.Warn("chan is full and new change has been dropped")
					}
				}()
			}
			d.mu.Unlock()
		},
	}

	err := d.namingClient.Subscribe(param)
	// if failed to Subscribe, retry
	if err != nil {
		var tempDelay time.Duration
		retry := d.RetriesAfterWatchFailed
		for d.RetriesAfterWatchFailed < 0 || retry >= 0 {
			err := d.namingClient.Subscribe(param)
			if err != nil {
				if d.RetriesAfterWatchFailed > 0 {
					retry--
				}
				if tempDelay == 0 {
					tempDelay = 1 * time.Second
				} else {
					tempDelay *= 2
				}
				if max := 30 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Warnf("can not subscribe (with retry %d, sleep %v): %s: %v", retry, tempDelay, d.servicePath, err)
				time.Sleep(tempDelay)
				continue
			}
			break
		}
	}
}

func (d *NacosDiscovery) Close() {
	close(d.stopCh)
}
