package stream

import (
	"context"

	"github.com/shopspring/decimal"
)

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
	Address  string
	From     string
	To       string
	Value    string
	TxHash   string
	BlockNum uint64
}

func (e TransferEvent) Amount(decimals int64) (decimal.Decimal, error) {
	a, err := decimal.NewFromString(e.Value)
	if err != nil {
		return decimal.Zero, err
	}
	d := decimal.NewFromInt(10).Pow(decimal.NewFromInt(decimals))
	return a.Div(d), nil
}
