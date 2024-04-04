package entity

// ConversationKind represents kind of conversation
type ConversationKind int

const (
	// Conversation2P2Kind represents person to person conversation kind
	Conversation2P2Kind ConversationKind = 0
	// ConversationGroupKind represents group conversation kind
	ConversationGroupKind ConversationKind = 1
)

type Conversation struct {
	ID               int64            `db:"id"                json:"id"`
	Title            string           `db:"title"             json:"title"`
	Color            string           `db:"color"             json:"color"`
	ConversationKind ConversationKind `db:"conversation_kind" json:"conversation_kind"`
	CreatorId        int64            `db:"creator_id"        json:"creator_id"`
}
