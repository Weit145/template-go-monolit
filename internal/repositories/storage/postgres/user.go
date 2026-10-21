package postgres

import (
	"context"
	"fmt"

	domain_user "github.com/Weit145/template-go-monolit/internal/domain/user"
	"github.com/Weit145/template-go-monolit/internal/logger"
	"github.com/Weit145/template-go-monolit/internal/repositories/storage/postgres/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserRepo interface {
	CreateUser(
		ctx context.Context,
		user domain_user.User,
	) (domain_user.User, error)

	DeleteUser(
		ctx context.Context,
		id uuid.UUID,
	) error
}

func (r *Repository) CreateUser(
	ctx context.Context,
	user domain_user.User,
) (domain_user.User, error) {
	row, err := r.ctxGetter.DefaultTrOrDB(ctx, r.db).Query(ctx,
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
	return result, nil
}

func (r *Repository) DeleteUser(
	ctx context.Context,
	id uuid.UUID,
) error {
	_, err := r.ctxGetter.DefaultTrOrDB(ctx, r.db).Exec(ctx,
		`DELETE FROM users
	 	WHERE users.id = $1`, id)
	if err != nil {
		return logger.WrapError(ctx, fmt.Errorf("delete user: %w", err))
	}
	return nil
}
