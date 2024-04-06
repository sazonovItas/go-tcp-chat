package hasher

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrMismatchedPasswords = errors.New("mismatched passwords")

// Hasher is interface for hashing and comparing passwords
type Hasher interface {
	// Password returns a hashed version of the password
	// Errors: unknown
	Password(password string) ([]byte, error)

	// Compare compares passwords for mathes
	// Errors: ErrMismatchedPasswords, unknown
	Compare(hashedPassword []byte, password []byte) error
}

// hasher is implementing Hasher interface
type hasher struct {
	cost int
}

func New(cost int) Hasher {
	return &hasher{cost: cost}
}

// Password is implementing Hasher interface
func (h *hasher) Password(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}

// Compare is implementing Hasher interface
func (h *hasher) Compare(hashedPassword []byte, password []byte) error {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err != nil {
		switch {
		case err == bcrypt.ErrMismatchedHashAndPassword:
			return ErrMismatchedPasswords
		default:
			return err
		}
	}

	return nil
}
