package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/entity"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/storage"
)

type MessageRepository interface {
	Create(ctx context.Context, msg *entity.Message) (uuid.UUID, error)
	FindById(ctx context.Context, id uuid.UUID) (*entity.Message, error)
	Update(ctx context.Context, message *entity.Message) error
	Delete(ctx context.Context, id uuid.UUID) error

	GetConvMessagesPrevTimestamp(
		ctx context.Context,
		convId int64,
		timestamp time.Time,
		limit int,
	) ([]entity.Message, error)

	GetConvMessagesNextTimestamp(
		ctx context.Context,
		convId int64,
		timestamp time.Time,
		limit int,
	) ([]entity.Message, error)

	GetConvMessagesBetweenTimestamp(
		ctx context.Context,
		convId int64,
		from, to time.Time,
	) ([]entity.Message, error)
}

type messageRepository struct {
	storage *storage.Storage
}

func NewMessageRepository(db *storage.Storage) MessageRepository {
	return &messageRepository{storage: db}
}

var (
	ErrGenerateUUIDFailed          = errors.New("failed generate uuid")
	ErrMessageCreateFailed         = errors.New("failed create message")
	ErrMessageUpdateFailed         = errors.New("failed update message")
	ErrMessageNotFound             = errors.New("message not found")
	ErrNoNewMessagesInConversation = errors.New("no new messages in conversation")
	ErrMessageDeleteFailed         = errors.New("failed to delete message")
)

// CreateMessage creates new message and returns it's id
func (ms *messageRepository) Create(
	ctx context.Context,
	msg *entity.Message,
) (id uuid.UUID, err error) {
	const op = "gochat.internal.domain.infastructure.datastore.Create"

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
		`
    INSERT INTO chat.messages (id, sender_id, conversation_id, message_kind, message, created_at, updated_at)
    VALUES (:id, :sender_id, :conversation_id, :message_kind, :message, :created_at, :updated_at)
    `,
		msg,
	)
	if err != nil {
		return id, fmt.Errorf("%s: %w", op, err)
	}

	res, err := result.RowsAffected()
	if err != nil {
		return id, fmt.Errorf("%s: %w", op, err)
	}

	if res != 1 {
		return id, ErrMessageCreateFailed
	}

	return
}

// GetMessageById returns message by id
func (ms *messageRepository) FindById(
	ctx context.Context,
	id uuid.UUID,
) (*entity.Message, error) {
	const op = "gochat.internal.domain.infastructure.datastore.message.FindById"

	var msg entity.Message
	err := ms.storage.Get(
		&msg,
		"SELECT * FROM chat.messages WHERE id=$1",
		id,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrMessageNotFound
		default:
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return &msg, nil
}

// UpdateMessage updates message
func (ms *messageRepository) Update(ctx context.Context, message *entity.Message) error {
	const op = "gochat.internal.domain.infastructure.datastore.message.Update"

	result, err := ms.storage.ExecContext(
		ctx,
		"UPDATE chat.messages SET message=$1, updated_at=$2 WHERE id=$3",
		message.Message,
		message.CreatedAt,
		message.ID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	res, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if res != 1 {
		return ErrMessageUpdateFailed
	}

	return nil
}

// DeleteMessageId deletes message by id
func (ms *messageRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	const op = "gochat.internal.domain.infastructure.datastore.message.Delete"

	result, err := ms.storage.ExecContext(
		ctx,
		"DELETE FROM chat.messages WHERE id=$1",
		id,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	res, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if res != 1 {
		return ErrMessageDeleteFailed
	}

	return nil
}

// GetConvMessagesPrevTimestamp returns limit or less messages
// previous to timestamp from conversation with id equals to convId
func (ms *messageRepository) GetConvMessagesPrevTimestamp(
	ctx context.Context,
	convId int64,
	timestamp time.Time,
	limit int,
) ([]entity.Message, error) {
	const op = "gochat.internal.domain.infastructure.datastore.message.GetConvMessagesPrevTimestamp"

	var messages []entity.Message
	err := ms.storage.SelectContext(
		ctx,
		&messages,
		`
    WITH ready_messages AS (
     SELECT id, sender_id, conversation_id, message_kind, message, created_at 
     FROM chat.messages 
     WHERE conversation_id=$1 AND created_at<$2 
		 ORDER BY created_at DESC
     LIMIT $3
    )
    SELECT * FROM ready_messages ORDER BY created_at ASC
    `,
		convId,
		timestamp,
		limit,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoNewMessagesInConversation
		default:
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return messages, nil
}

// GetConvMessagesNextTimestamp returns limit or less messages
// next to timestamp from conversation with id equals to convId
func (ms *messageRepository) GetConvMessagesNextTimestamp(
	ctx context.Context,
	convId int64,
	timestamp time.Time,
	limit int,
) ([]entity.Message, error) {
	const op = "gochat.internal.domain.infastructure.datastore.message.GetConvMessagesNextTimestamp"

	var messages []entity.Message
	err := ms.storage.SelectContext(ctx,
		&messages,
		`
    SELECT id, sender_id, conversation_id, message_kind, message, created_at 
    FROM chat.messages 
    WHERE conversation_id=$1 AND created_at>$2 
		ORDER BY created_at ASC
    LIMIT $3
    `,
		convId,
		timestamp,
		limit,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoNewMessagesInConversation
		default:
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return messages, nil
}

func (ms *messageRepository) GetConvMessagesBetweenTimestamp(
	ctx context.Context,
	convId int64,
	from, to time.Time,
) ([]entity.Message, error) {
	const op = "gochat.internal.domain.infastructure.datastore.message.GetConvMessagesBetweenTimestamp"

	var messages []entity.Message
	err := ms.storage.SelectContext(
		ctx,
		&messages,
		`
    SELECT id, sender_id, conversation_id, message_kind, message, created_at 
    FROM chat.messages 
    WHERE created_at BETWEEN $1 and $2
		ORDER BY created_at ASC
    `,
		from,
		to)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoNewMessagesInConversation
		default:
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return messages, nil
}
