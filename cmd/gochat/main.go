package main

import (
	"io"
	"log/slog"
	"os"

	"github.com/joho/godotenv"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/config"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/router"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/storage/postgres"
	"github.com/sazonovItas/gochat-tcp/internal/logger/sl"
	tcpws "github.com/sazonovItas/gochat-tcp/internal/server"
	"github.com/sazonovItas/gochat-tcp/internal/utils"
)

func main() {
	configEnv := utils.GetEnv()

	// Setup logger
	logger := NewLogger(configEnv, os.Stdout)
	_ = logger

	// Load env variable from file
	err := godotenv.Load("./configs/.env." + configEnv)
	if err != nil {
		logger.Error("error to load env variable from file", "error", err.Error())
		return
	}

	// Load storage config from env
	storageCfg, err := utils.LoadCfgFromEnv[config.Storage]()
	if err != nil {
		logger.Error("error to load storage config", "error", err.Error())
		return
	}

	// Load server config from env
	serverCfg, err := utils.LoadCfgFromEnv[config.TCPServer]()
	if err != nil {
		logger.Error("error to load server config from env", "error", err.Error())
		return
	}

	storage, err := postgres.New(storageCfg)
	if err != nil {
		logger.Error("error to init storage", "error", err.Error())
		return
	}
	defer storage.Close()

	// create new router for requests
	mux := router.New(&router.RouterOptions{
		Logger: logger,

		Timeout:     serverCfg.Timeout,
		IdleTimeout: serverCfg.IdleTimeout,
	})

	handlersSrv := tcpws.NewServer(serverCfg.Addr, mux)
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
