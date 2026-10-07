package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Weit145/template-go-monolit/internal/domain/health"
	"github.com/Weit145/template-go-monolit/internal/logger"
	utils "github.com/Weit145/template-go-monolit/internal/transport/api/utils/err"
)

type HealthService interface {
	CheckHealthService(ctx context.Context) (health health.Health, err error)
}

func Health(ser HealthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		health, err := ser.CheckHealthService(ctx)
		if err != nil || !health.IsPostgres() {
			logCtx := logger.ErrorCtx(ctx, err)
			slog.WarnContext(logCtx, "gateway is not ready", slog.Any("error", err))
			utils.MapErr(w, r, err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(utils.Success())
	}
}
