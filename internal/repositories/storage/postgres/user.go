package postgres

import (
	"context"
	"fmt"

	domain_user "github.com/Weit145/template-go-monolit/internal/domain/user"
	"github.com/Weit145/template-go-monolit/internal/logger"
	"github.com/Weit145/template-go-monolit/internal/repositories/storage/postgres/model"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) CreateUser(
	ctx context.Context,
	user domain_user.User,
) (domain_user.User, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return domain_user.User{}, logger.WrapError(ctx, fmt.Errorf("begin update email verification token tx: %w", err))
	}
	defer tx.Rollback(ctx)

	row, err := tx.Query(ctx,
		`INSERT INTO users (id, name)
		VALUES ($1, $2)
		RETURNING id, name`,
		user.ID(), user.Name(),
	)
	if err != nil {
		return domain_user.User{}, logger.WrapError(ctx, fmt.Errorf("insert user: %w", err))
	}
	UserBD, err := pgx.CollectOneRow(row, pgx.RowToStructByName[model.User])
	if err != nil {
		return domain_user.User{}, logger.WrapError(ctx, fmt.Errorf("parse user: %w", err))
	}

	result, err := UserBD.ToDomain(ctx)
	if err != nil {
		return domain_user.User{}, logger.WrapError(ctx, fmt.Errorf("to domain: %w", err))
	}

	if err := tx.Commit(ctx); err != nil {
		return domain_user.User{}, logger.WrapError(ctx, fmt.Errorf("commit user: %w", err))
	}

	return result, nil
}
