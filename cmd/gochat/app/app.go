package app

import (
	"log/slog"

	"github.com/redis/go-redis/v9"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/storage"
	tcpws "github.com/sazonovItas/gochat-tcp/internal/server"
)

type App interface {
	Logger() *slog.Logger
	Mux() *tcpws.MuxHandler
	Storage() *storage.Storage
	CacheStorage() *redis.Client
}

var _ App = (&app{})

func (a *app) Logger() *slog.Logger {
	return a.lg
}

func (a *app) Mux() *tcpws.MuxHandler {
	return a.mux
}

func (a *app) Storage() *storage.Storage {
	return a.storage
}

func (a *app) CacheStorage() *redis.Client {
	return a.cacheStorage
}

type app struct {
	lg           *slog.Logger
	mux          *tcpws.MuxHandler
	storage      *storage.Storage
	cacheStorage *redis.Client
}
