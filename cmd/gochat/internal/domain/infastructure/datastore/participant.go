package datastore

import (
	"context"
	"errors"
	"fmt"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/domain/model/entity"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/storage"
)

type ParticipantDatastore interface {
	Create(ctx context.Context, participant *entity.Participant) (int64, error)
	FindById(ctx context.Context, id int64) (*entity.Participant, error)
	FindByUserAndConvId(ctx context.Context, userId, convId int64) (*entity.Participant, error)
	Delete(ctx context.Context, id int64) error
}

type participantDatastore struct {
	storage *storage.Storage
}

func NewParticipantDatastore(db *storage.Storage) ParticipantDatastore {
	return &participantDatastore{storage: db}
}

var ErrParticipantDeleteFailed = errors.New("delete participant failed")

// CreateParticipant creates participant and return it's id or 0
func (ps *participantDatastore) Create(
	ctx context.Context,
	participant *entity.Participant,
) (int64, error) {
	const op = "gochat.internal.domain.infastructure.datastore.participant.Create"

	var id int64
	err := ps.storage.QueryRowContext(
		ctx,
		"INSERT INTO chat.participants (id, user_id, conversation_id) VALUES ($1, $2, $3) RETURNING id",
		participant.ID,
		participant.UserID,
		participant.ConversationID,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// FindById returns participant by id
func (ps *participantDatastore) FindById(
	ctx context.Context,
	id int64,
) (*entity.Participant, error) {
	const op = "gochat.internal.domain.infastructure.datastore.participant.FindById"

	var participant entity.Participant
	err := ps.storage.Get(
		&participant,
		"SELECT id, user_id, conversation_id FROM chat.conversations WHERE id=$1",
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &participant, nil
}

// FindByUserAndConvId deletes participant by user and conversation id
func (ps *participantDatastore) FindByUserAndConvId(
	ctx context.Context,
	userId, convId int64,
) (*entity.Participant, error) {
	const op = "gochat.internal.domain.infastructure.datastore.participant.FindByUserAndConvId"

	var participant entity.Participant
	err := ps.storage.GetContext(
		ctx,
		&participant,
		"SELECT id, user_id, conversation_id FROM chat.participants WHERE user_id=$1 AND conversation_id=$2",
		userId,
		convId,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &participant, nil
}

// Delete deletes participant by id
func (ps *participantDatastore) Delete(ctx context.Context, id int64) error {
	const op = "gochat.internal.storage.models.participant.DeleteParticipantById"

	result, err := ps.storage.ExecContext(ctx, "DELETE FROM chat.participants WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	res, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if res != 1 {
		return ErrParticipantDeleteFailed
	}

	return nil
}
