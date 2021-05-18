package wallet

import (
	"bank-system-go/internal/model"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var v202105190000 = &gormigrate.Migration{
	ID: "202105190000",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.WalletHistory{}); err != nil {
			return err
		}
		return nil
	},
	Rollback: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.WalletHistory{}); err != nil {
			return err
		}
		return nil
	},
}
