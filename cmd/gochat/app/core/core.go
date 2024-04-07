package core

import (
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/entity"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/repo"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/service"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/storage"
	"github.com/sazonovItas/gochat-tcp/pkg/cache"
)

type Core struct {
	// Logger for logging
	Logger *slog.Logger

	// Services that using by app
	MessageService service.MessageService
	UserService    service.UserService
	AuthService    service.AuthService
	EventService   service.EventService
}

func New(storage *storage.Storage, cacheStorage *redis.Client, lg *slog.Logger) *Core {
	var core Core

	core.Logger = lg

	// init message service
	core.MessageService = service.NewMessageService(repo.NewMessageRepository(storage), nil)

	// init user service
	core.UserService = service.NewUserService(
		repo.NewUserRepository(storage),
		&cache.CacheOpts{
			Client:            cacheStorage,
			KeyPrefix:         "user",
			DefaultExpiration: time.Minute * 10,
		})

	// init auth service
	tokenStorage := cache.NewCache[entity.Token](&cache.CacheOpts{
		Client:            cacheStorage,
		KeyPrefix:         "auth_token",
		DefaultExpiration: time.Minute * 30,
	})
	core.AuthService = service.NewAuthService(
		repo.NewTokenRepository(tokenStorage, core.UserService),
	)

	// init event service
	core.EventService = service.NewEventService()

	return &core
}
