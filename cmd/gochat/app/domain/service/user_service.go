package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/entity"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/repo"
	"github.com/sazonovItas/gochat-tcp/pkg/cache"
)

var ErrUserLoginAlreadyExists = errors.New("user login already exists")

type UserService interface {
	Create(ctx context.Context, user *entity.User) (int64, error)
	FindById(ctx context.Context, id int64) (*entity.User, error)
	FindByLogin(ctx context.Context, login string) (*entity.User, error)
	FindPublicUserById(ctx context.Context, id int64) (*entity.PublicUser, error)
	FindPublicUserByLogin(ctx context.Context, login string) (*entity.PublicUser, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id int64) error

	GetIdByLogin(ctx context.Context, login string) int64
	GetPublicUsersByConvId(ctx context.Context, convId int64) ([]entity.PublicUser, error)

	ValidateLogin(ctx context.Context, login string) error
}

type userService struct {
	repository repo.UserRepository
	cache      cache.Cache[entity.User]
}

func NewUserService(repository repo.UserRepository, opts *cache.CacheOpts) UserService {
	return &userService{
		repository: repository,
		cache:      cache.NewCache[entity.User](opts),
	}
}

func (us *userService) Create(ctx context.Context, user *entity.User) (int64, error) {
	id, err := us.repository.Create(ctx, user)
	if err != nil {
		return 0, err
	}

	key := fmt.Sprintf("%d", id)

	user.ID = id
	_ = us.cache.Set(ctx, key, *user, 0)
	return id, nil
}

func (us *userService) FindById(ctx context.Context, id int64) (*entity.User, error) {
	key := fmt.Sprintf("%d", id)

	cached, err := us.cache.Get(ctx, key)
	if err == nil {
		return &cached, nil
	}

	user, err := us.repository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	_ = us.cache.Set(ctx, key, *user, 0)
	return user, nil
}

func (us *userService) FindByLogin(ctx context.Context, login string) (*entity.User, error) {
	return us.repository.FindByLogin(ctx, login)
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
	key := fmt.Sprintf("%d", user.ID)

	if us.cache.Exists(ctx, key) {
		_ = us.cache.Set(ctx, key, *user, 0)
	}

	return us.repository.Update(ctx, user)
}

func (us *userService) Delete(ctx context.Context, Id int64) error {
	key := fmt.Sprintf("%d", Id)

	if us.cache.Exists(ctx, key) {
		_ = us.cache.Delete(ctx, key)
	}

	return us.repository.Delete(ctx, Id)
}

func (us *userService) GetIdByLogin(ctx context.Context, login string) int64 {
	return us.repository.GetIdByLogin(ctx, login)
}

func (us *userService) GetPublicUsersByConvId(
	ctx context.Context,
	convId int64,
) ([]entity.PublicUser, error) {
	return us.repository.GetPublicUsersByConvId(ctx, convId)
}

func (us *userService) ValidateLogin(ctx context.Context, login string) error {
	_, err := us.FindByLogin(ctx, login)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrUserNotFound):
			return nil
		default:
			return err
		}
	}

	return ErrUserLoginAlreadyExists
}
