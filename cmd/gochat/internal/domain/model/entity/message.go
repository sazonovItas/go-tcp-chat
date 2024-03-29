package entity

import (
	"time"

	"github.com/google/uuid"
)

// MessageKind represents message kind
type MessageKind int

const (
	// UserAddingMessage represents message of adding people to conversation
	UserAddingMessage MessageKind = 0
	// UserMessage represents text message from user
	UserTextMessage MessageKind = 1
)

type Message struct {
	ID             uuid.UUID   `db:"id"              json:"id"`
	SenderID       int64       `db:"sender_id"       json:"sender_id"`
	ConversationID int64       `db:"conversation_id" json:"conversation_id"`
	MessageKind    MessageKind `db:"message_kind"    json:"message_kind"`
	Message        string      `db:"message"         json:"message"`
	CreatedAt      time.Time   `db:"created_at"      json:"created_at"`
	UpdatedAt      time.Time   `db:"updated_at"      json:"updated_at"`
}
