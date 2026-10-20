package usecase

import (
	"context"
	"fmt"

	"github.com/Weit145/template-go-monolit/internal/domain/event"
	domain_user "github.com/Weit145/template-go-monolit/internal/domain/user"
	"github.com/Weit145/template-go-monolit/internal/logger"
)

func (s *Service) CreateAndDeleteUser(
	ctx context.Context,
	name string,
) (domain_user.User, error) {
	user, err := domain_user.NewUser(name)
	if err != nil {
		return domain_user.User{}, logger.WrapError(ctx, fmt.Errorf("create user domain: %w", err))
	}

	bdUser := domain_user.User{}

	err = s.trManager.Do(ctx, func(ctx context.Context) error {
		bdUser, err = s.repo.CreateUser(ctx, user)
		if err != nil {
			return logger.WrapError(ctx, fmt.Errorf("create user repo: %w", err))
		}

		err = s.repo.DeleteUser(ctx, bdUser.ID())
		if err != nil {
			return logger.WrapError(ctx, fmt.Errorf("delete user repo: %w", err))
		}
		return nil
	})
	if err != nil {
		return domain_user.User{}, logger.WrapError(ctx, fmt.Errorf("create and delete user transaction: %w", err))
	}
	s.events.Publish(event.Event{
		ID:   bdUser.ID(),
		Type: "user.created_and_deleted",
		Data: map[string]any{"id": bdUser.ID(), "name": bdUser.Name()},
	})
	return bdUser, nil
}
