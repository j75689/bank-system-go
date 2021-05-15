package service

import (
	"bank-system-go/internal/model"
	"context"
)

type UserService interface {
	Register(ctx context.Context, user model.User) (model.User, error)
}
