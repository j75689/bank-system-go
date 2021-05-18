package transation

import (
	"bank-system-go/internal/model"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var v202105181423 = &gormigrate.Migration{
	ID: "202105181423",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.Transation{}); err != nil {
			return err
		}
		return nil
	},
	Rollback: func(tx *gorm.DB) error {
		if err := tx.Migrator().DropTable(&model.Transation{}); err != nil {
			return err
		}
		return nil
	},
}
