package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/entity"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/storage"
)

type UserRepository interface {
	// Create creates new user and returns it's id
	// Errors: unknown
	Create(ctx context.Context, user *entity.User) (int64, error)

	// FindById returns user by id
	// Errors: ErrUserNotFound, unknown
	FindById(ctx context.Context, id int64) (*entity.User, error)

	// FindByLogin returns user by login
	// Errors: ErrUserNotFound, unknown
	FindByLogin(ctx context.Context, login string) (*entity.User, error)

	// Update updates user by id
	// Errors: ErrUserUpdateFailed
	Update(ctx context.Context, user *entity.User) error

	// Delete deletes user by id
	// Errors: ErrUserDeleteFailed
	Delete(ctx context.Context, id int64) error

	// GetIdByLogin returns user id by login
	GetIdByLogin(ctx context.Context, login string) int64

	// GetPublicUsersByConvId returns public users
	// Errors: unknown
	GetPublicUsers(ctx context.Context) ([]entity.PublicUser, error)
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
		case errors.Is(err, sql.ErrNoRows):
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
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrUserNotFound
		default:
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return &user, nil
}

// Update is implementing interface UserRepository
func (us *userRepository) Update(
	ctx context.Context,
	user *entity.User,
) error {
	const op = "gochat.internal.domain.repo.user_repo.Update"

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

// Delete is implementing interface UserRepository
func (us *userRepository) Delete(ctx context.Context, id int64) error {
	const op = "gochat.internal.domain.repo.user_repo.Delete"

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

// GetIdByLogin is implementing interface UserRepository
func (us *userRepository) GetIdByLogin(ctx context.Context, login string) int64 {
	var id int64

	err := us.storage.GetContext(ctx, &id, "SELECT id FROM chat.users WHERE login=$1", login)
	if err != nil {
		return 0
	}

	return id
}

// GetPublicUsers is implementing interface UserRepository
func (us *userRepository) GetPublicUsers(
	ctx context.Context,
) ([]entity.PublicUser, error) {
	const op = "gochat.internal.domain.repo.user_repo.GetPublicUsers"

	var users []entity.PublicUser
	err := us.storage.SelectContext(
		ctx,
		&users,
		`
      SELECT id, login, name, color FROM chat.users
    `,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}
