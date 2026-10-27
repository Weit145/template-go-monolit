package websocket

import (
	"context"
	"fmt"
	"time"

	"github.com/Weit145/template-go-monolit/internal/domain/event"
	"github.com/Weit145/template-go-monolit/internal/logger"
	"github.com/Weit145/template-go-monolit/internal/transport/api/v1/dto"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/google/uuid"
)

type Client struct {
	ID   uuid.UUID
	Conn *websocket.Conn
	Send chan event.Event
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		ID:   uuid.New(),
		Conn: conn,
		Send: make(chan event.Event, 32),
	}
}

func (c *Client) ClientID() uuid.UUID { return c.ID }

func (c *Client) Close() error { return c.Conn.CloseNow() }

func (c *Client) SendMessage(message event.Event) bool {
	select {
	case c.Send <- message:
		return true
	default:
		return false
	}
}

func (c *Client) WriteMessage(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-c.Send:
			if !ok {
				return nil
			}
			writeCtx, cancel := context.WithTimeout(
				ctx,
				5*time.Second,
			)
			err := wsjson.Write(writeCtx, c.Conn, dto.NewServerMessage(msg))
			cancel()
			if err != nil {
				return logger.WrapError(ctx, fmt.Errorf("write websocket event: %w", err))
			}
		}
	}
}

func (c *Client) ReadMessage(ctx context.Context) (event.Command, error) {
	var msg dto.ClientMessage
	err := wsjson.Read(ctx, c.Conn, &msg)
	if err != nil {
		return event.Command{}, logger.WrapError(ctx, fmt.Errorf("read json: %w", err))
	}
	return event.Command{ID: msg.EventID, Type: msg.Type}, nil
}
