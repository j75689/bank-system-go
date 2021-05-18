package transaction

import (
	"bank-system-go/internal/model"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var v202105190001 = &gormigrate.Migration{
	ID: "202105190001",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.Transaction{}); err != nil {
			return err
		}
		return nil
	},
	Rollback: func(tx *gorm.DB) error {
		if err := tx.Migrator().DropTable(&model.Transaction{}); err != nil {
			return err
		}
		return nil
	},
}
