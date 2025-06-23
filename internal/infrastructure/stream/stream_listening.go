package streamservice

import (
	"context"
	"ec-wallet/internal/domain/stream"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func (s *streamService) OnListening(ctx context.Context, transferEventChan chan stream.TransferEvent) {
	zapLogger := s.logger.With(
		zap.String("module", "streamService"),
	)

	lastID := "0" // 或 "$" 表示只收新的
	for {
		select {
		case <-ctx.Done():
			zapLogger.Info("🛑 OnListening 停止")
			return
		default:
			streams, err := s.cache.XRead(ctx, &redis.XReadArgs{
				Streams: []string{stream.TransferStreamKey, lastID},
				Block:   5 * time.Second,
				Count:   10,
			}).Result()

			if err != nil {
				if err == redis.Nil {
					continue // 無訊息
				}
				zapLogger.Error("❌ Redis XRead 錯誤", zap.Error(err))
				time.Sleep(time.Second)
				continue
			}

			for _, s := range streams {
				for _, msg := range s.Messages {
					lastID = msg.ID

					var transferEvent stream.TransferEvent

					// 从 msg.Values 中获取值
					transferEvent.From, _ = msg.Values["from"].(string)
					transferEvent.To, _ = msg.Values["to"].(string)
					transferEvent.Value, _ = msg.Values["value"].(string)
					transferEvent.TxHash, _ = msg.Values["txHash"].(string)
					transferEvent.BlockNum, _ = msg.Values["blockNum"].(uint64)

					transferEventChan <- transferEvent
				}
			}
		}
	}
}
