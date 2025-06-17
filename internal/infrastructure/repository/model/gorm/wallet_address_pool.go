package model

import (
	gormrepo "ec-wallet/internal/domain/gorm_repo"
	"time"
)

// WalletAddressPool 錢包地址資源池模型
type WalletAddressPool struct {
	ID            uint64     `gorm:"column:id;primaryKey"`
	Address       string     `gorm:"column:address;not null"`
	Chain         string     `gorm:"column:chain;not null"`
	Path          string     `gorm:"column:path;not null"`
	Index         int        `gorm:"column:index;not null"`
	CurrentStatus string     `gorm:"column:current_status;not null;default:AVAILABLE"`
	ReservedUntil *time.Time `gorm:"column:reserved_until"`
	CreatedAt     time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP"`
}

// TableName 指定資料表名稱
func (WalletAddressPool) TableName() string {
	return "wallet_address_pools"
}

func BatchWalletAddressPoolModelToDomain(ms []*WalletAddressPool) (ds []*gormrepo.WalletAddressPool) {
	for _, m := range ms {
		ds = append(ds, WalletAddressPoolModelToDomain(m))
	}
	return
}
func BatchWalletAddressPoolDomainToModel(ds []*gormrepo.WalletAddressPool) (ms []*WalletAddressPool) {
	for _, d := range ds {
		ms = append(ms, WalletAddressPoolDomainToModel(d))
	}
	return
}

func WalletAddressPoolModelToDomain(m *WalletAddressPool) (d *gormrepo.WalletAddressPool) {
	return &gormrepo.WalletAddressPool{
		ID:            m.ID,
		Address:       m.Address,
		Chain:         m.Chain,
		Path:          m.Path,
		Index:         m.Index,
		CurrentStatus: m.CurrentStatus,
		ReservedUntil: m.ReservedUntil,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func WalletAddressPoolDomainToModel(d *gormrepo.WalletAddressPool) (m *WalletAddressPool) {
	return &WalletAddressPool{
		ID:            d.ID,
		Address:       d.Address,
		Chain:         d.Chain,
		Path:          d.Path,
		Index:         d.Index,
		CurrentStatus: d.CurrentStatus,
		ReservedUntil: d.ReservedUntil,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
	}
}
