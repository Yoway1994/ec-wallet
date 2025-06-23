package streamservice

import (
	"context"
	"ec-wallet/internal/domain/stream"
	"ec-wallet/internal/errors"
	"ec-wallet/internal/infrastructure/logger"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	ttl = 24 * time.Hour
)

func NewStreamService(cache *redis.Client, logger *zap.Logger) stream.Stream {
	return &streamService{cache: cache, logger: logger}
}

type streamService struct {
	cache  *redis.Client
	logger *zap.Logger
}

func (s *streamService) WatchAddress(ctx context.Context, req *stream.WatchAddressRequest) error {
	zapLogger := logger.FromContext(ctx)
	//
	watchKey := s.addressWatchKey(req.Chain, req.Address)

	// 檢查是否已經在監聽
	exists, err := s.cache.Exists(ctx, watchKey).Result()
	if err != nil {
		err = errors.ErrStreamRedisCheckFailed.WithCause(err)
		zapLogger.Error(err.Error())
		return err
	}
	if exists > 0 {
		err = errors.ErrStreamAddressAlreadyWatched.WithMetadata(map[string]string{
			"chain":   req.Chain,
			"address": req.Address,
		})
		zapLogger.Error(err.Error())
		return err
	}

	// 將監聽請求加入 Redis Stream
	watchRequest := map[string]interface{}{
		"action":    stream.WatchActionStart,
		"address":   req.Address,
		"chain":     req.Chain,
		"timestamp": time.Now().Unix(),
	}

	streamID, err := s.cache.XAdd(ctx, &redis.XAddArgs{
		Stream: stream.AddressStreamKey,
		Values: watchRequest,
	}).Result()

	if err != nil {
		err = errors.ErrStreamAddWatchFailed.WithCause(err)
		zapLogger.Error(err.Error())
		return err
	}

	// 記錄監聽狀態
	watchData := map[string]interface{}{
		"stream_id":  streamID,
		"started_at": time.Now().Unix(),
		"status":     stream.DataStatusActive,
	}

	if err := s.cache.HSet(ctx, watchKey, watchData).Err(); err != nil {
		err = errors.ErrStreamAddWatchFailed.WithCause(err)
		zapLogger.Error(err.Error())
		return err
	}

	// 設置過期時間
	if err := s.cache.Expire(ctx, watchKey, ttl).Err(); err != nil {
		err = errors.ErrStreamSetExpiryFailed.WithCause(err)
		zapLogger.Error(err.Error())
		return err
	}

	return nil
}

func (s *streamService) addressWatchKey(chain, address string) string {
	return fmt.Sprintf("watch:%s:%s", chain, address)
}
