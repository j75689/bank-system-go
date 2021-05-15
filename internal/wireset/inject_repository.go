package wireset

import (
	"bank-system-go/internal/repository"
	"bank-system-go/internal/repository/gorm_repository"

	"github.com/google/wire"
)

var RepositorySet = wire.NewSet(
	gorm_repository.NewGORMRepository,
	InitUserRepository,
	InitWalletRepository,
	InitTransationRepository,
)

func InitUserRepository(gorm *gorm_repository.GORMRepository) repository.UserRepository {
	return gorm
}

func InitWalletRepository(gorm *gorm_repository.GORMRepository) repository.WalletRepository {
	return gorm
}

func InitTransationRepository(gorm *gorm_repository.GORMRepository) repository.TransationRepository {
	return gorm
}
