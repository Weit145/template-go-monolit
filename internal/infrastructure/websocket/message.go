package websocket

import (
	"github.com/Weit145/template-go-monolit/internal/domain/event"
	"github.com/google/uuid"
)

type ClientMessage struct {
	Type    string    `json:"type"`
	EventID uuid.UUID `json:"event_id"`
}

type ServerMessage struct {
	Type    string    `json:"type"`
	EventID uuid.UUID `json:"event_id"`
	Data    any       `json:"data"`
}

func NewServerMessage(e event.Event) ServerMessage {
	return ServerMessage{Type: e.Type, EventID: e.ID, Data: e.Data}
}
