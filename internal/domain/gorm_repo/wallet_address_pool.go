package gormrepo

import "time"

type WalletAddressPool struct {
	ID            uint64
	Address       string
	Chain         string
	Path          string
	Index         int
	CurrentStatus string
	ReservedUntil *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type QueryWalletAddressPoolsParams struct {
	ID            *uint64
	Address       *string
	Chain         *string
	CurrentStatus *string
}

type UpdateWalletAddressPoolsParams struct {
	Where         QueryWalletAddressPoolsParams
	CurrentStatus *string
	ReservedUntil *time.Time
}
