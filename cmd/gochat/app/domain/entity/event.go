package entity

import (
	"time"

	"github.com/gofrs/uuid"
)

type Event struct {
	ID        uuid.UUID
	Type      string
	Timestamp time.Time
	Payload   interface{}
}

type PublicEvent struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type NewMessageEvent struct {
	ID          string      `json:"id"`
	SenderID    int64       `json:"sender_id"`
	MessageKind MessageKind `json:"message_kind"`
	Message     string      `json:"message"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdateAt    time.Time   `json:"updated_at"`
}
