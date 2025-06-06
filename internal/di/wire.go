//go:build wireinject
// +build wireinject

package di

import (
	"ec-wallet/internal/domain"
	"ec-wallet/internal/infrastructure/database"
	"ec-wallet/internal/infrastructure/repository/gormRepo"
	"sync"

	"github.com/google/wire"
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

func NewRepository() (domain.Repository, error) {
	panic(wire.Build(gormRepo.NewRepository))
}
