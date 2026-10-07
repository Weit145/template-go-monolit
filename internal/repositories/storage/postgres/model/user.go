package model

import (
	"context"
	"fmt"

	domain_user "github.com/Weit145/template-go-monolit/internal/domain/user"
	"github.com/Weit145/template-go-monolit/internal/logger"
	"github.com/google/uuid"
)

type User struct {
	Id   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

func (u User) ToDomain(ctx context.Context) (domain_user.User, error) {
	user, err := domain_user.UserFromDB(u.Id, u.Name)
	if err != nil {
		return domain_user.User{}, logger.WrapError(ctx, fmt.Errorf("from db: %w", err))
	}
	return user, nil
}
