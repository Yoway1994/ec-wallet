package gormrepo

import "time"

type WalletAddressPool struct {
	ID            uint64
	Address       string
	Chain         string
	Path          string
	Index         int
	CurrentStatus string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type QueryWalletAddressPoolsParams struct {
	Address       *string
	Chain         *string
	CurrentStatus *string
}

const (
	AddressStatusAvailable   = "AVAILABLE"
	AddressStatusReserved    = "RESERVED"
	AddressStatusBlacklisted = "BLACKLISTED"
)
