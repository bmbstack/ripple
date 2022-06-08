package cache

import (
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

// MemoryItem represents a memory cache item.
type MemoryItem struct {
	val        string
	created    int64
	expiration time.Duration
}

// MemoryCache represents a memory cache adapter implementation.
type MemoryCache struct {
	lock  sync.RWMutex
	items map[string]*MemoryItem
}

// NewMemoryCache creates and returns a new memory cache.
func NewMemoryCache() *MemoryCache {
	return &MemoryCache{items: make(map[string]*MemoryItem)}
}

// Forever put value into cache with key forever save
func (c *MemoryCache) Forever(key, val string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.items[key] = &MemoryItem{
		val:        val,
		created:    time.Now().Unix(),
		expiration: 0,
	}
	return nil
}

func (c *MemoryCache) connect(opt Options) error {
	return nil
}

func (r *MemoryCache) Client() interface{} {
	return nil
}

func (r *MemoryCache) Key(key string) string {
	return key
}

func (c *MemoryCache) Set(key, val string, expiration time.Duration) {
}

func (c *MemoryCache) Get(key string) string {
	return ""
}

func (c *MemoryCache) Delete(key string) {
}

// DeleteByPrefix Delete deletes cached value by given prefix.
func (c *MemoryCache) DeleteByPrefix(prefix string) {
}

// Incr increases cached int-type value by given key as a counter.
func (c *MemoryCache) Incr(key string) int64 {
	return 0
}

func (c *MemoryCache) Decr(key string) int64 {
	return 0
}

func (c *MemoryCache) IsExist(key string) bool {
	return false
}

func (c *MemoryCache) Touch(key string) {
}

func (c *MemoryCache) Flush() {
}

func (c *MemoryCache) HGet(key, field string) string {
	//TODO implement me
	panic("implement me")
}

func (c *MemoryCache) HGetAll(key string) map[string]string {
	//TODO implement me
	panic("implement me")
}

func (c *MemoryCache) HSet(key string, values ...interface{}) int64 {
	//TODO implement me
	panic("implement me")
}

func (c *MemoryCache) HExists(key, field string) bool {
	//TODO implement me
	panic("implement me")
}

func (c *MemoryCache) HMGet(key string, fields ...string) []interface{} {
	//TODO implement me
	panic("implement me")
}

func (c *MemoryCache) HMSet(key string, values ...interface{}) bool {
	//TODO implement me
	panic("implement me")
}

func (c *MemoryCache) SCard(key string) int64 {
	//TODO implement me
	panic("implement me")
}

func (c *MemoryCache) SAdd(key string, members ...interface{}) int64 {
	//TODO implement me
	panic("implement me")
}

func (c *MemoryCache) SRem(key string, members ...interface{}) int64 {
	//TODO implement me
	panic("implement me")
}

func (c *MemoryCache) LRange(key string, start, stop int64) []string {
	//TODO implement me
	panic("implement me")
}

func (c *MemoryCache) LPush(key string, values ...interface{}) int64 {
	//TODO implement me
	panic("implement me")
}

func (c *MemoryCache) LLen(key string) int64 {
	//TODO implement me
	panic("implement me")
}

func (c *MemoryCache) ZAdd(key string, members ...*redis.Z) int64 {
	//TODO implement me
	panic("implement me")
}

func (c *MemoryCache) ZRange(key string, start, stop int64) []string {
	//TODO implement me
	panic("implement me")
}

func (c *MemoryCache) ZRank(key, member string) int64 {
	//TODO implement me
	panic("implement me")
}

func (c *MemoryCache) ZScore(key, member string) float64 {
	//TODO implement me
	panic("implement me")
}

func (c *MemoryCache) ZRem(key string, members ...interface{}) int64 {
	//TODO implement me
	panic("implement me")
}

func (c *MemoryCache) SetEX(key string, value interface{}, expiration time.Duration) string {
	//TODO implement me
	panic("implement me")
}

func (c *MemoryCache) SetNX(key string, value interface{}, expiration time.Duration) bool {
	//TODO implement me
	panic("implement me")
}

func (c *MemoryCache) HSetNX(key, field string, value interface{}) bool {
	//TODO implement me
	panic("implement me")
}

func (c *MemoryCache) SScan(key string, cursor uint64, match string, count int64) ([]string, uint64) {
	//TODO implement me
	panic("implement me")
}

func (c *MemoryCache) Del(keys ...string) int64 {
	//TODO implement me
	panic("implement me")
}

func (c *MemoryCache) HDel(key string, fields ...string) int64 {
	//TODO implement me
	panic("implement me")
}

func (c *MemoryCache) Type(key string) string {
	//TODO implement me
	panic("implement me")
}

func (c *MemoryCache) TTL(key string) time.Duration {
	//TODO implement me
	panic("implement me")
}

func (c *MemoryCache) Expire(key string, expiration time.Duration) bool {
	//TODO implement me
	panic("implement me")
}

func (c *MemoryCache) Exists(keys ...string) int64 {
	//TODO implement me
	panic("implement me")
}
