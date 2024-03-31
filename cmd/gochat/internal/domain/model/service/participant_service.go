package service

import (
	"context"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/domain/infastructure/datastore"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/domain/model/entity"
)

type ParticipantService interface {
	Create(ctx context.Context, participant *entity.Participant) (int64, error)
	FindById(ctx context.Context, id int64) (*entity.Participant, error)
	FindByUserAndConvId(ctx context.Context, userId, convId int64) (*entity.Participant, error)
	Update(ctx context.Context, participant *entity.Participant) error
	Delete(ctx context.Context, id int64) error
}

type participantService struct {
	datastore datastore.ParticipantDatastore
}

func NewParticipantService(datastore datastore.ParticipantDatastore) ParticipantService {
	return &participantService{datastore: datastore}
}

func (ps *participantService) Create(
	ctx context.Context,
	participant *entity.Participant,
) (int64, error) {
	return ps.datastore.Create(ctx, participant)
}

func (ps *participantService) FindById(ctx context.Context, id int64) (*entity.Participant, error) {
	return ps.datastore.FindById(ctx, id)
}

func (ps *participantService) FindByUserAndConvId(
	ctx context.Context,
	userId, convId int64,
) (*entity.Participant, error) {
	return ps.datastore.FindByUserAndConvId(ctx, userId, convId)
}

func (ps *participantService) Update(ctx context.Context, participant *entity.Participant) error {
	return ps.datastore.Update(ctx, participant)
}

func (ps *participantService) Delete(ctx context.Context, id int64) error {
	return ps.datastore.Delete(ctx, id)
}
