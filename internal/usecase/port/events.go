package port

import (
	"github.com/Weit145/template-go-monolit/internal/domain/event"
	"github.com/google/uuid"
)

type WebSocketHub interface {
	Publish(message event.Event)
	Register(event.WebSocketClient) bool
	Unregister(id_conn uuid.UUID)
	Subscribe(id_conn, event uuid.UUID)
	Unsubscribe(id_conn, event uuid.UUID)
}
type EventHub interface {
	WebSocketHub
}
