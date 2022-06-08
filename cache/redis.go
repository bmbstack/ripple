package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/gommon/color"
	"time"
)

var clients = make(map[string]*redis.Client)
var ctx = context.Background()

type RedisCache struct {
	alias     string
	prefix    string
	hasPrefix bool
}

// NewRedisCache creates and returns a new redis cache.
func NewRedisCache() *RedisCache {
	return &RedisCache{}
}

func (r *RedisCache) connect(opt Options) error {
	r.hasPrefix = len(opt.Section) > 0
	r.prefix = opt.Section
	r.alias = opt.Alias

	nc := redis.NewClient(&redis.Options{
		Addr:     opt.AdapterConfig.Addr,
		Password: opt.AdapterConfig.Password,
		DB:       opt.AdapterConfig.DB,
	})
	clients[opt.Alias] = nc

	fmt.Println(fmt.Sprintf("%s: %s, %s, db: %d", color.Green("Connect.redis"), opt.Section, opt.AdapterConfig.Addr, opt.AdapterConfig.DB))
	return nil
}

func (r *RedisCache) client() *redis.Client {
	if _, ok := clients[r.alias]; !ok {
		panic(fmt.Errorf("GetCache: cannot get cache alias '%s'", r.alias))
	}
	return clients[r.alias]
}

func (r *RedisCache) Client() interface{} {
	return r.client()
}

func (r *RedisCache) Key(key string) string {
	if r.hasPrefix {
		return r.prefix + ":" + key
	}
	return key
}

func (r *RedisCache) Set(key, val string, expiration time.Duration) {
	err := r.client().Set(ctx, r.Key(key), val, expiration).Err()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.Set error: %s", err.Error()))
	}
}

func (r *RedisCache) Get(key string) string {
	result, err := r.client().Get(ctx, r.Key(key)).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.Get error: %s", err.Error()))
		return ""
	}
	return result
}

func (r *RedisCache) Delete(key string) {
	err := r.client().Del(ctx, r.Key(key)).Err()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.Delete error: %s", err.Error()))
	}
}

// DeleteByPrefix Delete deletes cached value by given prefix key.
func (r *RedisCache) DeleteByPrefix(prefix string) {
	iter := r.client().Scan(ctx, 0, r.Key(prefix)+"*", 0).Iterator()
	for iter.Next(ctx) {
		err := r.client().Del(ctx, iter.Val()).Err()
		if err != nil {
			fmt.Println(fmt.Sprintf("Redis.DeleteByPrefix error: %s", err.Error()))
		}
	}
	if err := iter.Err(); err != nil {
		fmt.Println(fmt.Sprintf("Redis.DeleteByPrefix Iterator error: %s", err.Error()))
	}
}

func (r *RedisCache) Incr(key string) int64 {
	result, err := r.client().Incr(ctx, r.Key(key)).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.Incr error: %s", err.Error()))
		return 0
	}
	return result
}

func (r *RedisCache) Decr(key string) int64 {
	result, err := r.client().Decr(ctx, r.Key(key)).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.Decr error: %s", err.Error()))
		return 0
	}
	return result
}

func (r *RedisCache) IsExist(key string) bool {
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

func (r *RedisCache) Touch(key string) {
	err := r.client().Touch(ctx, r.Key(key)).Err()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.Touch error: %s", err.Error()))
	}
}

func (r *RedisCache) Flush() {
	err := r.client().FlushAll(ctx).Err()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.Flush error: %s", err.Error()))
	}
}

func (r *RedisCache) HGet(key, field string) string {
	result, err := r.client().HGet(context.Background(), r.Key(key), field).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.HGet error: %s", err.Error()))
		return ""
	}
	return result
}

func (r *RedisCache) HGetAll(key string) map[string]string {
	result, err := r.client().HGetAll(context.Background(), r.Key(key)).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.HGetAll error: %s", err.Error()))
		return map[string]string{}
	}
	return result
}

func (r *RedisCache) HSet(key string, values ...interface{}) int64 {
	result, err := r.client().HSet(context.Background(), r.Key(key), values...).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.HSet error: %s", err.Error()))
		return 0
	}
	return result
}

func (r *RedisCache) HExists(key, field string) bool {
	result, err := r.client().HExists(context.Background(), r.Key(key), field).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.HExists error: %s", err.Error()))
		return false
	}
	return result
}

func (r *RedisCache) HMGet(key string, fields ...string) []interface{} {
	result, err := r.client().HMGet(context.Background(), r.Key(key), fields...).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.HMGet error: %s", err.Error()))
		return []interface{}{}
	}
	return result
}

func (r *RedisCache) HMSet(key string, values ...interface{}) bool {
	result, err := r.client().HMSet(context.Background(), r.Key(key), values...).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.HMSet error: %s", err.Error()))
		return false
	}
	return result
}

func (r *RedisCache) SCard(key string) int64 {
	result, err := r.client().SCard(context.Background(), r.Key(key)).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.SCard error: %s", err.Error()))
		return 0
	}
	return result
}

func (r *RedisCache) SAdd(key string, members ...interface{}) int64 {
	result, err := r.client().SAdd(context.Background(), r.Key(key), members...).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.SAdd error: %s", err.Error()))
		return 0
	}
	return result
}

func (r *RedisCache) SRem(key string, members ...interface{}) int64 {
	result, err := r.client().SRem(context.Background(), r.Key(key), members...).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.SRem error: %s", err.Error()))
		return 0
	}
	return result
}

func (r *RedisCache) LRange(key string, start, stop int64) []string {
	result, err := r.client().LRange(context.Background(), r.Key(key), start, stop).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.Type error: %s", err.Error()))
		return []string{}
	}
	return result
}

func (r *RedisCache) LPush(key string, values ...interface{}) int64 {
	result, err := r.client().LPush(context.Background(), r.Key(key), values...).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.LPush error: %s", err.Error()))
		return 0
	}
	return result
}

func (r *RedisCache) LLen(key string) int64 {
	result, err := r.client().LLen(context.Background(), r.Key(key)).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.LLen error: %s", err.Error()))
		return 0
	}
	return result
}

func (r *RedisCache) ZAdd(key string, members ...*redis.Z) int64 {
	result, err := r.client().ZAdd(context.Background(), r.Key(key), members...).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.ZAdd error: %s", err.Error()))
		return 0
	}
	return result
}

func (r *RedisCache) ZRange(key string, start, stop int64) []string {
	result, err := r.client().ZRange(context.Background(), r.Key(key), start, stop).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.ZRange error: %s", err.Error()))
		return []string{}
	}
	return result
}

func (r *RedisCache) ZRank(key, member string) int64 {
	result, err := r.client().ZRank(context.Background(), r.Key(key), member).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.ZRank error: %s", err.Error()))
		return 0
	}
	return result
}

func (r *RedisCache) ZScore(key, member string) float64 {
	result, err := r.client().ZScore(context.Background(), r.Key(key), member).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.ZScore error: %s", err.Error()))
		return 0
	}
	return result
}

func (r *RedisCache) ZRem(key string, members ...interface{}) int64 {
	result, err := r.client().ZRem(context.Background(), r.Key(key), members...).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.ZRem error: %s", err.Error()))
		return 0
	}
	return result
}

func (r *RedisCache) SetEX(key string, value interface{}, expiration time.Duration) string {
	result, err := r.client().SetEX(context.Background(), r.Key(key), value, expiration).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.SetEX error: %s", err.Error()))
		return ""
	}
	return result
}

func (r *RedisCache) SetNX(key string, value interface{}, expiration time.Duration) bool {
	result, err := r.client().SetNX(context.Background(), r.Key(key), value, expiration).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.SetNX error: %s", err.Error()))
		return false
	}
	return result
}

func (r *RedisCache) HSetNX(key, field string, value interface{}) bool {
	result, err := r.client().HSetNX(context.Background(), r.Key(key), field, value).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.HSetNX error: %s", err.Error()))
		return false
	}
	return result
}

func (r *RedisCache) SScan(key string, cursor uint64, match string, count int64) ([]string, uint64) {
	keys, cursor, err := r.client().SScan(context.Background(), r.Key(key), cursor, match, count).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.SScan error: %s", err.Error()))
		return []string{}, 0
	}
	return keys, cursor
}

func (r *RedisCache) Del(keys ...string) int64 {
	var keyArray []string
	for _, key := range keys {
		keyArray = append(keyArray, r.Key(key))
	}
	result, err := r.client().Del(context.Background(), keyArray...).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.Del error: %s", err.Error()))
		return 0
	}
	return result
}

func (r *RedisCache) HDel(key string, fields ...string) int64 {
	result, err := r.client().HDel(context.Background(), r.Key(key), fields...).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.HDel error: %s", err.Error()))
		return 0
	}
	return result
}

func (r *RedisCache) Type(key string) string {
	result, err := r.client().Type(context.Background(), r.Key(key)).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.Type error: %s", err.Error()))
		return ""
	}
	return result
}

func (r *RedisCache) TTL(key string) time.Duration {
	result, err := r.client().TTL(context.Background(), r.Key(key)).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.TTL error: %s", err.Error()))
		return 0
	}
	return result
}

func (r *RedisCache) Expire(key string, expiration time.Duration) bool {
	result, err := r.client().Expire(context.Background(), r.Key(key), expiration).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.Expire error: %s", err.Error()))
		return false
	}
	return result
}

func (r *RedisCache) Exists(keys ...string) int64 {
	var keyArray []string
	for _, key := range keys {
		keyArray = append(keyArray, r.Key(key))
	}
	result, err := r.client().Exists(context.Background(), keyArray...).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis.Exists error: %s", err.Error()))
		return 0
	}
	return result
}
