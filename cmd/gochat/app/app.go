package app

import (
	"io"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/core"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/entity"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/repo"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/service"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/internal/hasher"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/storage"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/storage/postgres"
	rediscache "github.com/sazonovItas/gochat-tcp/cmd/gochat/app/storage/redis"
	"github.com/sazonovItas/gochat-tcp/internal/logger/sl"
	"github.com/sazonovItas/gochat-tcp/internal/middleware"
	tcpws "github.com/sazonovItas/gochat-tcp/internal/server"
	"github.com/sazonovItas/gochat-tcp/pkg/cache"
)

type Application struct {
	Logger *slog.Logger
	*core.Core

	listenAddr string
	mux        *tcpws.MuxHandler

	storage      *storage.Storage
	cacheStorage *redis.Client
}

func InitApp(cfg *AppConfig) (*Application, error) {
	var app Application

	// init logger
	app.Logger = NewLogger(cfg.Options.Env, cfg.Options.LogWriter)

	// init storage
	db, err := postgres.New(&cfg.Storage)
	if err != nil {
		return nil, err
	}
	app.storage = db

	// init cache storage
	cache, err := rediscache.New(&cfg.CacheStorage)
	if err != nil {
		return nil, err
	}
	app.cacheStorage = cache

	// init core
	app.Core = InitCore(db, cache)

	// setup server address and mux handler routes
	app.listenAddr = cfg.TCPServer.Addr
	app.mux = InitMux()

	// setup middlewares for mux handler
	app.mux.Use(middleware.RequestId())
	app.mux.Use(middleware.Logger(app.Logger))
	app.mux.Use(middleware.Timeout(cfg.TCPServer.Timeout))

	return &app, nil
}

func (app *Application) Run() error {
	defer func() {
		app.storage.Close()
		app.cacheStorage.Close()

		app.Logger.Info("server stopped")
	}()

	app.Logger.Info("server start running")
	return tcpws.ListenAndServe(app.listenAddr, app.mux)
}

func InitCore(storage *storage.Storage, cacheStorage *redis.Client) *core.Core {
	var core core.Core

	// init hasher
	core.Hasher = hasher.New(10)

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

	core.AuthService = service.NewAuthService(tokenStorage, core.UserService)

	return &core
}

func InitMux() *tcpws.MuxHandler {
	mux := tcpws.NewMuxHandler()

	return mux
}

// Create new logger that is specified by env
func NewLogger(env string, out io.Writer) *slog.Logger {
	var opts sl.HandlerOptions

	switch env {
	case "local":
		opts.SlogOpts.Level = slog.LevelDebug
	case "dev":
		opts.SlogOpts.Level = slog.LevelDebug
	case "prod":
		opts.SlogOpts.Level = slog.LevelInfo
	}

	return slog.New(sl.NewHandler(out, opts))
}
