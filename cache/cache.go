package cache

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
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
	Close() error

	Key(key string) string

	Set(key, val string, expiration time.Duration) error
	Get(key string) (string, error)

	HGet(key, field string) (string, error)
	HGetAll(key string) (map[string]string, error)
	HSet(key string, values ...interface{}) (int64, error)
	HExists(key, field string) (bool, error)

	HMGet(key string, fields ...string) ([]interface{}, error)
	HMSet(key string, values ...interface{}) (bool, error)

	SCard(key string) (int64, error)
	SAdd(key string, members ...interface{}) (int64, error)
	SRem(key string, members ...interface{}) (int64, error)

	LRange(key string, start, stop int64) ([]string, error)
	LPush(key string, values ...interface{}) (int64, error)
	LLen(key string) (int64, error)
	RPop(key string) (string, error)
	RPopCount(key string, count int) ([]string, error)
	RPopLPush(source, destination string) (string, error)
	RPush(key string, values ...interface{}) (int64, error)
	RPushX(key string, values ...interface{}) (int64, error)

	ZAdd(key string, members ...*redis.Z) (int64, error)
	ZRange(key string, start, stop int64) ([]string, error)
	ZRank(key, member string) (int64, error)
	ZScore(key, member string) (float64, error)
	ZRem(key string, members ...interface{}) (int64, error)

	SetEX(key string, value interface{}, expiration time.Duration) (string, error)
	SetNX(key string, value interface{}, expiration time.Duration) (bool, error)
	HSetNX(key, field string, value interface{}) (bool, error)
	SScan(key string, cursor uint64, match string, count int64) ([]string, uint64, error)
	RunScript(src string, keys []string, args ...interface{}) (interface{}, error)

	GetBit(key string, offset int64) (int64, error)
	SetBit(key string, offset int64, value int) (int64, error)
	BitCount(key string, bitCount *redis.BitCount) (int64, error)
	BitOpAnd(destKey string, keys ...string) (int64, error)
	BitOpOr(destKey string, keys ...string) (int64, error)
	BitOpXor(destKey string, keys ...string) (int64, error)
	BitOpNot(destKey string, key string) (int64, error)
	BitPos(key string, bit int64, pos ...int64) (int64, error)
	BitField(key string, args ...interface{}) ([]int64, error)

	Del(keys ...string) (int64, error)
	Delete(key string) error
	DeleteByPrefix(prefix string) error
	HDel(key string, fields ...string) (int64, error)

	Incr(key string) (int64, error)
	Decr(key string) (int64, error)

	Type(key string) (string, error)
	TTL(key string) (time.Duration, error)
	Expire(key string, expiration time.Duration) (bool, error)
	Exists(keys ...string) (int64, error)

	Exist(key string) bool
	Touch(key string) error
	Flush() error
}

type Cache struct {
	store Store
	Opt   Options
}

type AdapterConfig struct {
	Addr     string
	Password string
	DB       int

	// Type of connection pool.
	// true for FIFO pool, false for LIFO pool.
	// Note that fifo has higher overhead compared to lifo.
	PoolFIFO bool
	// Maximum number of socket connections.
	// Default is 10 connections per every available CPU as reported by runtime.GOMAXPROCS.
	PoolSize int
	// Minimum number of idle connections which is useful when establishing
	// new connection is slow.
	MinIdleConns int

	// redis://user:password@localhost:6789/3?dial_timeout=3&db=1&read_timeout=6s&max_retries=2
	// url, 如果使用url，其他配置无效, url的优先级高于其他配置
	Url string
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

func (this *Cache) Close() error {
	return this.store.Close()
}

func (this *Cache) Key(key string) string {
	return this.store.Key(key)
}

func (this *Cache) Set(key, val string, expiration time.Duration) error {
	return this.store.Set(key, val, expiration)
}

func (this *Cache) Get(key string) (string, error) {
	return this.store.Get(key)
}

func (this *Cache) HGet(key, field string) (string, error) {
	return this.store.HGet(key, field)
}

func (this *Cache) HGetAll(key string) (map[string]string, error) {
	return this.store.HGetAll(key)
}

func (this *Cache) HSet(key string, values ...interface{}) (int64, error) {
	return this.store.HSet(key, values...)
}

func (this *Cache) HExists(key, field string) (bool, error) {
	return this.store.HExists(key, field)
}

func (this *Cache) HMGet(key string, fields ...string) ([]interface{}, error) {
	return this.store.HMGet(key, fields...)
}

func (this *Cache) HMSet(key string, values ...interface{}) (bool, error) {
	return this.store.HMSet(key, values...)
}

func (this *Cache) SCard(key string) (int64, error) {
	return this.store.SCard(key)
}

func (this *Cache) SAdd(key string, members ...interface{}) (int64, error) {
	return this.store.SAdd(key, members...)
}

func (this *Cache) SRem(key string, members ...interface{}) (int64, error) {
	return this.store.SRem(key, members)
}

func (this *Cache) LRange(key string, start, stop int64) ([]string, error) {
	return this.store.LRange(key, start, stop)
}

func (this *Cache) LPush(key string, values ...interface{}) (int64, error) {
	return this.store.LPush(key, values...)
}

func (this *Cache) LLen(key string) (int64, error) {
	return this.store.LLen(key)
}

func (this *Cache) ZAdd(key string, members ...*redis.Z) (int64, error) {
	return this.store.ZAdd(key, members...)
}

func (this *Cache) ZRange(key string, start, stop int64) ([]string, error) {
	return this.store.ZRange(key, start, stop)
}

func (this *Cache) ZRank(key, member string) (int64, error) {
	return this.store.ZRank(key, member)
}

func (this *Cache) ZScore(key, member string) (float64, error) {
	return this.store.ZScore(key, member)
}

func (this *Cache) ZRem(key string, members ...interface{}) (int64, error) {
	return this.store.ZRem(key, members...)
}

func (this *Cache) SetEX(key string, value interface{}, expiration time.Duration) (string, error) {
	return this.store.SetEX(key, value, expiration)
}

func (this *Cache) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	return this.store.SetNX(key, value, expiration)
}

func (this *Cache) HSetNX(key, field string, value interface{}) (bool, error) {
	return this.store.HSetNX(key, field, value)
}

func (this *Cache) SScan(key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return this.store.SScan(key, cursor, match, count)
}

func (this *Cache) Del(keys ...string) (int64, error) {
	return this.store.Del(keys...)
}

func (this *Cache) Delete(key string) error {
	return this.store.Delete(key)
}

func (this *Cache) DeleteByPrefix(prefix string) error {
	return this.store.DeleteByPrefix(prefix)
}

func (this *Cache) HDel(key string, fields ...string) (int64, error) {
	return this.store.HDel(key, fields...)
}

func (this *Cache) Incr(key string) (int64, error) {
	return this.store.Incr(key)
}

func (this *Cache) Decr(key string) (int64, error) {
	return this.store.Decr(key)
}

func (this *Cache) Type(key string) (string, error) {
	return this.store.Type(key)
}

func (this *Cache) TTL(key string) (time.Duration, error) {
	return this.store.TTL(key)
}

func (this *Cache) Expire(key string, expiration time.Duration) (bool, error) {
	return this.store.Expire(key, expiration)
}

func (this *Cache) Exists(keys ...string) (int64, error) {
	return this.store.Exists(keys...)
}

func (this *Cache) Exist(key string) bool {
	return this.store.Exist(key)
}

func (this *Cache) Touch(key string) error {
	return this.store.Touch(key)
}

func (this *Cache) Flush() error {
	return this.store.Flush()
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
