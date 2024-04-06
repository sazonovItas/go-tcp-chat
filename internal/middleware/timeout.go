package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	tcpws "github.com/sazonovItas/gochat-tcp/internal/server"
)

func Timeout(timeout time.Duration) tcpws.Middleware {
	return func(next tcpws.HandlerFunc) tcpws.HandlerFunc {
		fn := func(resp *tcpws.Response, req *tcpws.Request) {
			if req.Proto == tcpws.ProtoHTTP {
				ctx, cancel := context.WithTimeout(req.Ctx(), timeout)
				defer func() {
					cancel()
					if errors.Is(ctx.Err(), context.DeadlineExceeded) {
						resp.StatusCode = http.StatusBadGateway
						resp.Status = http.StatusText(http.StatusBadGateway)
						resp.Body = ""
					}
				}()

				req = req.WithContext(ctx)
			}

			next.Serve(resp, req)
		}

		return fn
	}
}
