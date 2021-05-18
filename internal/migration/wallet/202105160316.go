package wallet

import (
	"bank-system-go/internal/model"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

var v202105162211 = &gormigrate.Migration{
	ID: "202105162211",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.Wallet{}); err != nil {
			return err
		}
		if err := tx.AutoMigrate(&model.Currency{}); err != nil {
			return err
		}
		if err := tx.Create(&model.Currency{
			Code:     "USD",
			Name:     "USD",
			FeeType:  model.FIXED,
			FeeValue: decimal.NewFromFloat(0.10),
		}).Error; err != nil {
			return err
		}
		return nil
	},
	Rollback: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.Wallet{}); err != nil {
			return err
		}
		if err := tx.AutoMigrate(&model.Currency{}); err != nil {
			return err
		}
		return nil
	},
}
