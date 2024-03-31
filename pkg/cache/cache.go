package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrCacheMissed = errors.New("cache missed")

// Cache is interface for redis based cache
type Cache[T any] interface {
	Set(ctx context.Context, key string, value T, expiration time.Duration) error
	Get(ctx context.Context, key string) (T, error)
	Delete(ctx context.Context, keys ...string) error
}

// CacheOpts represents options for redis based cache
type CacheOpts struct {
	Client *redis.Client

	DefaultExpiration time.Duration
}

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

func (c *cache[T]) Delete(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}
