package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// Erros
var (
	ErrKeyNotFound = errors.New("data with a key not found")
)

// Cache is interface for caching basic operations
type Cache[T any] interface {
	// Set sets data under a key with expiration time
	// Errors: unknown
	Set(ctx context.Context, key string, data T, expiration time.Duration) error

	// Get gets data with a key
	// Errors: ErrCacheMissed, unknown
	Get(ctx context.Context, key string) (T, error)

	// Delete deletes data with keys
	// Errors: unknown
	Delete(ctx context.Context, keys ...string) error

	// Exists check existence of a key
	Exists(ctx context.Context, key string) bool
}

// CacheOpts represents options for redis based cache
type CacheOpts struct {
	Client *redis.Client

	// key prefix using for avoid collision with other caches
	KeyPrefix         string
	DefaultExpiration time.Duration
}

// cache represents redis based cache and implements Cache interface
type cache[T any] struct {
	client *redis.Client

	keyPrefix         string
	defaultExpiration time.Duration
}

// NewCacheByOpts creates new cache by opts
func NewCache[T any](opts *CacheOpts) Cache[T] {
	return &cache[T]{
		client:            opts.Client,
		keyPrefix:         opts.KeyPrefix,
		defaultExpiration: opts.DefaultExpiration,
	}
}

// Set sets data under a key, if expiration is 0, then using default expiration
func (c *cache[T]) Set(
	ctx context.Context,
	key string,
	data T,
	expiration time.Duration,
) (err error) {
	if c.keyPrefix != "" {
		key = c.keyPrefix + ":" + key
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if expiration == 0 {
		expiration = c.defaultExpiration
	}

	return c.client.Set(ctx, key, payload, expiration).Err()
}

// Get gets data with a key, if key doesn't exist return error
func (c *cache[T]) Get(ctx context.Context, key string) (value T, err error) {
	if c.keyPrefix != "" {
		key = c.keyPrefix + ":" + key
	}

	payload, err := c.client.Get(ctx, key).Result()
	if err != nil {
		switch {
		case errors.Is(err, redis.Nil):
			return value, ErrKeyNotFound
		default:
			return value, err
		}
	}

	err = json.Unmarshal([]byte(payload), &value)
	if err != nil {
		return value, err
	}

	return value, nil
}

// Delete deletes keys from the cache
func (c *cache[T]) Delete(ctx context.Context, keys ...string) error {
	prefKeys := make([]string, len(keys))
	copy(prefKeys, keys)

	if c.keyPrefix != "" {
		for i := 0; i < len(prefKeys); i++ {
			prefKeys[i] = c.keyPrefix + ":" + prefKeys[i]
		}
	}

	return c.client.Del(ctx, prefKeys...).Err()
}

// Exists check existence of a key
func (c *cache[T]) Exists(ctx context.Context, key string) bool {
	if c.keyPrefix != "" {
		key = c.keyPrefix + ":" + key
	}

	res, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false
	}

	return res > 0
}
