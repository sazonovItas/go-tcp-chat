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
	Create(ctx context.Context, user *entity.User) (int64, error)

	GetById(ctx context.Context, id int64) (*entity.User, error)
	GetIdByLogin(ctx context.Context, login string) int64
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

func (aus *authService) Create(ctx context.Context, user *entity.User) (int64, error) {
	id, err := aus.datastore.Create(ctx, user)
	if err != nil {
		return 0, err
	}

	key := fmt.Sprintf("%s:%d", userCacheKey, id)

	user.ID = id
	_ = aus.cache.Set(ctx, key, *user, 0)
	return id, nil
}

func (aus *authService) GetIdByLogin(ctx context.Context, login string) int64 {
	return aus.datastore.GetIdByLogin(ctx, login)
}

func (aus *authService) GetById(ctx context.Context, id int64) (*entity.User, error) {
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
