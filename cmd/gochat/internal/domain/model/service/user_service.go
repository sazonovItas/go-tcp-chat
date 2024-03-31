package service

import (
	"context"
	"fmt"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/domain/infastructure/datastore"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/domain/model/entity"
	"github.com/sazonovItas/gochat-tcp/pkg/cache"
)

type UserService interface {
	FindById(ctx context.Context, id int64) (*entity.User, error)
	FindByLogin(ctx context.Context, login string) (*entity.User, error)
	FindPublicUserById(ctx context.Context, id int64) (*entity.PublicUser, error)
	FindPublicUserByLogin(ctx context.Context, login string) (*entity.PublicUser, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id int64) error

	GetPublicUsersByConvId(ctx context.Context, convId int64) ([]entity.PublicUser, error)
}

type userService struct {
	datastore datastore.UserDatastore
	cache     cache.Cache[entity.User]
}

func NewUserService(datastore datastore.UserDatastore, opts *cache.CacheOpts) UserService {
	return &userService{
		datastore: datastore,
		cache:     cache.NewCache[entity.User](opts),
	}
}

const userCacheKey = "user"

func (us *userService) FindById(ctx context.Context, id int64) (*entity.User, error) {
	key := fmt.Sprintf("%s:%d", userCacheKey, id)

	cached, err := us.cache.Get(ctx, key)
	if err == nil {
		return &cached, nil
	}

	user, err := us.datastore.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	_ = us.cache.Set(ctx, key, *user, 0)
	return user, nil
}

func (us *userService) FindByLogin(ctx context.Context, login string) (*entity.User, error) {
	return us.datastore.FindByLogin(ctx, login)
}

func (us *userService) FindPublicUserById(
	ctx context.Context,
	id int64,
) (*entity.PublicUser, error) {
	user, err := us.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &entity.PublicUser{
		ID:    user.ID,
		Login: user.Login,
		Name:  user.Name,
		Color: user.Color,
	}, nil
}

func (us *userService) FindPublicUserByLogin(
	ctx context.Context,
	login string,
) (*entity.PublicUser, error) {
	user, err := us.FindByLogin(ctx, login)
	if err != nil {
		return nil, err
	}

	return &entity.PublicUser{
		ID:    user.ID,
		Login: user.Login,
		Name:  user.Name,
		Color: user.Color,
	}, nil
}

func (us *userService) Update(
	ctx context.Context,
	user *entity.User,
) error {
	key := fmt.Sprintf("%s:%d", userCacheKey, user.ID)

	if us.cache.Exists(ctx, key) {
		_ = us.cache.Set(ctx, key, *user, 0)
	}

	return us.datastore.Update(ctx, user)
}

func (us *userService) Delete(ctx context.Context, Id int64) error {
	key := fmt.Sprintf("%s:%d", userCacheKey, Id)

	if us.cache.Exists(ctx, key) {
		_ = us.cache.Delete(ctx, key)
	}

	return us.datastore.Delete(ctx, Id)
}

func (us *userService) GetPublicUsersByConvId(
	ctx context.Context,
	convId int64,
) ([]entity.PublicUser, error) {
	return us.datastore.GetPublicUsersByConvId(ctx, convId)
}
