package models

import (
	"context"
	"errors"
	"fmt"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/storage"
)

type APIUser struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type User struct {
	ID           int64  `db:"id"`
	Login        string `db:"login"`
	Name         string `db:"name"`
	Color        string `db:"color"`
	PasswordHash string `db:"password_hash"`
}

type UpdateUser struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

var (
	ErrUserDoesNotExists = errors.New("user does not exists")
	ErrUserDeleteFailed  = errors.New("failed delete user")
	ErrUserUpdateFailed  = errors.New("failed update user")
)

// CreateUser creates new user and returns user id
func (us *UserStorage) CreateUser(ctx context.Context, user *User) (int64, error) {
	const op = "gochat.internal.storage.models.user.CreateUser"

	var id int64
	err := us.storage.QueryRowContext(
		ctx,
		"INSERT INTO chat.users (name, login, color, password_hash) VALUES ($1, $2, $3, $4) RETURNING id",
		user.Name,
		user.Login,
		user.Color,
		user.PasswordHash,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	return id, nil
}

// GetUserId returns user ID or 0 if user does not exists
func (us *UserStorage) GetUserId(ctx context.Context, login string) int64 {
	var id int64

	err := us.storage.GetContext(ctx, &id, "SELECT id FROM chat.users WHERE login=$1", login)
	if err != nil {
		return 0
	}

	return id
}

// GetUserById returns user model struct by id
func (us *UserStorage) GetUserById(ctx context.Context, id int64) (*User, error) {
	const op = "gochat.internal.storage.models.user.GetUserById"

	var user User
	err := us.storage.GetContext(
		ctx,
		&user,
		"SELECT id, login, name, color, password_hash FROM chat.users WHERE id=$1",
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return &user, nil
}

// GetUserByLogin returns user mode struct by login
func (us *UserStorage) GetUserByLogin(ctx context.Context, login string) (*User, error) {
	const op = "gochat.internal.storage.models.user.GetUserByLogin"

	var user User
	err := us.storage.GetContext(
		ctx,
		&user,
		"SELECT id, login, name, color, password_hash FROM chat.users WHERE login=$1",
		login,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return &user, nil
}

// UpdateUser updates user's data
// TODO: maybe return changed user
func (us *UserStorage) UpdateUser(
	ctx context.Context,
	updateUser *UpdateUser,
) error {
	const op = "gochat.internal.storage.models.user.UpdateUserById"

	result, err := us.storage.ExecContext(
		ctx,
		"UPDATE chat.users SET name=$1, color=$2 WHERE id=$3",
		updateUser.Name,
		updateUser.Color,
		updateUser.ID,
	)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	res, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	if res != 1 {
		return ErrUserUpdateFailed
	}

	return nil
}

// DeleteUserById deletes user by id
func (us *UserStorage) DeleteUserId(ctx context.Context, userId int64) error {
	const op = "gochat.internal.storage.models.user.DeleteUserById"

	result, err := us.storage.ExecContext(ctx, "DELETE FROM chat.users WHERE id=$1", userId)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	res, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	if res != 1 {
		return ErrUserDeleteFailed
	}

	return nil
}

type UserStorage struct {
	storage *storage.Storage
}

func NewUserStore(db *storage.Storage) *UserStorage {
	return &UserStorage{storage: db}
}
