//go:build wireinject
// +build wireinject

package wire

import (
	"ec-wallet/configs"
	gormrepo "ec-wallet/internal/domain/gorm_repo"
	"ec-wallet/internal/domain/stream"
	"ec-wallet/internal/domain/wallet"
	"ec-wallet/internal/infrastructure/cache"
	"ec-wallet/internal/infrastructure/database"
	gormRepoImpl "ec-wallet/internal/infrastructure/repository/gorm_repo"
	streamservice "ec-wallet/internal/infrastructure/stream"
	walletservice "ec-wallet/internal/infrastructure/wallet"
	"sync"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var db *gorm.DB
var dbOnce sync.Once

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

var config *configs.Config
var configOnce sync.Once

func NewConfig() *configs.Config {
	if config == nil {
		configOnce.Do(func() {
			config = configs.NewConfig()
		})
	}
	return config
}

var redisClient *redis.Client
var redisClientOnce sync.Once

func NewRedisClient() (*redis.Client, error) {
	var err error
	if redisClient == nil {
		redisClientOnce.Do(func() {
			redisClient, err = cache.NewRedisClient()
			if err != nil {
				return
			}
		})
	}
	return redisClient, nil
}

func NewRepository() (gormrepo.Repo, error) {
	panic(wire.Build(gormRepoImpl.NewRepository, NewDB))
}

func NewWallet() (wallet.Wallet, error) {
	panic(wire.Build(walletservice.NewWalletService, NewConfig, NewRepository))
}

func NewStreamService() (stream.Stream, error) {
	panic(wire.Build(streamservice.NewStreamService, NewRedisClient))
}
