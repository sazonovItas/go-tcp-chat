package entity

import "github.com/google/uuid"

type TokenID uuid.UUID

func (tid TokenID) String() string {
	return (uuid.UUID(tid)).String()
}

type Token struct {
	ID     TokenID `json:"id"`
	UserId int64   `json:"user_id"`
}
