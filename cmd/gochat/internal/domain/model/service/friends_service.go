package service

import (
	"context"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/domain/infastructure/datastore"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/domain/model/entity"
)

type FriendService interface {
	Create(ctx context.Context, friend *entity.Friend) (int64, error)
	FindById(ctx context.Context, id int64) (*entity.Friend, error)
	FindByUserAndFriendId(ctx context.Context, userId, friendId int64) (*entity.Friend, error)
	Delete(ctx context.Context, id int64) error
}

type friendService struct {
	datastore datastore.FriendDatastore
}

func NewFriendService(datastore datastore.FriendDatastore) FriendService {
	return &friendService{datastore: datastore}
}

func (fr *friendService) Create(ctx context.Context, friend *entity.Friend) (int64, error) {
	return fr.datastore.Create(ctx, friend)
}

func (fr *friendService) FindById(ctx context.Context, id int64) (*entity.Friend, error) {
	return fr.datastore.FindById(ctx, id)
}

func (fr *friendService) FindByUserAndFriendId(
	ctx context.Context,
	userId, friendId int64,
) (*entity.Friend, error) {
	return fr.datastore.FindByUserAndFriendId(ctx, userId, friendId)
}

func (fr *friendService) Delete(ctx context.Context, id int64) error {
	return fr.datastore.Delete(ctx, id)
}
