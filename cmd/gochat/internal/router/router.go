package router

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/sazonovItas/gochat-tcp/internal/middleware"
	tcpws "github.com/sazonovItas/gochat-tcp/internal/server"
)

type RouterOptions struct {
	// Logger logs some data of request and other things
	Logger *slog.Logger

	// Timeout using for timeout middleware
	Timeout     time.Duration
	IdleTimeout time.Duration
}

func New(rtOpt *RouterOptions) *tcpws.MuxHandler {
	mux := tcpws.NewMuxHandler()
	mux.Use(middleware.Timeout(rtOpt.Timeout))
	mux.Use(middleware.Logger(rtOpt.Logger))
	mux.Use(middleware.RequestId())

	mux.HandleFunc(
		"GET",
		"/user/{id}",
		func(resp *tcpws.Response, req *tcpws.Request) {
			resp.Status = http.StatusText(http.StatusOK)
			resp.StatusCode = http.StatusOK
			resp.Header["Content-Type"] = "application/json"
			resp.Body = "hello from server"
		},
	)

	return mux
}
