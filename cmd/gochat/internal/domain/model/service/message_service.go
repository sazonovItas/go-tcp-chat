package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/domain/infastructure/datastore"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/domain/model/entity"
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
	datastore datastore.MessageDatastore
	cache     cache.Cache[entity.Participant]
}

func NewMessageService(datastore datastore.MessageDatastore, opts *cache.CacheOpts) MessageService {
	return &messageService{
		datastore: datastore,
		cache:     nil,
	}
}

func (ms *messageService) Create(ctx context.Context, msg *entity.Message) (uuid.UUID, error) {
	return ms.datastore.Create(ctx, msg)
}

func (ms *messageService) FindById(ctx context.Context, id uuid.UUID) (*entity.Message, error) {
	return ms.datastore.FindById(ctx, id)
}

func (ms *messageService) Update(ctx context.Context, msg *entity.Message) error {
	return ms.datastore.Update(ctx, msg)
}

func (ms *messageService) Delete(ctx context.Context, id uuid.UUID) error {
	return ms.datastore.Delete(ctx, id)
}

func (ms *messageService) GetConvMessagesPrevTimestamp(
	ctx context.Context,
	convId int64,
	timestamp time.Time,
	limit int,
) ([]entity.Message, error) {
	return ms.datastore.GetConvMessagesPrevTimestamp(ctx, convId, timestamp, limit)
}

func (ms *messageService) GetConvMessagesNextTimestamp(
	ctx context.Context,
	convId int64,
	timestamp time.Time,
	limit int,
) ([]entity.Message, error) {
	return ms.datastore.GetConvMessagesNextTimestamp(ctx, convId, timestamp, limit)
}

func (ms *messageService) GetConvMessagesBetweenTimestamp(
	ctx context.Context,
	convId int64,
	from, to time.Time,
) ([]entity.Message, error) {
	return ms.datastore.GetConvMessagesBetweenTimestamp(ctx, convId, from, to)
}
