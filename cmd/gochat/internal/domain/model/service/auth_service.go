package service

import (
	"context"
	"fmt"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/domain/infastructure/datastore"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/domain/model/entity"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/session"
	"github.com/sazonovItas/gochat-tcp/pkg/cache"
)

type AuthService interface {
	CreateUser(ctx context.Context, user *entity.User) (int64, error)

	GetUserById(ctx context.Context, id int64) (*entity.User, error)
	GetUserIdByLogin(ctx context.Context, login string) int64

	CreateToken(ctx context.Context, login, passwordHash string) (session.SessionToken, error)
	CheckToken(ctx context.Context, token session.SessionToken) error
}

type authService struct {
	datastore datastore.UserDatastore
	manager   session.TokenManager
	cache     cache.Cache[entity.User]
}

func NewAuthService(
	datastore datastore.UserDatastore,
	cacheOpts *cache.CacheOpts,
	managerOpts *session.TokenManagerOptions,
) AuthService {
	return &authService{
		datastore: datastore,
		manager:   session.NewTokenManager(managerOpts),
		cache:     cache.NewCache[entity.User](cacheOpts),
	}
}

func (aus *authService) CreateToken(
	ctx context.Context,
	login, passwordHash string,
) (session.SessionToken, error) {
	return aus.manager.CreateToken(ctx, login, passwordHash)
}

func (aus *authService) CheckToken(ctx context.Context, token session.SessionToken) error {
	return aus.manager.CheckToken(ctx, token)
}

func (aus *authService) CreateUser(ctx context.Context, user *entity.User) (int64, error) {
	id, err := aus.datastore.Create(ctx, user)
	if err != nil {
		return 0, err
	}

	key := fmt.Sprintf("%s:%d", userCacheKey, id)

	user.ID = id
	_ = aus.cache.Set(ctx, key, *user, 0)
	return id, nil
}

func (aus *authService) GetUserIdByLogin(ctx context.Context, login string) int64 {
	return aus.datastore.GetIdByLogin(ctx, login)
}

func (aus *authService) GetUserById(ctx context.Context, id int64) (*entity.User, error) {
	key := fmt.Sprintf("%s:%d", userCacheKey, id)

	cached, err := aus.cache.Get(ctx, key)
	if err == nil {
		return &cached, nil
	}

	user, err := aus.datastore.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	_ = aus.cache.Set(ctx, key, *user, 0)
	return user, nil
}
