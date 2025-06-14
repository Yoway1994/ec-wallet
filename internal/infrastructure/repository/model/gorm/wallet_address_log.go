package model

import (
	"time"
)

// WalletAddressLog 地址狀態變更日誌模型
type WalletAddressLog struct {
	ID           uint       `gorm:"column:id;primaryKey"`
	AddressID    uint       `gorm:"column:address_id;not null"`
	Operation    string     `gorm:"column:operation;not null"`
	StatusAfter  string     `gorm:"column:status_after;not null"`
	StatusBefore string     `gorm:"column:status_before"`
	OperationAt  time.Time  `gorm:"column:operation_at;not null;default:CURRENT_TIMESTAMP"`
	ValidUntil   *time.Time `gorm:"column:valid_until"`
	OrderID      string     `gorm:"column:order_id"`
	UserID       string     `gorm:"column:user_id"`
}

// TableName 指定資料表名稱
func (WalletAddressLog) TableName() string {
	return "wallet_address_logs"
}
