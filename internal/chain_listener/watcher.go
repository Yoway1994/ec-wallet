package chainlistener

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func BnbBlockWatcher() *BlockWatcher {
	return &BlockWatcher{
		handler: func(ctx context.Context, block *types.Block) error {
			// 輸出日誌
			fmt.Printf("📦 新區塊: %d\n", block.NumberU64())
			return nil
		},
	}
}

func BnbTransferWatcher() *TransactionWatcher {
	targetAddress := common.HexToAddress("0x906cAD3F4350CD7d3474CBd5f5DFe056e3BD7908")
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
			fmt.Printf("💰 收到 BNB: 到 %s 來自 %s, 數量: %s\n",
				targetAddress.Hex(), sender.Hex(), amount.String())

			return nil
		},
	}
}

// 創建 EventWatcher 實例
func Erc20TransferWatcher() *EventWatcher {
	return &EventWatcher{
		contractAddress: common.HexToAddress("0x0dEb24A269C09CADA1DdA15bE5E6b8B928596c13"),
		eventSignature:  "Transfer(address,address,uint256)",
		handler: func(ctx context.Context, vLog types.Log) error {
			// 創建 ABI 解析器（可以全局定義一次）
			erc20ABI, err := abi.JSON(strings.NewReader(ERC20ABI))
			if err != nil {
				log.Println("解析 ABI 失敗:", err)
				return err
			}

			// 解析事件數據
			var transferEvent struct {
				From  common.Address
				To    common.Address
				Value *big.Int
			}

			if err := erc20ABI.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data); err != nil {
				log.Printf("解析事件失敗: %v", err)
				return err
			}

			// 從 Topics 提取發送方和接收方地址
			transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
			transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())

			// 輸出日誌
			fmt.Printf("💸 Transfer: 從 %s 到 %s 價值 %s\n",
				transferEvent.From.Hex(), transferEvent.To.Hex(), transferEvent.Value.String())
			return nil
		},
	}
}

// 你可以放你自己的 ERC20 ABI
const ERC20ABI = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},
{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],
"name":"Transfer","type":"event"}]`
