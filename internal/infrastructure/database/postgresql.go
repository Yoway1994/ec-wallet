package database

import (
	"ec-wallet/configs"
	"ec-wallet/internal/errors"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 目前沒有多db需求, 將來再從domain層使用
func PostgresqlConnection() (*gorm.DB, error) {
	config := configs.NewConfig()
	return PostgresqlConnectionWithConfig(config)

}

func PostgresqlConnectionWithConfig(config *configs.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.DB.Host, config.DB.Port, config.DB.User, config.DB.Password, config.DB.Database, config.DB.SSLMode,
	)

	// 配置GORM日誌
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // 慢SQL閾值
			LogLevel:                  logger.Info, // 日誌級別
			IgnoreRecordNotFoundError: true,        // 忽略記錄未找到錯誤
			Colorful:                  true,        // 啟用彩色
		},
	)

	// 建立連接
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		return nil, errors.ErrDatabaseUnavailable.WithCause(err)
	}

	// 獲取原生SQL連接以配置連接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.ErrDatabaseUnavailable.WithCause(err)
	}

	// 配置連接池
	sqlDB.SetMaxIdleConns(config.DB.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.DB.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.DB.MaxLifetime)

	return db, nil
}
