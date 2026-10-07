package domain_validator

import (
	"net/mail"
	"strings"

	domain_error "github.com/Weit145/template-go-monolit/internal/domain/error"
)

func NormalizeValue(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", domain_error.ErrInvalidInput
	}
	return value, nil
}

func NormalizeEmail(email string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return "", domain_error.ErrInvalidInput
	}
	address, err := mail.ParseAddress(email)
	if err != nil || address.Address != email {
		return "", domain_error.ErrInvalidInput
	}
	return email, nil
}

func NormalizePassword(password string) (string, error) {
	if strings.TrimSpace(password) == "" || len(password) < 6 || len(password) > 72 {
		return "", domain_error.ErrInvalidInput
	}
	return password, nil
}

func ValidateOpaqueValue(value string) (string, error) {
	if strings.TrimSpace(value) == "" {
		return "", domain_error.ErrInvalidInput
	}
	return value, nil
}
