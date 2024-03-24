package models

import (
	"context"
	"errors"
	"fmt"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/storage"
)

// ConversationKind represents kind of conversation
type ConversationKind int

const (
	// Conversation2P2Kind represents person to person conversation kind
	Conversation2P2Kind ConversationKind = 0
	// ConversationGroupKind represents group conversation kind
	ConversationGroupKind ConversationKind = 1
)

type Conversation struct {
	ID               int64            `db:"id"`
	Title            string           `db:"title"`
	ConversationType ConversationKind `db:"conversation_type"`
	CreatorId        int64            `db:"creator_id"`
}

type UpdateConversation struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

var (
	ErrConversationUpdateFailed = errors.New("failed update conversation")
	ErrConversationDeleteFailed = errors.New("failed delete conversation")
)

// CreateConversation creates new conversation and returns conversation id
func (cs *ConversationStorage) CreateConversation(
	ctx context.Context,
	conversation *Conversation,
) (int64, error) {
	const op = "gochat.internal.storage.models.conversation.CreateConversation"

	var id int64
	err := cs.storage.QueryRowContext(
		ctx,
		"INSERT INTO chat.conversations (title, conversation_type, creator_id) VALUES ($1, $2, $3) RETURNING id",
		conversation.Title,
		conversation.ConversationType,
		conversation.CreatorId,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	return id, nil
}

// GetConversationById returns conversation by id
func (cs *ConversationStorage) GetConversationById(
	ctx context.Context,
	id int64,
) (*Conversation, error) {
	const op = "gochat.internal.storage.models.conversation.GetConversationById"

	var conversation Conversation
	err := cs.storage.Get(
		&conversation,
		"SELECT id, title, conversation_type, creator_id FROM chat.conversations WHERE id=$1",
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return &conversation, err
}

// GetConversationsByTitle returns slice of conversation by title
func (cs *ConversationStorage) GetConversationsByTitle(
	ctx context.Context,
	title string,
) ([]Conversation, error) {
	const op = "gochat.internal.storage.models.conversation.GetConversationsByTitle"

	var conversations []Conversation
	err := cs.storage.SelectContext(
		ctx,
		conversations,
		"SELECT id, title, conversation_type, creator_id FROM conversations WHERE title=$1",
		title,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return conversations, nil
}

// UpdateConversation updates conversation's data
// TODO: maybe return changed conversation
func (cs *ConversationStorage) UpdateConversation(
	ctx context.Context,
	updateConversation *UpdateConversation,
) error {
	const op = "gochat.internal.storage.models.conversation.UpdateConversation"

	result, err := cs.storage.ExecContext(
		ctx,
		"UPDATE chat.users SET title=$1 WHERE id=$2",
		updateConversation.Title,
		updateConversation.ID,
	)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	res, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	if res != 1 {
		return ErrConversationUpdateFailed
	}

	return nil
}

// DeleteConversation deletes conversation by id
func (cs *ConversationStorage) DeleteConversationById(ctx context.Context, id int64) error {
	const op = "gochat.internal.storage.models.conversation.DeleteConversation"

	result, err := cs.storage.ExecContext(ctx, "DELETE FROM chat.conversations WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	res, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	if res != 1 {
		return ErrConversationDeleteFailed
	}

	return nil
}

type ConversationStorage struct {
	storage *storage.Storage
}

func NewConversationStorage(db *storage.Storage) *ConversationStorage {
	return &ConversationStorage{storage: db}
}
