package websocket

import (
	"context"
	"sync"

	domain_event "github.com/Weit145/template-go-monolit/internal/domain/event"
	"github.com/Weit145/template-go-monolit/internal/usecase/port"
	"github.com/google/uuid"
)

type HubInfastructure interface {
	Registre(client port.WebSocketClient) bool
	UnRegistre(id_conn uuid.UUID)
	Subscribe(id_conn, event uuid.UUID)
	UnSubscribe(id_conn, event uuid.UUID)
	Publish(message domain_event.Event)
	Shutdown(ctx context.Context) error
}

type Hub struct {
	mu      sync.RWMutex
	clients map[uuid.UUID]port.WebSocketClient
	event   map[uuid.UUID]map[uuid.UUID]struct{}
	closed  bool
	active  sync.WaitGroup
}

func NewHub() *Hub {
	return &Hub{
		mu:      sync.RWMutex{},
		clients: make(map[uuid.UUID]port.WebSocketClient),
		event:   make(map[uuid.UUID]map[uuid.UUID]struct{}),
	}
}

func (h *Hub) Registre(client port.WebSocketClient) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.closed {
		return false
	}
	if _, exists := h.clients[client.ClientID()]; exists {
		return false
	}
	h.clients[client.ClientID()] = client
	h.active.Add(1)
	return true
}

func (h *Hub) UnRegistre(id_conn uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, exists := h.clients[id_conn]; !exists {
		return
	}
	delete(h.clients, id_conn)
	for id_event, subscr := range h.event {
		delete(subscr, id_conn)
		if len(subscr) == 0 {
			delete(h.event, id_event)
		}
	}
	h.active.Done()
}

func (h *Hub) Subscribe(id_conn, event uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[id_conn]; !ok {
		return
	}

	if _, ok := h.event[event]; !ok {
		h.event[event] = make(map[uuid.UUID]struct{})
	}

	h.event[event][id_conn] = struct{}{}
}

func (h *Hub) UnSubscribe(id_conn, event uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.event[event][id_conn]; !ok {
		return
	}
	delete(h.event[event], id_conn)
	if len(h.event[event]) == 0 {
		delete(h.event, event)
	}
}

func (h *Hub) Publish(message domain_event.Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for connID := range h.event[message.ID] {
		client, ok := h.clients[connID]
		if !ok {
			continue
		}

		client.SendMessage(message)
	}
}

func (h *Hub) Shutdown(ctx context.Context) error {
	h.mu.Lock()
	h.closed = true
	clients := make([]port.WebSocketClient, 0, len(h.clients))
	for _, client := range h.clients {
		clients = append(clients, client)
	}
	h.mu.Unlock()

	for _, client := range clients {
		_ = client.Close()
	}
	done := make(chan struct{})
	go func() {
		h.active.Wait()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
