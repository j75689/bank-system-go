package http

import (
	"bank-system-go/internal/config"
	"bank-system-go/internal/delivery/http"
	"fmt"

	"github.com/rs/zerolog"
)

type Application struct {
	logger     zerolog.Logger
	config     config.Config
	httpServer *http.HttpServer
}

func (application Application) Start() error {
	application.logger.Info().Msgf("http server listen :%d", application.config.HTTP.Port)
	return application.httpServer.Run(fmt.Sprintf(":%d", application.config.HTTP.Port))
}

func newApplication(
	logger zerolog.Logger,
	config config.Config,
	httpServer *http.HttpServer,
) Application {
	return Application{
		logger:     logger,
		config:     config,
		httpServer: httpServer,
	}
}
