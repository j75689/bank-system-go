package gorm_repository

import (
	"bank-system-go/internal/model"
	"bank-system-go/internal/repository"
	"context"
	"errors"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

var (
	_ repository.UserRepository       = (*GORMRepository)(nil)
	_ repository.WalletRepository     = (*GORMRepository)(nil)
	_ repository.TransationRepository = (*GORMRepository)(nil)
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

func (repo *GORMRepository) GetWallet(ctx context.Context, filter model.Wallet) (model.Wallet, error) {
	wallet := model.Wallet{}
	return wallet, repo.db.WithContext(ctx).Where(filter).First(&wallet).Error
}

func (repo *GORMRepository) CreateWallet(ctx context.Context, value model.Wallet) error {
	return repo.db.WithContext(ctx).Create(&value).Error
}

func (repo *GORMRepository) UpdateWallet(ctx context.Context, filter model.Wallet, value model.Wallet) error {
	return repo.db.WithContext(ctx).Where(filter).Updates(value).Error
}

func (repo *GORMRepository) UpdateBalance(ctx context.Context, filter model.Wallet, amount decimal.Decimal) error {
	operator := "+"
	isLock := false
	if amount.LessThan(decimal.Zero) {
		isLock = true
		operator = "-"
	}
	db := repo.db.WithContext(ctx).
		Where(filter)

	if isLock {
		db = db.Where("balance >= ?", amount)
	}

	db = db.UpdateColumn("balance", gorm.Expr("balance "+operator+" ?", amount))
	if err := db.Error; err != nil {
		return err
	}
	if db.RowsAffected <= 0 {
		return errors.New("no record be affect")
	}
	return nil
}

func (repo *GORMRepository) DeleteWallet(ctx context.Context, filter model.Wallet) error {
	return repo.db.WithContext(ctx).Delete(filter).Error
}

func (repo *GORMRepository) GetTransation(ctx context.Context, filter model.Transation) (model.Transation, error) {
	transation := model.Transation{}
	return transation, repo.db.WithContext(ctx).Where(filter).First(&transation).Error
}

func (repo *GORMRepository) ListTransation(
	ctx context.Context,
	filter model.Transation,
	pagination model.Pagination, sorting model.Sorting) ([]model.Transation, int64, error) {
	var (
		transations = []model.Transation{}
		total       int64
	)
	db := repo.db.WithContext(ctx).Where(filter)
	err := db.Count(&total).Error
	if err != nil {
		return transations, total, err
	}

	return transations, total,
		db.Scopes(pagination.LimitAndOffset, sorting.Sort).Find(&transations).Error
}

func (repo *GORMRepository) CreateTransation(ctx context.Context, value model.Transation) error {
	return repo.db.WithContext(ctx).Create(&value).Error
}

func (repo *GORMRepository) DeleteTransation(ctx context.Context, filter model.Transation) error {
	return repo.db.WithContext(ctx).Delete(filter).Error
}
