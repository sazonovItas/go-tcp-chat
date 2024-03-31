package service

import (
	"context"
	"fmt"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/domain/infastructure/datastore"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/domain/model/entity"
	"github.com/sazonovItas/gochat-tcp/pkg/cache"
)

type ConversationService interface {
	Create(ctx context.Context, conversation *entity.Conversation) (int64, error)
	FindById(ctx context.Context, id int64) (*entity.Conversation, error)
	Update(ctx context.Context, conversation *entity.Conversation) error
	Delete(ctx context.Context, id int64) error
}

type conversationService struct {
	datastore datastore.ConversationDatastore
	cache     cache.Cache[entity.Conversation]
}

func NewConversationService(
	datastore datastore.ConversationDatastore,
	opts *cache.CacheOpts,
) ConversationService {
	return &conversationService{
		datastore: datastore,
		cache:     cache.NewCache[entity.Conversation](opts),
	}
}

const conversationCacheKey = "conversation"

func (cr *conversationService) Create(
	ctx context.Context,
	conversation *entity.Conversation,
) (int64, error) {
	id, err := cr.datastore.Create(ctx, conversation)
	if err != nil {
		return 0, err
	}

	key := fmt.Sprintf("%s:%d", conversationCacheKey, id)

	conversation.ID = id
	_ = cr.cache.Set(ctx, key, *conversation, 0)
	return id, nil
}

func (cr *conversationService) FindById(
	ctx context.Context,
	id int64,
) (*entity.Conversation, error) {
	key := fmt.Sprintf("%s:%d", conversationCacheKey, id)

	cached, err := cr.cache.Get(ctx, key)
	if err == nil {
		return &cached, nil
	}

	conversation, err := cr.datastore.FindById(ctx, id)
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
	key := fmt.Sprintf("%s:%d", conversationCacheKey, conversation.ID)

	if cr.cache.Exists(ctx, key) {
		_ = cr.cache.Set(ctx, key, *conversation, 0)
	}

	return cr.datastore.Update(ctx, conversation)
}

func (cr *conversationService) Delete(ctx context.Context, id int64) error {
	key := fmt.Sprintf("%s:%d", conversationCacheKey, id)

	if cr.cache.Exists(ctx, key) {
		_ = cr.cache.Delete(ctx, key)
	}

	return cr.datastore.Delete(ctx, id)
}
