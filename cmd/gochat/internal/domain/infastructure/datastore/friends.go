package datastore

import (
	"context"
	"errors"
	"fmt"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/domain/model/entity"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/storage"
)

type FriendDatastore interface {
	Create(ctx context.Context, friend *entity.Friend) (int64, error)
	FindById(ctx context.Context, id int64) (*entity.Friend, error)
	FindByUserAndFriendId(ctx context.Context, userId, friendId int64) (*entity.Friend, error)
	Delete(ctx context.Context, id int64) error
}

type friendDatastore struct {
	storage *storage.Storage
}

func NewFriendDatastore(db *storage.Storage) FriendDatastore {
	return &friendDatastore{storage: db}
}

var ErrFriendDeleteFailed = errors.New("delete friend failed")

// CreateFriend creates friend and return it's id
func (fs *friendDatastore) Create(
	ctx context.Context,
	friend *entity.Friend,
) (int64, error) {
	const op = "gochat.internal.domain.infastructure.datastore.participant.CreateFriend"

	var id int64
	err := fs.storage.QueryRowContext(
		ctx,
		"INSERT INTO chat.friends (id, user_id, friend_id) VALUES ($1, $2, $3) RETURNING id",
		friend.ID,
		friend.UserID,
		friend.FriendID,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// FindById returns friend by unique friend id
func (fs *friendDatastore) FindById(ctx context.Context, id int64) (*entity.Friend, error) {
	const op = "gochat.internal.domain.infastructure.datastore.participant.FindById"

	var friend entity.Friend
	err := fs.storage.Get(
		&friend,
		"SELECT id, user_id, friend_id FROM chat.friends WHERE id=$1",
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &friend, nil
}

// FindByUserAndFriendId returns friend by user id and friend id
func (fs *friendDatastore) FindByUserAndFriendId(
	ctx context.Context,
	userId, friendId int64,
) (*entity.Friend, error) {
	const op = "gochat.internal.domain.infastructure.datastore.participant.FindByUserAndFriendId"

	var friend entity.Friend
	err := fs.storage.Get(
		&friend,
		"SELECT id, user_id, friend_id FROM chat.friends WHERE user_id=$1 AND friend_id=$2",
		userId, friendId,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &friend, nil
}

// Delete deletes friend by id
func (fs *friendDatastore) Delete(ctx context.Context, id int64) error {
	const op = "gochat.internal.domain.infastructure.datastore.participant.Delete"

	result, err := fs.storage.ExecContext(ctx, "DELETE FROM chat.friends WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	res, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if res != 1 {
		return ErrFriendDeleteFailed
	}

	return nil
}
