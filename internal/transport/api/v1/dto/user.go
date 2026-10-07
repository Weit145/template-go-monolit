package dto

import (
	"context"

	domain_user "github.com/Weit145/template-go-monolit/internal/domain/user"

	utils "github.com/Weit145/template-go-monolit/internal/transport/api/utils/err"
	"github.com/google/uuid"
)

type UserRequest struct {
	Name string `json:"name" validate:"required,min=2,max=32"`
}

func (r *UserRequest) Validate(ctx context.Context) error {
	if err := utils.Struct(ctx, r); err != nil {
		return err
	}
	return nil
}

type UserResponse struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func NewUserResponse(req domain_user.User) *UserResponse {
	return &UserResponse{
		Id:   req.ID(),
		Name: req.Name(),
	}
}
