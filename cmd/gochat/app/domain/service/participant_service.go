package service

import (
	"context"
	"fmt"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/entity"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/repo"
	"github.com/sazonovItas/gochat-tcp/pkg/cache"
)

type ParticipantService interface {
	Create(ctx context.Context, participant *entity.Participant) (int64, error)
	FindById(ctx context.Context, id int64) (*entity.Participant, error)
	FindByUserAndConvId(ctx context.Context, userId, convId int64) (*entity.Participant, error)
	Update(ctx context.Context, participant *entity.Participant) error
	Delete(ctx context.Context, id int64) error
}

type participantService struct {
	datastore repo.ParticipantRepository
	cache     cache.Cache[entity.Participant]
}

func NewParticipantService(
	datastore repo.ParticipantRepository,
	opts *cache.CacheOpts,
) ParticipantService {
	return &participantService{
		datastore: datastore,
		cache:     cache.NewCache[entity.Participant](opts),
	}
}

const participantCacheKey = "participant"

func (ps *participantService) Create(
	ctx context.Context,
	participant *entity.Participant,
) (int64, error) {
	id, err := ps.datastore.Create(ctx, participant)
	if err != nil {
		return 0, err
	}

	key := fmt.Sprintf("%s:%d", participantCacheKey, id)

	participant.ID = id
	_ = ps.cache.Set(ctx, key, *participant, 0)
	return id, nil
}

func (ps *participantService) FindById(ctx context.Context, id int64) (*entity.Participant, error) {
	key := fmt.Sprintf("%s:%d", participantCacheKey, id)

	cached, err := ps.cache.Get(ctx, key)
	if err == nil {
		return &cached, nil
	}

	participant, err := ps.datastore.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	_ = ps.cache.Set(ctx, key, *participant, 0)
	return participant, nil
}

func (ps *participantService) FindByUserAndConvId(
	ctx context.Context,
	userId, convId int64,
) (*entity.Participant, error) {
	return ps.datastore.FindByUserAndConvId(ctx, userId, convId)
}

func (ps *participantService) Update(ctx context.Context, participant *entity.Participant) error {
	key := fmt.Sprintf("%s:%d", participantCacheKey, participant.ID)

	if ps.cache.Exists(ctx, key) {
		_ = ps.cache.Set(ctx, key, *participant, 0)
	}

	return ps.datastore.Update(ctx, participant)
}

func (ps *participantService) Delete(ctx context.Context, id int64) error {
	key := fmt.Sprintf("%s:%d", participantCacheKey, id)

	if ps.cache.Exists(ctx, key) {
		_ = ps.cache.Delete(ctx, key)
	}

	return ps.datastore.Delete(ctx, id)
}
