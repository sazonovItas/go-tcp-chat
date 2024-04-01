package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/entity"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/repo"
	"github.com/sazonovItas/gochat-tcp/pkg/cache"
)

type MessageService interface {
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

type messageService struct {
	repository repo.MessageRepository
	cache      cache.Cache[entity.Participant]
}

func NewMessageService(datastore repo.MessageRepository, opts *cache.CacheOpts) MessageService {
	return &messageService{
		repository: datastore,
		cache:      nil,
	}
}

func (ms *messageService) Create(ctx context.Context, msg *entity.Message) (uuid.UUID, error) {
	return ms.repository.Create(ctx, msg)
}

func (ms *messageService) FindById(ctx context.Context, id uuid.UUID) (*entity.Message, error) {
	return ms.repository.FindById(ctx, id)
}

func (ms *messageService) Update(ctx context.Context, msg *entity.Message) error {
	return ms.repository.Update(ctx, msg)
}

func (ms *messageService) Delete(ctx context.Context, id uuid.UUID) error {
	return ms.repository.Delete(ctx, id)
}

func (ms *messageService) GetConvMessagesPrevTimestamp(
	ctx context.Context,
	convId int64,
	timestamp time.Time,
	limit int,
) ([]entity.Message, error) {
	return ms.repository.GetConvMessagesPrevTimestamp(ctx, convId, timestamp, limit)
}

func (ms *messageService) GetConvMessagesNextTimestamp(
	ctx context.Context,
	convId int64,
	timestamp time.Time,
	limit int,
) ([]entity.Message, error) {
	return ms.repository.GetConvMessagesNextTimestamp(ctx, convId, timestamp, limit)
}

func (ms *messageService) GetConvMessagesBetweenTimestamp(
	ctx context.Context,
	convId int64,
	from, to time.Time,
) ([]entity.Message, error) {
	return ms.repository.GetConvMessagesBetweenTimestamp(ctx, convId, from, to)
}
