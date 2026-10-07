package port

import "context"

type Repository interface {
	CheckHealth(ctx context.Context) error
}
