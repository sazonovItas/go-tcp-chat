package service

import (
	"context"
	"errors"
	"time"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/entity"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/repo"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/internal/color"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/internal/hasher"
)

const DefaultTokenExpiration = time.Minute * 30

var (
	ErrMismatchedTokens = errors.New("mismatched tokens")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrInvalidToken     = errors.New("invalid token")
)

// AuthService is interface for managing user's authorization tokens
type AuthService interface {
	// SignUp sign up user and returns new user entity
	// Errors: unknown
	SignUp(ctx context.Context, authUser *entity.AuthUser) (*entity.User, error)

	// SignIn sign in user by login and password
	// Errors: ErrInvalidPassword, unknown
	SignIn(ctx context.Context, authUser *entity.AuthUser, user *entity.User) (entity.Token, error)

	// SignInByToken sign in user by token
	// Errors: ErrInvalidToken, unknown
	SignInByToken(ctx context.Context, authToken entity.Token) error

	// ValidateToken validates token
	// if token exists but not the same, then token would be deleted
	// Errors: ErrMismatchedTokens, ErrTokenNotFound, unknown
	ValidateToken(ctx context.Context, authToken entity.Token) error
}

type authService struct {
	hasher          hasher.Hasher
	tokenRepository repo.TokenRepository
}

func NewAuthService(tokenRepository repo.TokenRepository) AuthService {
	return &authService{
		tokenRepository: tokenRepository,
		hasher:          hasher.New(10),
	}
}

// SignUp is implementing interface AuthService
func (aus *authService) SignUp(
	ctx context.Context,
	authUser *entity.AuthUser,
) (*entity.User, error) {
	passwordHash, err := aus.hasher.Password(authUser.Password)
	if err != nil {
		return nil, err
	}

	err = aus.hasher.Compare(passwordHash, []byte(authUser.Password))
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Login:        authUser.Login,
		Name:         authUser.Login,
		Color:        color.GetRandomColorInHex(),
		PasswordHash: string(passwordHash),
	}
	return user, nil
}

// SignIn is implementing interface AuthService
func (aus *authService) SignIn(
	ctx context.Context,
	authUser *entity.AuthUser,
	user *entity.User,
) (entity.Token, error) {
	err := aus.hasher.Compare([]byte(user.PasswordHash), []byte(authUser.Password))
	if err != nil {
		switch {
		case errors.Is(err, hasher.ErrMismatchedPasswords):
			return entity.Token{}, ErrInvalidPassword
		default:
			return entity.Token{}, nil
		}
	}

	tk, err := aus.tokenRepository.CreateToken(ctx, user.ID)
	if err != nil {
		return entity.Token{}, err
	}

	err = aus.tokenRepository.SaveToken(ctx, tk, DefaultTokenExpiration)
	if err != nil {
		return entity.Token{}, nil
	}

	return tk, nil
}

// SignInByToken is implementing interface AuthService
func (aus *authService) SignInByToken(ctx context.Context, authToken entity.Token) error {
	err := aus.ValidateToken(ctx, authToken)
	if err != nil {
		switch {
		case errors.Is(err, ErrMismatchedTokens) || errors.Is(err, repo.ErrTokenNotFound):
			return ErrInvalidToken
		default:
			return err
		}
	}

	return nil
}

// ValidateToken is implementing interface AuthService
func (aus *authService) ValidateToken(ctx context.Context, authToken entity.Token) error {
	tk, err := aus.tokenRepository.TokenById(ctx, authToken.ID)
	if err != nil {
		return err
	}

	if tk != authToken {
		_ = aus.tokenRepository.DeleteToken(ctx, authToken.ID)
		return ErrMismatchedTokens
	}

	return nil
}
