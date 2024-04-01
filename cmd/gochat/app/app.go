package app

import (
	"log/slog"

	"github.com/redis/go-redis/v9"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/storage"
	tcpws "github.com/sazonovItas/gochat-tcp/internal/server"
)

type Application interface {
	Logger() *slog.Logger
	Mux() *tcpws.MuxHandler
}

var _ Application = (&application{})

func (a *application) Logger() *slog.Logger {
	return a.lg
}

func (a *application) Mux() *tcpws.MuxHandler {
	return a.mux
}

type application struct {
	lg           *slog.Logger
	mux          *tcpws.MuxHandler
	storage      *storage.Storage
	cacheStorage *redis.Client
}

func (a *application) Run() error {
	defer func() {
		a.storage.Close()
		a.cacheStorage.Close()
	}()

	return nil
}
