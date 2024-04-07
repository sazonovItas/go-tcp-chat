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
	// Create creates new user and returns it's id
	// Errors: unknown
	Create(ctx context.Context, user *entity.User) (int64, error)

	// FindById returns user by id
	// Errors: ErrUserNotFound, unknown
	FindById(ctx context.Context, id int64) (*entity.User, error)

	// FindByLogin returns user by login
	// Errors: ErrUserNotFound, unknown
	FindByLogin(ctx context.Context, login string) (*entity.User, error)

	// FindById returns user by id
	// Errors: ErrUserNotFound, unknown
	FindPublicUserById(ctx context.Context, id int64) (*entity.PublicUser, error)

	// FindByLogin returns user by login
	// Errors: ErrUserNotFound, unknown
	FindPublicUserByLogin(ctx context.Context, login string) (*entity.PublicUser, error)

	// Update updates user by id
	// Errors: ErrUserUpdateFailed
	Update(ctx context.Context, user *entity.User) error

	// Delete deletes user by id
	// Errors: ErrUserDeleteFailed
	Delete(ctx context.Context, id int64) error

	// GetIdByLogin returns user id by login
	GetIdByLogin(ctx context.Context, login string) int64

	// GetPublicUsersByConvId returns public users
	// Errors: unknown
	GetPublicUsers(ctx context.Context) ([]entity.PublicUser, error)

	// ValidateLogin validates login
	// Errors: ErrUserLoginAlreadyExists, unknown
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

// Create is implementing interface UserService
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

// FindById is implementing interface UserService
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

// FindByLogin is implementing interface UserService
func (us *userService) FindByLogin(ctx context.Context, login string) (*entity.User, error) {
	return us.repository.FindByLogin(ctx, login)
}

// FindPublicUserById is implementing interface UserService
func (us *userService) FindPublicUserById(
	ctx context.Context,
	id int64,
) (*entity.PublicUser, error) {
	key := fmt.Sprintf("%d", id)

	cached, err := us.cache.Get(ctx, key)
	if err == nil {
		return &entity.PublicUser{
			ID:    cached.ID,
			Login: cached.Login,
			Name:  cached.Name,
			Color: cached.Color,
		}, nil
	}

	user, err := us.repository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	_ = us.cache.Set(ctx, key, *user, 0)
	return &entity.PublicUser{
		ID:    user.ID,
		Login: user.Login,
		Name:  user.Name,
		Color: user.Color,
	}, nil
}

// FindPublicUserByLogin is implementing interface UserService
func (us *userService) FindPublicUserByLogin(
	ctx context.Context,
	login string,
) (*entity.PublicUser, error) {
	user, err := us.repository.FindByLogin(ctx, login)
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

// Update is implementing interface UserService
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

// Delete is implementing interface UserService
func (us *userService) Delete(ctx context.Context, Id int64) error {
	key := fmt.Sprintf("%d", Id)

	if us.cache.Exists(ctx, key) {
		_ = us.cache.Delete(ctx, key)
	}

	return us.repository.Delete(ctx, Id)
}

// Delete is implementing interface UserService
func (us *userService) GetIdByLogin(ctx context.Context, login string) int64 {
	return us.repository.GetIdByLogin(ctx, login)
}

// GetPublicUsers is implementing interface UserService
func (us *userService) GetPublicUsers(
	ctx context.Context,
) ([]entity.PublicUser, error) {
	return us.repository.GetPublicUsers(ctx)
}

// ValidateLogin is implementing interface ValidateLogin
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
