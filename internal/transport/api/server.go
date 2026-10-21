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

func NewServer(handler http.Handler, settings Settings) *http.Server {
	return &http.Server{
		Addr:              settings.Addr,
		Handler:           handler,
		ReadHeaderTimeout: settings.ReadHeaderTimeout,
		ReadTimeout:       settings.ReadTimeout,
		WriteTimeout:      settings.WriteTimeout,
		IdleTimeout:       settings.IdleTimeout,
	}
}

func StartServer(ctx context.Context, server *http.Server, shutdownTimeout time.Duration) error {
	slog.Info("http server started", "addr", server.Addr)
	serveErr := serveHTTP(ctx, server, shutdownTimeout)
	return logger.WrapError(ctx, serveErr)
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
