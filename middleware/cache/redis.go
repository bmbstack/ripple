package cache

import (
	"encoding/json"
	"fmt"
	"time"

	redigo "github.com/garyburd/redigo/redis"
)

const GC_HASH_KEY = "TagCache:CacheGCKeys"

var _  CacheStore = new(RedisCache)


type RedisConfig struct {
	Addr        string
	Passwd      string
	SelectDB    int
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
	Wait        bool
	OccupyMode  bool // use whole db
}

func prepareConfig(conf RedisConfig) RedisConfig {
	if conf.MaxIdle == 0 {
		conf.MaxIdle = 10
	}
	if conf.MaxActive == 0 {
		conf.MaxActive = 10
	}

	if conf.IdleTimeout == 0 {
		conf.IdleTimeout = 60
	}

	return conf
}

type RedisCache struct {
	pool       *redigo.Pool
	occupyMode bool
	interval   int
	prefix     string
	hasPrefix  bool
}

func (r *RedisCache) key(key string) string {
	if r.hasPrefix {
		return r.prefix + ":" + key
	}

	return key
}

func (r *RedisCache) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	c := r.pool.Get()
	defer c.Close()
	return c.Do(commandName, args...)
}

func (r *RedisCache) Put(key, val string, timeout int64) (err error) {
	var timeoutUnix int64 = 0
	if timeout == 0 {
		_, err = r.do("SET", r.key(key), val)
	} else {
		_, err = r.do("SETEX", r.key(key), timeout, val)
		timeoutUnix = time.Now().Unix() + timeout
	}
	if err != nil {
		return
	}

	if r.occupyMode {
		return
	}

	_, err = r.do("HSET", r.key(GC_HASH_KEY), r.key(key), timeoutUnix)
	return
}

func (r *RedisCache) Get(key string) string {
	v, _ := redigo.String(r.do("GET", r.key(key)))
	return v
}

// Delete deletes cached value by given key.
func (r *RedisCache) Delete(key string) (err error) {
	if _, err = r.do("DEL", r.key(key)); err != nil {
		return
	}
	if r.occupyMode {
		return
	}

	_, err = r.do("HDEL", r.key(GC_HASH_KEY), r.key(key))
	return
}

// Incr increases cached int-type value by given key as a counter.
func (r *RedisCache) Incr(key string) (int64, error) {
	return redigo.Int64(r.do("INCR", key))
}

// Decr decreases cached int-type value by given key as a counter.
func (r *RedisCache) Decr(key string) (int64, error) {
	return redigo.Int64(r.do("DECR", key))
}

// IsExist returns true if cached value exists.
func (r *RedisCache) IsExist(key string) bool {
	v, err := redigo.Bool(r.do("EXISTS", key))
	if err != nil {
		return false
	}
	return v
}

// Flush deletes all cached data.
func (r *RedisCache) Flush() (err error) {
	if r.occupyMode {
		_, err = r.do("FLUSHDB")
		return
	}

	keys, err := redigo.MultiBulk(r.do("HKEYS", r.key(GC_HASH_KEY)))
	if err != nil {
		return
	}

	fmt.Println(keys)

	conn := r.pool.Get()
	defer conn.Close()

	_, err = conn.Do("DEL", keys...)

	return
}

func (r *RedisCache) startGC() {
	if r.occupyMode {
		return
	}

	kvs, err := redigo.Int64Map(r.do("HGETALL", r.key(GC_HASH_KEY)))
	if err != nil {
		return
	}

	nowUnix := time.Now().Unix()

	outKeys := make([]interface{}, 0)

	for k, v := range kvs {
		if v == 0 {
			continue
		}

		if v < nowUnix {
			outKeys = append(outKeys, k)
		}
	}

	if len(outKeys) > 0 {
		_, err = r.do("DEL", outKeys...)
		if err != nil {
			fmt.Println(err)
		}

		args := make([]interface{}, len(outKeys)+1)
		args[0] = r.key(GC_HASH_KEY)
		copy(args[1:], outKeys)
		_, err = r.do("HDEL", args...)
		if err != nil {
			fmt.Println(err)
		}
	}

	time.AfterFunc(time.Duration(r.interval)*time.Second, func() { r.startGC() })
}

// StartAndGC starts GC routine based on config string settings.
func (r *RedisCache) StartAndGC(opt Options) error {
	var conf RedisConfig
	err := json.Unmarshal([]byte(opt.AdapterConfig), &conf)
	if err != nil {
		return fmt.Errorf("RedisConfig parse err %v", err)
	}

	conf = prepareConfig(conf)

	r.occupyMode = conf.OccupyMode
	r.interval = opt.Interval
	r.hasPrefix = len(opt.Section) > 0
	r.prefix = opt.Section
	r.pool = newRedisPool(conf)

	conn := r.pool.Get()

	_, err = conn.Do("PING")
	if err != nil {
		return fmt.Errorf("redis conn err %v", err)
	}
	conn.Close()

	go r.startGC()

	return nil
}

func newRedisPool(conf RedisConfig) *redigo.Pool {
	return &redigo.Pool{
		MaxIdle:     conf.MaxIdle,
		IdleTimeout: time.Duration(conf.IdleTimeout) * time.Second,
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", conf.Addr)
			if err != nil {
				return nil, err
			}
			if len(conf.Passwd) > 0 {
				if _, err := c.Do("AUTH", conf.Passwd); err != nil {
					c.Close()
					return nil, err
				}
			}
			_, err = c.Do("SELECT", conf.SelectDB)
			if err != nil {
				c.Close()
				return nil, err
			}

			return c, err
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

// update expire time
func (r *RedisCache) Touch(key string, expire int64) (err error) {
	if _, err = r.do("EXPIRE", key, expire); err != nil {
		return
	}

	if r.occupyMode {
		return
	}

	_, err = r.do("HSET", r.key(GC_HASH_KEY), r.key(key), (time.Now().Unix() + expire))

	return nil
}
