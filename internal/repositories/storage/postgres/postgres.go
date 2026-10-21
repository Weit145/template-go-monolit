package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/Weit145/template-go-monolit/internal/logger"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	trmmanager "github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Settings struct {
	Address         string
	MaxConns        *int32
	MinConns        *int32
	MaxConnLifetime *time.Duration
	MaxConnIdleTime *time.Duration
}

type RepositoryInterfase interface {
	Close()
	CheckHealth(ctx context.Context) error
	UserRepo
}

type TrManagerInterfase interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) (err error)
}

type Repository struct {
	db        *pgxpool.Pool
	ctxGetter *trmpgx.CtxGetter
}

func NewRepository(ctx context.Context, settigs Settings) (*Repository, error) {
	if strings.TrimSpace(settigs.Address) == "" {
		err := fmt.Errorf("database dsn is empty")
		slog.ErrorContext(ctx, "database dsn is empty", slog.Any("error", err))
		return nil, logger.WrapError(ctx, err)
	}

	cfg, err := pgxpool.ParseConfig(settigs.Address)
	if err != nil {
		return nil, logger.WrapError(ctx, fmt.Errorf("parse pgxpool config: %w", err))
	}

	if settigs.MaxConns != nil {
		cfg.MaxConns = *settigs.MaxConns
	}
	if settigs.MinConns != nil {
		cfg.MinConns = *settigs.MinConns
	}
	if settigs.MaxConnLifetime != nil {
		cfg.MaxConnLifetime = *settigs.MaxConnLifetime
	}
	if settigs.MaxConnIdleTime != nil {
		cfg.MaxConnIdleTime = *settigs.MaxConnIdleTime
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, logger.WrapError(ctx, fmt.Errorf("create database pool: %w", err))
	}
	if err := pingWithRetry(ctx, pool.Ping, 10, 2*time.Second, 2*time.Second); err != nil {
		pool.Close()
		return nil, logger.WrapError(ctx, err)
	}

	return &Repository{db: pool, ctxGetter: trmpgx.DefaultCtxGetter}, nil
}

func pingWithRetry(ctx context.Context, ping func(context.Context) error, attempts int, timeout, pause time.Duration) error {
	var err error
	for i := 1; i <= attempts; i++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		pingCtx, cancel := context.WithTimeout(ctx, timeout)
		err = ping(pingCtx)
		cancel()
		if err == nil {
			return nil
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		slog.WarnContext(ctx, "database ping failed", "try", i, slog.Any("error", err))
		if i == attempts {
			break
		}
		timer := time.NewTimer(pause)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
	return fmt.Errorf("ping database after %d attempts: %w", attempts, err)
}

func NewTrManager(repo *Repository) (TrManagerInterfase, error) {
	trManager, err := trmmanager.New(trmpgx.NewDefaultFactory(repo.db))
	if err != nil {
		return nil, fmt.Errorf("create trmanager: %w", err)
	}
	return trManager, nil
}

func (r *Repository) Close() {
	if r.db != nil {
		r.db.Close()
	}
}

func (r *Repository) CheckHealth(ctx context.Context) error {
	if err := r.db.Ping(ctx); err != nil {
		return logger.WrapError(ctx, fmt.Errorf("ping database: %w", err))
	}
	return nil
}
