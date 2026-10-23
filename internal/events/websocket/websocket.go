package websocket

import (
	"context"
	"sync"

	"github.com/Weit145/template-go-monolit/internal/domain/event"
	domain_event "github.com/Weit145/template-go-monolit/internal/domain/event"
	"github.com/google/uuid"
)

type Hub interface {
	Registre(client event.WebSocketClient) bool
	UnRegistre(id_conn uuid.UUID)
	Subscribe(id_conn, event uuid.UUID)
	UnSubscribe(id_conn, event uuid.UUID)
	Publish(message domain_event.Event)
	Shutdown(ctx context.Context) error
}

type hub struct {
	mu      sync.RWMutex
	clients map[uuid.UUID]event.WebSocketClient
	event   map[uuid.UUID]map[uuid.UUID]struct{}
	closed  bool
	active  sync.WaitGroup
}

func NewHub() Hub {
	return &hub{
		mu:      sync.RWMutex{},
		clients: make(map[uuid.UUID]event.WebSocketClient),
		event:   make(map[uuid.UUID]map[uuid.UUID]struct{}),
	}
}

func (h *hub) Registre(client event.WebSocketClient) bool {
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

func (h *hub) UnRegistre(id_conn uuid.UUID) {
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

func (h *hub) Subscribe(id_conn, event uuid.UUID) {
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

func (h *hub) UnSubscribe(id_conn, event uuid.UUID) {
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

func (h *hub) Publish(message domain_event.Event) {
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

func (h *hub) Shutdown(ctx context.Context) error {
	h.mu.Lock()
	h.closed = true
	clients := make([]event.WebSocketClient, 0, len(h.clients))
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
