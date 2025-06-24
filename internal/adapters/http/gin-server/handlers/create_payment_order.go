package handlers

import (
	"ec-wallet/internal/adapters/http/gin-server/utils"
	"ec-wallet/internal/domain/order"
	"ec-wallet/internal/domain/stream"
	"ec-wallet/internal/domain/wallet"
	"ec-wallet/internal/errors"
	"ec-wallet/internal/infrastructure/logger"
	"ec-wallet/internal/wire"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type PaymentAddressRequest struct {
	OrderID   string `json:"order_id" binding:"required"`
	Chain     string `json:"chain" binding:"required"`
	AmountUsd string `json:"amount_usd" binding:"required"`
	Token     string `json:"token" binding:"required"`
}

// PaymentAddressResponse 支付地址的響應格式
type PaymentAddressResponse struct {
	OrderID       string    `json:"order_id" example:"ORD12345678"`
	ReservationID string    `json:"reservation_id"`
	Address       string    `json:"address" example:"0x6C318c04Ed42cEe76a61870543bf70F55aEf1fdb"`
	Chain         string    `json:"chain" example:"BSC"`
	CreatedAt     time.Time `json:"created_at"`
	ExpireTime    time.Time `json:"expire_time"`
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
// @Router /v1/payment-orders [post]
func CreatePaymentOrder(c *gin.Context) {
	zapLogger := logger.FromGinContext(c)
	//
	zapLogger.Debug("開始獲取支付地址")
	var req PaymentAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleError(c, errors.ErrInvalidParameter.WithCause(err))
		return
	}
	zapLogger.Debug("解析Request Body",
		zap.String("收款鏈別", req.Chain),
		zap.String("付款數量", req.AmountUsd),
	)

	// 轉換amount usd型別
	amountUsd, err := decimal.NewFromString(req.AmountUsd)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// 分配付款地址
	walletService, err := wire.NewWalletService()
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	reservation, err := walletService.AcquireAddress(c,
		wallet.WithOrderID(req.OrderID))
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	zapLogger.Debug("分配地址", zap.String("到期時間", reservation.ExpiresAt.String()))

	// 通知開始監聽特定地址
	streamService, err := wire.NewStreamService()
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	err = streamService.WatchAddress(c, &stream.WatchAddressRequest{
		Address: reservation.Address,
		Chain:   req.Chain,
	})
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// 保存訂單和地址的對應關係到資料庫
	orderService, err := wire.NewOrderService()
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	newOrder := order.NewPaymentOrder(&order.NewPaymentOrderParams{
		OrderID:    req.OrderID,
		Address:    reservation.Address,
		Chain:      req.Chain,
		Token:      req.Token,
		AmountUSD:  amountUsd,
		ExpireTime: reservation.ExpiresAt,
	})
	err = orderService.CreatePaymentOrder(c, newOrder)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	//
	c.JSON(http.StatusOK, PaymentAddressResponse{
		OrderID:       req.OrderID,
		ReservationID: reservation.ReservationID,
		Address:       reservation.Address,
		Chain:         req.Chain,
		CreatedAt:     reservation.ReservedAt,
		ExpireTime:    reservation.ExpiresAt,
	})
}
