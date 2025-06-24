package model

import (
	"time"
)

// PaymentToken represents the payment_tokens table in the database
type PaymentToken struct {
	ID              uint64     `gorm:"primaryKey"`
	Symbol          string     `gorm:"column:symbol;type:varchar(64);not null"`
	Chain           string     `gorm:"column:chain;type:varchar(64);not null"`
	ContractAddress *string    `gorm:"column:contract_address;type:varchar(1024)"`
	Decimals        int64      `gorm:"column:decimals;type:integer;not null"`
	IsActive        bool       `gorm:"column:is_active;type:boolean;not null;default:true"`
	CreatedAt       time.Time  `gorm:"column:created_at;type:timestamp with time zone;autoCreateTime;not null"`
	UpdatedAt       *time.Time `gorm:"column:updated_at;type:timestamp with time zone"`
}

// TableName overrides the table name
func (p *PaymentToken) TableName() string {
	return "payment_tokens"
}
