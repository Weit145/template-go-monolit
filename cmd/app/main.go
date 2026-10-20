package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Weit145/template-go-monolit/internal/config"
	"github.com/Weit145/template-go-monolit/internal/infrastructure"
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

	//Init transaction manager
	trManager, err := postgres.NewTrManager(ctx, repo)
	if err != nil {
		return logger.WrapError(ctx, fmt.Errorf("TrManager: %w", err))
	}
	//Init infastructure
	inf := infrastructure.NewInfrastructure()

	//Init usecase
	ser := usecase.NewService(repo, trManager, inf)

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
	serveErr := api.StartServer(ctx, ser, settings)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer cancel()
	if err := inf.Shutdown(shutdownCtx); err != nil {
		serveErr = errors.Join(serveErr, fmt.Errorf("shutdown websocket hub: %w", err))
	}
	return logger.WrapError(ctx, serveErr)

}
