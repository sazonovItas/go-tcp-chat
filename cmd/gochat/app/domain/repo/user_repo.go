package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/entity"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/storage"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) (int64, error)
	FindById(ctx context.Context, id int64) (*entity.User, error)
	FindByLogin(ctx context.Context, login string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id int64) error

	GetIdByLogin(ctx context.Context, login string) int64
	GetPublicUsersByConvId(ctx context.Context, convId int64) ([]entity.PublicUser, error)
}

type userRepository struct {
	storage *storage.Storage
}

func NewUserRepository(db *storage.Storage) UserRepository {
	return &userRepository{storage: db}
}

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserDeleteFailed = errors.New("failed delete user")
	ErrUserUpdateFailed = errors.New("failed update user")
)

// CreateUser creates new user and returns user id
func (us *userRepository) Create(ctx context.Context, user *entity.User) (int64, error) {
	const op = "gochat.internal.domain.infastructure.datastore.user.Create"

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
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// FindById returns user model struct by id
func (us *userRepository) FindById(ctx context.Context, id int64) (*entity.User, error) {
	const op = "gochat.internal.domain.infastructure.datastore.user.FindById"

	var user entity.User
	err := us.storage.GetContext(
		ctx,
		&user,
		"SELECT id, login, name, color, password_hash FROM chat.users WHERE id=$1",
		id,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrUserNotFound
		default:
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return &user, nil
}

// GetUserByLogin returns user mode struct by login
func (us *userRepository) FindByLogin(ctx context.Context, login string) (*entity.User, error) {
	const op = "gochat.internal.domain.infastructure.datastore.user.FindByLogin"

	var user entity.User
	err := us.storage.GetContext(
		ctx,
		&user,
		"SELECT id, login, name, color, password_hash FROM chat.users WHERE login=$1",
		login,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrUserNotFound
		default:
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return &user, nil
}

// UpdateUser updates user's data
func (us *userRepository) Update(
	ctx context.Context,
	user *entity.User,
) error {
	const op = "gochat.internal.domain.infastructure.datastore.user.Update"

	result, err := us.storage.ExecContext(
		ctx,
		"UPDATE chat.users SET name=$1, color=$2 WHERE id=$3",
		user.Name,
		user.Color,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	res, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if res != 1 {
		return ErrUserUpdateFailed
	}

	return nil
}

// Delete deletes user byid
func (us *userRepository) Delete(ctx context.Context, id int64) error {
	const op = "gochat.internal.domain.infastructure.datastore.user.Delete"

	result, err := us.storage.ExecContext(ctx, "DELETE FROM chat.users WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	res, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if res != 1 {
		return ErrUserDeleteFailed
	}

	return nil
}

// GetUserId returns user ID or 0 if user does not exists
func (us *userRepository) GetIdByLogin(ctx context.Context, login string) int64 {
	var id int64

	err := us.storage.GetContext(ctx, &id, "SELECT id FROM chat.users WHERE login=$1", login)
	if err != nil {
		return 0
	}

	return id
}

func (us *userRepository) GetPublicUsersByConvId(
	ctx context.Context,
	convId int64,
) ([]entity.PublicUser, error) {
	const op = "gochat.internal.domain.infastructure.datastore.user.GetPublicUsersByConvId"

	var users []entity.PublicUser
	err := us.storage.SelectContext(
		ctx,
		&users,
		`
    WITH (
      SELECT user_id FROM chat.participants WHERE conversation_id=$1
    ) AS users_ids
    SELECT U.id, U.login, U.name, U.color 
    FROM chat.users AS U
      JOIN users_ids AS UID ON U.id=UID.user_id
    `,
		convId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}
