package chainlistener

import (
	"context"
	"ec-wallet/internal/wire"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/redis/go-redis/v9"
)

const (
	BscTestnetWebsocketURL = "wss://bsc-testnet-rpc.publicnode.com"
)

type EVMChainListener struct {
	cache               *redis.Client
	client              *ethclient.Client
	eventWatchers       map[string]*EventWatcher
	transactionWatchers map[string]*TransactionWatcher
	blockWatchers       map[string]*BlockWatcher
	ctx                 context.Context
	cancel              context.CancelFunc
	wsURL               string // 保存 URL 用於重連
}

// 事件監聽器（合約事件）
type EventWatcher struct {
	contractAddress common.Address
	eventSignature  string
	handler         func(ctx context.Context, log types.Log) error
}

// 交易監聽器（原生幣轉帳）
type TransactionWatcher struct {
	fromAddress *common.Address // nil 表示監聽所有
	toAddress   *common.Address // nil 表示監聽所有
	handler     func(ctx context.Context, tx *types.Transaction) error
}

// 區塊監聽器
type BlockWatcher struct {
	handler func(ctx context.Context, block *types.Block) error
}

func NewEVMChainListener(wsURL string) (*EVMChainListener, error) {
	client, err := ethclient.Dial(wsURL)
	if err != nil {
		return nil, fmt.Errorf("初始連接失敗: %w", err)
	}

	cache, err := wire.NewRedisClient()
	if err != nil {
		return nil, fmt.Errorf("redis init fail: %w", err)
	}
	return &EVMChainListener{
		client:              client,
		eventWatchers:       make(map[string]*EventWatcher),
		transactionWatchers: make(map[string]*TransactionWatcher),
		blockWatchers:       make(map[string]*BlockWatcher),
		wsURL:               wsURL,
		cache:               cache,
	}, nil
}

func (l *EVMChainListener) RegisterEventWatcher(name string, watcher *EventWatcher) {
	l.eventWatchers[name] = watcher
}

func (l *EVMChainListener) RegisterTransactionWatcher(name string, watcher *TransactionWatcher) {
	l.transactionWatchers[name] = watcher
}

func (l *EVMChainListener) RegisterBlockWatcher(name string, watcher *BlockWatcher) {
	l.blockWatchers[name] = watcher
}

func (l *EVMChainListener) Start() error {
	l.ctx, l.cancel = context.WithCancel(context.Background())
	for {
		select {
		case <-l.ctx.Done():
			return nil
		default:
			err := l.runListeners()
			if err != nil {
				log.Printf("監聽器錯誤: %v，3秒後重啟", err)
				l.reconnect() // 重新建立連接
				time.Sleep(3 * time.Second)
				continue
			}
		}
	}
}

func (l *EVMChainListener) runListeners() error {
	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	// WaitGroup 確保我們等到所有監聽器都結束
	// 才進行重連邏輯

	if len(l.eventWatchers) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := l.startEventSubscription(); err != nil {
				select {
				case errChan <- err:
				default:
				}
			}
		}()
	}

	if len(l.blockWatchers) > 0 || len(l.transactionWatchers) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := l.startBlockSubscription(); err != nil {
				select {
				case errChan <- err:
				default:
				}
			}
		}()
	}

	// 等待所有監聽器結束
	wg.Wait()

	// 檢查是否有錯誤
	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}

// 啟動事件訂閱
func (l *EVMChainListener) startEventSubscription() error {
	query := ethereum.FilterQuery{
		Addresses: l.getWatchedEventAddresses(),
		Topics:    l.getWatchedEventTopics(),
	}

	logs := make(chan types.Log)
	sub, err := l.client.SubscribeFilterLogs(l.ctx, query, logs)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			return err
		case vLog := <-logs:
			l.processEventLog(vLog)
		case <-l.ctx.Done():
			return nil
		}
	}
}

// 啟動區塊訂閱
func (l *EVMChainListener) startBlockSubscription() error {
	headers := make(chan *types.Header)
	sub, err := l.client.SubscribeNewHead(l.ctx, headers)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			return err
		case header := <-headers:
			l.processBlockHeader(header)
		case <-l.ctx.Done():
			return nil
		}
	}
}

// 處理事件日誌
func (l *EVMChainListener) processEventLog(vLog types.Log) {
	for name, watcher := range l.eventWatchers {
		// 檢查合約地址是否匹配
		if vLog.Address != watcher.contractAddress {
			continue
		}

		// 檢查事件簽名是否匹配
		eventHash := crypto.Keccak256Hash([]byte(watcher.eventSignature))
		if len(vLog.Topics) > 0 && vLog.Topics[0] == eventHash {
			if err := watcher.handler(l.ctx, vLog); err != nil {
				log.Printf("事件處理器 %s 處理失敗: %v", name, err)
			}
		}
	}
}

// 處理區塊頭
func (l *EVMChainListener) processBlockHeader(header *types.Header) {
	// 獲取完整區塊
	block, err := l.client.BlockByHash(l.ctx, header.Hash())
	if err != nil {
		log.Printf("獲取區塊失敗: %v", err)
		return
	}

	// 處理區塊監聽器
	for name, watcher := range l.blockWatchers {
		if err := watcher.handler(l.ctx, block); err != nil {
			log.Printf("區塊處理器 %s 處理失敗: %v", name, err)
		}
	}

	// 處理交易監聽器
	if len(l.transactionWatchers) > 0 {
		l.processBlockTransactions(block)
	}
}

// 處理區塊中的交易
func (l *EVMChainListener) processBlockTransactions(block *types.Block) {
	for _, tx := range block.Transactions() {
		for name, watcher := range l.transactionWatchers {
			// 檢查接收地址
			to := tx.To()
			if watcher.toAddress != nil && (to == nil || *to != *watcher.toAddress) {
				continue
			}

			// 檢查發送地址
			if watcher.fromAddress != nil {
				signer := types.LatestSignerForChainID(tx.ChainId())
				sender, err := types.Sender(signer, tx)
				if err != nil || sender != *watcher.fromAddress {
					continue
				}
			}

			// 處理匹配的交易
			if err := watcher.handler(l.ctx, tx); err != nil {
				log.Printf("交易處理器 %s 處理失敗: %v", name, err)
			}
		}
	}
}

func (l *EVMChainListener) reconnect() error {
	l.client.Close()
	client, err := ethclient.Dial(l.wsURL)
	if err != nil {
		return err
	}
	l.client = client
	return nil
}

// 獲取需要監聽的合約地址列表
func (l *EVMChainListener) getWatchedEventAddresses() []common.Address {
	addressMap := make(map[common.Address]bool)
	for _, watcher := range l.eventWatchers {
		addressMap[watcher.contractAddress] = true
	}

	addresses := make([]common.Address, 0, len(addressMap))
	for addr := range addressMap {
		addresses = append(addresses, addr)
	}
	return addresses
}

// 獲取需要監聽的事件主題列表
func (l *EVMChainListener) getWatchedEventTopics() [][]common.Hash {
	// 簡化實現: 僅使用所有事件簽名的 hash 作為第一個主題
	var eventSignatures []common.Hash
	for _, watcher := range l.eventWatchers {
		eventSignatures = append(eventSignatures, crypto.Keccak256Hash([]byte(watcher.eventSignature)))
	}
	return [][]common.Hash{eventSignatures}
}
