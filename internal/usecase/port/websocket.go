package port

import (
	"context"

	"github.com/Weit145/template-go-monolit/internal/domain/event"
	"github.com/google/uuid"
)

type WebSocketClient interface {
	ClientID() uuid.UUID
	ReadeMessage(context.Context) (event.Command, error)
	WriteMessage(context.Context) error
	SendMessage(event.Event) bool
	Close() error
}

type EventHub interface {
	EventPublisher
	Registre(WebSocketClient) bool
	UnRegistre(id_conn uuid.UUID)
	Subscribe(id_conn, event uuid.UUID)
	UnSubscribe(id_conn, event uuid.UUID)
}
