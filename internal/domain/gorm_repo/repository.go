package gormrepo

import (
	"context"

	"gorm.io/gorm"
)

type Repo interface {
	//
	Begin() *gorm.DB
	Commit(tx *gorm.DB) *gorm.DB
	RollBack(tx *gorm.DB) *gorm.DB
	//
	QueryWalletAddressPools(ctx context.Context, tx *gorm.DB, params *QueryWalletAddressPoolsParams) ([]*WalletAddressPool, error)
	CreateWalletAddressPools(ctx context.Context, tx *gorm.DB, pools []*WalletAddressPool) ([]uint64, error)
}
