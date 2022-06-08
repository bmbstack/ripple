package cache

import (
	"fmt"
	"github.com/go-redis/redis/v8"
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
	connect(opt Options) error

	Client() interface{}

	Key(key string) string
	Set(key, val string, expiration time.Duration)
	Get(key string) string

	HGet(key, field string) string
	HGetAll(key string) map[string]string
	HSet(key string, values ...interface{}) int64
	HExists(key, field string) bool

	HMGet(key string, fields ...string) []interface{}
	HMSet(key string, values ...interface{}) bool

	SCard(key string) int64
	SAdd(key string, members ...interface{}) int64
	SRem(key string, members ...interface{}) int64

	LRange(key string, start, stop int64) []string
	LPush(key string, values ...interface{}) int64
	LLen(key string) int64

	ZAdd(key string, members ...*redis.Z) int64
	ZRange(key string, start, stop int64) []string
	ZRank(key, member string) int64
	ZScore(key, member string) float64
	ZRem(key string, members ...interface{}) int64

	SetEX(key string, value interface{}, expiration time.Duration) string
	SetNX(key string, value interface{}, expiration time.Duration) bool
	HSetNX(key, field string, value interface{}) bool
	SScan(key string, cursor uint64, match string, count int64) ([]string, uint64)

	Del(keys ...string) int64
	Delete(key string)
	DeleteByPrefix(prefix string)
	HDel(key string, fields ...string) int64

	Incr(key string) int64
	Decr(key string) int64

	Type(key string) string
	TTL(key string) time.Duration
	Expire(key string, expiration time.Duration) bool
	Exists(keys ...string) int64

	IsExist(key string) bool
	Touch(key string)
	Flush()
}

type Cache struct {
	store Store
	Opt   Options
}

type AdapterConfig struct {
	Addr     string
	Password string
	DB       int
}

type Options struct {
	// alias
	Alias string
	// Name of adapter. Default is "redis".
	Adapter string
	// Adapter configuration, it's corresponding to adapter.
	AdapterConfig AdapterConfig
	// key prefix Default is ""
	Section string
}

func prepareOptions(opt Options) Options {
	if len(opt.Adapter) == 0 {
		opt.Adapter = "redis"
	}
	return opt
}

func NewCache(options Options) (*Cache, error) {
	opt := prepareOptions(options)
	adapterKey := fmt.Sprintf("%s_%s", opt.Adapter, opt.Alias)
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
	return newCache, store.connect(opt)
}

func (this *Cache) Client() interface{} {
	return this.store.Client()
}

func (this *Cache) Key(key string) string {
	return this.store.Key(key)
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
