package cache

import (
	"time"

	"github.com/go-redis/redis/v8"
)

const VersionName = "0.1.0"

func Version() string {
	return VersionName
}

// Store is a thin key-value cache abstraction.
// Redis-specific commands (hash/list/set/zset/bitmap/script) are intentionally
// NOT part of this interface; use Client() to access them directly.
type Store interface {
	Client() *redis.Client
	Close() error

	Key(key string) string

	Set(key, val string, expiration time.Duration) error
	Get(key string) (string, error)

	Del(keys ...string) (int64, error)
	Delete(key string) error
	DeleteByPrefix(prefix string) error

	Expire(key string, expiration time.Duration) (bool, error)
	TTL(key string) (time.Duration, error)
	Exists(keys ...string) (int64, error)
	Exist(key string) bool

	Incr(key string) (int64, error)
	Decr(key string) (int64, error)
}

// Cache is a redis-backed cache. It wraps a Store and delegates all operations.
type Cache struct {
	store Store
	Opt   Options
}

// AdapterConfig holds redis connection settings.
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

// Options configures a Cache instance.
type Options struct {
	// alias
	Alias string
	// Adapter configuration.
	AdapterConfig AdapterConfig
	// key prefix Default is ""
	Section string
}

// NewCache creates a new redis-backed cache.
func NewCache(options Options) (*Cache, error) {
	store := NewRedisCache()
	newCache := &Cache{store: store, Opt: options}
	return newCache, store.connect(options)
}

func (this *Cache) Client() *redis.Client {
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

func (this *Cache) Del(keys ...string) (int64, error) {
	return this.store.Del(keys...)
}

func (this *Cache) Delete(key string) error {
	return this.store.Delete(key)
}

func (this *Cache) DeleteByPrefix(prefix string) error {
	return this.store.DeleteByPrefix(prefix)
}

func (this *Cache) Expire(key string, expiration time.Duration) (bool, error) {
	return this.store.Expire(key, expiration)
}

func (this *Cache) TTL(key string) (time.Duration, error) {
	return this.store.TTL(key)
}

func (this *Cache) Exists(keys ...string) (int64, error) {
	return this.store.Exists(keys...)
}

func (this *Cache) Exist(key string) bool {
	return this.store.Exist(key)
}

func (this *Cache) Incr(key string) (int64, error) {
	return this.store.Incr(key)
}

func (this *Cache) Decr(key string) (int64, error) {
	return this.store.Decr(key)
}
