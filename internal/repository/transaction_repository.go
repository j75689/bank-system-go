package repository

import (
	"bank-system-go/internal/model"
	"context"
)

type TransactionRepository interface {
	GetTransaction(ctx context.Context, filter model.Transaction) (model.Transaction, error)
	ListTransaction(ctx context.Context, filter model.Transaction,
		baseFilter model.BaseFilter, pagination model.Pagination, sorting model.Sorting) ([]model.Transaction, int64, error)
	CreateTransaction(ctx context.Context, value *model.Transaction) error
	DeleteTransaction(ctx context.Context, filter model.Transaction) error
}
