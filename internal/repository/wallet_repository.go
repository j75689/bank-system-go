package repository

import (
	"bank-system-go/internal/model"
	"context"

	"github.com/shopspring/decimal"
)

type WalletRepository interface {
	GetWallet(ctx context.Context, filter model.Wallet) (model.Wallet, error)
	CreateWallet(ctx context.Context, value model.Wallet) error
	UpdateWallet(ctx context.Context, filter model.Wallet, value model.Wallet) error
	UpdateBalance(ctx context.Context, filter model.Wallet, amount decimal.Decimal) error
	DeleteWallet(ctx context.Context, filter model.Wallet) error
}
