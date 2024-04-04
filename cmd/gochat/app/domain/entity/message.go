package entity

import (
	"time"

	"github.com/gofrs/uuid"
)

// MessageKind represents message kind
type MessageKind int

const (
	// CreateConversationMessage represetns a message of creating conversation
	CreateConversationMessage MessageKind = 0
	// AddingUserMessage represents a message of adding people to conversation
	AddingUserMessage MessageKind = 1
	// UserTextMessage represents a text message from user
	UserTextMessage MessageKind = 2
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

type MarshalMessage struct {
	ID             string      `json:"id"`
	SenderID       int64       `json:"sender_id"`
	ConversationID int64       `json:"conversation_id"`
	MessageKind    MessageKind `json:"message_kind"`
	Message        string      `json:"message"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}
