package main

import (
	"context"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/config"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/storage/models"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/storage/postgres"
	"github.com/sazonovItas/gochat-tcp/internal/logger/sl"
	"github.com/sazonovItas/gochat-tcp/internal/middleware"
	tcpws "github.com/sazonovItas/gochat-tcp/internal/server"
	"github.com/sazonovItas/gochat-tcp/internal/utils"
)

func main() {
	// Load config from env variable
	cfg, err := utils.LoadCfgFromEnv[config.Config]("CONFIG_PATH")
	if err != nil {
		log.Fatalf("error to load config: %s", err.Error())
	}
	_ = cfg

	// Setup logger
	logger := NewLogger(cfg.Env, os.Stdout)
	_ = logger

	storage, err := postgres.New(cfg.Storage)
	if err != nil {
		logger.Error("error to init storage", "error", err)
		return
	}
	defer storage.Close()

	mux := tcpws.NewMuxHandler()
	mux.Use(middleware.Timeout(cfg.TCPServer.Timeout))
	mux.Use(middleware.Logger(logger))
	mux.Use(middleware.RequestId())

	mux.HandleFunc(
		"GET",
		"/user/{id}",
		func(resp *tcpws.Response, req *tcpws.Request) {
			logger.Info("handling request", "request", req)

			resp.Status = http.StatusText(http.StatusOK)
			resp.StatusCode = http.StatusOK
			resp.Header["Content-Type"] = "application/json"
			resp.Body = "hello from server"
		},
	)

	db := models.NewModelStorage(storage)

	// Test create user
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte("itas124"),
		bcrypt.DefaultCost,
	)
	if err != nil {
		logger.Error("error to hash password", "error", err.Error())
		return
	}

	logger.Info(
		"generated password hash",
		"password_hash",
		string(passwordHash),
		"password_hash_len",
		len(string(passwordHash)),
		"password",
		"itas124",
	)

	user_id, err := db.CreateUser(context.Background(), &models.User{
		Login:        "itas",
		Name:         "Alex",
		Color:        "#ef1512",
		PasswordHash: string(passwordHash),
	})
	if err != nil {
		logger.Error("error to create user", "error", err.Error())
		return
	}

	user, err := db.GetUserById(context.Background(), user_id)
	if err != nil {
		logger.Error("error to get user", "error", err.Error(), "user_id", user_id)
		return
	}
	logger.Info("get new user from db", "user", user)

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("itas124"))
	if err != nil {
		logger.Error(
			"error to compare password_hash and password",
			"error",
			err,
			"password_hash",
			passwordHash,
			"password",
			"itas124",
		)
		return
	}

	err = db.UpdateUser(context.Background(), &models.UpdateUser{
		ID:    user_id,
		Name:  "Itas",
		Color: "#ef1",
	})
	if err != nil {
		logger.Error("error to update user", "error", err.Error())
		return
	}
	logger.Info("user updated", "user_id", user_id)

	// Test create conversation
	conv_id, err := db.CreateConversation(context.Background(), &models.Conversation{
		Title:            "test",
		ConversationType: models.Conversation2P2Kind,
		CreatorId:        1,
	})
	if err != nil {
		logger.Error("error to create user", "error", err.Error())
		return
	}

	var conversation models.Conversation
	err = storage.Get(
		&conversation,
		"SELECT id, title, conversation_type, creator_id FROM chat.conversations WHERE id=$1",
		conv_id,
	)
	if err != nil {
		logger.Error("error to conversation", "error", err.Error(), "conversation_id", conv_id)
		return
	}
	logger.Info("get new conversation from db", "conversation", conversation)

	// Test create message
	msg_id, err := db.CreateMessage(context.Background(), &models.Message{
		SenderID:       user_id,
		ConversationID: conv_id,
		MessageType:    models.UserTextMessage,
		Message:        "hello",
		CreatedAt:      time.Now(),
	})
	if err != nil {
		logger.Error("error to create message", "error", err.Error())
		return
	}

	message, err := db.GetMessageById(context.Background(), msg_id)
	if err != nil {
		logger.Error("error to get message", "error", err.Error(), "message_id", msg_id)
		return
	}
	logger.Info("get new message from db", "message", message)

	err = db.UpdateMessage(context.Background(), &models.UpdateMessage{
		ID:        msg_id,
		Message:   "hi, updated message",
		CreatedAt: time.Now(),
	})
	if err != nil {
		logger.Error("error to update message", "error", err.Error())
	}

	message, err = db.GetMessageById(context.Background(), msg_id)
	if err != nil {
		logger.Error("error to get message", "error", err.Error(), "message_id", msg_id)
		return
	}
	logger.Info("get updated message from db", "message", message)

	err = db.DeleteMessageId(context.Background(), msg_id)
	if err != nil {
		logger.Error("error to delete message", "error", err.Error(), "msg_id", msg_id)
		return
	}

	// handlersSrv := tcpws.NewServer(cfg.TCPServer.Addr, mux)
	// logger.Error("server stoped", "error", handlersSrv.ListenAndServe())
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
