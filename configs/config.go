package configs

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DB DBConfig
}

type DBConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	Database     string
	SSLMode      string
	MaxIdleConns int
	MaxOpenConns int
	MaxLifetime  time.Duration
}

func NewConfig() *Config {
	return &Config{
		DB: DBConfig{
			Host:         viper.GetString("db.host"),
			Port:         viper.GetString("db.port"),
			User:         viper.GetString("db.user"),
			Password:     viper.GetString("db.password"),
			Database:     viper.GetString("db.database"),
			SSLMode:      viper.GetString("db.sslmode"),
			MaxIdleConns: viper.GetInt("db.maxidleconns"),
			MaxOpenConns: viper.GetInt("db.maxopenconns"),
			MaxLifetime:  viper.GetDuration("db.maxlifetime"),
		},
	}
}
