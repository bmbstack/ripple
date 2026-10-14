package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/gommon/color"
)

var clients = make(map[string]*redis.Client)
var ctx = context.Background()

// RedisCache implements the Store interface backed by redis.
type RedisCache struct {
	alias     string
	prefix    string
	hasPrefix bool
}

// NewRedisCache creates and returns a new redis cache.
func NewRedisCache() *RedisCache {
	return &RedisCache{}
}

func (r *RedisCache) connect(opt Options) (err error) {
	r.hasPrefix = len(opt.Section) > 0
	r.prefix = opt.Section
	r.alias = opt.Alias

	var redisOpt *redis.Options
	if opt.AdapterConfig.Url != "" {
		redisOpt, err = redis.ParseURL(opt.AdapterConfig.Url)
		if err != nil {
			return err
		}
	} else {
		redisOpt = &redis.Options{
			Addr:         opt.AdapterConfig.Addr,
			Password:     opt.AdapterConfig.Password,
			DB:           opt.AdapterConfig.DB,
			PoolFIFO:     opt.AdapterConfig.PoolFIFO,
			PoolSize:     opt.AdapterConfig.PoolSize,
			MinIdleConns: opt.AdapterConfig.MinIdleConns,
		}
	}

	nc := redis.NewClient(redisOpt)
	clients[opt.Alias] = nc

	fmt.Println(fmt.Sprintf("%s: %s, %s, db: %d, poolSize=%d, MinIdleConns=%d, poolFIFO=%v", color.Green("Connect.redis"), opt.Section, redisOpt.Addr, redisOpt.DB, redisOpt.PoolSize, redisOpt.MinIdleConns, redisOpt.PoolFIFO))
	return nil
}

func (r *RedisCache) client() *redis.Client {
	if _, ok := clients[r.alias]; !ok {
		panic(fmt.Errorf("GetCache: cannot get cache alias '%s'", r.alias))
	}
	return clients[r.alias]
}

func (r *RedisCache) Client() *redis.Client {
	return r.client()
}

func (r *RedisCache) Close() error {
	return r.client().Close()
}

func (r *RedisCache) Key(key string) string {
	if r.hasPrefix {
		return r.prefix + ":" + key
	}
	return key
}

func (r *RedisCache) Set(key, val string, expiration time.Duration) error {
	return r.client().Set(ctx, r.Key(key), val, expiration).Err()
}

func (r *RedisCache) Get(key string) (string, error) {
	return r.client().Get(ctx, r.Key(key)).Result()
}

func (r *RedisCache) Del(keys ...string) (int64, error) {
	var keyArray []string
	for _, key := range keys {
		keyArray = append(keyArray, r.Key(key))
	}
	return r.client().Del(context.Background(), keyArray...).Result()
}

func (r *RedisCache) Delete(key string) error {
	return r.client().Del(ctx, r.Key(key)).Err()
}

// DeleteByPrefix deletes cached values by given prefix key.
func (r *RedisCache) DeleteByPrefix(prefix string) error {
	iter := r.client().Scan(ctx, 0, r.Key(prefix)+"*", 0).Iterator()
	for iter.Next(ctx) {
		err := r.client().Del(ctx, iter.Val()).Err()
		if err != nil {
			fmt.Println(fmt.Sprintf("Redis.DeleteByPrefix error: %s", err.Error()))
		}
	}
	return iter.Err()
}

func (r *RedisCache) Expire(key string, expiration time.Duration) (bool, error) {
	return r.client().Expire(context.Background(), r.Key(key), expiration).Result()
}

func (r *RedisCache) TTL(key string) (time.Duration, error) {
	return r.client().TTL(context.Background(), r.Key(key)).Result()
}

func (r *RedisCache) Exists(keys ...string) (int64, error) {
	var keyArray []string
	for _, key := range keys {
		keyArray = append(keyArray, r.Key(key))
	}
	return r.client().Exists(context.Background(), keyArray...).Result()
}

func (r *RedisCache) Exist(key string) bool {
	result, err := r.client().Exists(ctx, r.Key(key)).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.IsExist error: %s", err.Error()))
		return false
	}
	if result == 1 {
		return true
	}
	return false
}

func (r *RedisCache) Incr(key string) (int64, error) {
	return r.client().Incr(ctx, r.Key(key)).Result()
}

func (r *RedisCache) Decr(key string) (int64, error) {
	return r.client().Decr(ctx, r.Key(key)).Result()
}
