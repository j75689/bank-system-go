package v1

import (
	"bank-system-go/internal/model"
	"bank-system-go/internal/repository"
	"bank-system-go/internal/service"
	"context"
)

var _ service.TransactionService = (*TransactionService)(nil)

func NewTransactionService(transactionRepo repository.TransactionRepository) service.TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
	}
}

type TransactionService struct {
	transactionRepo repository.TransactionRepository
}

func (svc *TransactionService) CreateTransaction(ctx context.Context, transaction model.Transaction) (model.Transaction, error) {
	return transaction, svc.transactionRepo.CreateTransaction(ctx, &transaction)
}

func (svc *TransactionService) ListTransaction(ctx context.Context, filter model.Transaction, baseFilter model.BaseFilter,
	pagination model.Pagination, sorting model.Sorting) ([]model.Transaction, int64, error) {
	return svc.transactionRepo.ListTransaction(ctx, filter, baseFilter, pagination, sorting)
}
