package datastore

import (
	"context"
	"errors"
	"fmt"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/domain/model/entity"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/storage"
)

type ConversationDatastore interface {
	Create(ctx context.Context, conversation *entity.Conversation) (int64, error)
	FindById(ctx context.Context, id int64) (*entity.Conversation, error)
	Update(ctx context.Context, conversation *entity.Conversation) error
	Delete(ctx context.Context, id int64) error
}

type conversationDatastore struct {
	storage *storage.Storage
}

func NewConversationDatastore(db *storage.Storage) ConversationDatastore {
	return &conversationDatastore{storage: db}
}

var (
	ErrConversationUpdateFailed = errors.New("failed update conversation")
	ErrConversationDeleteFailed = errors.New("failed delete conversation")
)

// Create creates new conversation and returns conversation id
func (cs *conversationDatastore) Create(
	ctx context.Context,
	conversation *entity.Conversation,
) (int64, error) {
	const op = "gochat.internal.storage.models.conversation.CreateConversation"

	var id int64
	err := cs.storage.QueryRowContext(
		ctx,
		"INSERT INTO chat.conversations (title, conversation_kind, creator_id) VALUES ($1, $2, $3) RETURNING id",
		conversation.Title,
		conversation.ConversationKind,
		conversation.CreatorId,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// GetById returns conversation by id
func (cs *conversationDatastore) FindById(
	ctx context.Context,
	id int64,
) (*entity.Conversation, error) {
	const op = "gochat.internal.storage.models.conversation.GetConversationById"

	var conversation entity.Conversation
	err := cs.storage.Get(
		&conversation,
		"SELECT id, title, conversation_kind, creator_id FROM chat.conversations WHERE id=$1",
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &conversation, err
}

// UpdateConversation updates conversation's data
func (cs *conversationDatastore) Update(
	ctx context.Context,
	conversation *entity.Conversation,
) error {
	const op = "gochat.internal.storage.models.conversation.UpdateConversation"

	result, err := cs.storage.ExecContext(
		ctx,
		"UPDATE chat.users SET title=$1 WHERE id=$2",
		conversation.Title,
		conversation.ID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	res, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if res != 1 {
		return ErrConversationUpdateFailed
	}

	return nil
}

// DeleteConversation deletes conversation by id
func (cs *conversationDatastore) Delete(ctx context.Context, id int64) error {
	const op = "gochat.internal.storage.models.conversation.DeleteConversation"

	result, err := cs.storage.ExecContext(ctx, "DELETE FROM chat.conversations WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	res, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if res != 1 {
		return ErrConversationDeleteFailed
	}

	return nil
}
