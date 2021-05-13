package repository

import (
	"bank-system-go/internal/model"
	"context"
)

type TransationRepository interface {
	GetTransation(ctx context.Context, filter model.Transation) (model.Transation, error)
	ListTransation(ctx context.Context, filter model.Transation,
		pagination model.Pagination, sorting model.Sorting) ([]model.Transation, int64, error)
	CreateTransation(ctx context.Context, value model.Transation) error
	DeleteTransation(ctx context.Context, filter model.Transation) error
}
