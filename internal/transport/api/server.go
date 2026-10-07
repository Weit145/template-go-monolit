package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/netip"
	"time"

	"github.com/Weit145/template-go-monolit/internal/logger"
	v1 "github.com/Weit145/template-go-monolit/internal/transport/api/v1"
)

type Settings struct {
	Addr              string
	TrustedProxies    []netip.Prefix
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
}

type httpServer interface {
	ListenAndServe() error
	Shutdown(context.Context) error
}

func StartServer(ctx context.Context, ser v1.UseCase, settings Settings) error {
	router := v1.NewRouter(ser, settings.TrustedProxies)
	slog.Info("Init api")

	server := &http.Server{
		Addr:              settings.Addr,
		Handler:           router,
		ReadHeaderTimeout: settings.ReadHeaderTimeout,
		ReadTimeout:       settings.ReadTimeout,
		WriteTimeout:      settings.WriteTimeout,
		IdleTimeout:       settings.IdleTimeout,
	}
	slog.Info("http server started", "addr", settings.Addr)
	if err := serveHTTP(ctx, server, settings.ShutdownTimeout); err != nil {
		return logger.WrapError(ctx, err)
	}
	return nil
}

func serveHTTP(ctx context.Context, server httpServer, shutdownTimeout time.Duration) error {
	shutdownSignal, stopShutdown := context.WithCancel(ctx)
	defer stopShutdown()

	shutdownDone := make(chan error, 1)
	go func() {
		<-shutdownSignal.Done()
		if ctx.Err() == nil {
			shutdownDone <- nil
			return
		}

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		shutdownDone <- server.Shutdown(shutdownCtx)
	}()

	serveErr := server.ListenAndServe()
	if ctx.Err() != nil {
		if err := <-shutdownDone; err != nil {
			return fmt.Errorf("shutdown http server: %w", err)
		}
	}
	if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
		return fmt.Errorf("listen and serve: %w", serveErr)
	}
	return nil
}
