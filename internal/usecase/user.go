package usecase

import (
	"context"
	"errors"
	"fmt"

	domain_error "github.com/Weit145/template-go-monolit/internal/domain/error"
	domain_user "github.com/Weit145/template-go-monolit/internal/domain/user"
	"github.com/Weit145/template-go-monolit/internal/logger"
)

type UserService interface {
	CreateAndDeleteUser(
		ctx context.Context,
		name string,
	) (domain_user.User, error)
}

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
			return logger.WrapError(ctx, fmt.Errorf("create user repo: %w", errors.Join(err, domain_error.ErrDataBase)))
		}

		err = s.repo.DeleteUser(ctx, bdUser.ID())
		if err != nil {
			return logger.WrapError(ctx, fmt.Errorf("delete user repo: %w", errors.Join(err, domain_error.ErrDataBase)))
		}
		return nil
	})
	if err != nil {
		return domain_user.User{}, logger.WrapError(ctx, fmt.Errorf("data base transaction: %w", errors.Join(err, domain_error.ErrDataBase)))
	}
	return bdUser, nil
}
