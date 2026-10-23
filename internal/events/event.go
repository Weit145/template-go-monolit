package events

import (
	"context"

	"github.com/Weit145/template-go-monolit/internal/domain/event"
	"github.com/Weit145/template-go-monolit/internal/events/websocket"
	"github.com/google/uuid"
)

type EventBus interface {
	websocket.Hub
}

type eventBus struct {
	hub websocket.Hub
}

func NewEventBus() EventBus {
	return &eventBus{
		hub: websocket.NewHub(),
	}
}

func (e *eventBus) Registre(client event.WebSocketClient) bool {
	return e.hub.Registre(client)
}

func (e *eventBus) UnRegistre(id_conn uuid.UUID) {
	e.hub.UnRegistre(id_conn)
}

func (e *eventBus) Subscribe(id_conn, event uuid.UUID) {
	e.hub.Subscribe(id_conn, event)
}

func (e *eventBus) UnSubscribe(id_conn, event uuid.UUID) {
	e.hub.UnSubscribe(id_conn, event)
}

func (e *eventBus) Publish(message event.Event) {
	e.hub.Publish(message)
}

func (e *eventBus) Shutdown(ctx context.Context) error {
	return e.hub.Shutdown(ctx)
}
