package wallet

import (
	"bank-system-go/internal/config"
	"bank-system-go/internal/controller"
	"bank-system-go/internal/migration/wallet"
	"bank-system-go/pkg/logger"
	"context"
	"fmt"

	"github.com/go-gormigrate/gormigrate/v2"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type Application struct {
	logger     logger.Logger
	config     config.Config
	db         *gorm.DB
	controller *controller.WalletController
}

func (application Application) Migrate() error {
	m := gormigrate.New(application.db, gormigrate.DefaultOptions, wallet.Migrations)
	if err := m.Migrate(); err != nil {
		return err
	}
	fmt.Println("migration complete")
	return nil
}

func (application Application) Start() error {

	ctx := context.Background()
	errg := errgroup.Group{}
	errg.Go(func() error {
		return application.controller.CreateWallet(ctx)
	})
	errg.Go(func() error {
		return application.controller.ListWallet(ctx)
	})
	errg.Go(func() error {
		return application.controller.UpdateWalletBalance(ctx)
	})
	errg.Go(func() error {
		return application.controller.Transfer(ctx)
	})

	return errg.Wait()
}

func newApplication(
	logger logger.Logger,
	config config.Config,
	db *gorm.DB,
	controller *controller.WalletController,
) Application {
	return Application{
		logger:     logger,
		config:     config,
		db:         db,
		controller: controller,
	}
}
