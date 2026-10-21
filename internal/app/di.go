package app

import (
	"context"
	"fmt"

	"github.com/Weit145/template-go-monolit/internal/closer"
	"github.com/Weit145/template-go-monolit/internal/config"
	"github.com/Weit145/template-go-monolit/internal/infrastructure"
	"github.com/Weit145/template-go-monolit/internal/repositories/storage/postgres"
	handlers "github.com/Weit145/template-go-monolit/internal/transport/api/v1/handler"
	"github.com/Weit145/template-go-monolit/internal/usecase"
)

type diContainer struct {
	cfg    config.Config
	closer closer.Closer

	db        *postgres.Repository
	trManager postgres.TrManagerInterfase
	inf       infrastructure.InfrastructureInterfase
	service   usecase.ServiceInterfase
	handler   handlers.Handler
}

func newDIContainer(cfg config.Config) *diContainer {
	return &diContainer{cfg: cfg}
}

func (d *diContainer) DB(ctx context.Context) (*postgres.Repository, error) {
	if d.db == nil {
		repo, err := postgres.NewRepository(ctx, postgres.Settings{
			Address:         d.cfg.DB.Address,
			MaxConns:        d.cfg.DB.MaxConns,
			MinConns:        d.cfg.DB.MinConns,
			MaxConnLifetime: d.cfg.DB.MaxConnLifetime,
			MaxConnIdleTime: d.cfg.DB.MaxConnIdleTime,
		})
		if err != nil {
			return nil, fmt.Errorf("init postgres: %w", err)
		}
		d.closer.Add("postgres", func(context.Context) error {
			repo.Close()
			return nil
		})
		d.db = repo
	}
	return d.db, nil
}

func (d *diContainer) TrManager(ctx context.Context) (postgres.TrManagerInterfase, error) {
	if d.trManager == nil {
		db, err := d.DB(ctx)
		if err != nil {
			return nil, err
		}
		trManager, err := postgres.NewTrManager(db)
		if err != nil {
			return nil, err
		}
		d.trManager = trManager
	}
	return d.trManager, nil
}

func (d *diContainer) Infrastructure() infrastructure.InfrastructureInterfase {
	if d.inf == nil {
		inf := infrastructure.NewInfrastructure()
		d.closer.Add("websocket hub", inf.Shutdown)
		d.inf = inf
	}
	return d.inf
}

func (d *diContainer) Service(ctx context.Context) (usecase.ServiceInterfase, error) {
	if d.service == nil {
		db, err := d.DB(ctx)
		if err != nil {
			return nil, err
		}
		trManager, err := d.TrManager(ctx)
		if err != nil {
			return nil, err
		}
		d.service = usecase.NewService(db, trManager, d.Infrastructure())
	}
	return d.service, nil
}

func (d *diContainer) Handler(ctx context.Context) (handlers.Handler, error) {
	if d.handler == nil {
		service, err := d.Service(ctx)
		if err != nil {
			return nil, err
		}
		d.handler = handlers.NewHandler(service)
	}
	return d.handler, nil
}
