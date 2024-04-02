package repo

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/entity"
)

var (
	ErrGenerateUUID  = errors.New("failed to generate uuid")
	ErrTokenNotFound = errors.New("token not found")
)

type TokenRepository interface {
	CreateToken(ctx context.Context, userId int64) (entity.Token, error)
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

type tokenRepository struct {
	tokenStorage TokenStorage
	userStorage  UserStorage
}

func NewTokenRepository(tokenStorage TokenStorage, userStorage UserStorage) TokenRepository {
	return &tokenRepository{
		tokenStorage: tokenStorage,
		userStorage:  userStorage,
	}
}

func (tr *tokenRepository) CreateToken(
	ctx context.Context,
	userId int64,
) (token entity.Token, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ErrGenerateUUID
		}
	}()

	token.ID = entity.TokenID(uuid.New().String())
	token.UserId = userId
	return
}

func (tr *tokenRepository) SaveToken(
	ctx context.Context,
	Token entity.Token,
	expiration time.Duration,
) error {
	return tr.tokenStorage.Set(ctx, Token.ID.String(), Token, expiration)
}

func (tr *tokenRepository) TokenById(ctx context.Context, id entity.TokenID) (entity.Token, error) {
	return tr.tokenStorage.Get(ctx, id.String())
}

func (tr *tokenRepository) UserByTokenId(
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

func (tr *tokenRepository) DeleteToken(ctx context.Context, id entity.TokenID) error {
	return tr.tokenStorage.Delete(ctx, id.String())
}
