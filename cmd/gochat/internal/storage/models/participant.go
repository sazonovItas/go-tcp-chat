package models

import (
	"context"
	"errors"
	"fmt"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/storage"
)

type Participant struct {
	ID             int64 `db:"id"`
	UserID         int64 `db:"user_id"`
	ConversationID int64 `db:"conversation_id"`
}

var ErrParticipantDeleteFailed = errors.New("delete participant failed")

// CreateParticipant creates participant and return it's id
func (ps *ParticipantStorage) CreateParticipant(
	ctx context.Context,
	participant Participant,
) (int64, error) {
	const op = "gochat.internal.storage.models.participant.CreateParticipant"

	var id int64
	err := ps.storage.QueryRowContext(
		ctx,
		"INSERT INTO chat.participants (id, user_id, conversation_id) VALUES ($1, $2, $3) RETURNING id",
		participant.ID,
		participant.UserID,
		participant.ConversationID,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	return id, nil
}

func (ps *ParticipantStorage) GetParticipantById(
	ctx context.Context,
	id int64,
) (*Participant, error) {
	const op = "gochat.internal.storage.models.participant.GetParticipantById"

	var participant Participant
	err := ps.storage.Get(
		&participant,
		"SELECT id, user_id, conversation_id FROM chat.conversations WHERE id=$1",
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return &participant, nil
}

func (ps *ParticipantStorage) DeleteParticipantById(ctx context.Context, id int64) error {
	const op = "gochat.internal.storage.models.participant.DeleteParticipantById"

	result, err := ps.storage.ExecContext(ctx, "DELETE FROM chat.participants WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	res, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	if res != 1 {
		return ErrParticipantDeleteFailed
	}

	return nil
}

// DeleteParticipantByUserConvId deletes participant by user and conversation id
func (ps *ParticipantStorage) DeleteParticipantByUserConvId(
	ctx context.Context,
	userId int64,
	convId int64,
) error {
	const op = "gochat.internal.storage.models.participant.DeleteParticipantByUserConvId"

	result, err := ps.storage.ExecContext(
		ctx,
		"DELETE FROM chat.participants WHERE user_id=$1 AND conversation_id=$2",
		userId,
		convId,
	)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	res, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	if res != 1 {
		return ErrParticipantDeleteFailed
	}

	return nil
}

type ParticipantStorage struct {
	storage *storage.Storage
}

func NewParticipantStore(db *storage.Storage) *ParticipantStorage {
	return &ParticipantStorage{storage: db}
}
