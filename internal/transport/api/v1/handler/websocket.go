package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	client "github.com/Weit145/template-go-monolit/internal/infrastructure/websocket"
	"github.com/Weit145/template-go-monolit/internal/usecase/port"
	"github.com/coder/websocket"
)

type WebSocketService interface {
	WebSocket(context.Context, port.WebSocketClient) error
}

func WebSocket(ser WebSocketService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			slog.DebugContext(ctx, "failed to connect websocket", slog.Any("error", err))
			return
		}
		defer conn.CloseNow()
		c := client.NewClient(conn)
		logWebSocketError(ctx, "websocket stopped", ser.WebSocket(ctx, c))
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
