package websocket

import (
	"sync"
	"testing"
	"time"

	"github.com/Weit145/template-go-monolit/internal/domain/event"
	"github.com/google/uuid"
)

func TestPublishWithFullQueueAndConcurrentSubscriptions(t *testing.T) {
	hub := NewHub()
	c := NewClient(nil)
	if !hub.Registre(c) {
		t.Fatal("registration failed")
	}
	eventID := uuid.New()
	publish := func(data any) {
		hub.Publish(event.Event{ID: eventID, Type: "update", Data: data})
	}
	hub.Subscribe(c.ID, eventID)
	for i := 0; i < cap(c.Send); i++ {
		publish(i)
	}
	done := make(chan struct{})
	go func() {
		publish("queue is full")
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("a full queue blocked publication")
	}
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				hub.Subscribe(c.ID, eventID)
				publish(j)
				hub.UnSubscribe(c.ID, eventID)
			}
		}()
	}
	hub.UnRegistre(c.ID)
	wg.Wait()
	if len(hub.clients) != 0 || len(hub.event) != 0 {
		t.Fatal("unregistered client or subscriptions remained in hub")
	}
}
