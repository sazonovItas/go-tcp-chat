package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	tcpws "github.com/sazonovItas/gochat-tcp/internal/server"
)

func Timeout(timeout time.Duration) func(next tcpws.HandlerFunc) tcpws.HandlerFunc {
	return func(next tcpws.HandlerFunc) tcpws.HandlerFunc {
		fn := func(resp *tcpws.Response, req *tcpws.Request) {
			ctx, cancel := context.WithTimeout(req.Context(), timeout)
			defer func() {
				cancel()
				if errors.Is(ctx.Err(), context.DeadlineExceeded) {
					resp.Header["Status"] = http.StatusBadGateway
				}
			}()

			req = req.WithContext(ctx)
			next.Serve(resp, req)
		}

		return tcpws.HandlerFunc(fn)
	}
}
