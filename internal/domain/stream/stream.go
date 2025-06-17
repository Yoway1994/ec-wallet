package stream

import "context"

type Stream interface {
	WatchAddress(ctx context.Context, req *WatchAddressRequest) error
}

const (
	WatchActionStart = "start_watch"

	DataStatusActive = "active"

	AddressStreamKey = "stream:address"
)

type WatchAddressRequest struct {
	Address string
	Chain   string
}
