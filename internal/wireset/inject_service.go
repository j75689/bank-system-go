package wireset

import (
	v1 "bank-system-go/internal/service/v1"

	"github.com/google/wire"
)

var ServiceV1Set = wire.NewSet(
	v1.NewUserService,
	v1.NewWalletService,
)
