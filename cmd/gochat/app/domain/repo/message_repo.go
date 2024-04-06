package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/entity"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/storage"
)

type MessageRepository interface {
	// Create creates new message and returns it's id
	// Errors: ErrGenerateUUIDFailed, ErrMessageCreateFailed, unknown
	Create(ctx context.Context, msg *entity.Message) (uuid.UUID, error)

	// FindById finds message by id
	// Errors: ErrMessageNotFound, unknown
	FindById(ctx context.Context, id uuid.UUID) (*entity.Message, error)

	// Update updates message by id
	// Errors: ErrMessageUpdateFailed, unknown
	Update(ctx context.Context, message *entity.Message) error

	// Delete deletes message by id
	// Errors: ErrMessageDeleteFailed, unknown
	Delete(ctx context.Context, id uuid.UUID) error

	// GetConvMessagesPrevTimestamp returns limits count of messages previous to timestamp
	// Errors: ErrNoMessages, unknown
	GetConvMessagesPrevTimestamp(
		ctx context.Context,
		timestamp time.Time,
		limit int,
	) ([]entity.Message, error)

	// GetConvMessagesNextTimestamp returns limits count of messages next to timestamp
	// Errors: ErrNoMessages, unknown
	GetConvMessagesNextTimestamp(
		ctx context.Context,
		timestamp time.Time,
		limit int,
	) ([]entity.Message, error)

	// GetConvMessagesBetweenTimestamp returns limits count of messages next to timestamp
	// Errors: ErrNoMessages, unknown
	GetConvMessagesBetweenTimestamp(
		ctx context.Context,
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
	ErrGenerateUUIDFailed  = errors.New("failed generate uuid")
	ErrMessageCreateFailed = errors.New("failed create message")
	ErrMessageUpdateFailed = errors.New("failed update message")
	ErrMessageNotFound     = errors.New("message not found")
	ErrNoMessages          = errors.New("no messages")
	ErrMessageDeleteFailed = errors.New("failed to delete message")
)

// Create is implementing interface MessageRepository
func (ms *messageRepository) Create(
	ctx context.Context,
	msg *entity.Message,
) (uuid.UUID, error) {
	const op = "gochat.internal.domain.infastructure.datastore.Create"

	// Generate new uuid for message
	id, err := uuid.NewV4()
	if err != nil {
		return id, ErrGenerateUUIDFailed
	}

	msg.ID = id
	result, err := ms.storage.NamedExecContext(
		ctx,
		`
    INSERT INTO chat.messages (id, sender_id, message_kind, message, created_at, updated_at)
    VALUES (:id, :sender_id, :message_kind, :message, :created_at, :updated_at)
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

	return id, nil
}

// FindById is implementing interface MessageRepository
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

// Update is implementing interface MessageRepository
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

// Delete is implementing interface MessageRepository
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

// GetConvMessagesPrevTimestamp is implementing MessageRepository interface
func (ms *messageRepository) GetConvMessagesPrevTimestamp(
	ctx context.Context,
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
     SELECT id, sender_id, message_kind, message, created_at 
     FROM chat.messages 
     WHERE created_at<$1 
		 ORDER BY created_at DESC
     LIMIT $2
    )
    SELECT * FROM ready_messages ORDER BY created_at ASC
    `,
		timestamp,
		limit,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoMessages
		default:
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return messages, nil
}

// GetConvMessagesNextTimestamp is implementing interface MessageRepository
func (ms *messageRepository) GetConvMessagesNextTimestamp(
	ctx context.Context,
	timestamp time.Time,
	limit int,
) ([]entity.Message, error) {
	const op = "gochat.internal.domain.infastructure.datastore.message.GetConvMessagesNextTimestamp"

	var messages []entity.Message
	err := ms.storage.SelectContext(ctx,
		&messages,
		`
    SELECT id, sender_id, message_kind, message, created_at 
    FROM chat.messages 
    WHERE created_at>$1 
		ORDER BY created_at ASC
    LIMIT $2
    `,
		timestamp,
		limit,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoMessages
		default:
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return messages, nil
}

// GetConvMessagesBetweenTimestamp is implementing interface MessageRepository
func (ms *messageRepository) GetConvMessagesBetweenTimestamp(
	ctx context.Context,
	from, to time.Time,
) ([]entity.Message, error) {
	const op = "gochat.internal.domain.infastructure.datastore.message.GetConvMessagesBetweenTimestamp"

	var messages []entity.Message
	err := ms.storage.SelectContext(
		ctx,
		&messages,
		`
    SELECT id, sender_id, message_kind, message, created_at 
    FROM chat.messages 
    WHERE created_at BETWEEN $1 and $2
		ORDER BY created_at ASC
    `,
		from,
		to)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoMessages
		default:
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return messages, nil
}
