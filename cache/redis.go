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

func (r *RedisCache) Client() interface{} {
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

func (r *RedisCache) HGet(key, field string) (string, error) {
	return r.client().HGet(context.Background(), r.Key(key), field).Result()
}

func (r *RedisCache) HGetAll(key string) (map[string]string, error) {
	return r.client().HGetAll(context.Background(), r.Key(key)).Result()
}

func (r *RedisCache) HSet(key string, values ...interface{}) (int64, error) {
	return r.client().HSet(context.Background(), r.Key(key), values...).Result()
}

func (r *RedisCache) HExists(key, field string) (bool, error) {
	return r.client().HExists(context.Background(), r.Key(key), field).Result()
}

func (r *RedisCache) HMGet(key string, fields ...string) ([]interface{}, error) {
	return r.client().HMGet(context.Background(), r.Key(key), fields...).Result()
}

func (r *RedisCache) HMSet(key string, values ...interface{}) (bool, error) {
	return r.client().HMSet(context.Background(), r.Key(key), values...).Result()
}

func (r *RedisCache) SCard(key string) (int64, error) {
	return r.client().SCard(context.Background(), r.Key(key)).Result()
}

func (r *RedisCache) SAdd(key string, members ...interface{}) (int64, error) {
	return r.client().SAdd(context.Background(), r.Key(key), members...).Result()
}

func (r *RedisCache) SRem(key string, members ...interface{}) (int64, error) {
	return r.client().SRem(context.Background(), r.Key(key), members...).Result()
}

func (r *RedisCache) LRange(key string, start, stop int64) ([]string, error) {
	return r.client().LRange(context.Background(), r.Key(key), start, stop).Result()
}

func (r *RedisCache) LPush(key string, values ...interface{}) (int64, error) {
	return r.client().LPush(context.Background(), r.Key(key), values...).Result()
}

func (r *RedisCache) LLen(key string) (int64, error) {
	return r.client().LLen(context.Background(), r.Key(key)).Result()
}

func (r *RedisCache) RPop(key string) (string, error) {
	return r.client().RPop(context.Background(), r.Key(key)).Result()
}

func (r *RedisCache) RPopCount(key string, count int) ([]string, error) {
	return r.client().RPopCount(context.Background(), r.Key(key), count).Result()
}

func (r *RedisCache) RPopLPush(source, destination string) (string, error) {
	return r.client().RPopLPush(context.Background(), r.Key(source), r.Key(destination)).Result()
}

func (r *RedisCache) RPush(key string, values ...interface{}) (int64, error) {
	return r.client().RPush(context.Background(), r.Key(key), values...).Result()
}

func (r *RedisCache) RPushX(key string, values ...interface{}) (int64, error) {
	return r.client().RPushX(context.Background(), r.Key(key), values...).Result()
}

func (r *RedisCache) ZAdd(key string, members ...*redis.Z) (int64, error) {
	return r.client().ZAdd(context.Background(), r.Key(key), members...).Result()
}

func (r *RedisCache) ZRange(key string, start, stop int64) ([]string, error) {
	return r.client().ZRange(context.Background(), r.Key(key), start, stop).Result()
}

func (r *RedisCache) ZRank(key, member string) (int64, error) {
	return r.client().ZRank(context.Background(), r.Key(key), member).Result()
}

func (r *RedisCache) ZScore(key, member string) (float64, error) {
	return r.client().ZScore(context.Background(), r.Key(key), member).Result()
}

func (r *RedisCache) ZRem(key string, members ...interface{}) (int64, error) {
	return r.client().ZRem(context.Background(), r.Key(key), members...).Result()
}

func (r *RedisCache) SetEX(key string, value interface{}, expiration time.Duration) (string, error) {
	return r.client().SetEX(context.Background(), r.Key(key), value, expiration).Result()
}

func (r *RedisCache) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	return r.client().SetNX(context.Background(), r.Key(key), value, expiration).Result()
}

func (r *RedisCache) HSetNX(key, field string, value interface{}) (bool, error) {
	return r.client().HSetNX(context.Background(), r.Key(key), field, value).Result()
}

func (r *RedisCache) SScan(key string, cursor uint64, match string, count int64) (keys []string, cursorOut uint64, err error) {
	keys, cursorOut, err = r.client().SScan(context.Background(), r.Key(key), cursor, match, count).Result()
	return keys, cursorOut, err
}

func (r *RedisCache) RunScript(src string, keys []string, args ...interface{}) (interface{}, error) {
	script := redis.NewScript(src)
	return script.Run(context.Background(), r.client(), keys, args).Result()
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

// DeleteByPrefix Delete deletes cached value by given prefix key.
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

func (r *RedisCache) HDel(key string, fields ...string) (int64, error) {
	return r.client().HDel(context.Background(), r.Key(key), fields...).Result()
}

func (r *RedisCache) Incr(key string) (int64, error) {
	return r.client().Incr(ctx, r.Key(key)).Result()
}

func (r *RedisCache) Decr(key string) (int64, error) {
	return r.client().Decr(ctx, r.Key(key)).Result()
}

func (r *RedisCache) Type(key string) (string, error) {
	return r.client().Type(context.Background(), r.Key(key)).Result()
}

func (r *RedisCache) TTL(key string) (time.Duration, error) {
	return r.client().TTL(context.Background(), r.Key(key)).Result()
}

func (r *RedisCache) Expire(key string, expiration time.Duration) (bool, error) {
	return r.client().Expire(context.Background(), r.Key(key), expiration).Result()
}

func (r *RedisCache) Exists(keys ...string) (int64, error) {
	var keyArray []string
	for _, key := range keys {
		keyArray = append(keyArray, r.Key(key))
	}
	return r.client().Exists(context.Background(), keyArray...).Result()
}

func (r *RedisCache) GetBit(key string, offset int64) (int64, error) {
	return r.client().GetBit(context.Background(), r.Key(key), offset).Result()
}

func (r *RedisCache) SetBit(key string, offset int64, value int) (int64, error) {
	return r.client().SetBit(context.Background(), r.Key(key), offset, value).Result()
}

func (r *RedisCache) BitCount(key string, bitCount *redis.BitCount) (int64, error) {
	return r.client().BitCount(context.Background(), r.Key(key), bitCount).Result()
}

func (r *RedisCache) BitOpAnd(destKey string, keys ...string) (int64, error) {
	var keyArray []string
	for _, key := range keys {
		keyArray = append(keyArray, r.Key(key))
	}
	return r.client().BitOpAnd(context.Background(), r.Key(destKey), keyArray...).Result()
}

func (r *RedisCache) BitOpOr(destKey string, keys ...string) (int64, error) {
	var keyArray []string
	for _, key := range keys {
		keyArray = append(keyArray, r.Key(key))
	}
	return r.client().BitOpOr(context.Background(), r.Key(destKey), keyArray...).Result()
}

func (r *RedisCache) BitOpXor(destKey string, keys ...string) (int64, error) {
	var keyArray []string
	for _, key := range keys {
		keyArray = append(keyArray, r.Key(key))
	}
	return r.client().BitOpXor(context.Background(), r.Key(destKey), keyArray...).Result()
}

func (r *RedisCache) BitOpNot(destKey string, key string) (int64, error) {
	return r.client().BitOpNot(context.Background(), r.Key(destKey), r.Key(key)).Result()
}

func (r *RedisCache) BitPos(key string, bit int64, pos ...int64) (int64, error) {
	return r.client().BitPos(context.Background(), r.Key(key), bit, pos...).Result()
}

func (r *RedisCache) BitField(key string, args ...interface{}) ([]int64, error) {
	return r.client().BitField(context.Background(), r.Key(key), args...).Result()
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

func (r *RedisCache) Touch(key string) error {
	return r.client().Touch(ctx, r.Key(key)).Err()
}

func (r *RedisCache) Flush() error {
	return r.client().FlushAll(ctx).Err()
}
