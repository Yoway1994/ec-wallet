package order

import (
	"context"
	"time"
)

type Order interface {
	CreatePaymentOrder(ctx context.Context, order *PaymentOrder) error
}

const (
	PaymentStatusPending   = "PENDING"
	PaymentStatusCompleted = "COMPLETED"
	PaymentStatusExpired   = "EXPIRED"
	PaymentStatusCancelled = "CANCELLED"
	PaymentStatusFailed    = "FAILED"
)

type PaymentOrder struct {
	OrderID    string
	Address    string
	Chain      string
	Token      string
	AmountUSD  float64
	Status     string
	TxHash     string
	ExpireTime time.Time
	PaidAt     *time.Time
}

func NewPaymentOrder(params *NewPaymentOrderParams) *PaymentOrder {
	return &PaymentOrder{
		OrderID:    params.OrderID,
		Address:    params.Address,
		Chain:      params.Chain,
		Token:      params.Token,
		AmountUSD:  params.AmountUSD,
		Status:     PaymentStatusPending, // 預設PENDING
		ExpireTime: params.ExpireTime,
	}
}

type NewPaymentOrderParams struct {
	OrderID    string
	Address    string
	Chain      string
	Token      string
	AmountUSD  float64
	ExpireTime time.Time
}
