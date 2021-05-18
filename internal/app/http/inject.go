package http

import (
	"bank-system-go/internal/config"
	"bank-system-go/internal/delivery/http"
	"bank-system-go/pkg/logger"
	"bank-system-go/pkg/mq"
	"fmt"
)

type Application struct {
	logger     logger.Logger
	config     config.Config
	httpServer *http.HttpServer
	mq         mq.MQ
}

func (application Application) Start() error {
	application.logger.Info().Msgf("http server listen :%d", application.config.HTTP.Port)
	return application.httpServer.Run(fmt.Sprintf(":%d", application.config.HTTP.Port))
}

func (application Application) Stop() error {
	if err := application.httpServer.Shutdown(); err != nil {
		return err
	}
	return application.mq.Close()
}

func newApplication(
	logger logger.Logger,
	config config.Config,
	httpServer *http.HttpServer,
	mq mq.MQ,
) Application {
	return Application{
		logger:     logger,
		config:     config,
		httpServer: httpServer,
		mq:         mq,
	}
}
