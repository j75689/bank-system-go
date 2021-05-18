package repository

import (
	"bank-system-go/internal/model"
	"context"

	"github.com/shopspring/decimal"
)

type WalletRepository interface {
	GetWallet(ctx context.Context, filter model.Wallet) (model.Wallet, error)
	ListWallet(ctx context.Context, filter model.Wallet,
		pagination model.Pagination, sorting model.Sorting) ([]model.Wallet, int64, error)
	CreateWallet(ctx context.Context, value *model.Wallet) error
	UpdateWallet(ctx context.Context, filter model.Wallet, value *model.Wallet) error
	UpdateBalance(ctx context.Context, filter model.Wallet, request_id string, transactionType model.TransactionType, amount, fee decimal.Decimal) error
	DeleteWallet(ctx context.Context, filter model.Wallet) error
	Revert(ctx context.Context, requestID, accountNumber string, transactionType model.TransactionType) error
	Transfer(ctx context.Context, requestID string, from, to model.Wallet, transactionType model.TransactionType, amount, fee decimal.Decimal) error
}
