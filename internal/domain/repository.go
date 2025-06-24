package domain

import (
	"context"
	"ec-wallet/internal/domain/order"
	"ec-wallet/internal/domain/wallet"

	"gorm.io/gorm"
)

type Repo interface {
	// TX
	Begin() *gorm.DB
	Commit(tx *gorm.DB) *gorm.DB
	RollBack(tx *gorm.DB) *gorm.DB
	// WalletAddressPool
	GetWalletAddressPool(ctx context.Context, tx *gorm.DB, params *wallet.QueryWalletAddressPoolsParams) (*wallet.WalletAddressPool, error)
	QueryWalletAddressPools(ctx context.Context, tx *gorm.DB, params *wallet.QueryWalletAddressPoolsParams) ([]*wallet.WalletAddressPool, error)
	CreateWalletAddressPools(ctx context.Context, tx *gorm.DB, pools []*wallet.WalletAddressPool) ([]uint64, error)
	UpdateWalletAddressPools(ctx context.Context, tx *gorm.DB, updates *wallet.UpdateWalletAddressPoolsParams) (int64, error)
	// WalletAddressLog
	CreateWalletAddressLogs(ctx context.Context, tx *gorm.DB, logs []*wallet.WalletAddressLog) ([]uint64, error)
	// PaymentOrder
	CreatePaymentOrders(ctx context.Context, tx *gorm.DB, orders []*order.PaymentOrder) ([]uint64, error)
}
