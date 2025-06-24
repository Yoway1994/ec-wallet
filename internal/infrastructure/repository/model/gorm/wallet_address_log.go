package model

import (
	"ec-wallet/internal/domain/wallet"
	"time"
)

// WalletAddressLog 地址狀態變更日誌模型
type WalletAddressLog struct {
	ID           uint64     `gorm:"column:id;primaryKey"`
	AddressID    uint64     `gorm:"column:address_id;not null"`
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

// WalletAddressLogDomainToModel 轉換domain為model
func WalletAddressLogDomainToModel(d *wallet.WalletAddressLog) (m *WalletAddressLog) {
	return &WalletAddressLog{
		ID:           d.ID,
		AddressID:    d.AddressID,
		Operation:    d.Operation,
		StatusAfter:  d.StatusAfter,
		StatusBefore: d.StatusBefore,
		OperationAt:  d.OperationAt,
		ValidUntil:   d.ValidUntil,
		OrderID:      d.OrderID,
		UserID:       d.UserID,
	}
}

// BatchWalletAddressLogDomainToModel 批量轉換domain為model
func BatchWalletAddressLogDomainToModel(ds []*wallet.WalletAddressLog) (ms []*WalletAddressLog) {
	for _, d := range ds {
		ms = append(ms, WalletAddressLogDomainToModel(d))
	}
	return
}
