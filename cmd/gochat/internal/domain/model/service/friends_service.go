package service

import (
	"context"
	"fmt"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/domain/infastructure/datastore"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/domain/model/entity"
	"github.com/sazonovItas/gochat-tcp/pkg/cache"
)

type FriendService interface {
	Create(ctx context.Context, friend *entity.Friend) (int64, error)
	FindById(ctx context.Context, id int64) (*entity.Friend, error)
	FindByUserAndFriendId(ctx context.Context, userId, friendId int64) (*entity.Friend, error)
	Delete(ctx context.Context, id int64) error
}

type friendService struct {
	datastore datastore.FriendDatastore
	cache     cache.Cache[entity.Friend]
}

func NewFriendService(datastore datastore.FriendDatastore, opts *cache.CacheOpts) FriendService {
	return &friendService{
		datastore: datastore,
		cache:     cache.NewCache[entity.Friend](opts),
	}
}

const friendCacheKey = "friend"

func (fr *friendService) Create(ctx context.Context, friend *entity.Friend) (int64, error) {
	id, err := fr.datastore.Create(ctx, friend)
	if err != nil {
		return 0, err
	}

	key := fmt.Sprintf("%s:%d", friendCacheKey, id)

	friend.ID = id
	_ = fr.cache.Set(ctx, key, *friend, 0)
	return id, nil
}

func (fr *friendService) FindById(ctx context.Context, id int64) (*entity.Friend, error) {
	key := fmt.Sprintf("%s:%d", friendCacheKey, id)

	cached, err := fr.cache.Get(ctx, key)
	if err == nil {
		return &cached, nil
	}

	friend, err := fr.datastore.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	_ = fr.cache.Set(ctx, key, *friend, 0)
	return friend, err
}

func (fr *friendService) FindByUserAndFriendId(
	ctx context.Context,
	userId, friendId int64,
) (*entity.Friend, error) {
	return fr.datastore.FindByUserAndFriendId(ctx, userId, friendId)
}

func (fr *friendService) Delete(ctx context.Context, id int64) error {
	key := fmt.Sprintf("%s:%d", friendCacheKey, id)

	if fr.cache.Exists(ctx, key) {
		_ = fr.cache.Delete(ctx, key)
	}

	return fr.datastore.Delete(ctx, id)
}
