package wallet

import (
	"time"

	"github.com/gin-gonic/gin"
)

// WalletResponse 錢包信息的響應格式
type WalletResponse struct {
	ID        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Address   string    `json:"address" example:"0x71C7656EC7ab88b098defB751B7401B5f6d8976F"`
	Type      string    `json:"type" example:"ETH"`
	Balance   string    `json:"balance" example:"0.2451"`
	CreatedAt time.Time `json:"created_at" example:"2025-06-07T02:33:34+08:00"`
}

// ErrorResponse 錯誤響應的格式
type ErrorResponse struct {
	Code    int    `json:"code" example:"404"`
	Message string `json:"message" example:"Not found"`
}

// Get godoc
// @Summary Get wallet details
// @Description Get details of a specific wallet by ID
// @Tags wallets
// @Accept json
// @Produce json
// @Param id path string true "Wallet ID"
// @Success 200 {object} WalletResponse
// @Failure 404 {object} ErrorResponse
// @Router /wallets/{id} [get]
func Get(c *gin.Context) {
}
