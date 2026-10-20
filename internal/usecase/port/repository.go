package port

import (
	"context"

	domain_user "github.com/Weit145/template-go-monolit/internal/domain/user"
	"github.com/google/uuid"
)

type TrManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) (err error)
}

type Repository interface {
	CheckHealth(ctx context.Context) error
	CreateUser(ctx context.Context, user domain_user.User) (domain_user.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
