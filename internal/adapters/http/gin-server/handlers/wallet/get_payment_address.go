package wallet

import (
	"ec-wallet/internal/adapters/http/gin-server/utils"
	"ec-wallet/internal/errors"
	"ec-wallet/internal/wire"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type PaymentAddressRequest struct {
	OrderID string `form:"order_id" binding:"required"`
}

// PaymentAddressResponse 支付地址的響應格式
type PaymentAddressResponse struct {
	OrderID    string    `json:"order_id" example:"ORD12345678"`
	Address    string    `json:"address" example:"bnb1w7jflwesfnrp0lfgnthkvq55m8gzrlav5ktmyk"`
	Chain      string    `json:"chain" example:"BNB"`
	CreatedAt  time.Time `json:"created_at"`
	ExpireTime time.Time `json:"expire_time"`
}

// ErrorResponse 錯誤響應的格式
type ErrorResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Invalid order ID"`
}

// GetPaymentAddress godoc
// @Summary Get payment address for an order
// @Description Allocate a BNB address for a specific order number
// @Tags payment
// @Accept json
// @Produce json
// @Param order_id query string true "Order ID"
// @Success 200 {object} PaymentAddressResponse
// @Failure 400 {object} ErrorResponse "Invalid order ID"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /wallet/v1/payment-address [get]
func GetPaymentAddress(c *gin.Context) {
	var req PaymentAddressRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.HandleError(c, errors.ErrInvalidParameter.WithCause(err))
		return
	}

	// 驗證訂單編號格式是否有效
	if len(req.OrderID) < 5 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid order ID format",
		})
		return
	}

	// 檢查是否已經為此訂單分配了地址（保證冪等性）
	// TODO: 實現從資料庫查詢是否已存在訂單-地址映射

	wallet, err := wire.NewWallet()
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// 4. 如果沒有現有地址，從地址池中分配一個
	// TODO: 實現地址池及分配邏輯
	// 這裡僅作為示例
	reservation, err := wallet.AcquireAddress(c)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// 5. 保存訂單和地址的對應關係到資料庫
	// TODO: 實現資料庫寫入邏輯

	// 7. 構建並返回響應
	now := time.Now()
	expireTime := now.Add(24 * time.Hour) // 支付地址有效期為 24 小時

	c.JSON(http.StatusOK, PaymentAddressResponse{
		OrderID:    req.OrderID,
		Address:    reservation.Address,
		Chain:      "BNB",
		CreatedAt:  now,
		ExpireTime: expireTime,
	})
}
