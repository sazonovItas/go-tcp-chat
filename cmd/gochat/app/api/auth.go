package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/entity"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/repo"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/service"
	tcpws "github.com/sazonovItas/gochat-tcp/internal/server"
)

var (
	SuccessfulSignUp = "successful sign up"
	SuccessfulSignIn = "successful sign in"
)

// /api/v1/signup
func (api *Api) SignUp(resp *tcpws.Response, req *tcpws.Request) {
	const op = "gochat.app.api.auth.SignUp"

	var authUser entity.AuthUser

	// unmurshal auth data
	err := json.Unmarshal([]byte(req.Body), &authUser)
	if err != nil {
		resp.StatusCode = http.StatusBadRequest
		resp.Status = fmt.Errorf("%s: %w", op, err).Error()
		return
	}

	// validate login for existence
	err = api.app.UserService.ValidateLogin(req.Ctx(), authUser.Login)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserLoginAlreadyExists):
			resp.StatusCode = http.StatusBadRequest
			resp.Status = service.ErrUserLoginAlreadyExists.Error()
		default:
			resp.StatusCode = http.StatusInternalServerError
			resp.Status = fmt.Errorf("%s: %w", op, err).Error()
		}
		return
	}

	// prepare user
	user, err := api.app.AuthService.SignUp(req.Ctx(), &authUser)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.Status = fmt.Errorf("%s: %w", op, err).Error()
		return
	}

	_, err = api.app.UserService.Create(req.Ctx(), user)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.Status = fmt.Errorf("%s: %w", op, err).Error()
		return
	}

	resp.StatusCode = http.StatusOK
	resp.Status = SuccessfulSignUp
}

// /api/v1/signup
func (api *Api) SignIn(resp *tcpws.Response, req *tcpws.Request) {
	const op = "gochat.app.api.auth.SignIn"

	var authUser entity.AuthUser

	// unmurshal auth data
	err := json.Unmarshal([]byte(req.Body), &authUser)
	if err != nil {
		resp.StatusCode = http.StatusBadRequest
		resp.Status = fmt.Errorf("%s: %w", op, err).Error()
		return
	}

	// find user by login
	user, err := api.app.UserService.FindByLogin(req.Ctx(), authUser.Login)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrUserNotFound):
			resp.StatusCode = http.StatusNotFound
			resp.Status = repo.ErrUserNotFound.Error()
		default:
			resp.StatusCode = http.StatusInternalServerError
			resp.Status = fmt.Errorf("%s: %w", op, err).Error()
		}
		return
	}

	tk, err := api.app.AuthService.SignIn(req.Ctx(), &authUser, user)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidPassword):
			resp.StatusCode = http.StatusBadRequest
			resp.Status = service.ErrInvalidPassword.Error()
		default:
			resp.StatusCode = http.StatusInternalServerError
			resp.Status = fmt.Errorf("%s: %w", op, err).Error()
		}
		return
	}

	type AuthData struct {
		AuthToken entity.Token `json:"auth_token"`
		User      entity.User  `json:"user"`
	}

	response, err := json.Marshal(AuthData{
		AuthToken: tk,
		User:      *user,
	})
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.Status = fmt.Errorf("%s: %w", op, err).Error()
		return
	}

	resp.StatusCode = http.StatusOK
	resp.Status = SuccessfulSignIn
	resp.Body = string(response)
}

// /api/v1/signin/token
func (api *Api) SignInByToken(resp *tcpws.Response, req *tcpws.Request) {
	const op = "gochat.app.api.auth.SignInByToken"

	var token entity.Token
	err := json.Unmarshal([]byte(req.Body), &token)
	if err != nil {
		resp.StatusCode = http.StatusBadRequest
		resp.Status = fmt.Errorf("%s: %w", op, err).Error()
		return
	}

	err = api.app.AuthService.SignInByToken(req.Ctx(), token)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidToken):
			resp.StatusCode = http.StatusUnauthorized
			resp.Status = err.Error()
		default:
			resp.StatusCode = http.StatusInternalServerError
			resp.Status = fmt.Errorf("%s: %w", op, err).Error()
		}

		return
	}

	user, err := api.app.UserService.FindById(req.Ctx(), token.UserId)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.Status = fmt.Errorf("%s: %w", op, err).Error()
		return
	}

	response, err := json.Marshal(*user)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.Status = fmt.Errorf("%s: %w", op, err).Error()
		return
	}

	resp.StatusCode = http.StatusOK
	resp.Status = SuccessfulSignIn
	resp.Body = string(response)
}
