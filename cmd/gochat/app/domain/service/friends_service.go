package service

import (
	"context"
	"fmt"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/entity"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/repo"
	"github.com/sazonovItas/gochat-tcp/pkg/cache"
)

type FriendService interface {
	Create(ctx context.Context, friend *entity.Friend) (int64, error)
	FindById(ctx context.Context, id int64) (*entity.Friend, error)
	FindByUserAndFriendId(ctx context.Context, userId, friendId int64) (*entity.Friend, error)
	Delete(ctx context.Context, id int64) error
}

type friendService struct {
	repository repo.FriendRepository
	cache      cache.Cache[entity.Friend]
}

func NewFriendService(repository repo.FriendRepository, opts *cache.CacheOpts) FriendService {
	return &friendService{
		repository: repository,
		cache:      cache.NewCache[entity.Friend](opts),
	}
}

func (fr *friendService) Create(ctx context.Context, friend *entity.Friend) (int64, error) {
	id, err := fr.repository.Create(ctx, friend)
	if err != nil {
		return 0, err
	}

	key := fmt.Sprintf("%d", id)

	friend.ID = id
	_ = fr.cache.Set(ctx, key, *friend, 0)
	return id, nil
}

func (fr *friendService) FindById(ctx context.Context, id int64) (*entity.Friend, error) {
	key := fmt.Sprintf("%d", id)

	cached, err := fr.cache.Get(ctx, key)
	if err == nil {
		return &cached, nil
	}

	friend, err := fr.repository.FindById(ctx, id)
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
	return fr.repository.FindByUserAndFriendId(ctx, userId, friendId)
}

func (fr *friendService) Delete(ctx context.Context, id int64) error {
	key := fmt.Sprintf("%d", id)

	if fr.cache.Exists(ctx, key) {
		_ = fr.cache.Delete(ctx, key)
	}

	return fr.repository.Delete(ctx, id)
}
