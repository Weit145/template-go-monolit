package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/Weit145/template-go-monolit/internal/domain/event"
	client "github.com/Weit145/template-go-monolit/internal/transport/api/v1/websocket"
	"github.com/coder/websocket"
)

type WebSocketService interface {
	RegistreWebSocket(context.Context, event.WebSocketClient) error
	UnRegistreWebSocket(context.Context, uuid.UUID) error
	HandleWebSocketMessage(
		ctx context.Context,
		id_conn uuid.UUID,
		message event.Command,
	) error
}

func (h *handler) WebSocket(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		slog.DebugContext(ctx, "failed to connect websocket", slog.Any("error", err))
		return
	}
	defer conn.CloseNow()
	c := client.NewClient(conn)
	err = h.ser.RegistreWebSocket(ctx, c)
	if err != nil {
		logWebSocketError(ctx, "websocket error", err)
		return
	}

	writerDone := make(chan error, 1)
	defer func() {
		cancel()
		_ = c.Close()
		writeErr := <-writerDone
		err = h.ser.UnRegistreWebSocket(ctx, c.ClientID())
		if writeErr != nil && !errors.Is(writeErr, context.Canceled) {
			if errors.Is(err, context.Canceled) {
				err = writeErr
			} else {
				err = errors.Join(err, writeErr)
			}
		}
	}()

	go func() {
		writerDone <- c.WriteMessage(ctx)
	}()
	for {
		message, err := c.ReadeMessage(ctx)
		if err != nil {
			logWebSocketError(ctx, "websocket stopped", err)
			return
		}
		if err := h.ser.HandleWebSocketMessage(ctx, c.ClientID(), message); err != nil {
			logWebSocketError(ctx, "websocket stopped", err)
			return
		}
	}
}

func logWebSocketError(ctx context.Context, message string, err error) {
	if err == nil || errors.Is(err, context.Canceled) {
		return
	}
	status := websocket.CloseStatus(err)
	if status == websocket.StatusNormalClosure || status == websocket.StatusGoingAway {
		return
	}
	slog.WarnContext(ctx, message, slog.Any("error", err))
}
