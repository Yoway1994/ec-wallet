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

	// 註冊區塊更新
	listener.RegisterBlockWatcher("NewBlock", chainlistener.BscBlockWatcher())

	// 監聽地址
	go listener.OnWatchAddress()

	// 啟動監聽器
	if err := listener.Start(); err != nil {
		log.Fatal("監聽器錯誤:", err)
	}

}
