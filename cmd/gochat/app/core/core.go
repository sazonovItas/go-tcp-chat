package core

import (
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/service"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/internal/hasher"
)

type Core struct {
	// Services that using by app
	ConversationService service.ConversationService
	FriendService       service.FriendService
	MessageService      service.MessageService
	ParticipantService  service.ParticipantService
	UserService         service.UserService
	AuthService         service.AuthService

	// Additional functionality for the app
	Hasher hasher.Hasher
}
