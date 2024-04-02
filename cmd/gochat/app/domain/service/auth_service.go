package service

import (
	"context"
	"errors"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/entity"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/repo"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/internal/hasher"
)

var ErrMismatchedTokens = errors.New("mismatched tokens")

// AuthService is interface for managing user's authorization tokens
// TODO: Divide auth service on auth service and token repo
type AuthService interface {
	SignUp(ctx context.Context, authUser *entity.AuthUser) (*entity.User, error)
	SignIn(ctx context.Context, authUser *entity.AuthUser) (entity.Token, error)
}

type authService struct {
	hasher.Hasher
	tokenRepository repo.TokenRepository
}

func NewAuthService(tokenRepository repo.TokenRepository) AuthService {
	return &authService{
		tokenRepository: tokenRepository,
		Hasher:          hasher.New(10),
	}
}

func (aus *authService) SignUp(
	ctx context.Context,
	authUser *entity.AuthUser,
) (*entity.User, error) {
	return nil, nil
}

func (aus *authService) SignIn(
	ctx context.Context,
	authUser *entity.AuthUser,
) (entity.Token, error) {
	return entity.Token{}, nil
}

func (aus *authService) Validate(ctx context.Context, authToken entity.Token) error {
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
