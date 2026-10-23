package usecase

import (
	"context"
	"fmt"

	domain_error "github.com/Weit145/template-go-monolit/internal/domain/error"
	"github.com/Weit145/template-go-monolit/internal/domain/event"
	"github.com/Weit145/template-go-monolit/internal/logger"
	"github.com/google/uuid"
)

type WebSocketService interface {
	RegistreWebSocket(
		ctx context.Context,
		client event.WebSocketClient,
	) (err error)
	UnRegistreWebSocket(
		ctx context.Context,
		clientID uuid.UUID,
	) (err error)
	HandleWebSocketMessage(
		ctx context.Context,
		id_conn uuid.UUID,
		message event.Command,
	) error
}

func (s *Service) RegistreWebSocket(
	ctx context.Context,
	client event.WebSocketClient,
) (err error) {
	if !s.events.Registre(client) {
		return logger.WrapError(ctx, fmt.Errorf("register websocket: %w", domain_error.ErrUnavailable))
	}
	return nil
}

func (s *Service) UnRegistreWebSocket(
	ctx context.Context,
	clientID uuid.UUID,
) (err error) {
	s.events.UnRegistre(clientID)
	return nil
}

func (s *Service) HandleWebSocketMessage(
	ctx context.Context,
	id_conn uuid.UUID,
	message event.Command,
) error {
	if message.ID == uuid.Nil {
		return logger.WrapError(ctx, fmt.Errorf("empty event id: %w", domain_error.ErrInvalidInput))
	}
	switch message.Type {
	case "subscribe":
		s.events.Subscribe(id_conn, message.ID)
	case "unsubscribe":
		s.events.UnSubscribe(id_conn, message.ID)
	default:
		return logger.WrapError(ctx, fmt.Errorf("unknown websocket command %q: %w", message.Type, domain_error.ErrInvalidInput))
	}
	return nil
}
