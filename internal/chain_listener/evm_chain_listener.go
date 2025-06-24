package chainlistener

import (
	"context"
	"ec-wallet/internal/wire"
	"fmt"
	"log"
	"math/big"
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
	// infra
	cache  *redis.Client
	client *ethclient.Client
	ctx    context.Context
	cancel context.CancelFunc
	wsURL  string // 保存 URL 用於重連
	// watcher
	eventWatchers       map[string]*EventWatcher
	transactionWatchers map[string]*TransactionWatcher
	blockWatchers       map[string]*BlockWatcher
	// 地址轉帳監聽
	watchedAddress        map[common.Address]bool
	watchedAddressMu      sync.RWMutex
	watchedAddressChanged chan struct{}
	//
	// txTrackers   map[common.Hash]*TxConfirmationTracker
	// trackerMutex sync.RWMutex
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
		client:                client,
		wsURL:                 wsURL,
		cache:                 cache,
		eventWatchers:         make(map[string]*EventWatcher),
		transactionWatchers:   make(map[string]*TransactionWatcher),
		blockWatchers:         make(map[string]*BlockWatcher),
		watchedAddress:        make(map[common.Address]bool),
		watchedAddressChanged: make(chan struct{}, 100),
	}, nil
}

func (l *EVMChainListener) Start() error {
	l.ctx, l.cancel = context.WithCancel(context.Background())
	for {
		select {
		case <-l.ctx.Done():
			return nil
		default:
			log.Println("✅ 監聽器啟動...")
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
	wg.Add(1)
	go func() {
		defer wg.Done()
		// 僅處理ERC20轉帳事件
		log.Println("✅ 訂閱ERC20轉帳事件...")
		if err := l.startTransferEventSubscription(); err != nil {
			select {
			case errChan <- err:
			default:
			}
		}
	}()

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
func (l *EVMChainListener) startTransferEventSubscription() error {
	// 使用緩衝通道
	logs := make(chan types.Log, 1000)
	// dispatcher
	go l.dispatchLogs(logs)

	for {
		log.Println("✅ 動態更新訂閱...")
		// 動態訂閱topic
		query := ethereum.FilterQuery{
			FromBlock: big.NewInt(55667838),
			ToBlock:   nil,
			Addresses: l.getStableCoinAddresses(),
			Topics:    l.getWatchedTransferEventTopics(),
		}
		log.Println("開始訂閱...")
		sub, err := l.client.SubscribeFilterLogs(l.ctx, query, logs)
		if err != nil {
			fmt.Println("訂閱失敗:", err)
			return err
		}
		subscriptionDone := false
		log.Println("✅ 訂閱成功，開始處理日誌...")
		// 處理日誌更新
		for !subscriptionDone { // 成功訂閱時進入for迴圈
			select {
			case err := <-sub.Err():
				return err
			case <-l.watchedAddressChanged: // 沒有地址更新時阻塞
				log.Println("✅ 檢測到地址更新，重新啟動訂閱...")
				sub.Unsubscribe()
				subscriptionDone = true // 跳出for迴圈並重新訂閱
			case <-l.ctx.Done():
				return nil
			}
		}
	}
}

func (l *EVMChainListener) dispatchLogs(logs chan types.Log) {
	for {
		select {
		case log := <-logs:
			fmt.Println("✅ 收到日誌...")
			go l.processTransferEventLog(log)
		case <-l.ctx.Done():
			return
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
			go l.processBlockHeader(header)
		case <-l.ctx.Done():
			return nil
		}
	}
}

type transferEvent struct {
	From        common.Address
	To          common.Address
	Address     common.Address
	Value       *big.Int
	BlockNumber uint64
	TxHash      common.Hash
}

func (l *EVMChainListener) processTransferEventLog(vLog types.Log) {
	var transferEvent transferEvent
	// 解析事件數據
	if err := erc20ABI.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data); err != nil {
		log.Printf("解析事件失敗: %v", err)
	}

	// 從 Topics 提取發送方和接收方地址
	transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
	transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())
	transferEvent.BlockNumber = vLog.BlockNumber
	transferEvent.TxHash = vLog.TxHash
	log.Printf("🔔 Transfer Event: %s | From: %s | To: %s | Value: %s | Block: %d | TxHash: %s",
		transferEvent.Address.Hex(),
		transferEvent.From.Hex(),
		transferEvent.To.Hex(),
		transferEvent.Value.String(),
		transferEvent.BlockNumber,
		transferEvent.TxHash.Hex(),
	)
	// 通知ec wallet service發現轉帳
	l.notifyTransferEvent(&transferEvent)
	// 通知tx tracker開始追蹤

}

// 處理事件日誌
// func (l *EVMChainListener) processEventLog(vLog types.Log) {
// 	for name, watcher := range l.eventWatchers {
// 		// 檢查合約地址是否匹配
// 		if vLog.Address != watcher.contractAddress {
// 			continue
// 		}

// 		// 檢查事件簽名是否匹配
// 		eventHash := crypto.Keccak256Hash([]byte(watcher.eventSignature))
// 		if len(vLog.Topics) > 0 && vLog.Topics[0] == eventHash {
// 			if err := watcher.handler(l.ctx, vLog); err != nil {
// 				log.Printf("事件處理器 %s 處理失敗: %v", name, err)
// 			}
// 		}
// 	}
// }

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
// func (l *EVMChainListener) getWatchedEventAddresses() []common.Address {
// 	// 先用map去除冗餘
// 	addressMap := make(map[common.Address]bool)
// 	for _, watcher := range l.eventWatchers {
// 		addressMap[watcher.contractAddress] = true
// 	}

// 	addresses := make([]common.Address, 0, len(addressMap))
// 	for addr := range addressMap {
// 		addresses = append(addresses, addr)
// 	}
// 	return addresses
// }

// 獲取穩定幣合約地址
func (l *EVMChainListener) getStableCoinAddresses() []common.Address {
	return []common.Address{
		common.HexToAddress("0x0dEb24A269C09CADA1DdA15bE5E6b8B928596c13"), // USDC-bsc-test
	}
}

// 修正版本：正確獲取需要監聽的轉帳事件主題列表
func (l *EVMChainListener) getWatchedTransferEventTopics() [][]common.Hash {
	// 準備 topics[2] - to 地址過濾
	toTopics := make([]common.Hash, 0, len(l.watchedAddress))
	if len(l.watchedAddress) > 0 {
		for addr := range l.watchedAddress {
			paddedAddr := common.BytesToHash(common.LeftPadBytes(addr.Bytes(), 32))
			toTopics = append(toTopics, paddedAddr)
		}
	}

	topics := [][]common.Hash{
		{crypto.Keccak256Hash([]byte(TransferEventSignature))}, // topics[0]: Transfer(address,address,uint256)
		{}, // topics[1]: 不過濾 from 地址 (空切片表示不過濾)
		toTopics,
	}

	return topics
}

func (l *EVMChainListener) getTransferEventTopics() [][]common.Hash {
	return [][]common.Hash{{crypto.Keccak256Hash([]byte(TransferEventSignature))}}
}
