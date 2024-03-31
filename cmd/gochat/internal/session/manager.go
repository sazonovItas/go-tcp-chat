package session

import (
	"context"
	"errors"
	"time"
)

var (
	ErrTokenExpired = errors.New("session token expired")
	ErrCreateToken  = errors.New("cannot create token")
)

type TokenManager interface {
	CreateToken(ctx context.Context, login, passwordHash string) (SessionToken, error)
	CheckToken(ctx context.Context, token SessionToken) error
}

type TokenManagerOptions struct {
	SessionStorage SessionStorage[SessionToken]
	Expiration     time.Duration
}

type tokenManager struct {
	sessionStorage SessionStorage[SessionToken]
	expiration     time.Duration
}

func NewTokenManager(managerOpts *TokenManagerOptions) TokenManager {
	return &tokenManager{
		sessionStorage: managerOpts.SessionStorage,
		expiration:     managerOpts.Expiration,
	}
}

func (tm *tokenManager) CreateToken(
	ctx context.Context,
	login, passwordHash string,
) (SessionToken, error) {
	tk := NewSessionToken(login, passwordHash)

	err := tm.sessionStorage.Set(ctx, tk.UUID.String(), tk, tm.expiration)
	if err != nil {
		return SessionToken{}, err
	}

	return tk, nil
}

func (tm *tokenManager) CheckToken(ctx context.Context, token SessionToken) error {
	tk, err := tm.sessionStorage.Get(ctx, token.UUID.String())
	if err != nil {
		return ErrTokenExpired
	}

	if tk != token {
		_ = tm.sessionStorage.Delete(ctx, token.UUID.String())
		return ErrTokenExpired
	}

	return nil
}
