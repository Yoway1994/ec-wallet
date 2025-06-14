package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
)

func HandleError(c *gin.Context, err error) {
	// 將 Kratos 錯誤轉換為 HTTP 響應
	if e, ok := err.(*errors.Error); ok {
		c.JSON(int(e.Code), gin.H{
			"code":    e.Code,
			"message": e.Message,
		})
		return
	}

	// 未知錯誤
	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    http.StatusInternalServerError,
		"message": "內部服務錯誤",
	})
}
