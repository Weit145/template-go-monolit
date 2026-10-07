package domain_user

import (
	"fmt"

	domain_validator "github.com/Weit145/template-go-monolit/internal/domain/validator"
	"github.com/google/uuid"
)

type User struct {
	id   uuid.UUID
	name string
}

func (u User) ID() uuid.UUID { return u.id }
func (u User) Name() string  { return u.name }

func NewUser(name string) (User, error) {
	name, err := domain_validator.ValidateOpaqueValue(name)
	if err != nil {
		return User{}, fmt.Errorf("invalid name: %w", err)
	}
	return User{
		id:   uuid.New(),
		name: name,
	}, nil
}

func UserFromDB(id uuid.UUID, name string) (User, error) {
	if id == uuid.Nil {
		return User{}, fmt.Errorf("invalid id")
	}
	name, err := domain_validator.ValidateOpaqueValue(name)
	if err != nil {
		return User{}, fmt.Errorf("invalid name: %w", err)
	}
	return User{
		id:   id,
		name: name,
	}, nil
}
