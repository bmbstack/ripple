package cache

import (
	"fmt"
	"strings"
	"time"
)

const VersionName = "0.1.0"

func Version() string {
	return VersionName
}

var stores = make(map[string]Store)

// Store Cache is the interface that operates the cache data.
type Store interface {
	// Connect based on config string settings.
	Connect(opt Options) error
	// Client get client
	Client() interface{}
	// Set puts value into cache with key and expire time.
	Set(key, val string, expiration time.Duration)
	// Get gets cached value by given key.
	Get(key string) string
	// Delete deletes cached value by given key.
	Delete(key string)
	// DeleteByPrefix deletes cached value by given prefix.
	DeleteByPrefix(prefix string)
	// Incr increases cached int-type value by given key as a counter.
	Incr(key string) int64
	// Decr decreases cached int-type value by given key as a counter.
	Decr(key string) int64
	// IsExist returns true if cached value exists.
	IsExist(key string) bool
	// Touch touch
	Touch(key string)
	// Flush deletes all cached data.
	Flush()
}

type Cache struct {
	store Store
	Opt   Options
}

type AdapterConfig struct {
	Addr     string
	Password string
}

type Options struct {
	// Name of adapter. Default is "redis".
	Adapter string
	// Adapter configuration, it's corresponding to adapter.
	AdapterConfig AdapterConfig
	// key prefix Default is ""
	Section string
}

func prepareOptions(options []Options) Options {
	var opt Options
	if len(options) > 0 {
		opt = options[0]
	}

	if len(opt.Adapter) == 0 {
		opt.Adapter = "redis"
	}
	return opt
}

func NewCache(alias string, options ...Options) (*Cache, error) {
	opt := prepareOptions(options)
	adapterKey := fmt.Sprintf("%s_%s", opt.Adapter, alias)
	if strings.EqualFold(opt.Adapter, "redis") {
		Register(adapterKey, NewRedisCache())
	} else if strings.EqualFold(opt.Adapter, "memory") {
		Register(adapterKey, NewMemoryCache())
	} else {
		return nil, fmt.Errorf("cache: unknown adapter type: %s", opt.Adapter)
	}

	store, ok := stores[adapterKey]
	if !ok {
		return nil, fmt.Errorf("cache: unknown adapter '%s'(forgot to import?)", opt.Adapter)
	}

	newCache := &Cache{}
	newCache.store = store
	newCache.Opt = opt
	return newCache, store.Connect(opt)
}

func (this *Cache) Connect(opt Options) error {
	return this.store.Connect(opt)
}

func (this *Cache) Set(key, val string, expiration time.Duration) {
	this.store.Set(key, val, expiration)
}

func (this *Cache) Get(key string) string {
	return this.store.Get(key)
}

func (this *Cache) Delete(key string) {
	this.store.Delete(key)
}

func (this *Cache) DeleteByPrefix(prefix string) {
	this.store.DeleteByPrefix(prefix)
}

func (this *Cache) Incr(key string) int64 {
	return this.store.Incr(key)
}

func (this *Cache) Decr(key string) int64 {
	return this.store.Decr(key)
}

func (this *Cache) IsExist(key string) bool {
	return this.store.IsExist(key)
}

func (this *Cache) Touch(key string) {
	this.store.Touch(key)
}

func (this *Cache) Flush() {
	this.store.Flush()
}

// Register registers a store.
func Register(name string, store Store) {
	if store == nil {
		panic("cache: cannot register store with nil value")
	}
	if _, dup := stores[name]; dup {
		return
	}
	stores[name] = store
}
