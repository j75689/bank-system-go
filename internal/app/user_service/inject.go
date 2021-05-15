package user

import (
	"bank-system-go/internal/config"
	"bank-system-go/internal/controller"
	"bank-system-go/internal/migration/user"
	"bank-system-go/pkg/logger"
	"context"
	"fmt"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type Application struct {
	logger     logger.Logger
	config     config.Config
	db         *gorm.DB
	controller *controller.UserController
}

func (application Application) Migrate() error {
	m := gormigrate.New(application.db, gormigrate.DefaultOptions, user.Migrations)
	if err := m.Migrate(); err != nil {
		return err
	}
	fmt.Println("migration complete")
	return nil
}

func (application Application) Start() error {
	return application.controller.CreateUser(context.Background())
}

func newApplication(
	logger logger.Logger,
	config config.Config,
	db *gorm.DB,
	controller *controller.UserController,
) Application {
	return Application{
		logger:     logger,
		config:     config,
		db:         db,
		controller: controller,
	}
}
