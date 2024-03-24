package middleware

import (
	"context"
	"sync/atomic"

	tcpws "github.com/sazonovItas/gochat-tcp/internal/server"
)

// Type for the context
type CtxKeyRequestId int

const RequestIdKey CtxKeyRequestId = 0

var reqid uint64

func RequestId() tcpws.Middleware {
	return func(next tcpws.HandlerFunc) tcpws.HandlerFunc {
		fn := func(resp *tcpws.Response, req *tcpws.Request) {
			myId := NextRequestId()
			ctx := context.WithValue(req.Context(), RequestIdKey, myId)
			req = req.WithContext(ctx)
			next.Serve(resp, req)
		}

		return fn
	}
}

func NextRequestId() uint64 {
	return atomic.AddUint64(&reqid, 1)
}
