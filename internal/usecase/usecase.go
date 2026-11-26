package usecase

import (
	"context"
	"errors"

	domain_error "github.com/Weit145/template-go-monolit/internal/domain/error"
	domain_health "github.com/Weit145/template-go-monolit/internal/domain/health"
	"github.com/Weit145/template-go-monolit/internal/usecase/port"
)

type ServiceInterface interface {
	CheckHealthService(
		ctx context.Context,
	) (result domain_health.Health, err error)
	UserService
	WebSocketService
}

type Service struct {
	repo      port.Repository
	trManager port.TrManager
	events    port.EventHub
}

func NewService(repo port.Repository, trManager port.TrManager, events port.EventHub) ServiceInterface {
	return &Service{repo: repo, trManager: trManager, events: events}
}

func (s *Service) CheckHealthService(
	ctx context.Context,
) (result domain_health.Health, err error) {
	err = s.repo.CheckHealth(ctx)
	if err != nil {
		return domain_health.NewHealth(false), errors.Join(err, domain_error.ErrUnavailable)
	}
	return domain_health.NewHealth(true), nil

}
