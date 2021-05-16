package service

import (
	"bank-system-go/internal/model"
	"context"
)

type WalletService interface {
	CreateWallet(ctx context.Context, wallet model.Wallet) (model.Wallet, error)
}
