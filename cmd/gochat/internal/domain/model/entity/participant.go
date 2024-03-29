package entity

import "time"

type Participant struct {
	ID             int64     `db:"id"`
	UserID         int64     `db:"user_id"`
	ConversationID int64     `db:"conversation_id"`
	UpdatedAt      time.Time `db:"updated_at"`
}
