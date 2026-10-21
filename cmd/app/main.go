package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Weit145/template-go-monolit/internal/app"
)

func main() {
	if err := run(); err != nil {
		slog.Error("Ошибка приложения", "err", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	resetSignals := context.AfterFunc(ctx, stop)
	defer resetSignals()

	a, err := app.New(ctx)
	if err != nil {
		return fmt.Errorf("Инициализация приложения: %w", err)
	}
	return a.Run(ctx)
}
