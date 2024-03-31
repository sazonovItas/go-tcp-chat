package service

import (
	"context"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/domain/infastructure/datastore"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/domain/model/entity"
)

type ConversationService interface {
	Create(ctx context.Context, conversation *entity.Conversation) (int64, error)
	FindById(ctx context.Context, id int64) (*entity.Conversation, error)
	Update(ctx context.Context, conversation *entity.Conversation) error
	Delete(ctx context.Context, id int64) error
}

type conversationService struct {
	datastore datastore.ConversationDatastore
}

func NewConversationService(datastore datastore.ConversationDatastore) ConversationService {
	return &conversationService{datastore: datastore}
}

func (cr *conversationService) Create(
	ctx context.Context,
	conversation *entity.Conversation,
) (int64, error) {
	return cr.datastore.Create(ctx, conversation)
}

func (cr *conversationService) FindById(
	ctx context.Context,
	id int64,
) (*entity.Conversation, error) {
	return cr.datastore.FindById(ctx, id)
}

func (cr *conversationService) Update(
	ctx context.Context,
	conversation *entity.Conversation,
) error {
	return cr.datastore.Update(ctx, conversation)
}

func (cr *conversationService) Delete(ctx context.Context, id int64) error {
	return cr.datastore.Delete(ctx, id)
}
