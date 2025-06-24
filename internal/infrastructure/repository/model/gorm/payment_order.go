package model

import (
	"ec-wallet/internal/domain/order"
	"time"

	"github.com/shopspring/decimal"
)

// PaymentOrder represents the payment_orders table in the database
type PaymentOrder struct {
	ID         uint64          `gorm:"primaryKey"`
	OrderID    string          `gorm:"column:order_id;type:varchar(64);unique;not null"`
	Address    string          `gorm:"column:address;type:varchar(128);not null"`
	Chain      string          `gorm:"column:chain;type:varchar(20);not null"`
	Token      string          `gorm:"column:token;type:varchar(20);not null"`
	AmountUSD  decimal.Decimal `gorm:"column:amount_usd;type:decimal(20,8);not null"`
	Status     string          `gorm:"column:status;type:varchar(20);default:pending;not null"`
	TxHash     string          `gorm:"column:tx_hash;type:varchar(128)"`
	CreatedAt  time.Time       `gorm:"column:created_at;type:timestamp with time zone;autoCreateTime;not null"`
	UpdatedAt  time.Time       `gorm:"column:updated_at;type:timestamp with time zone"`
	ExpireTime time.Time       `gorm:"column:expire_time;type:timestamp with time zone;not null"`
	PaidAt     *time.Time      `gorm:"column:paid_at;type:timestamp with time zone"`
}

// TableName overrides the table name
func (p *PaymentOrder) TableName() string {
	return "payment_orders"
}

// PaymentOrderDomainToModel 轉換domain為model
func PaymentOrderDomainToModel(d *order.PaymentOrder) (m *PaymentOrder) {
	return &PaymentOrder{
		// ID:         d.ID,
		OrderID:   d.OrderID,
		Address:   d.Address,
		Chain:     d.Chain,
		Token:     d.Token,
		AmountUSD: d.AmountUSD,
		Status:    d.Status,
		TxHash:    d.TxHash,
		CreatedAt: d.CreatedAt,
		// UpdatedAt:  d.UpdatedAt,
		ExpireTime: d.ExpireTime,
		PaidAt:     d.PaidAt,
	}
}

// BatchPaymentOrderDomainToModel 批量轉換domain為model
func BatchPaymentOrderDomainToModel(ds []*order.PaymentOrder) (ms []*PaymentOrder) {
	ms = make([]*PaymentOrder, 0, len(ds))
	for _, d := range ds {
		ms = append(ms, PaymentOrderDomainToModel(d))
	}
	return
}

// PaymentOrderModelToDomain 轉換model為domain
func PaymentOrderModelToDomain(m *PaymentOrder) (d *order.PaymentOrder) {
	return &order.PaymentOrder{
		OrderID:    m.OrderID,
		Address:    m.Address,
		Chain:      m.Chain,
		Token:      m.Token,
		AmountUSD:  m.AmountUSD,
		Status:     m.Status,
		TxHash:     m.TxHash,
		ExpireTime: m.ExpireTime,
		PaidAt:     m.PaidAt,
		CreatedAt:  m.CreatedAt,
	}
}

// BatchPaymentOrderModelToDomain 批量轉換model為domain
func BatchPaymentOrderModelToDomain(ms []*PaymentOrder) (ds []*order.PaymentOrder) {
	ds = make([]*order.PaymentOrder, 0, len(ms))
	for _, m := range ms {
		ds = append(ds, PaymentOrderModelToDomain(m))
	}
	return
}
