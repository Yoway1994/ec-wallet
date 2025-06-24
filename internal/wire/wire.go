//go:build wireinject
// +build wireinject

package wire

import (
	"ec-wallet/configs"
	"ec-wallet/internal/domain"
	"ec-wallet/internal/domain/order"
	"ec-wallet/internal/domain/stream"
	"ec-wallet/internal/domain/wallet"
	"ec-wallet/internal/infrastructure/cache"
	"ec-wallet/internal/infrastructure/database"
	"ec-wallet/internal/infrastructure/logger"
	orderservice "ec-wallet/internal/infrastructure/order"
	gormRepoImpl "ec-wallet/internal/infrastructure/repository/gorm_repo"
	streamservice "ec-wallet/internal/infrastructure/stream"
	walletservice "ec-wallet/internal/infrastructure/wallet"

	"sync"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
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

var zLog *zap.Logger
var zLogOnce sync.Once

func NewLogger() *zap.Logger {
	if zLog == nil {
		zLogOnce.Do(func() {
			config := NewConfig()
			zLog = logger.NewLogger(config)
		})
	}
	return zLog
}

var tokens map[string]*order.PaymentToken
var tokensOnce sync.Once

func NewTokens() map[string]*order.PaymentToken {
	if tokens == nil {
		tokensOnce.Do(func() {
			db, _ := NewDB()
			tokens = orderservice.ProvideTokens(db)
		})
	}
	return tokens
}

func NewRepository() (domain.Repo, error) {
	panic(wire.Build(gormRepoImpl.NewRepository, NewDB))
}

func NewWalletService() (wallet.Wallet, error) {
	panic(wire.Build(walletservice.NewWalletService, NewConfig, NewRepository))
}

func NewStreamService() (stream.Stream, error) {
	panic(wire.Build(streamservice.NewStreamService, NewRedisClient, NewLogger))
}

func NewOrderService() (order.Order, error) {
	panic(wire.Build(orderservice.NewOrderService, NewRepository, NewTokens))
}
