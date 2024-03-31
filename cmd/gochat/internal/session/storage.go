package session

import (
	"context"
	"time"
)

type SessionStorage[T any] interface {
	Set(ctx context.Context, key string, value T, expiration time.Duration) error
	Get(ctx context.Context, key string) (T, error)
	Delete(ctx context.Context, keys ...string) error
}
