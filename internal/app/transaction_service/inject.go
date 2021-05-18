package transaction

import (
	"bank-system-go/internal/config"
	"bank-system-go/internal/controller"
	"bank-system-go/internal/migration/transaction"
	"bank-system-go/pkg/logger"
	"bank-system-go/pkg/mq"
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
	controller *controller.TransactionController
	mq         mq.MQ
}

func (application Application) Migrate() error {
	m := gormigrate.New(application.db, gormigrate.DefaultOptions, transaction.Migrations)
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
		return application.controller.CreateTransaction(ctx)
	})
	errg.Go(func() error {
		return application.controller.ListTransaction(ctx)
	})

	return errg.Wait()
}

func (application Application) Stop() error {
	if err := application.mq.Close(); err != nil {
		return err
	}
	db, err := application.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func newApplication(
	logger logger.Logger,
	config config.Config,
	db *gorm.DB,
	controller *controller.TransactionController,
	mq mq.MQ,
) Application {
	return Application{
		logger:     logger,
		config:     config,
		db:         db,
		controller: controller,
		mq:         mq,
	}
}
