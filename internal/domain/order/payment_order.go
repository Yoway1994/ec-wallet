package order

import (
	"context"
	"ec-wallet/internal/domain/stream"
	"time"

	"github.com/shopspring/decimal"
)

type Order interface {
	CreatePaymentOrder(ctx context.Context, order *PaymentOrder) error
	ValidateOrder(ctx context.Context, event *stream.TransferEvent) (bool, error)
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
	AmountUSD  decimal.Decimal
	Status     string
	TxHash     string
	ExpireTime time.Time
	PaidAt     *time.Time
	CreatedAt  time.Time
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
	AmountUSD  decimal.Decimal
	ExpireTime time.Time
}

type QueryPaymentOrdersParams struct {
	ID      *uint64
	Address *string
}
