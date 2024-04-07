package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/entity"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/service"
	tcpws "github.com/sazonovItas/gochat-tcp/internal/server"
)

const (
	ProtoNotSupported   = "not supported protocol"
	ReadyForMessages    = "ready for messages"
	UnauthorizedMessage = "token expired"
)

// /api/v1/chatting
func (api *Api) Chatting(resp *tcpws.Response, req *tcpws.Request) {
	const op = "gochat.app.api.chatting.Chatting"

	if req.Proto != tcpws.ProtoWS {
		resp.StatusCode = http.StatusBadRequest
		resp.Status = ProtoNotSupported
		return
	}

	var token entity.Token
	if err := json.Unmarshal([]byte(req.Body), &token); err != nil {
		resp.StatusCode = http.StatusBadRequest
		resp.Status = ProtoNotSupported
		return
	}

	if err := api.app.AuthService.ValidateToken(req.Ctx(), token); err != nil {
		resp.StatusCode = http.StatusUnauthorized
		resp.Status = UnauthorizedMessage
		return
	}

	resp.StatusCode = http.StatusOK
	resp.Status = ReadyForMessages
	if err := resp.Write(); err != nil {
		return
	}

	eventch := make(chan entity.Event, 5)
	subscriberId := api.app.EventService.Subscribe(service.NewMessageEventType, eventch)
	defer func() {
		api.app.EventService.Unsubscribe(service.NewMessageEventType, subscriberId)
		close(eventch)
	}()

	api.app.Logger.Info(
		"new subscriber on NewMessageEvent",
		"subscriber_id",
		subscriberId,
		"user_id",
		token.UserId,
	)

	stopch := make(chan struct{})
	defer close(stopch)

	// send events
	go func() {
		for {
			select {
			case <-stopch:
				return
			case event := <-eventch:
				api.app.Logger.Info("new event", "event", event)
				sendEvent := api.app.EventService.CreatePublicEvent(&event)

				msg, err := json.Marshal(*sendEvent)
				if err != nil {
					api.app.Logger.Error(
						"json marshal send event",
						"error",
						fmt.Errorf("%s: %w", op, err).Error(),
					)
					continue
				}

				_, err = resp.Conn.Write(msg)
				if err != nil {
					api.app.Logger.Error(
						"message send",
						"error",
						fmt.Errorf("%s: %w", op, err).Error(),
					)
					continue
				}
			}
		}
	}()

	// read events
	api.app.Logger.Info("start read messages")
	for {
		frame, err := resp.Conn.ReadFrame()
		if err != nil {
			switch {
			case errors.Is(err, io.EOF):
				api.app.Logger.Info("disconnection from user", "user_id", token.UserId)
			default:
				api.app.Logger.Error("read frame", "error", err.Error())
			}
			stopch <- struct{}{}
			break
		}

		var receivedEvent entity.PublicEvent
		if err := json.Unmarshal(frame, &receivedEvent); err != nil {
			api.app.Logger.Error("json unmarshal", "error", fmt.Errorf("%s: %w", op, err).Error())
			continue
		}

		event, err := api.app.EventService.CreateEvent(&receivedEvent)
		if err != nil {
			api.app.Logger.Error("create event", "error", fmt.Errorf("%s: %w", op, err).Error())
			continue
		}

		if msg, ok := event.Payload.(entity.NewMessageEvent); ok {
			msg.CreatedAt = time.Now()
			msg.UpdateAt = time.Now()
			id, err := api.app.MessageService.Create(context.Background(), &entity.Message{
				SenderID:    msg.SenderID,
				MessageKind: msg.MessageKind,
				Message:     msg.Message,
				CreatedAt:   msg.CreatedAt,
			})
			if err != nil {
				api.app.Logger.Error("create msg", "error", fmt.Errorf("%s: %w", op, err).Error())
				continue
			}

			msg.ID = id.String()
			event.Payload = msg
		}

		api.app.EventService.Publish(*event)
	}
}
