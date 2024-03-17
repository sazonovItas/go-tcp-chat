package main

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

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

			resp.Status = http.StatusText(http.StatusOK)
			resp.StatusCode = http.StatusOK
			resp.Header["Content-Type"] = "application/json"
			resp.Header["Content-Length"] = 100
			resp.Body = "hello, that's me"
		}),
	)

	// urlExample := "postgres://username:password@localhost:5432/database_name"
	connUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Storage.User,
		cfg.Storage.Password,
		cfg.Storage.Host,
		cfg.Storage.Port,
		cfg.Storage.Name,
	)

	db, err := sqlx.Connect("pgx", connUrl)
	if err != nil {
		logger.Error("error connect db", "error", err.Error())
	}

	if err := db.Ping(); err != nil {
		logger.Error("error to ping db", "error", err.Error())
	}

	type Message struct {
		Guid            string `db:"guid"`
		Sender_id       int    `db:"sender_id"`
		Conversation_id int    `db:"conversation_id"`
		Message         string `db:"message"`
		Created_at      string `db:"created_at"`
	}

	var messages []Message
	err = db.Select(&messages, "SELECT * FROM chat.messages")
	if err != nil {
		logger.Error("error scan row", "error", err.Error())
	}

	logger.Info(
		"message from query row",
		"messages", messages,
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
