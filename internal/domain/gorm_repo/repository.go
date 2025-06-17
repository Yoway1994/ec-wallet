package gormrepo

import (
	"context"

	"gorm.io/gorm"
)

type Repo interface {
	// TX
	Begin() *gorm.DB
	Commit(tx *gorm.DB) *gorm.DB
	RollBack(tx *gorm.DB) *gorm.DB
	// WalletAddressPool
	GetWalletAddressPool(ctx context.Context, tx *gorm.DB, params *QueryWalletAddressPoolsParams) (*WalletAddressPool, error)
	QueryWalletAddressPools(ctx context.Context, tx *gorm.DB, params *QueryWalletAddressPoolsParams) ([]*WalletAddressPool, error)
	CreateWalletAddressPools(ctx context.Context, tx *gorm.DB, pools []*WalletAddressPool) ([]uint64, error)
	UpdateWalletAddressPools(ctx context.Context, tx *gorm.DB, updates *UpdateWalletAddressPoolsParams) (int64, error)
	// WalletAddressLog
	CreateWalletAddressLogs(ctx context.Context, tx *gorm.DB, logs []*WalletAddressLog) ([]uint64, error)
}
