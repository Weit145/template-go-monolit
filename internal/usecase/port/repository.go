package port

import (
	"context"

	domain_user "github.com/Weit145/template-go-monolit/internal/domain/user"
)

type Repository interface {
	CheckHealth(ctx context.Context) error
	CreateUser(ctx context.Context, user domain_user.User) (domain_user.User, error)
}
