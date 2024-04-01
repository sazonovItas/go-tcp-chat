package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/domain/model/entity"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/domain/model/service"
	tcpws "github.com/sazonovItas/gochat-tcp/internal/server"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// TODO: add user to common chat for testing
func (auh *AuthHandler) SignUp(resp *tcpws.Response, req *tcpws.Request) {
	const op = "gochat.internal.domain.infastructure.transport.http.auth_handler.SignUp"

	var authUser entity.AuthUser

	err := json.Unmarshal([]byte(req.Body), &authUser)
	if err != nil {
		resp.StatusCode = http.StatusBadRequest
		resp.Status = fmt.Sprintf("%s: %s", op, err.Error())
		return
	}

	id := auh.authService.GetUserIdByLogin(req.Ctx(), authUser.Login)
	if id != 0 {
		resp.StatusCode = http.StatusBadRequest
		resp.Status = "User with this login is alredy exists"
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(authUser.Password), bcrypt.DefaultCost)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.Status = fmt.Sprintf("%s: %s", op, err.Error())
		return
	}

	newUser := &entity.User{
		Login:        authUser.Login,
		Name:         authUser.Login,
		Color:        "#737BBB",
		PasswordHash: string(passwordHash),
	}

	_, err = auh.authService.CreateUser(req.Ctx(), newUser)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.Status = fmt.Sprintf("%s: %s", op, err.Error())
		return
	}

	resp.StatusCode = http.StatusOK
	resp.Status = http.StatusText(http.StatusOK)
}

// SignIn sign in user and return token
func (auh *AuthHandler) SignIn(resp *tcpws.Response, req *tcpws.Request) {
	const op = "gochat.internal.domain.infastructure.transport.http.auth_handler.SignIn"

	var authUser entity.AuthUser

	err := json.Unmarshal([]byte(req.Body), &authUser)
	if err != nil {
		resp.StatusCode = http.StatusBadRequest
		resp.Status = fmt.Sprintf("%s: %s", op, err.Error())
		return
	}

	id := auh.authService.GetUserIdByLogin(req.Ctx(), authUser.Login)
	if id == 0 {
		resp.StatusCode = http.StatusNotFound
		resp.Status = "Login is incorrect"
		return
	}

	user, err := auh.authService.GetUserById(req.Ctx(), id)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.Status = fmt.Sprintf("%s: %s", op, err.Error())
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(authUser.Password))
	if err != nil {
		resp.StatusCode = http.StatusNotFound
		resp.Status = "Password is incorrect"
		return
	}

	token, err := auh.authService.CreateToken(req.Ctx(), user.Login, user.PasswordHash)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.Status = fmt.Sprintf("%s: %s", op, err.Error())
		return
	}

	tk, err := json.Marshal(token)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.Status = fmt.Sprintf("%s: %s", op, err.Error())
		return
	}

	resp.StatusCode = http.StatusOK
	resp.Status = http.StatusText(http.StatusOK)
	resp.Body = string(tk)
}
