package app

import (
	"io"
	"log/slog"

	"github.com/redis/go-redis/v9"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/api"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/core"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/storage"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/storage/postgres"
	rediscache "github.com/sazonovItas/gochat-tcp/cmd/gochat/app/storage/redis"
	"github.com/sazonovItas/gochat-tcp/internal/logger/sl"
	"github.com/sazonovItas/gochat-tcp/internal/middleware"
	tcpws "github.com/sazonovItas/gochat-tcp/internal/server"
)

type Application struct {
	Logger *slog.Logger
	*core.Core

	listenAddr string
	mux        *tcpws.MuxHandler

	storage      *storage.Storage
	cacheStorage *redis.Client
}

func InitApp(cfg *Config) (*Application, error) {
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
	app.Core = core.New(db, cache, app.Logger)

	// setup server address and mux handler routes
	app.listenAddr = cfg.TCPServer.Addr
	app.mux = tcpws.NewMuxHandler()

	// setup middlewares for mux handler
	app.mux.Use(middleware.RequestId())
	app.mux.Use(middleware.Logger(app.Logger))
	app.mux.Use(middleware.Timeout(cfg.TCPServer.Timeout))

	// init routes for app
	InitRoutes(app.mux, app.Core)

	return &app, nil
}

func (app *Application) Run() error {
	defer func() {
		app.storage.Close()
		app.cacheStorage.Close()

		app.Logger.Info("server stopped")
	}()

	app.Logger.Info("server start running", "address", app.listenAddr)
	return tcpws.ListenAndServe(app.listenAddr, app.mux)
}

func InitRoutes(mux *tcpws.MuxHandler, core *core.Core) *tcpws.MuxHandler {
	handlers := api.NewApi(core)

	mux.HandleFunc("POST", "/api/v1/signup", handlers.SignUp)
	mux.HandleFunc("POST", "/api/v1/signin", handlers.SignIn)
	mux.HandleFunc(tcpws.ProtoWS, "/api/v1/chatting", handlers.Chatting)

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
