package service

import (
	"context"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/domain/infastructure/datastore"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/domain/model/entity"
)

type UserService interface {
	Create(ctx context.Context, user *entity.User) (int64, error)
	FindById(ctx context.Context, id int64) (*entity.User, error)
	FindByLogin(ctx context.Context, login string) (*entity.User, error)
	FindByLoginAndPasswordHash(
		ctx context.Context,
		login, passwordHash string,
	) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id int64) error

	GetIdByLogin(ctx context.Context, login string) int64
	GetPublicUsersByConvId(ctx context.Context, convId int64) ([]entity.PublicUser, error)
}

type userService struct {
	datastore datastore.UserDatastore
}

func NewUserService(datastore datastore.UserDatastore) UserService {
	return &userService{datastore: datastore}
}

func (us *userService) Create(ctx context.Context, user *entity.User) (int64, error) {
	return us.datastore.Create(ctx, user)
}

func (us *userService) FindById(ctx context.Context, id int64) (*entity.User, error) {
	return us.datastore.FindById(ctx, id)
}

func (us *userService) FindByLogin(ctx context.Context, login string) (*entity.User, error) {
	return us.datastore.FindByLogin(ctx, login)
}

func (us *userService) FindByLoginAndPasswordHash(
	ctx context.Context,
	login, passwordHash string,
) (*entity.User, error) {
	return us.datastore.FindByLoginAndPasswordHash(ctx, login, passwordHash)
}

func (us *userService) Update(
	ctx context.Context,
	user *entity.User,
) error {
	return us.datastore.Update(ctx, user)
}

func (us *userService) Delete(ctx context.Context, userId int64) error {
	return us.datastore.Delete(ctx, userId)
}

func (us *userService) GetIdByLogin(ctx context.Context, login string) int64 {
	return us.datastore.GetIdByLogin(ctx, login)
}

func (us *userService) GetPublicUsersByConvId(
	ctx context.Context,
	convId int64,
) ([]entity.PublicUser, error) {
	return us.datastore.GetPublicUsersByConvId(ctx, convId)
}
