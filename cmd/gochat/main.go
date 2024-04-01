package main

import (
	"io"
	"log/slog"
	"os"

	"github.com/joho/godotenv"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/config"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/storage/postgres"
	rediscache "github.com/sazonovItas/gochat-tcp/cmd/gochat/app/storage/redis"
	"github.com/sazonovItas/gochat-tcp/internal/logger/sl"
	"github.com/sazonovItas/gochat-tcp/internal/utils"
)

func main() {
	configEnv := utils.GetEnv()

	// Setup logger
	logger := NewLogger(configEnv, os.Stdout)
	_ = logger

	// Load env variables from file
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
	_ = serverCfg

	redisCfg, err := utils.LoadCfgFromEnv[config.Redis]()
	if err != nil {
		logger.Error("error to load redis config", "error", err.Error())
		return
	}

	storage, err := postgres.New(storageCfg)
	if err != nil {
		logger.Error("error to init storage", "error", err.Error())
		return
	}
	defer storage.Close()

	redis, err := rediscache.New(redisCfg)
	if err != nil {
		logger.Error("error to init redis", "error", err.Error())
		return
	}
	defer redis.Close()
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
