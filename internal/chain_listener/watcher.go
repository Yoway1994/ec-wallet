package chainlistener

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const TransferEventSignature = "Transfer(address,address,uint256)"

var erc20ABI abi.ABI

func init() {
	ERC20ABI := `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},
{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],
"name":"Transfer","type":"event"}]`
	var err error
	erc20ABI, err = abi.JSON(strings.NewReader(ERC20ABI))
	if err != nil {
		log.Fatal("解析 ABI 失敗:", err)
	}
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

func (l *EVMChainListener) RegisterEventWatcher(name string, watcher *EventWatcher) {
	l.eventWatchers[name] = watcher
}

func (l *EVMChainListener) RegisterBlockWatcher(name string, watcher *BlockWatcher) {
	l.blockWatchers[name] = watcher
}

func (l *EVMChainListener) RegisterTransactionWatcher(name string, watcher *TransactionWatcher) {
	l.transactionWatchers[name] = watcher
}

func BnbBlockWatcher() *BlockWatcher {
	return &BlockWatcher{
		handler: func(ctx context.Context, block *types.Block) error {
			// 輸出日誌
			fmt.Printf("📦 新區塊: %d\n", block.NumberU64())
			return nil
		},
	}
}

func NativeTokenTransferWatcher(targetAddress common.Address) *TransactionWatcher {
	return &TransactionWatcher{
		fromAddress: nil,            // 不限制發送方
		toAddress:   &targetAddress, // 只監聽發送到目標地址的交易
		handler: func(ctx context.Context, tx *types.Transaction) error {
			// 獲取交易金額
			amount := tx.Value()
			// 獲取發送方地址
			signer := types.LatestSignerForChainID(tx.ChainId())
			sender, err := types.Sender(signer, tx)
			if err != nil {
				log.Printf("無法獲取發送人地址: %v", err)
				return err
			}

			// 輸出日誌
			fmt.Printf("💰 收到 原生幣: 到 %s 來自 %s, 數量: %s\n",
				targetAddress.Hex(), sender.Hex(), amount.String())

			return nil
		},
	}
}

// 創建 EventWatcher 實例
func ERC20TransferWatcher(contractAddr common.Address) *EventWatcher {
	var transferEvent transferEvent
	return &EventWatcher{
		contractAddress: contractAddr,
		eventSignature:  TransferEventSignature,
		handler: func(ctx context.Context, vLog types.Log) error {
			// 解析事件數據, 這只會解碼 data 裡面的內容，也就是 value，不會去讀 topics！
			if err := erc20ABI.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data); err != nil {
				log.Printf("解析事件失敗: %v", err)
				return err
			}

			// 從 Topics 提取發送方和接收方地址
			transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
			transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())

			// 輸出日誌, 發現監聽到特定事件
			fmt.Printf("💸 Transfer: 從 %s 到 %s 價值 %s\n",
				transferEvent.From.Hex(), transferEvent.To.Hex(), transferEvent.Value.String())
			//
			return nil
		},
	}
}
