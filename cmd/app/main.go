package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Weit145/template-go-monolit/internal/config"
	"github.com/Weit145/template-go-monolit/internal/logger"
	"github.com/Weit145/template-go-monolit/internal/repositories/storage/postgres"
	"github.com/Weit145/template-go-monolit/internal/transport/api"
	"github.com/Weit145/template-go-monolit/internal/usecase"
)

func main() {
	if err := run(); err != nil {
		slog.ErrorContext(logger.ErrorCtx(context.Background(), err), "service failed", slog.Any("error", err))
		os.Exit(1)
	}
}

func run() error {
	//Init config
	cfg, err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	//Init logger
	logger.InitLogging(cfg.Log.Level)
	slog.Info("Init logger")

	//Init context
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	//Init storage
	repo, err := postgres.NewRepository(ctx,
		postgres.Settings{
			Address:         cfg.DB.Address,
			MaxConns:        cfg.DB.MaxConns,
			MinConns:        cfg.DB.MinConns,
			MaxConnLifetime: cfg.DB.MaxConnLifetime,
			MaxConnIdleTime: cfg.DB.MaxConnIdleTime,
		})
	if err != nil {
		return logger.WrapError(ctx, fmt.Errorf("Connect storage: %w", err))
	}
	defer repo.Close()
	slog.Info("Init storage")

	//Init usecase
	ser := usecase.NewService(repo)

	//Init api
	settings := api.Settings{
		Addr:              cfg.HTTP.Address,
		TrustedProxies:    cfg.HTTP.TrustedProxies,
		ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
		ReadTimeout:       cfg.HTTP.ReadTimeout,
		WriteTimeout:      cfg.HTTP.WriteTimeout,
		IdleTimeout:       cfg.HTTP.IdleTimeout,
		ShutdownTimeout:   cfg.HTTP.ShutdownTimeout,
	}
	if err := api.StartServer(ctx, ser, settings); err != nil {
		return logger.WrapError(ctx, err)
	}
	return nil

}
