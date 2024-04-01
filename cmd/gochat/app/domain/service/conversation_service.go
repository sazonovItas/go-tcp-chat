package service

import (
	"context"
	"fmt"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/entity"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/repo"
	"github.com/sazonovItas/gochat-tcp/pkg/cache"
)

type ConversationService interface {
	Create(ctx context.Context, conversation *entity.Conversation) (int64, error)
	FindById(ctx context.Context, id int64) (*entity.Conversation, error)
	Update(ctx context.Context, conversation *entity.Conversation) error
	Delete(ctx context.Context, id int64) error
}

type conversationService struct {
	repository repo.ConversationRepository
	cache      cache.Cache[entity.Conversation]
}

func NewConversationService(
	repository repo.ConversationRepository,
	opts *cache.CacheOpts,
) ConversationService {
	return &conversationService{
		repository: repository,
		cache:      cache.NewCache[entity.Conversation](opts),
	}
}

func (cr *conversationService) Create(
	ctx context.Context,
	conversation *entity.Conversation,
) (int64, error) {
	id, err := cr.repository.Create(ctx, conversation)
	if err != nil {
		return 0, err
	}

	key := fmt.Sprintf("%d", id)

	conversation.ID = id
	_ = cr.cache.Set(ctx, key, *conversation, 0)
	return id, nil
}

func (cr *conversationService) FindById(
	ctx context.Context,
	id int64,
) (*entity.Conversation, error) {
	key := fmt.Sprintf("%d", id)

	cached, err := cr.cache.Get(ctx, key)
	if err == nil {
		return &cached, nil
	}

	conversation, err := cr.repository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	_ = cr.cache.Set(ctx, key, *conversation, 0)
	return conversation, nil
}

func (cr *conversationService) Update(
	ctx context.Context,
	conversation *entity.Conversation,
) error {
	key := fmt.Sprintf("%d", conversation.ID)

	if cr.cache.Exists(ctx, key) {
		_ = cr.cache.Set(ctx, key, *conversation, 0)
	}

	return cr.repository.Update(ctx, conversation)
}

func (cr *conversationService) Delete(ctx context.Context, id int64) error {
	key := fmt.Sprintf("%d", id)

	if cr.cache.Exists(ctx, key) {
		_ = cr.cache.Delete(ctx, key)
	}

	return cr.repository.Delete(ctx, id)
}
