package usecase

import (
	"context"
	"errors"
	"fmt"

	domain_error "github.com/Weit145/template-go-monolit/internal/domain/error"
	"github.com/Weit145/template-go-monolit/internal/domain/event"
	"github.com/Weit145/template-go-monolit/internal/logger"
	"github.com/Weit145/template-go-monolit/internal/usecase/port"
	"github.com/google/uuid"
)

type WebSocketService interface {
	WebSocket(
		ctx context.Context,
		client port.WebSocketClient,
	) (resultErr error)
	HandleWebSocketMessage(
		ctx context.Context,
		id_conn uuid.UUID,
		message event.Command,
	) error
}

func (s *Service) WebSocket(
	ctx context.Context,
	client port.WebSocketClient,
) (resultErr error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if !s.events.Registre(client) {
		return logger.WrapError(ctx, fmt.Errorf("register websocket: %w", domain_error.ErrUnavailable))
	}
	writerDone := make(chan error, 1)
	defer func() {
		cancel()
		_ = client.Close()
		writeErr := <-writerDone
		s.events.UnRegistre(client.ClientID())
		if writeErr != nil && !errors.Is(writeErr, context.Canceled) {
			if errors.Is(resultErr, context.Canceled) {
				resultErr = writeErr
			} else {
				resultErr = errors.Join(resultErr, writeErr)
			}
		}
	}()
	go func() {
		defer cancel()
		writerDone <- client.WriteMessage(ctx)
	}()
	for {
		message, err := client.ReadeMessage(ctx)
		if err != nil {
			return err
		}
		if err := s.HandleWebSocketMessage(ctx, client.ClientID(), message); err != nil {
			return err
		}
	}
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
