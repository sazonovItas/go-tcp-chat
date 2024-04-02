package core

import (
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/entity"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/repo"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/service"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/storage"
	"github.com/sazonovItas/gochat-tcp/pkg/cache"
)

type Core struct {
	// Services that using by app
	ConversationService service.ConversationService
	FriendService       service.FriendService
	MessageService      service.MessageService
	ParticipantService  service.ParticipantService
	UserService         service.UserService
	AuthService         service.AuthService
}

func New(storage *storage.Storage, cacheStorage *redis.Client) *Core {
	var core Core

	// init conversation service
	core.ConversationService = service.NewConversationService(
		repo.NewConversationRepository(storage),
		&cache.CacheOpts{
			Client:            cacheStorage,
			KeyPrefix:         "conversation",
			DefaultExpiration: time.Minute * 5,
		},
	)

	// init friend service
	core.FriendService = service.NewFriendService(
		repo.NewFriendRepository(storage),
		&cache.CacheOpts{
			Client:            cacheStorage,
			KeyPrefix:         "friend",
			DefaultExpiration: time.Minute * 5,
		},
	)

	// init message service
	core.MessageService = service.NewMessageService(repo.NewMessageRepository(storage), nil)

	// init participant service
	core.ParticipantService = service.NewParticipantService(
		repo.NewParticipantRepository(storage),
		&cache.CacheOpts{
			Client:            cacheStorage,
			KeyPrefix:         "participant",
			DefaultExpiration: time.Minute * 5,
		})

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

	return &core
}
