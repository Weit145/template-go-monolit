package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	domain_user "github.com/Weit145/template-go-monolit/internal/domain/user"
	"github.com/Weit145/template-go-monolit/internal/logger"
	utils "github.com/Weit145/template-go-monolit/internal/transport/api/utils/err"
	"github.com/Weit145/template-go-monolit/internal/transport/api/utils/request"
	"github.com/Weit145/template-go-monolit/internal/transport/api/v1/dto"
)

type UserService interface {
	CreateUser(ctx context.Context, name string) (domain_user.User, error)
}

func CreateUser(ser UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		var req dto.UserRequest
		if err := request.DecodeJSON(w, r, &req); err != nil {
			slog.DebugContext(ctx, "failed to decode request body", slog.Any("error", err))
			utils.DecodeFailed(w, err)
			return
		}

		if err := req.Validate(ctx); err != nil {
			slog.DebugContext(ctx, "failed to validate request body", slog.Any("error", err))
			utils.MapErr(w, r, err)
			return
		}

		user, err := ser.CreateUser(ctx, req.Name)
		if err != nil {
			logCtx := logger.ErrorCtx(ctx, err)
			slog.WarnContext(logCtx, "error service", slog.Any("error", err))
			utils.MapErr(w, r, err)
			return
		}

		response := dto.NewUserResponse(user)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(response)
	}
}
