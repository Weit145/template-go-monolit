package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/Weit145/template-go-monolit/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Settings struct {
	Address         string
	MaxConns        *int32
	MinConns        *int32
	MaxConnLifetime *time.Duration
	MaxConnIdleTime *time.Duration
}

type Repository struct {
	db *pgxpool.Pool
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

	var pool *pgxpool.Pool

	for i := 1; i <= 10; i++ {
		pool, err = pgxpool.NewWithConfig(ctx, cfg)
		if err == nil {
			break
		}
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "Storage no connect", "try", i, slog.Any("error", err))
		select {
		case <-ctx.Done():
			return nil, logger.WrapError(ctx, ctx.Err())
		case <-time.After(time.Second * 2):
		}
	}
	if err != nil {
		return nil, logger.WrapError(ctx, fmt.Errorf("connect to database: %w", err))
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, logger.WrapError(ctx, fmt.Errorf("ping database: %w", err))
	}

	return &Repository{db: pool}, nil
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
