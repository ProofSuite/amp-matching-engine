package types

import "fmt"

type EventType string

type Event struct {
	eventType EventType `json:"eventType"`
	payload   interface{}
}

func (e *Event) String() string {
	return fmt.Sprintf("\nEvent Type: %v\nPayload: %v\n\n", e.eventType, e.payload)
}
