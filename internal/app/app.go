package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Weit145/template-go-monolit/internal/config"
	"github.com/Weit145/template-go-monolit/internal/logger"
	"github.com/Weit145/template-go-monolit/internal/transport/api"
)

type App struct {
	diContainer *diContainer
	httpServer  *http.Server
}

func New(ctx context.Context) (*App, error) {
	a := &App{}
	if err := a.initDeps(ctx); err != nil {
		return nil, errors.Join(err, a.shutdown())
	}
	return a, nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func() error{
		a.initConfig,
		a.initLogger,
		func() error { return a.initHTTPServer(ctx) },
	}
	for _, init := range inits {
		if err := init(); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initConfig() error {
	cfg, err := config.InitConfig()
	if err != nil {
		return fmt.Errorf("init config: %w", err)
	}
	a.diContainer = newDIContainer(cfg)
	return nil
}

func (a *App) initLogger() error {
	logger.InitLogging(a.diContainer.cfg.Log.Level)
	slog.Info("Init logger")
	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	handler, err := a.diContainer.Handler(ctx)
	if err != nil {
		return err
	}
	cfg := a.diContainer.cfg.HTTP
	a.httpServer = api.NewServer(handler.Routes(), api.Settings{
		Addr:              cfg.Address,
		TrustedProxies:    cfg.TrustedProxies,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		ShutdownTimeout:   cfg.ShutdownTimeout,
	})
	return nil
}

func (a *App) Run(ctx context.Context) error {
	serveErr := make(chan error, 1)
	go func() {
		err := a.httpServer.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		serveErr <- err
	}()

	var err error
	select {
	case err = <-serveErr:
	case <-ctx.Done():
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), a.diContainer.cfg.HTTP.ShutdownTimeout)
		defer shutdownCancel()
		err = a.httpServer.Shutdown(shutdownCtx)
		if err != nil {
			err = errors.Join(err, a.httpServer.Close())
		}
		err = errors.Join(err, <-serveErr)
	}
	return errors.Join(err, a.shutdown())
}

func (a *App) shutdown() error {
	if a.diContainer == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), a.diContainer.cfg.HTTP.ShutdownTimeout)
	defer cancel()
	return a.diContainer.closer.CloseAll(ctx)
}
