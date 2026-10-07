package usecase

import (
	"context"
	"fmt"

	domain_user "github.com/Weit145/template-go-monolit/internal/domain/user"
	"github.com/Weit145/template-go-monolit/internal/logger"
)

func (s *Service) CreateUser(
	ctx context.Context,
	name string,
) (domain_user.User, error) {
	user, err := domain_user.NewUser(name)
	if err != nil {
		return domain_user.User{}, logger.WrapError(ctx, fmt.Errorf("create user domain: %w", err))
	}
	bdUser, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return domain_user.User{}, logger.WrapError(ctx, fmt.Errorf("create user repo: %w", err))
	}
	return bdUser, nil
}
