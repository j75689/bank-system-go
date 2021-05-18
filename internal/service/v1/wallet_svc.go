package v1

import (
	"bank-system-go/internal/config"
	"bank-system-go/internal/model"
	"bank-system-go/internal/repository"
	"bank-system-go/internal/service"
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var _ service.WalletService = (*WalletService)(nil)

func NewWalletService(config config.Config, walletRepo repository.WalletRepository) service.WalletService {
	return &WalletService{
		config:     config,
		walletRepo: walletRepo,
	}
}

type WalletService struct {
	config     config.Config
	walletRepo repository.WalletRepository
}

func (svc *WalletService) CreateWallet(ctx context.Context, wallet model.Wallet) (model.Wallet, error) {
	wallet.AccountNumber = uuid.New().String()
	return wallet, svc.walletRepo.CreateWallet(ctx, &wallet)
}

func (svc *WalletService) GetWallet(ctx context.Context, filter model.Wallet) (model.Wallet, error) {
	return svc.walletRepo.GetWallet(ctx, filter)
}

func (svc *WalletService) ListWallet(ctx context.Context, filter model.Wallet,
	pagination model.Pagination, sorting model.Sorting) ([]model.Wallet, int64, error) {
	return svc.walletRepo.ListWallet(ctx, filter, pagination, sorting)
}

func (svc *WalletService) UpdateBalance(ctx context.Context, filter model.Wallet, requestID string, transactionType model.TransactionType, amount decimal.Decimal) error {
	return svc.walletRepo.UpdateBalance(ctx, filter, requestID, transactionType, amount)
}
