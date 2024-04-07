package service

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/gofrs/uuid"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/domain/entity"
)

const (
	NewMessageEventType = "NewMessageEvent"
)

var ErrUnknownEventType = errors.New("unknown event type")

type EventService interface {
	EventBus

	// CreateEvent creates event from public event and returns it
	// Errors: ErrUnknownEventType, unknown
	CreateEvent(event *entity.PublicEvent) (*entity.Event, error)

	// CreatePublicEvent creates public event from event
	CreatePublicEvent(event *entity.Event) *entity.PublicEvent
}

type EventBus interface {
	// Subscribe adds a new subscriber for a given event type
	Subscribe(eventType string, subscriber chan<- entity.Event) int

	// Unsubscribe deletes a subscriber by id from a given event type
	Unsubscribe(eventType string, id int)

	// Publish sends an event to all subscribers of a given event type
	Publish(event entity.Event)
}

type eventService struct {
	EventBus
}

func NewEventService() EventService {
	return &eventService{
		EventBus: NewEventBus(),
	}
}

func (es *eventService) CreatePublicEvent(event *entity.Event) *entity.PublicEvent {
	return &entity.PublicEvent{
		Type:    event.Type,
		Payload: event.Payload,
	}
}

// CreateEvent is implementing interface EventService
func (es *eventService) CreateEvent(event *entity.PublicEvent) (*entity.Event, error) {
	data, err := json.Marshal(event.Payload)
	if err != nil {
		return nil, err
	}

	var payload interface{}
	switch event.Type {
	case NewMessageEventType:
		var newMessageEvent entity.NewMessageEvent
		if err := json.Unmarshal(data, &newMessageEvent); err != nil {
			return nil, err
		}
		payload = newMessageEvent
	default:
		return nil, ErrUnknownEventType
	}

	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	return &entity.Event{
		ID:        id,
		Type:      event.Type,
		Timestamp: time.Now(),
		Payload:   payload,
	}, nil
}

type eventBus struct {
	mu          sync.Mutex
	subscribers map[string][]chan<- entity.Event
}

func NewEventBus() EventBus {
	return &eventBus{
		subscribers: make(map[string][]chan<- entity.Event),
	}
}

// Subscribe is implementing interface EventBus
func (eb *eventBus) Subscribe(eventType string, subscriber chan<- entity.Event) int {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	for i := 0; i < len(eb.subscribers[eventType]); i++ {
		if eb.subscribers[eventType][i] == nil {
			eb.subscribers[eventType][i] = subscriber
			return i
		}
	}

	id := len(eb.subscribers[eventType])
	eb.subscribers[eventType] = append(eb.subscribers[eventType], subscriber)
	return id
}

// Unsubscribe is implementing interface EventBus
func (eb *eventBus) Unsubscribe(eventType string, id int) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.subscribers[eventType][id] = nil
}

// Publish is implementing interface EventBus
func (eb *eventBus) Publish(event entity.Event) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	subscribers := eb.subscribers[event.Type]
	for _, subscriber := range subscribers {
		subscriber <- event
	}
}
