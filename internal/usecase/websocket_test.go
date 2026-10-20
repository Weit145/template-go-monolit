package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	domain_error "github.com/Weit145/template-go-monolit/internal/domain/error"
	"github.com/Weit145/template-go-monolit/internal/domain/event"
	"github.com/google/uuid"
)

func TestHandleWebSocketMessage(t *testing.T) {
	connID, objectID := uuid.New(), uuid.New()
	for _, test := range []struct {
		name    string
		command event.Command
		wantErr bool
	}{
		{name: "subscribe", command: event.Command{ID: objectID, Type: "subscribe"}},
		{name: "unsubscribe", command: event.Command{ID: objectID, Type: "unsubscribe"}},
		{name: "unknown command", command: event.Command{ID: objectID, Type: "unknown"}, wantErr: true},
		{name: "empty object ID", command: event.Command{Type: "subscribe"}, wantErr: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			hub := &eventTestPublisher{}
			service := NewService(nil, nil, hub)
			err := service.HandleWebSocketMessage(context.Background(), connID, test.command)
			if test.wantErr {
				if !errors.Is(err, domain_error.ErrInvalidInput) || len(hub.commands) != 0 {
					t.Fatalf("invalid command reached infrastructure: err=%v, commands=%v", err, hub.commands)
				}
				return
			}
			if err != nil || len(hub.commands) != 1 || hub.commands[0] != test.command || hub.connID != connID {
				t.Fatalf("command was not forwarded correctly: err=%v, commands=%v", err, hub.commands)
			}
		})
	}
}

type websocketTestClient struct {
	id         uuid.UUID
	readErr    error
	writeErr   error
	closed     bool
	writerDone chan struct{}
}

func (c *websocketTestClient) ClientID() uuid.UUID          { return c.id }
func (c *websocketTestClient) SendMessage(event.Event) bool { return true }
func (c *websocketTestClient) Close() error {
	c.closed = true
	return nil
}
func (c *websocketTestClient) ReadeMessage(ctx context.Context) (event.Command, error) {
	if c.readErr != nil {
		return event.Command{}, c.readErr
	}
	<-ctx.Done()
	return event.Command{}, ctx.Err()
}
func (c *websocketTestClient) WriteMessage(ctx context.Context) error {
	defer close(c.writerDone)
	if c.writeErr != nil {
		return c.writeErr
	}
	<-ctx.Done()
	return ctx.Err()
}

func TestWebSocketStopsBothLoopsOnFailure(t *testing.T) {
	for _, stage := range []string{"read", "write"} {
		t.Run(stage, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			failure := errors.New("connection failed")
			client := &websocketTestClient{id: uuid.New(), writerDone: make(chan struct{})}
			if stage == "read" {
				client.readErr = failure
			} else {
				client.writeErr = failure
			}
			hub := &eventTestPublisher{}
			service := NewService(nil, nil, hub)
			err := service.WebSocket(ctx, client)
			if !errors.Is(err, failure) || !client.closed || !hub.unregistered {
				t.Fatalf("session cleanup failed: err=%v, closed=%v, unregistered=%v", err, client.closed, hub.unregistered)
			}
			select {
			case <-client.writerDone:
			default:
				t.Fatal("writer was not stopped before returning")
			}
		})
	}
}
