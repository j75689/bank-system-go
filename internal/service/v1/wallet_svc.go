package v1

import (
	"bank-system-go/internal/config"
	"bank-system-go/internal/model"
	"bank-system-go/internal/repository"
	"bank-system-go/internal/service"
	"context"

	"github.com/google/uuid"
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
