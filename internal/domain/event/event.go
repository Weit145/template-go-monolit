package event

import (
	"context"

	"github.com/google/uuid"
)

type WebSocketClient interface {
	ClientID() uuid.UUID
	ReadeMessage(context.Context) (Command, error)
	WriteMessage(context.Context) error
	SendMessage(Event) bool
	Close() error
}

type Event struct {
	ID   uuid.UUID
	Type string
	Data any
}

type Command struct {
	ID   uuid.UUID
	Type string
}
