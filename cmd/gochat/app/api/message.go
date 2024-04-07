package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/entity"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/repo"
	tcpws "github.com/sazonovItas/gochat-tcp/internal/server"
)

const (
	NoMoreMessages = "no more messages"
)

// /api/v1/messages
func (api *Api) MessagesPrevTimestamp(resp *tcpws.Response, req *tcpws.Request) {
	const op = "gochat.app.api.messages.MessagesPrevTimestamp"

	type request struct {
		Token     entity.Token `json:"auth_token"`
		Timestamp time.Time    `json:"timestamp"`
		Limit     int          `json:"limit"`
	}

	var r request
	if err := json.Unmarshal([]byte(req.Body), &r); err != nil {
		resp.StatusCode = http.StatusBadRequest
		resp.Status = fmt.Errorf("%s: %w", op, err).Error()
		return
	}

	if err := api.app.AuthService.ValidateToken(req.Ctx(), r.Token); err != nil {
		resp.StatusCode = http.StatusUnauthorized
		resp.Status = UnauthorizedMessage
		return
	}

	messages, err := api.app.MessageService.GetConvMessagesPrevTimestamp(
		req.Ctx(),
		r.Timestamp,
		r.Limit,
	)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrNoMessages):
			resp.StatusCode = http.StatusOK
			resp.Status = NoMoreMessages
			resp.Body = "{messages:[]}"
		default:
			resp.StatusCode = http.StatusInternalServerError
			resp.Status = fmt.Errorf("%s: %w", op, err).Error()
			return
		}
	}

	type response struct {
		Messages []entity.Message `json:"messages"`
	}

	data, err := json.Marshal(response{Messages: messages})
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.Status = fmt.Errorf("%s: %w", op, err).Error()
		return
	}

	resp.StatusCode = http.StatusOK
	resp.Status = http.StatusText(http.StatusOK)
	resp.Body = string(data)
}
