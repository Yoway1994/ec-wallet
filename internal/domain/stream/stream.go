package stream

import "context"

type Stream interface {
	OnListening(ctx context.Context, transferEventChan chan TransferEvent)
	WatchAddress(ctx context.Context, req *WatchAddressRequest) error
}

const (
	// action
	WatchActionStart = "start_watch"
	// status
	DataStatusActive = "active"
	// stream
	TransferStreamKey = "stream:transfer"
	AddressStreamKey  = "stream:address"
)

type WatchAddressRequest struct {
	Address string
	Chain   string
}

type TransferEvent struct {
	From     string
	To       string
	Value    string
	TxHash   string
	BlockNum uint64
}
