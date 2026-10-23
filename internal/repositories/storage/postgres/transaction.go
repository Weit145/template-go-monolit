package postgres

import (
	"context"
	"fmt"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	trmmanager "github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

type TrManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) (err error)
}

func NewTrManager(repo *Repository) (TrManager, error) {
	trManager, err := trmmanager.New(trmpgx.NewDefaultFactory(repo.db))
	if err != nil {
		return nil, fmt.Errorf("create trmanager: %w", err)
	}
	return trManager, nil
}
