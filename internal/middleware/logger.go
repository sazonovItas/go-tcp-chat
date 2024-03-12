package middleware

import tcpws "github.com/sazonovItas/gochat-tcp/internal/server"

func Logger(next tcpws.HandlerFunc) tcpws.HandlerFunc {
	return RequestLogger(next)
}

func RequestLogger(next tcpws.HandlerFunc) tcpws.HandlerFunc {
	return next.Serve
}
