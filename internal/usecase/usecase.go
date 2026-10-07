package usecase

import (
	"context"

	domain_health "github.com/Weit145/template-go-monolit/internal/domain/health"
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
) (result domain_health.Health, err error) {
	err = s.repo.CheckHealth(ctx)
	if err != nil {
		return domain_health.NewHealth(false), err
	}
	return domain_health.NewHealth(true), nil

}
