package main

import (
	chainlistener "ec-wallet/internal/chain_listener"
	"log"
)

func main() {
	// 初始化監聽器
	listener, err := chainlistener.NewEVMChainListener("wss://bsc-testnet-rpc.publicnode.com")
	if err != nil {
		log.Fatal("初始化監聽器失敗:", err)
	}

	// 設置 ERC20 Transfer 事件監聽器
	listener.RegisterEventWatcher("TokenTransfer", chainlistener.Erc20TransferWatcher())
	// 設置 BNB 轉帳監聽器
	listener.RegisterTransactionWatcher("BNBTransfer", chainlistener.BnbTransferWatcher())
	// 設置區塊監聽器 (可選)
	listener.RegisterBlockWatcher("NewBlock", chainlistener.BnbBlockWatcher())

	// 啟動監聽器
	if err := listener.Start(); err != nil {
		log.Fatal("監聽器錯誤:", err)
	}
}
