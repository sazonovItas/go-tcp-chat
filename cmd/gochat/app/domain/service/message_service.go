package service

import (
	"context"
	"time"

	"github.com/gofrs/uuid"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/entity"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/repo"
	"github.com/sazonovItas/gochat-tcp/pkg/cache"
)

type MessageService interface {
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

	// GetConvMessagesBetweenTimestamp returns messages between timestamp
	// Errors: ErrNoMessages, unknown
	GetConvMessagesBetweenTimestamp(
		ctx context.Context,
		from, to time.Time,
	) ([]entity.Message, error)
}

type messageService struct {
	repository repo.MessageRepository
	cache      cache.Cache[entity.Message]
}

func NewMessageService(repository repo.MessageRepository, opts *cache.CacheOpts) MessageService {
	return &messageService{
		repository: repository,
		cache:      nil,
	}
}

// Create is implementing interface MessageService
func (ms *messageService) Create(ctx context.Context, msg *entity.Message) (uuid.UUID, error) {
	return ms.repository.Create(ctx, msg)
}

// FindById is implementing interface MessageService
func (ms *messageService) FindById(ctx context.Context, id uuid.UUID) (*entity.Message, error) {
	return ms.repository.FindById(ctx, id)
}

// Update is implementing interface MessageService
func (ms *messageService) Update(ctx context.Context, msg *entity.Message) error {
	return ms.repository.Update(ctx, msg)
}

// Delete is implementing interface MessageService
func (ms *messageService) Delete(ctx context.Context, id uuid.UUID) error {
	return ms.repository.Delete(ctx, id)
}

// GetConvMessagesPrevTimestamp is implementing interface MessageService
func (ms *messageService) GetConvMessagesPrevTimestamp(
	ctx context.Context,
	timestamp time.Time,
	limit int,
) ([]entity.Message, error) {
	return ms.repository.GetConvMessagesPrevTimestamp(ctx, timestamp, limit)
}

// GetConvMessagesNextTimestamp is implementing interface MessageService
func (ms *messageService) GetConvMessagesNextTimestamp(
	ctx context.Context,
	timestamp time.Time,
	limit int,
) ([]entity.Message, error) {
	return ms.repository.GetConvMessagesNextTimestamp(ctx, timestamp, limit)
}

// GetConvMessagesBetweenTimestamp is implementing interface MessageService
func (ms *messageService) GetConvMessagesBetweenTimestamp(
	ctx context.Context,
	from, to time.Time,
) ([]entity.Message, error) {
	return ms.repository.GetConvMessagesBetweenTimestamp(ctx, from, to)
}
