package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/Weit145/template-go-monolit/internal/domain/event"
	domain_user "github.com/Weit145/template-go-monolit/internal/domain/user"
	"github.com/Weit145/template-go-monolit/internal/usecase/port"
	"github.com/google/uuid"
)

type eventTestRepository struct {
	createErr error
	deleteErr error
}

func (*eventTestRepository) CheckHealth(context.Context) error { return nil }
func (r *eventTestRepository) CreateUser(_ context.Context, user domain_user.User) (domain_user.User, error) {
	return user, r.createErr
}
func (r *eventTestRepository) DeleteUser(context.Context, uuid.UUID) error { return r.deleteErr }

type eventTestTransaction struct {
	commitErr error
	committed bool
}

func (m *eventTestTransaction) Do(ctx context.Context, fn func(context.Context) error) error {
	if err := fn(ctx); err != nil {
		return err
	}
	if m.commitErr != nil {
		return m.commitErr
	}
	m.committed = true
	return nil
}

type eventTestPublisher struct {
	t            *testing.T
	tx           *eventTestTransaction
	calls        int
	eventID      uuid.UUID
	commands     []event.Command
	connID       uuid.UUID
	unregistered bool
}

func (*eventTestPublisher) Registre(port.WebSocketClient) bool { return true }
func (p *eventTestPublisher) UnRegistre(uuid.UUID)             { p.unregistered = true }
func (p *eventTestPublisher) Subscribe(id_conn, eventID uuid.UUID) {
	p.connID = id_conn
	p.commands = append(p.commands, event.Command{ID: eventID, Type: "subscribe"})
}
func (p *eventTestPublisher) UnSubscribe(id_conn, eventID uuid.UUID) {
	p.connID = id_conn
	p.commands = append(p.commands, event.Command{ID: eventID, Type: "unsubscribe"})
}

func (p *eventTestPublisher) Publish(message event.Event) {
	if !p.tx.committed {
		p.t.Error("event published before transaction committed")
	}
	if message.Type != "user.created_and_deleted" {
		p.t.Errorf("unexpected event type: %s", message.Type)
	}
	p.calls++
	p.eventID = message.ID
}

func TestCreateUserPublishesOnlyAfterCommit(t *testing.T) {
	failure := errors.New("database failure")
	for _, stage := range []string{"success", "create", "delete", "commit"} {
		t.Run(stage, func(t *testing.T) {
			repo := &eventTestRepository{}
			tx := &eventTestTransaction{}
			switch stage {
			case "create":
				repo.createErr = failure
			case "delete":
				repo.deleteErr = failure
			case "commit":
				tx.commitErr = failure
			}
			publisher := &eventTestPublisher{t: t, tx: tx}
			service := NewService(repo, tx, publisher)
			user, err := service.CreateUser(context.Background(), "Alice")
			if stage != "success" {
				if !errors.Is(err, failure) || publisher.calls != 0 {
					t.Fatalf("failed transaction: err=%v, published=%d", err, publisher.calls)
				}
				return
			}
			if err != nil || publisher.calls != 1 || publisher.eventID != user.ID() || user.ID() == uuid.Nil {
				t.Fatalf("successful transaction: err=%v, published=%d, eventID=%s", err, publisher.calls, publisher.eventID)
			}
		})
	}
}
