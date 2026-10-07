package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	domain_error "github.com/Weit145/template-go-monolit/internal/domain/error"
	"github.com/Weit145/template-go-monolit/internal/logger"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func MapErr(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, domain_error.ErrUnauthorized):
		writeJSON(w, http.StatusUnauthorized, Error("unauthorized"))
	case errors.Is(err, domain_error.ErrForbidden):
		writeJSON(w, http.StatusForbidden, Error("forbidden"))
	case errors.Is(err, domain_error.ErrInvalidInput):
		writeJSON(w, http.StatusBadRequest, Error("invalid request"))
	case errors.Is(err, domain_error.ErrNotFound):
		writeJSON(w, http.StatusNotFound, Error("not found"))
	case errors.Is(err, domain_error.ErrAlreadyExists):
		writeJSON(w, http.StatusConflict, Error("already exists"))
	case errors.Is(err, domain_error.ErrRateLimit):
		writeJSON(w, http.StatusTooManyRequests, Error("rate limit exceeded"))
	case errors.Is(err, domain_error.ErrUnavailable):
		writeJSON(w, http.StatusServiceUnavailable, Error("service unavailable"))
	case errors.Is(err, domain_error.ErrTimeout):
		writeJSON(w, http.StatusGatewayTimeout, Error("gateway timeout"))
	default:
		writeJSON(w, http.StatusInternalServerError, Error("internal server error"))
	}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func DecodeFailed(w http.ResponseWriter, err error) {
	var maxBytesErr *http.MaxBytesError
	if errors.As(err, &maxBytesErr) {
		writeJSON(w, http.StatusRequestEntityTooLarge, Error("request body too large"))
		return
	}

	writeJSON(w, http.StatusBadRequest, Error("invalid request"))
}

func CookieFaild(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(Response{Status: "error unauthorized"})
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func Success() *Response {
	return &Response{
		Status: "success",
	}
}

func Error(msg string) *Response {
	return &Response{
		Status: "error",
		Error:  msg,
	}
}

func Struct(ctx context.Context, value any) error {
	if err := validate.Struct(value); err != nil {
		return logger.WrapError(ctx, fmt.Errorf("validate request: %w: %w", err, domain_error.ErrInvalidInput))
	}

	return nil
}
