package ginserver

import (
	"ec-wallet/internal/adapters/http/gin-server/handlers/user"
	"ec-wallet/internal/adapters/http/gin-server/handlers/wallet"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 基本middleware, 將來可替換
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	api := r.Group("/api")
	{
		userV1 := api.Group("/user/v1")
		{
			userV1.GET("", user.Get)
		}

		walletV1 := api.Group("wallet/v1")
		{
			walletV1.GET("", wallet.Get)
		}
	}

	return r
}
