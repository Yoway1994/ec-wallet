package ginserver

import (
	"ec-wallet/internal/adapters/http/gin-server/handlers"
	"ec-wallet/internal/adapters/http/gin-server/middleware"
	"ec-wallet/internal/wire"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 基本middleware, 將來可替換
	r.Use(gin.Recovery())
	// r.Use(gin.Logger())

	//
	r.Use(middleware.RequestIDMiddleware(""))
	baseLogger, err := wire.NewLogger()
	if err != nil {
		panic(err)
	}
	r.Use(middleware.LoggerMiddleware(baseLogger, middleware.LoggerConfig{
		LogQueryParams: false,
		LogUserAgent:   false,
	}))

	// health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	api := r.Group("/api")
	{
		api.POST("/v1/payment-orders", handlers.CreatePaymentOrder)
	}

	return r
}
