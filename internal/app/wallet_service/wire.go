//+build wireinject

//The build tag makes sure the stub is not built in the final build.

package wallet

import (
	"bank-system-go/internal/config"
	"bank-system-go/internal/controller"
	"bank-system-go/internal/wireset"

	"github.com/google/wire"
)

func Initialize(configPath string) (Application, error) {
	wire.Build(
		newApplication,
		config.NewConfig,
		wireset.InitLogger,
		wireset.InitDatabase,
		wireset.InitMQ,
		wireset.RepositorySet,
		wireset.ServiceV1Set,
		controller.NewWalletController,
	)
	return Application{}, nil
}
