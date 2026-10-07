package usecase

import (
	"context"

	"github.com/Weit145/template-go-monolit/internal/domain/health"
	"github.com/Weit145/template-go-monolit/internal/usecase/port"
)

type Service struct {
	repo port.Repository
}

func NewService(repo port.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CheckHealthService(
	ctx context.Context,
) (result health.Health, err error) {
	err = s.repo.CheckHealth(ctx)
	if err != nil {
		return health.NewHealth(false), err
	}
	return health.NewHealth(true), nil

}
