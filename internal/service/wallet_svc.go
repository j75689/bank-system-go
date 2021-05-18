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
	UpdateBalance(ctx context.Context, filter model.Wallet, requestID string, transationType model.TransationType, amount decimal.Decimal) error
}
