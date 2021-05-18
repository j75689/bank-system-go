package service

import (
	"bank-system-go/internal/model"
	"context"

	"github.com/shopspring/decimal"
)

type WalletService interface {
	CreateWallet(ctx context.Context, wallet model.Wallet) (model.Wallet, error)
	GetWallet(ctx context.Context, filter model.Wallet) (model.Wallet, error)
	ListWallet(ctx context.Context, filter model.Wallet,
		pagination model.Pagination, sorting model.Sorting) ([]model.Wallet, int64, error)
	UpdateBalance(ctx context.Context, filter model.Wallet, requestID string, transactionType model.TransactionType, amount, fee decimal.Decimal) error
	Revert(ctx context.Context, requestID, accountNumber string, transactionType model.TransactionType) error
	Transfer(ctx context.Context, requestID string, from, to model.Wallet, transactionType model.TransactionType, amount, fee decimal.Decimal) error
}
