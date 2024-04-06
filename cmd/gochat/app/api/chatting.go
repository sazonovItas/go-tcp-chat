package api

import (
	"net/http"
	"time"

	tcpws "github.com/sazonovItas/gochat-tcp/internal/server"
)

const (
	ProtoNotSupported = "not supported protocol"
	ReadyForMessages  = "ready for messages"
)

// /api/v1/chatting
func (api *Api) Chatting(resp *tcpws.Response, req *tcpws.Request) {
	// const op = "gochat.app.api.chatting.Chatting"

	if req.Proto != tcpws.ProtoWS {
		resp.StatusCode = http.StatusBadRequest
		resp.Status = ProtoNotSupported
		return
	}

	resp.StatusCode = http.StatusOK
	resp.Status = ReadyForMessages
	err := resp.Write()
	if err != nil {
		return
	}

	api.app.Logger.Info("Handling chatting connection",
		"remote_addr", resp.Conn.RemoteAddr(),
	)

	// TODO: handle chatting connection
	time.Sleep(time.Minute * 2)
}
