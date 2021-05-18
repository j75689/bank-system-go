package wallet

import (
	"bank-system-go/internal/model"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var v202105181826 = &gormigrate.Migration{
	ID: "202105181826",
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
