package main

import (
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app"
	"github.com/sazonovItas/gochat-tcp/internal/logger/sl"
	"github.com/sazonovItas/gochat-tcp/internal/utils"
)

func main() {
	configEnv := utils.GetEnv()

	// Load env variables from file
	err := godotenv.Load("./configs/.env." + configEnv)
	if err != nil {
		log.Fatalf("%s: %s", "error to load env variables from file", err.Error())
		return
	}

	// Setup logger
	logger := NewLogger(configEnv, os.Stdout)
	_ = logger

	cfg, err := app.InitAppConfig()
	if err != nil {
		logger.Error("error to init app config", "error", err.Error())
		return
	}
	_ = cfg
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
