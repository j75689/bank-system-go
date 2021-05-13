package wireset

import (
	"bank-system-go/internal/config"
	"bank-system-go/pkg/logger"
)

func InitLogger(config config.Config) (logger.Logger, error) {
	return logger.NewLogger(config.Logger.Level, config.Logger.Format)
}
