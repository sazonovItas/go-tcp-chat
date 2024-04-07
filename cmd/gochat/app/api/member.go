package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/repo"
	tcpws "github.com/sazonovItas/gochat-tcp/internal/server"
)

// /api/v1/member/{id}
func (api *Api) GetChatMemberById(resp *tcpws.Response, req *tcpws.Request) {
	const op = "gochat.app.api.user.GetChatMemberById"

	userId, err := strconv.ParseInt(req.ParamByName("id"), 10, 64)
	if err != nil {
		resp.StatusCode = http.StatusBadRequest
		resp.Status = fmt.Errorf("%s: %w", op, err).Error()
		return
	}

	user, err := api.app.UserService.FindPublicUserById(req.Ctx(), userId)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrUserNotFound):
			resp.StatusCode = http.StatusNotFound
			resp.Status = fmt.Errorf("%s: %w", op, err).Error()
		default:
			resp.StatusCode = http.StatusInternalServerError
			resp.Status = fmt.Errorf("%s: %w", op, err).Error()
		}

		return
	}

	response, err := json.Marshal(*user)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.Status = fmt.Errorf("%s: %w", op, err).Error()
		return
	}

	resp.StatusCode = http.StatusOK
	resp.Status = http.StatusText(http.StatusOK)
	resp.Body = string(response)
}

// /api/v1/member
func (api *Api) GetChatMembers(resp *tcpws.Response, req *tcpws.Request) {
	const op = "gochat.app.api.user.GetChatMembers"

	users, err := api.app.UserService.GetPublicUsers(req.Ctx())
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.Status = fmt.Errorf("%s: %w", op, err).Error()
		return
	}

	response, err := json.Marshal(users)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.Status = fmt.Errorf("%s: %w", op, err).Error()
		return
	}

	resp.StatusCode = http.StatusOK
	resp.Status = http.StatusText(http.StatusOK)
	resp.Body = string(response)
}
