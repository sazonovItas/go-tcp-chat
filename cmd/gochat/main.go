package main

import (
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/config"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/router"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/storage/postgres"
	"github.com/sazonovItas/gochat-tcp/internal/logger/sl"
	tcpws "github.com/sazonovItas/gochat-tcp/internal/server"
	"github.com/sazonovItas/gochat-tcp/internal/utils"
)

func main() {
	// Load config from env variable
	cfg, err := utils.LoadCfgFromFile[config.Config](os.Getenv("CONFIG_PATH"))
	if err != nil {
		log.Fatalf("error to load config: %s", err.Error())
	}
	_ = cfg

	// Setup logger
	logger := NewLogger(cfg.Env, os.Stdout)
	_ = logger

	storage, err := postgres.New(cfg.Storage)
	if err != nil {
		logger.Error("error to init storage", "error", err.Error())
		return
	}
	defer storage.Close()

	// create new router for requests
	mux := router.New(&router.RouterOptions{
		Logger: logger,

		Timeout:     cfg.TCPServer.Timeout,
		IdleTimeout: cfg.TCPServer.IdleTimeout,
	})

	handlersSrv := tcpws.NewServer(cfg.TCPServer.Addr, mux)
	logger.Error("server stoped", "error", handlersSrv.ListenAndServe())
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
