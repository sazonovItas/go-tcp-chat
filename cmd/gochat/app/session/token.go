package session

import "github.com/google/uuid"

type SessionToken struct {
	UUID         uuid.UUID `json:"uuid"`
	ID           int64     `json:"id"`
	Login        string    `json:"login"`
	PasswordHash string    `json:"password_hash"`
}

func NewSessionToken(login, passwordHash string) SessionToken {
	return SessionToken{
		UUID:         uuid.New(),
		Login:        login,
		PasswordHash: passwordHash,
	}
}
