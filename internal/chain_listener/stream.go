package chainlistener

import (
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	WatchActionStart = "start_watch"
	DataStatusActive = "active"
	AddressStreamKey = "stream:address"
)

func (l *EVMChainListener) OnWatchAddress() {
	lastID := "0" // 或 "$" 表示只收新的
	for {
		select {
		case <-l.ctx.Done():
			log.Println("🛑 OnWatchAddress 停止")
			return
		default:
			streams, err := l.cache.XRead(l.ctx, &redis.XReadArgs{
				Streams: []string{AddressStreamKey, lastID},
				Block:   5 * time.Second,
				Count:   10,
			}).Result()

			if err != nil {
				if err == redis.Nil {
					continue // 無訊息
				}
				log.Printf("❌ Redis XRead 錯誤: %v", err)
				time.Sleep(time.Second)
				continue
			}

			for _, stream := range streams {
				for _, msg := range stream.Messages {
					lastID = msg.ID

					action, ok1 := msg.Values["action"].(string)
					address, ok2 := msg.Values["address"].(string)

					if !ok1 || !ok2 {
						log.Printf("⚠️ 訊息格式錯誤: %+v", msg.Values)
						continue
					}

					switch action {
					case WatchActionStart:
						l.addWatchAddress(address)
					case "remove":
						l.removeWatchAddress(address)
					default:
						log.Printf("⚠️ 未知 action: %s", action)
					}
				}
			}
		}
	}

}

func (l *EVMChainListener) addWatchAddress(address string) {

}

func (l *EVMChainListener) removeWatchAddress(address string) {
	//
}
