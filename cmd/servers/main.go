package main

import (
	_ "ec-wallet/docs/swagger"
	ginserver "ec-wallet/internal/adapters/http/gin-server"
	streamadapter "ec-wallet/internal/adapters/stream"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title EC-Wallet API
// @version 1.0
// @description API Server for EC-Wallet application
// @BasePath /api

func main() {
	// background
	go streamadapter.StartListeningHandler()

	// client
	router := ginserver.SetupRouter()

	// Swagger 文檔路徑
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// start listening...
	router.Run()
}
