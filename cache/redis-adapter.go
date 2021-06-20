package cache

import (
	"github.com/MagalixTechnologies/core/logger"
	redisCache "github.com/go-redis/cache"
	"github.com/go-redis/redis"
	cache "github.com/victorspringer/http-cache"
	"github.com/vmihailenco/msgpack"
	"time"
)

// Adapter is the memory adapter data structure.
type RedisAdapter struct {
	store *redisCache.Codec
}

// RedisOptions exports go-redis RedisOptions type.
type RedisOptions redis.Options

// Get implements the cache Adapter interface Get method.
func (a *RedisAdapter) Get(key uint64) ([]byte, bool) {

	var c []byte
	if err := a.store.Get(cache.KeyAsString(key), &c); err == nil {
		logger.Debugw("found a expected in the cache", "key", key)
		return c, true
	} else {
		if err.Error() != "cache: key is missing" {
			logger.Errorw("failed to  get the cache expected", "err", err, "key", key)
		}
	}

	return nil, false
}

// Set implements the cache Adapter interface Set method.
func (a *RedisAdapter) Set(key uint64, response []byte, expiration time.Time) {
	err := a.store.Set(&redisCache.Item{
		Key:        cache.KeyAsString(key),
		Object:     response,
		Expiration: expiration.Sub(time.Now()),
	})
	if err != nil {
		logger.Errorw("failed to cache the expected", "key", key)
	}
}

// Release implements the cache Adapter interface Release method.
func (a *RedisAdapter) Release(key uint64) {
	a.store.Delete(cache.KeyAsString(key))
}

// NewRedisAdapter initializes Redis adapter.
func NewRedisAdapter(opt *RedisOptions) cache.Adapter {
	ropt := redis.Options(*opt)
	return &RedisAdapter{
		&redisCache.Codec{
			Redis: redis.NewClient(&ropt),
			Marshal: func(v interface{}) ([]byte, error) {
				return msgpack.Marshal(v)

			},
			Unmarshal: func(b []byte, v interface{}) error {
				return msgpack.Unmarshal(b, v)
			},
		},
	}
}
