package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/config"
)

// TODO: Addd optional redis config
// TODO: Need to know more about redis connections
func New(cfg *config.Redis) (*redis.Client, error) {
	const op = "gochat.internal.storage.redis.New"

	// new client for redis
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// check ping redis
	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return client, nil
}
