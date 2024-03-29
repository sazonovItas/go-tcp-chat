package event

import "time"

type (
	EventPayload interface{}

	Event interface {
		IDer
		EventName() string
		Payload() EventPayload
		Metadata() Metadata
		OccuredAt() time.Time
	}

	event struct {
		Entity
		payload   EventPayload
		metadata  Metadata
		occuredAt time.Time
	}
)

func (e event) EventName() string          { return e.name }
func (e event) EventPayload() EventPayload { return e.payload }
func (e event) EventMetadata() Metadata    { return e.metadata }
func (e event) OccuredAt() time.Time       { return e.occuredAt }
