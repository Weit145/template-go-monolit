package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	domain_error "github.com/Weit145/template-go-monolit/internal/domain/error"
	domain_health "github.com/Weit145/template-go-monolit/internal/domain/health"
	"github.com/Weit145/template-go-monolit/internal/logger"
	utils "github.com/Weit145/template-go-monolit/internal/transport/api/utils/err"
)

type HealthService interface {
	CheckHealthService(ctx context.Context) (health domain_health.Health, err error)
}

func Health(ser HealthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		health, err := ser.CheckHealthService(ctx)
		if err != nil || !health.IsPostgres() {
			logCtx := logger.ErrorCtx(ctx, err)
			slog.WarnContext(logCtx, "gateway is not ready", slog.Any("error", err))
			if err != nil {
				utils.MapErr(w, r, err)
			} else {
				utils.MapErr(w, r, domain_error.ErrUnavailable)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(utils.Success())
	}
}
