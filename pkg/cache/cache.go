package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// Erros for upper layer
var (
	ErrCacheMissed    = errors.New("cache missed")
	ErrCacheSetFailed = errors.New("cache set failed")
	ErrCacheGetFailed = errors.New("cache get failed")
)

// Cache is interface implementing basic operations on cache
type Cache[T any] interface {
	Set(ctx context.Context, key string, value T, expiration time.Duration) error
	Get(ctx context.Context, key string) (T, error)
	Delete(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, key string) bool
}

// CacheOpts represents options for redis based cache
type CacheOpts struct {
	Client *redis.Client

	DefaultExpiration time.Duration
}

// cache represents redis based cache and implements Cache interface
type cache[T any] struct {
	client *redis.Client

	defaultExpiration time.Duration
}

// NewCacheByOpts creates new cache by opts
func NewCache[T any](opts *CacheOpts) Cache[T] {
	return &cache[T]{
		client:            opts.Client,
		defaultExpiration: opts.DefaultExpiration,
	}
}

// Set sets data under a key, if expiration is 0, then usgin default expiration
func (c *cache[T]) Set(
	ctx context.Context,
	key string,
	value T,
	expiration time.Duration,
) (err error) {
	payload, err := json.Marshal(value)
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
	payload, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return value, err
	}

	err = json.Unmarshal([]byte(payload), &value)
	if err != nil {
		return value, err
	}

	return value, nil
}

// Delete deletes keys from the cache
func (c *cache[T]) Delete(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

// Exists check existence of a key
func (c *cache[T]) Exists(ctx context.Context, key string) bool {
	res, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false
	}

	return res > 0
}
