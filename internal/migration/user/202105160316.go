package user

import (
	"bank-system-go/internal/model"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var v202105160316 = &gormigrate.Migration{
	ID: "202105160316",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.User{}); err != nil {
			return err
		}
		if err := tx.AutoMigrate(&model.AccessLog{}); err != nil {
			return err
		}
		return nil
	},
	Rollback: func(tx *gorm.DB) error {
		if err := tx.Migrator().DropTable(&model.User{}); err != nil {
			return err
		}
		if err := tx.Migrator().DropTable(&model.AccessLog{}); err != nil {
			return err
		}
		return nil
	},
}
