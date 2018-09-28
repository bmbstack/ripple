package cache

import (
	"fmt"
	"strings"
)

const _VERSION = "0.1.0"

func Version() string {
	return _VERSION
}

var _ Cache = new(Engine)

// Cache is the interface that operates the cache data.
type CacheStore interface {
	// Put puts value into cache with key and expire time.
	Put(key, val string, timeout int64) error
	// Get gets cached value by given key.
	Get(key string) string
	// Delete deletes cached value by given key.
	Delete(key string) error
	// Delete deletes cached value by given prefix.
	DeleteByPrefix(prefix string) error
	// Incr increases cached int-type value by given key as a counter.
	Incr(key string) (int64, error)
	// Decr decreases cached int-type value by given key as a counter.
	Decr(key string) (int64, error)
	// IsExist returns true if cached value exists.
	IsExist(key string) bool
	// Flush deletes all cached data.
	Flush() error
	// StartAndGC starts GC routine based on config string settings.
	StartAndGC(opt Options) error
	// update expire time
	Touch(key string, expire int64) error
}

type Cache interface {
	CacheStore
	Tags(tags []string) Cache
}

type Options struct {
	// Name of adapter. Default is "memory".
	Adapter string
	// Adapter configuration, it's corresponding to adapter.
	AdapterConfig string
	// GC interval time in seconds. Default is 60.
	Interval int
	// key prefix Default is ""
	Section string
}

func prepareOptions(options []Options) Options {
	var opt Options
	if len(options) > 0 {
		opt = options[0]
	}

	// if len(opt.Section) == 0 {
	// 	opt.Section = ""
	// }

	if len(opt.Adapter) == 0 {
		opt.Adapter = "memory"
	}
	if opt.Interval == 0 {
		opt.Interval = 60
	}

	return opt
}

func New(alias string, options ...Options) (Cache, error) {
	opt := prepareOptions(options)
	adapterKey := fmt.Sprintf("%s_%s", opt.Adapter, alias)
	if strings.EqualFold(opt.Adapter, "redis") {
		Register(adapterKey, &RedisCache{})
	} else if strings.EqualFold(opt.Adapter, "memory") {
		Register(adapterKey, NewMemoryCacher())
	} else {
		return nil, fmt.Errorf("cache: unknown adapter type: %s", opt.Adapter)
	}

	adapter, ok := adapters[adapterKey]
	if !ok {
		return nil, fmt.Errorf("cache: unknown adapter '%s'(forgot to import?)", opt.Adapter)
	}

	engine := &Engine{}
	engine.Opt = opt
	engine.store = adapter

	return engine, adapter.StartAndGC(opt)
}

type Engine struct {
	Opt   Options
	store CacheStore
}

func (this *Engine) Put(key, val string, timeout int64) error {
	return this.store.Put(key, val, timeout)
}

func (this *Engine) Get(key string) string {
	return this.store.Get(key)
}

func (this *Engine) Delete(key string) error {
	return this.store.Delete(key)
}

func (this *Engine) DeleteByPrefix(prefix string) error {
	return this.store.DeleteByPrefix(prefix)
}

func (this *Engine) Incr(key string) (int64, error) {
	return this.store.Incr(key)
}

func (this *Engine) Decr(key string) (int64, error) {
	return this.store.Decr(key)
}

func (this *Engine) IsExist(key string) bool {
	return this.store.IsExist(key)
}

func (this *Engine) Flush() error {
	return this.store.Flush()
}

func (this *Engine) StartAndGC(opt Options) error {
	return this.store.StartAndGC(opt)
}

func (this *Engine) Touch(key string, expire int64) error {
	return this.store.Touch(key, expire)
}

func (this *Engine) Tags(tags []string) Cache {
	return NewTagCache(this.store, tags...)
}

var adapters = make(map[string]CacheStore)

// Register registers a adapter.
func Register(name string, adapter CacheStore) {
	if adapter == nil {
		panic("cache: cannot register adapter with nil value")
	}
	if _, dup := adapters[name]; dup {
		return
	}
	adapters[name] = adapter
}
