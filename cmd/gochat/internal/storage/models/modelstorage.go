package models

import "github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/storage"

type ModelStorage struct {
	*UserStorage
	*ConversationStorage
	*MessageStorage
	*ParticipantStorage
	*FriendStorage
}

func NewModelStorage(db *storage.Storage) *ModelStorage {
	return &ModelStorage{
		UserStorage:         NewUserStore(db),
		ConversationStorage: NewConversationStorage(db),
		MessageStorage:      NewMessageStorage(db),
		ParticipantStorage:  NewParticipantStore(db),
		FriendStorage:       NewFriendStorage(db),
	}
}
