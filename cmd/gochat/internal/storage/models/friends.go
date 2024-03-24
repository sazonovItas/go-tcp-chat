package models

import (
	"context"
	"errors"
	"fmt"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/storage"
)

type Friend struct {
	ID       int64 `db:"id"`
	UserID   int64 `db:"user_id"`
	FriendID int64 `db:"friend_id"`
}

var ErrFriendDeleteFailed = errors.New("delete friend failed")

// CreateFriend creates friend and return it's id
func (fs *FriendStorage) CreateFriend(
	ctx context.Context,
	friend Friend,
) (int64, error) {
	const op = "gochat.internal.storage.models.participant.CreateFriend"

	var id int64
	err := fs.storage.QueryRowContext(
		ctx,
		"INSERT INTO chat.friends (id, user_id, friend_id) VALUES ($1, $2, $3) RETURNING id",
		friend.ID,
		friend.UserID,
		friend.FriendID,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	return id, nil
}

func (fs *FriendStorage) GetFriendById(ctx context.Context, id int64) (*Friend, error) {
	const op = "gochat.internal.storage.models.participant.GetFriendById"

	var friend Friend
	err := fs.storage.Get(
		&friend,
		"SELECT id, user_id, friend_id FROM chat.friends WHERE id=$1",
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return &friend, nil
}

func (fs *FriendStorage) DeleteFriendById(ctx context.Context, id int64) error {
	const op = "gochat.internal.storage.models.participant.GetFriendById"

	result, err := fs.storage.ExecContext(ctx, "DELETE FROM chat.friends WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	res, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	if res != 1 {
		return ErrFriendDeleteFailed
	}

	return nil
}

type FriendStorage struct {
	storage *storage.Storage
}

func NewFriendStorage(db *storage.Storage) *FriendStorage {
	return &FriendStorage{storage: db}
}
