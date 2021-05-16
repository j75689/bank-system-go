package repository

import (
	"bank-system-go/internal/model"
	"context"
)

type UserRepository interface {
	GetUser(ctx context.Context, filter model.User) (model.User, error)
	CreateUser(ctx context.Context, value *model.User) error
	UpdateUser(ctx context.Context, filter model.User, value *model.User) error
	DeleteUser(ctx context.Context, filter model.User) error
	CreateAccessLog(ctx context.Context, value *model.AccessLog) error
	ListAccessLog(ctx context.Context, filter model.AccessLog,
		pagination model.Pagination, sorting model.Sorting) ([]model.AccessLog, int64, error)
}
