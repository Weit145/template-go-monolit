package v1_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Weit145/template-go-monolit/internal/domain/event"
	domain_user "github.com/Weit145/template-go-monolit/internal/domain/user"
	"github.com/Weit145/template-go-monolit/internal/infrastructure"
	socket "github.com/Weit145/template-go-monolit/internal/infrastructure/websocket"
	v1 "github.com/Weit145/template-go-monolit/internal/transport/api/v1"
	"github.com/Weit145/template-go-monolit/internal/usecase"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/google/uuid"
)

type testRepository struct {
	id uuid.UUID
}

func (testRepository) CheckHealth(context.Context) error { return nil }

func (r testRepository) CreateUser(_ context.Context, user domain_user.User) (domain_user.User, error) {
	return domain_user.UserFromDB(r.id, user.Name())
}

func (testRepository) DeleteUser(context.Context, uuid.UUID) error { return nil }

type testTransaction struct{}

func (testTransaction) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

type observedHub struct {
	*infrastructure.Infrastructure
	changes chan string
	left    chan struct{}
}

func (h *observedHub) Subscribe(clientID, eventID uuid.UUID) {
	h.Infrastructure.Subscribe(clientID, eventID)
	h.changes <- "subscribe"
}

func (h *observedHub) UnSubscribe(clientID, eventID uuid.UUID) {
	h.Infrastructure.UnSubscribe(clientID, eventID)
	h.changes <- "unsubscribe"
}

func (h *observedHub) UnRegistre(id uuid.UUID) {
	h.Infrastructure.UnRegistre(id)
	h.left <- struct{}{}
}

func TestWebSocketLifecycle(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	hub := &observedHub{Infrastructure: infrastructure.NewInfrastructure(), changes: make(chan string, 8), left: make(chan struct{}, 2)}
	watchedUserID := uuid.New()
	service := usecase.NewService(testRepository{id: watchedUserID}, testTransaction{}, hub)
	server := httptest.NewServer(v1.NewRouter(service, nil))
	defer server.Close()
	defer hub.Shutdown(ctx)

	conn, response, err := websocket.Dial(ctx, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws", nil)
	if err != nil {
		t.Fatalf("WebSocket upgrade failed: %v", err)
	}
	defer conn.CloseNow()
	if _, err := uuid.Parse(response.Header.Get("X-Request-ID")); err != nil {
		t.Fatalf("missing request ID after upgrade: %v", err)
	}
	first, second := uuid.New(), uuid.New()
	command := func(kind string, eventID uuid.UUID) {
		t.Helper()
		if err := wsjson.Write(ctx, conn, socket.ClientMessage{Type: kind, EventID: eventID}); err != nil {
			t.Fatal(err)
		}
		select {
		case got := <-hub.changes:
			if got != kind {
				t.Fatalf("command = %s, want %s", got, kind)
			}
		case <-ctx.Done():
			t.Fatal("subscription was not processed")
		}
	}
	command("subscribe", first)
	command("subscribe", second)
	publish := func(id uuid.UUID, data any) {
		hub.Publish(event.Event{ID: id, Type: "update", Data: data})
	}
	publish(uuid.New(), "not subscribed")
	publish(first, "delivered")
	var message socket.ServerMessage
	if err := wsjson.Read(ctx, conn, &message); err != nil {
		t.Fatal(err)
	}
	if message.EventID != first || message.Data != "delivered" || message.Type != "update" {
		t.Fatalf("unexpected message: %+v", message)
	}
	command("unsubscribe", first)
	publish(first, "must not arrive")
	publish(second, "still subscribed")
	if err := wsjson.Read(ctx, conn, &message); err != nil {
		t.Fatal(err)
	}
	if message.EventID != second || message.Data != "still subscribed" {
		t.Fatalf("unsubscribe did not take effect: %+v", message)
	}
	command("subscribe", watchedUserID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, server.URL+"/api/v1/users", strings.NewReader(`{"name":"Alice"}`))
	if err != nil {
		t.Fatal(err)
	}
	response, err = server.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	_ = response.Body.Close()
	if response.StatusCode != http.StatusCreated {
		t.Fatalf("create user status: %d", response.StatusCode)
	}
	if err := wsjson.Read(ctx, conn, &message); err != nil {
		t.Fatal(err)
	}
	data, ok := message.Data.(map[string]any)
	if !ok || data["name"] != "Alice" || message.EventID != watchedUserID || message.Type != "user.created_and_deleted" {
		t.Fatalf("API event was not delivered: %+v", message)
	}
	_ = conn.CloseNow()
	select {
	case <-hub.left:
	case <-ctx.Done():
		t.Fatal("disconnected client handler did not finish")
	}

	// A live idle connection must not prevent application shutdown.
	idle, _, err := websocket.Dial(ctx, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer idle.CloseNow()
	if err := hub.Shutdown(ctx); err != nil {
		t.Fatalf("shutdown did not finish: %v", err)
	}
	if _, _, err := idle.Read(ctx); err == nil {
		t.Fatal("connection remained open after shutdown")
	}
}
