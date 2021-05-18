package service

import (
	"bank-system-go/internal/model"
	"context"
)

type TransactionService interface {
	CreateTransaction(ctx context.Context, transaction model.Transaction) (model.Transaction, error)
	ListTransaction(ctx context.Context, filter model.Transaction, baseFilter model.BaseFilter,
		pagination model.Pagination, sorting model.Sorting) ([]model.Transaction, int64, error)
}
