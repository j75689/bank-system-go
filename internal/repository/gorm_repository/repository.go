package gorm_repository

import (
	pkgErr "bank-system-go/internal/errors"
	"bank-system-go/internal/model"
	"bank-system-go/internal/repository"
	"context"
	"errors"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

var (
	_ repository.UserRepository        = (*GORMRepository)(nil)
	_ repository.WalletRepository      = (*GORMRepository)(nil)
	_ repository.TransactionRepository = (*GORMRepository)(nil)
)

func NewGORMRepository(db *gorm.DB) *GORMRepository {
	return &GORMRepository{
		db: db,
	}
}

type GORMRepository struct {
	db *gorm.DB
}

func (repo *GORMRepository) GetUser(ctx context.Context, filter model.User) (model.User, error) {
	user := model.User{}
	return user, repo.db.WithContext(ctx).Where(filter).First(&user).Error
}

func (repo *GORMRepository) CreateUser(ctx context.Context, value *model.User) error {
	return repo.db.WithContext(ctx).Create(value).Error
}

func (repo *GORMRepository) UpdateUser(ctx context.Context, filter model.User, value *model.User) error {
	return repo.db.WithContext(ctx).Where(filter).Updates(value).Error
}

func (repo *GORMRepository) DeleteUser(ctx context.Context, filter model.User) error {
	return repo.db.WithContext(ctx).Delete(filter).Error
}

func (repo *GORMRepository) CreateAccessLog(ctx context.Context, value *model.AccessLog) error {
	return repo.db.WithContext(ctx).Create(value).Error
}

func (repo *GORMRepository) ListAccessLog(ctx context.Context, filter model.AccessLog,
	pagination model.Pagination, sorting model.Sorting) ([]model.AccessLog, int64, error) {
	var (
		accesslogs = []model.AccessLog{}
		total      int64
	)
	db := repo.db.WithContext(ctx).Where(filter)
	err := db.Count(&total).Error
	if err != nil {
		return accesslogs, total, err
	}

	return accesslogs, total,
		db.Scopes(pagination.LimitAndOffset, sorting.Sort).Find(&accesslogs).Error
}

func (repo *GORMRepository) GetWallet(ctx context.Context, filter model.Wallet) (model.Wallet, error) {
	wallet := model.Wallet{}
	db := repo.db.WithContext(ctx).Where(filter).Scopes(wallet.Preload).First(&wallet)
	if db.Statement.RowsAffected <= 0 {
		return wallet, pkgErr.ErrWalletAccountNotFound
	}
	return wallet, nil
}

func (repo *GORMRepository) ListWallet(ctx context.Context, filter model.Wallet,
	pagination model.Pagination, sorting model.Sorting) ([]model.Wallet, int64, error) {
	var (
		wallets = []model.Wallet{}
		total   int64
	)
	db := repo.db.WithContext(ctx).Model(model.Wallet{}).Where(filter)
	err := db.Count(&total).Error
	if err != nil {
		return wallets, total, err
	}

	return wallets, total,
		db.Scopes(model.Wallet{}.Preload, pagination.LimitAndOffset, sorting.Sort).Find(&wallets).Error
}

func (repo *GORMRepository) CreateWallet(ctx context.Context, value *model.Wallet) error {
	return repo.db.WithContext(ctx).Create(value).Error
}

func (repo *GORMRepository) UpdateWallet(ctx context.Context, filter model.Wallet, value *model.Wallet) error {
	return repo.db.WithContext(ctx).Where(filter).Updates(value).Error
}

func (repo *GORMRepository) UpdateBalance(ctx context.Context, filter model.Wallet, request_id string, transactionType model.TransactionType, amount decimal.Decimal) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		history := model.WalletHistory{
			RequestID:       request_id,
			UserID:          filter.UserID,
			TransactionType: transactionType,
			AccountNumber:   filter.AccountNumber,
			Amount:          amount,
		}
		err := tx.Where(history).First(&history).Error
		if err == nil {
			// already existed
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		isLock := false
		if amount.LessThan(decimal.Zero) {
			isLock = true
		}
		db := tx.WithContext(ctx).
			Where(filter)

		if isLock {
			db = db.Where("cast(balance as decimal) >= ?", amount.Abs())
		}

		db = db.Model(model.Wallet{}).UpdateColumn("balance", gorm.Expr("cast(balance as decimal) + cast(? as decimal)", amount))
		if err := db.Error; err != nil {
			return err
		}
		if db.RowsAffected <= 0 {
			return pkgErr.ErrWalletBalanceInsufficient
		}
		return tx.Model(model.WalletHistory{}).Create(&history).Error
	})
}

func (repo *GORMRepository) DeleteWallet(ctx context.Context, filter model.Wallet) error {
	return repo.db.WithContext(ctx).Delete(filter).Error
}

func (repo *GORMRepository) GetTransaction(ctx context.Context, filter model.Transaction) (model.Transaction, error) {
	transaction := model.Transaction{}
	return transaction, repo.db.WithContext(ctx).Where(filter).First(&transaction).Error
}

func (repo *GORMRepository) ListTransaction(
	ctx context.Context,
	filter model.Transaction,
	baseFilter model.BaseFilter, pagination model.Pagination, sorting model.Sorting) ([]model.Transaction, int64, error) {
	var (
		transactions = []model.Transaction{}
		total        int64
	)
	db := repo.db.WithContext(ctx).Model(model.Transaction{}).Where(filter)
	err := db.Count(&total).Error
	if err != nil {
		return transactions, total, err
	}

	return transactions, total,
		db.Scopes(baseFilter.Filter, pagination.LimitAndOffset, sorting.Sort).Find(&transactions).Error
}

func (repo *GORMRepository) CreateTransaction(ctx context.Context, value *model.Transaction) error {
	return repo.db.WithContext(ctx).Create(value).Error
}

func (repo *GORMRepository) DeleteTransaction(ctx context.Context, filter model.Transaction) error {
	return repo.db.WithContext(ctx).Delete(filter).Error
}
