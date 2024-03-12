package main

import (
	"io"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/sazonovItas/gochat-tcp/internal/config"
	"github.com/sazonovItas/gochat-tcp/internal/logger/sl"
	"github.com/sazonovItas/gochat-tcp/internal/middleware"
	tcpws "github.com/sazonovItas/gochat-tcp/internal/server"
)

func main() {
	// Load config from env variable
	cfg, err := config.LoadCfgFromEnv("CONFIG_PATH")
	if err != nil {
		log.Fatalf("error to load config: %s", err.Error())
	}
	_ = cfg

	// Setup logger
	logger := NewLogger(cfg.Env, os.Stdout)
	_ = logger

	mux := tcpws.NewMuxHandler()
	mux.HandleFunc(
		"GET",
		"/user/{id}",
		middleware.Timeout(time.Second)(func(resp *tcpws.Response, req *tcpws.Request) {
			logger.Info("handling request", "request", req)
			time.Sleep(time.Second)

			resp.Header["Status"] = 200
			resp.Header["Status-Code"] = "OK"
			resp.Header["Content-Type"] = "text/json"
			resp.Body = "request accepting"
			_ = resp.Write()
		}),
	)

	logger.Error("server stoped", "error", tcpws.ListenAndServe(cfg.TCPServer.Addr, mux))
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
