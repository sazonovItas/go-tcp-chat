package middleware

import (
	"log/slog"
	"time"

	tcpws "github.com/sazonovItas/gochat-tcp/internal/server"
)

func Logger(log *slog.Logger) tcpws.Middleware {
	return func(next tcpws.HandlerFunc) tcpws.HandlerFunc {
		fn := func(resp *tcpws.Response, req *tcpws.Request) {
			t1 := time.Now()

			var requestId uint64
			if reqId, ok := req.Ctx().Value(RequestIdKey).(uint64); ok {
				requestId = reqId
			}

			defer func() {
				log.Info("request completed",
					slog.Uint64("request_id", requestId),
					slog.String("method", req.Method),
					slog.String("path", req.Url),
					slog.String("remote_addr", resp.Conn.RemoteAddr().String()),
					slog.Any("header", req.Header),
					slog.Int("status_code", resp.StatusCode),
					slog.String("status", resp.Status),
					slog.Int("response_size", len(resp.Body)),
					slog.String("duration", time.Since(t1).String()),
				)
			}()

			next.Serve(resp, req)
		}

		return fn
	}
}
