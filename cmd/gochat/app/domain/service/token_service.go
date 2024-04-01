package service

import (
	"context"
	"errors"
	"time"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/entity"
)

var ErrTokenNotFound = errors.New("token not found")

// TokenService is interface for managing user's authorization tokens
type TokenService interface {
	SaveToken(ctx context.Context, Token entity.Token, expiration time.Duration) error
	TokenById(ctx context.Context, id entity.TokenID) (entity.Token, error)
	UserByTokenId(ctx context.Context, id entity.TokenID) (*entity.User, error)
	DeleteToken(ctx context.Context, id entity.TokenID) error
}

// TokenStorage is interface for store session token
type TokenStorage interface {
	Set(ctx context.Context, key string, value entity.Token, expiration time.Duration) error
	Get(ctx context.Context, key string) (entity.Token, error)
	Delete(ctx context.Context, keys ...string) error
}

// UserStorage is interface for user storage
type UserStorage interface {
	FindById(ctx context.Context, id int64) (*entity.User, error)
}

type tokenService struct {
	tokenStorage TokenStorage
	userStorage  UserStorage
}

func NewTokenRepository(tokenStorage TokenStorage, userStorage UserStorage) TokenService {
	return &tokenService{
		tokenStorage: tokenStorage,
		userStorage:  userStorage,
	}
}

func (tr *tokenService) SaveToken(
	ctx context.Context,
	Token entity.Token,
	expiration time.Duration,
) error {
	return tr.tokenStorage.Set(ctx, Token.UUID.String(), Token, expiration)
}

func (tr *tokenService) TokenById(ctx context.Context, id entity.TokenID) (entity.Token, error) {
	return tr.tokenStorage.Get(ctx, id.String())
}

func (tr *tokenService) UserByTokenId(
	ctx context.Context,
	id entity.TokenID,
) (*entity.User, error) {
	tk, err := tr.tokenStorage.Get(ctx, id.String())
	if err != nil {
		return nil, err
	}

	user, err := tr.userStorage.FindById(ctx, tk.UserId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (tr *tokenService) DeleteToken(ctx context.Context, id entity.TokenID) error {
	return tr.tokenStorage.Delete(ctx, id.String())
}
