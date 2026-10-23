package port

import (
	"github.com/Weit145/template-go-monolit/internal/domain/event"
	"github.com/google/uuid"
)

type WebSocketHub interface {
	Publish(message event.Event)
	Registre(event.WebSocketClient) bool
	UnRegistre(id_conn uuid.UUID)
	Subscribe(id_conn, event uuid.UUID)
	UnSubscribe(id_conn, event uuid.UUID)
}
type EventHub interface {
	WebSocketHub
}
