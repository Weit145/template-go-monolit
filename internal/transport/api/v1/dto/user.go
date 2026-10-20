package dto

import (
	"context"

	domain_user "github.com/Weit145/template-go-monolit/internal/domain/user"

	utils "github.com/Weit145/template-go-monolit/internal/transport/api/utils/err"
	"github.com/google/uuid"
)

type CreateAndDeleteUserRequest struct {
	Name string `json:"name" validate:"required,min=2,max=32"`
}

func (r *CreateAndDeleteUserRequest) Validate(ctx context.Context) error {
	if err := utils.Struct(ctx, r); err != nil {
		return err
	}
	return nil
}

type CreateAndDeleteUserResponse struct {
	Id     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Status string    `json:"status"`
}

func NewCreateAndDeleteUserResponse(req domain_user.User) *CreateAndDeleteUserResponse {
	return &CreateAndDeleteUserResponse{
		Id:     req.ID(),
		Name:   req.Name(),
		Status: "created_and_deleted",
	}
}
