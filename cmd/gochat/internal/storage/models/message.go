package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/storage"
)

// MessageKind represents message kind
type MessageKind int

const (
	// UserAddingMessage represents message of adding people to conversation
	UserAddingMessage MessageKind = 0
	// UserMessage represents text message from user
	UserTextMessage MessageKind = 1
)

var (
	ErrGenerateUUIDFailed  = errors.New("failed generate uuid")
	ErrMessageCreateFailed = errors.New("failed create message")
	ErrMessageUpdateFailed = errors.New("failed update message")
	ErrMessageDeleteFailed = errors.New("failed to delete message")
)

type Message struct {
	ID             uuid.UUID   `db:"id"`
	SenderID       int64       `db:"sender_id"`
	ConversationID int64       `db:"conversation_id"`
	MessageType    MessageKind `db:"message_type"`
	Message        string      `db:"message"`
	CreatedAt      time.Time   `db:"created_at"`
}

type UpdateMessage struct {
	ID        uuid.UUID `json:"id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateMessage creates new message and returns it's id
func (ms *MessageStorage) CreateMessage(
	ctx context.Context,
	msg *Message,
) (id uuid.UUID, err error) {
	const op = "gochat.internal.storage.models.message.CreateMessage"

	defer func() {
		if r := recover(); r != nil {
			err = ErrGenerateUUIDFailed
		}
	}()

	// Generate new uuid for message
	id = uuid.New()

	msg.ID = id
	result, err := ms.storage.NamedExecContext(
		ctx,
		"INSERT INTO chat.messages (id, sender_id, conversation_id, message_type, message, created_at) VALUES (:id, :sender_id, :conversation_id, :message_type, :message, :created_at)",
		msg,
	)
	if err != nil {
		return id, fmt.Errorf("%s: %v", op, err)
	}

	res, err := result.RowsAffected()
	if err != nil {
		return id, fmt.Errorf("%s: %v", op, err)
	}

	if res != 1 {
		return id, ErrMessageCreateFailed
	}

	return
}

// GetMessageById returns message by id
func (ms *MessageStorage) GetMessageById(
	ctx context.Context,
	id uuid.UUID,
) (*Message, error) {
	const op = "gochat.internal.storage.models.message.GetMessageById"

	var msg Message
	err := ms.storage.Get(
		&msg,
		"SELECT id, sender_id, conversation_id, message_type, message, created_at FROM chat.messages WHERE id=$1",
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return &msg, nil
}

// UpdateMessage updates message
// TODO: maybe return changed message
func (ms *MessageStorage) UpdateMessage(ctx context.Context, updateMessage *UpdateMessage) error {
	const op = "gochat.internal.storage.models.message.UpdateMessage"

	result, err := ms.storage.ExecContext(
		ctx,
		"UPDATE chat.messages SET message=$1, created_at=$2 WHERE id=$3",
		updateMessage.Message,
		updateMessage.CreatedAt,
		updateMessage.ID,
	)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	res, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	if res != 1 {
		return ErrMessageUpdateFailed
	}

	return nil
}

// DeleteMessageId deletes message by id
func (ms *MessageStorage) DeleteMessageId(
	ctx context.Context,
	id uuid.UUID,
) error {
	const op = "gochat.internal.storage.models.message.UpdateMessage"

	result, err := ms.storage.ExecContext(
		ctx,
		"DELETE FROM chat.messages WHERE id=$1",
		id,
	)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	res, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	if res != 1 {
		return ErrMessageDeleteFailed
	}

	return nil
}

type APIMessage struct {
	ID          uuid.UUID   `json:"id"`
	Sender      APIUser     `json:"user"`
	MessageType MessageKind `json:"message_kind"`
	Message     string      `json:"message"`
	CreatedAt   time.Time   `json:"created_at"`
}

// APIMessagesFromTimestampByConvId returns sorted asc slice of api messages by timestamp
// with conversation id = convId from given timestamp to latest messages with limit
func (ms *MessageStorage) APIMessagesFromTimestampByConvId(
	ctx context.Context,
	convId int64,
	from time.Time,
	limit int,
) ([]APIMessage, error) {
	const op = "gochat.internal.storage.models.message.APIMessagesFromTimestampByConvId"

	rows, err := ms.storage.QueryContext(ctx, `
    WITH msgs AS (
      SELECT id AS msg_id, sender_id, message_type, message, created_at FROM chat.messages 
      WHERE conversation_id=$1 ORDER BY created_at ASC LIMIT $2)
    SELECT  msgs.msg_id, 
            users.id, 
            users.name, 
            users.color, 
            msgs.message_type, 
            msgs.message, 
            msgs.created_at 
    FROM msgs
      JOIN chat.users AS users
        ON msgs.sender_id = users.id 
    `, convId, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	defer rows.Close()

	var (
		messages []APIMessage
		msg      APIMessage
	)
	for rows.Next() {
		err := rows.Scan(
			&msg.ID,
			&msg.Sender.ID,
			&msg.Sender.Name,
			&msg.Sender.Color,
			&msg.MessageType,
			&msg.Message,
			&msg.CreatedAt,
		)
		if err != nil {
			return messages, err
		}

		messages = append(messages, msg)
	}

	return messages, nil
}

// GetLastAPIMessagesByConvId returns sorted asc slice of last api messages in conversation
func (ms *MessageStorage) GetLastAPIMessagesByConvId(
	ctx context.Context,
	convId int64,
	count int,
) ([]APIMessage, error) {
	const op = "gochat.internal.storage.models.message.APIMessagesByConvIdAndTimestamp"

	rows, err := ms.storage.QueryContext(ctx, `
    WITH msgs AS (
      SELECT id AS msg_id, sender_id, message_type, message, created_at FROM chat.messages 
      WHERE conversation_id=$1 ORDER BY created_at DESC LIMIT $2)
    SELECT  msgs.msg_id, 
            users.id, 
            users.name, 
            users.color, 
            msgs.message_type, 
            msgs.message, 
            msgs.created_at 
    FROM msgs
      JOIN chat.users AS users
        ON msgs.sender_id = users.id 
    ORDER BY msgs.created_at ASC
    `, convId, count)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	defer rows.Close()

	var (
		messages []APIMessage
		msg      APIMessage
	)
	for rows.Next() {
		err := rows.Scan(
			&msg.ID,
			&msg.Sender.ID,
			&msg.Sender.Name,
			&msg.Sender.Color,
			&msg.MessageType,
			&msg.Message,
			&msg.CreatedAt,
		)
		if err != nil {
			return messages, err
		}

		messages = append(messages, msg)
	}

	return messages, nil
}

type MessageStorage struct {
	storage *storage.Storage
}

func NewMessageStorage(db *storage.Storage) *MessageStorage {
	return &MessageStorage{storage: db}
}
