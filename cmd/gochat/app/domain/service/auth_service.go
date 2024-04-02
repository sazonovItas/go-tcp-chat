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
)

// AuthService is interface for managing user's authorization tokens
type AuthService interface {
	SignUp(ctx context.Context, authUser *entity.AuthUser) (*entity.User, error)
	SignIn(ctx context.Context, authUser *entity.AuthUser, user *entity.User) (entity.Token, error)
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

// TODO: add color generation
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

func (aus *authService) ValidateToken(ctx context.Context, authToken entity.Token) error {
	tk, err := aus.tokenRepository.TokenById(ctx, authToken.ID)
	if err != nil {
		return err
	}

	if tk != authToken {
		// TODO: add deletion of the token
		return ErrMismatchedTokens
	}

	return nil
}
