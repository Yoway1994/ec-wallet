package chainlistener

import (
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/redis/go-redis/v9"
)

const (
	AddressStreamKey  = "stream:address"
	TransferStreamKey = "stream:transfer"
	//
	WatchActionStart = "start_watch"
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
						l.handleNewWatchAddress(address)
					default:
						log.Printf("⚠️ 未知 action: %s", action)
					}
				}
			}
		}
	}
}

// 新增監控地址
func (l *EVMChainListener) handleNewWatchAddress(address string) {
	l.watchedAddressMu.Lock()
	defer l.watchedAddressMu.Unlock()
	// 更新地址並通知
	addr := common.HexToAddress(address)
	l.watchedAddress[addr] = true
	log.Println("✅ 新增監控地址:", address)
	l.watchedAddressChanged <- struct{}{}
}

func (l *EVMChainListener) notifyTransferEvent(event *transferEvent) {
	eventData := map[string]interface{}{
		"address":  event.Address.String(),
		"from":     event.From.String(),
		"to":       event.To.String(),
		"value":    event.Value.String(),
		"txHash":   event.TxHash.String(),
		"blockNum": event.BlockNumber,
	}

	// Add to Redis stream
	_, err := l.cache.XAdd(l.ctx, &redis.XAddArgs{
		Stream: TransferStreamKey,
		Values: eventData,
	}).Result()

	if err != nil {
		log.Printf("❌ Redis XAdd 錯誤: %v", err)
		return
	}
}
