//+build wireinject

//The build tag makes sure the stub is not built in the final build.

package http

import (
	"bank-system-go/internal/config"
	"bank-system-go/internal/delivery/http"
	"bank-system-go/internal/wireset"

	"github.com/google/wire"
)

func Initialize(configPath string) (Application, error) {
	wire.Build(
		newApplication,
		config.NewConfig,
		wireset.InitLogger,
		wireset.InitMQ,
		http.NewHttpServer,
	)
	return Application{}, nil
}
