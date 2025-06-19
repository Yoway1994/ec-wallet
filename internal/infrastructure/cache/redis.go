package cache

import (
	"context"
	"ec-wallet/configs"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient() (*redis.Client, error) {
	config := configs.NewConfig()
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Cache.Host, config.Cache.Port),
		Password: config.Cache.Auth,
		DB:       config.Cache.Database,
		PoolSize: config.Cache.MaxActive,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
