package handlers

import (
	"ec-wallet/internal/adapters/http/gin-server/utils"
	"ec-wallet/internal/domain/stream"
	"ec-wallet/internal/errors"
	"ec-wallet/internal/wire"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type PaymentAddressRequest struct {
	OrderID string `json:"order_id" binding:"required"`
	Chain   string `json:"chain" binding:"required"`
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

// CreatePaymentOrder godoc
// @Summary Create a payment order with cryptocurrency address
// @Description Create a payment order and allocate a cryptocurrency address for payment collection
// @Tags payment
// @Accept json
// @Produce json
// @Param payment_order body PaymentAddressRequest true "Payment Order Details"
// @Success 200 {object} PaymentAddressResponse
// @Failure 400 {object} ErrorResponse "Invalid order information"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /v1/payment-orders [post]
func CreatePaymentOrder(c *gin.Context) {
	var req PaymentAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleError(c, errors.ErrInvalidParameter.WithCause(err))
		return
	}

	// 分配付款地址
	walletService, err := wire.NewWallet()
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	reservation, err := walletService.AcquireAddress(c)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// 保存訂單和地址的對應關係到資料庫

	// 通知開始監聽特定地址
	streamService, err := wire.NewStreamService()
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	err = streamService.WatchAddress(c, &stream.WatchAddressRequest{
		Address: reservation.Address,
		Chain:   "BNB",
	})
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	//
	c.JSON(http.StatusOK, PaymentAddressResponse{
		OrderID:    reservation.ReservationID,
		Address:    reservation.Address,
		Chain:      "",
		CreatedAt:  reservation.ReservedAt,
		ExpireTime: reservation.ExpiresAt,
	})
}
