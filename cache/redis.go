package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/gommon/color"
	"time"
)

var client *redis.Client
var ctx = context.Background()

type RedisCache struct {
	prefix    string
	hasPrefix bool
}

// NewRedisCache creates and returns a new redis cache.
func NewRedisCache() *RedisCache {
	return &RedisCache{}
}

func (r *RedisCache) Connect(opt Options) error {
	r.hasPrefix = len(opt.Section) > 0
	r.prefix = opt.Section

	client = redis.NewClient(&redis.Options{
		Addr:     opt.AdapterConfig.Addr,
		Password: opt.AdapterConfig.Password,
	})

	fmt.Println(fmt.Sprintf("%s: %s, %s", color.Green("Connect.redis"), opt.Section, opt.AdapterConfig.Addr))
	return nil
}

func (r *RedisCache) Client() interface{} {
	return client
}

func (r *RedisCache) Key(key string) string {
	if r.hasPrefix {
		return r.prefix + ":" + key
	}
	return key
}

func (r *RedisCache) Set(key, val string, expiration time.Duration) {
	err := client.Set(ctx, r.Key(key), val, expiration).Err()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.Set error: %s", err.Error()))
	}
}

func (r *RedisCache) Get(key string) string {
	result, err := client.Get(ctx, r.Key(key)).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.Get error: %s", err.Error()))
		return ""
	}
	return result
}

// Delete deletes cached value by given key.
func (r *RedisCache) Delete(key string) {
	err := client.Del(ctx, r.Key(key)).Err()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.Delete error: %s", err.Error()))
	}
}

// Delete deletes cached value by given prefix key.
func (r *RedisCache) DeleteByPrefix(prefix string) {
	iter := client.Scan(ctx, 0, r.Key(prefix)+"*", 0).Iterator()
	for iter.Next(ctx) {
		err := client.Del(ctx, iter.Val()).Err()
		if err != nil {
			fmt.Println(fmt.Sprintf("Redis.DeleteByPrefix error: %s", err.Error()))
		}
	}
	if err := iter.Err(); err != nil {
		fmt.Println(fmt.Sprintf("Redis.DeleteByPrefix Iterator error: %s", err.Error()))
	}
}

// Incr increases cached int-type value by given key as a counter.
func (r *RedisCache) Incr(key string) int64 {
	result, err := client.Incr(ctx, r.Key(key)).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.Incr error: %s", err.Error()))
		return 0
	}
	return result
}

// Decr decreases cached int-type value by given key as a counter.
func (r *RedisCache) Decr(key string) int64 {
	result, err := client.Decr(ctx, r.Key(key)).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.Decr error: %s", err.Error()))
		return 0
	}
	return result
}

// IsExist returns true if cached value exists.
func (r *RedisCache) IsExist(key string) bool {
	result, err := client.Exists(ctx, r.Key(key)).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.IsExist error: %s", err.Error()))
		return false
	}
	if result == 1 {
		return true
	}
	return false
}

// update expire time
func (r *RedisCache) Touch(key string) {
	err := client.Touch(ctx, r.Key(key)).Err()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.Touch error: %s", err.Error()))
	}
}

// Flush deletes all cached data.
func (r *RedisCache) Flush() {
	err := client.FlushAll(ctx).Err()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.Flush error: %s", err.Error()))
	}
}
