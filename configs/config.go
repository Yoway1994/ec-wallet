package configs

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DB     DBConfig
	Cache  RedisConfig
	Wallet WalletConfig
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

type RedisConfig struct {
	Host       string
	Port       string
	Auth       string
	Database   int
	Max_active int
}

type WalletConfig struct {
	Mnemonic string
}

// 初始化配置，在 NewConfig 之前調用
func init() {
	// 設置配置文件名稱、類型和路徑
	viper.SetConfigName(".env") // 配置文件名稱
	viper.SetConfigType("env")  // 支持的擴展名

	// 添加搜索路徑
	viper.AddConfigPath(".")   // 先查找項目根目錄
	viper.AddConfigPath("../") // 如果在子目錄運行，還可以查找上一級目錄

	// 讀取環境變量
	viper.AutomaticEnv()

	// 讀取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Warning: Config file not found, using default values")
		} else {
			log.Printf("Error reading config file: %s", err)
		}
	} else {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func NewConfig() *Config {
	// 嘗試從配置文件讀取，如果沒有則使用默認值
	return &Config{
		DB: DBConfig{
			Host:         getStringWithDefault("db.host", "localhost"),
			Port:         getStringWithDefault("db.port", "5432"),
			User:         getStringWithDefault("db.user", "postgres"),
			Password:     getStringWithDefault("db.password", "password"),
			Database:     getStringWithDefault("db.database", "ecwallet"),
			SSLMode:      getStringWithDefault("db.sslmode", "disable"),
			MaxIdleConns: getIntWithDefault("db.maxidleconns", 5),
			MaxOpenConns: getIntWithDefault("db.maxopenconns", 10),
			MaxLifetime:  getDurationWithDefault("db.maxlifetime", 5*time.Minute),
		},
		Cache: RedisConfig{
			Host:       getStringWithDefault("redis.host", "localhost"),
			Port:       getStringWithDefault("redis.port", "6379"),
			Auth:       getStringWithDefault("redis.auth", ""),
			Database:   getIntWithDefault("redis.database", 0),
			Max_active: getIntWithDefault("redis.max_active", 10),
		},
		Wallet: WalletConfig{
			Mnemonic: getStringWithDefault("wallet.mnemonic", "test test test test test test test test test test test junk"),
		},
	}
}

// 帶默認值的輔助函數
func getStringWithDefault(key string, defaultValue string) string {
	if !viper.IsSet(key) {
		return defaultValue
	}
	return viper.GetString(key)
}

func getIntWithDefault(key string, defaultValue int) int {
	if !viper.IsSet(key) {
		return defaultValue
	}
	return viper.GetInt(key)
}

func getDurationWithDefault(key string, defaultValue time.Duration) time.Duration {
	if !viper.IsSet(key) {
		return defaultValue
	}
	return viper.GetDuration(key)
}
