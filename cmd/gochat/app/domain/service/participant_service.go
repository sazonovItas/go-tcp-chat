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
	repository repo.ParticipantRepository
	cache      cache.Cache[entity.Participant]
}

func NewParticipantService(
	datastore repo.ParticipantRepository,
	opts *cache.CacheOpts,
) ParticipantService {
	return &participantService{
		repository: datastore,
		cache:      cache.NewCache[entity.Participant](opts),
	}
}

func (ps *participantService) Create(
	ctx context.Context,
	participant *entity.Participant,
) (int64, error) {
	id, err := ps.repository.Create(ctx, participant)
	if err != nil {
		return 0, err
	}

	key := fmt.Sprintf("%d", id)

	participant.ID = id
	_ = ps.cache.Set(ctx, key, *participant, 0)
	return id, nil
}

func (ps *participantService) FindById(ctx context.Context, id int64) (*entity.Participant, error) {
	key := fmt.Sprintf("%d", id)

	cached, err := ps.cache.Get(ctx, key)
	if err == nil {
		return &cached, nil
	}

	participant, err := ps.repository.FindById(ctx, id)
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
	return ps.repository.FindByUserAndConvId(ctx, userId, convId)
}

func (ps *participantService) Update(ctx context.Context, participant *entity.Participant) error {
	key := fmt.Sprintf("%d", participant.ID)

	if ps.cache.Exists(ctx, key) {
		_ = ps.cache.Set(ctx, key, *participant, 0)
	}

	return ps.repository.Update(ctx, participant)
}

func (ps *participantService) Delete(ctx context.Context, id int64) error {
	key := fmt.Sprintf("%d", id)

	if ps.cache.Exists(ctx, key) {
		_ = ps.cache.Delete(ctx, key)
	}

	return ps.repository.Delete(ctx, id)
}
