package wireset

import (
	"bank-system-go/internal/config"
	"bank-system-go/pkg/database"
	"bank-system-go/pkg/logger"

	"gorm.io/gorm"
)

func InitDatabase(config config.Config, logger logger.Logger) (*gorm.DB, error) {
	return database.NewDataBase(
		config.DataBase.Driver,
		config.DataBase.Host,
		config.DataBase.Port,
		config.DataBase.Database,
		config.DataBase.InstanceName,
		config.DataBase.User,
		config.DataBase.Password,
		config.DataBase.SSLMode,
		config.DataBase.ConnectTimeout,
		config.DataBase.ReadTimeout,
		config.DataBase.WriteTimeout,
		config.DataBase.DialTimeout,
		database.SetConnMaxLifetime(config.DataBase.MaxLifetime),
		database.SetMaxIdleConns(config.DataBase.MaxIdleConn),
		database.SetMaxOpenConns(config.DataBase.MaxOpenConn),
		database.SetConnMaxIdleTime(config.DataBase.MaxIdleTime),
		database.SetLogger(logger.WarpedGormLogger()),
	)
}
