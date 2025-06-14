//go:build wireinject
// +build wireinject

package wire

import (
	"ec-wallet/configs"
	gormrepo "ec-wallet/internal/domain/gorm_repo"
	"ec-wallet/internal/domain/wallet"
	"ec-wallet/internal/infrastructure/database"
	gormRepoImpl "ec-wallet/internal/infrastructure/repository/gorm_repo"
	walletservice "ec-wallet/internal/infrastructure/wallet"
	"sync"

	"github.com/google/wire"
	"gorm.io/gorm"
)

var db *gorm.DB
var dbOnce sync.Once

var config *configs.Config
var configOnce sync.Once

func NewDB() (*gorm.DB, error) {
	var err error
	if db == nil {
		dbOnce.Do(func() {
			db, err = database.PostgresqlConnection()
			if err != nil {
				return
			}
		})
	}
	return db, err
}

func NewConfig() *configs.Config {
	if config == nil {
		configOnce.Do(func() {
			config = configs.NewConfig()
		})
	}
	return config
}

func NewRepository() (gormrepo.Repo, error) {
	panic(wire.Build(gormRepoImpl.NewRepository, NewDB))
}

func NewWallet() (wallet.Wallet, error) {
	panic(wire.Build(walletservice.NewWalletService, NewConfig, NewRepository))
}
