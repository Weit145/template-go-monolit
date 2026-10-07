package domain_error

import "errors"

var (
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("forbidden")
	ErrInvalidInput  = errors.New("invalid input")
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrRateLimit     = errors.New("rate limit exceeded")
	ErrUnavailable   = errors.New("service unavailable")
	ErrTimeout       = errors.New("upstream timeout")
	ErrInternal      = errors.New("internal error")
)
