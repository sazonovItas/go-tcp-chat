package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/entity"
	tcpws "github.com/sazonovItas/gochat-tcp/internal/server"
)

const (
	ProtoNotSupported   = "not supported protocol"
	ReadyForMessages    = "ready for messages"
	UnauthorizedMessage = "token expired"
)

// /api/v1/chatting
func (api *Api) Chatting(resp *tcpws.Response, req *tcpws.Request) {
	// const op = "gochat.app.api.chatting.Chatting"

	if req.Proto != tcpws.ProtoWS {
		resp.StatusCode = http.StatusBadRequest
		resp.Status = ProtoNotSupported
		return
	}

	var token entity.Token
	if err := json.Unmarshal([]byte(req.Body), &token); err != nil {
		resp.StatusCode = http.StatusBadRequest
		resp.Status = ProtoNotSupported
		return
	}

	if err := api.app.AuthService.ValidateToken(req.Ctx(), token); err != nil {
		resp.StatusCode = http.StatusUnauthorized
		resp.Status = UnauthorizedMessage
		return
	}

	resp.StatusCode = http.StatusOK
	resp.Status = ReadyForMessages
	if err := resp.Write(); err != nil {
		return
	}

	api.app.Logger.Info("Handling chatting connection",
		"remote_addr", resp.Conn.RemoteAddr(),
	)

	// TODO: handle chatting connection
	time.Sleep(time.Second * 10)
}
