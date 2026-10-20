package port

import "github.com/Weit145/template-go-monolit/internal/domain/event"

type EventPublisher interface {
	Publish(event.Event)
}
