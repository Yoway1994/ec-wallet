package gormrepo

import "time"

func NewWalletAddressLog(params *NewWalletAddressLogParams) *WalletAddressLog {
	log := &WalletAddressLog{
		AddressID:    params.AddressID,
		Operation:    params.Operation,
		StatusAfter:  params.StatusAfter,
		StatusBefore: params.StatusBefore,
		OrderID:      params.OrderID,
		UserID:       params.UserID,
	}

	if log.ValidUntil == nil {
		log.ValidUntil = log.CalculateExpiryTime(params.ExpiresIn)
	}
	return log
}

type NewWalletAddressLogParams struct {
	AddressID    uint64
	Operation    string
	StatusAfter  string
	StatusBefore string
	OperationAt  time.Time
	ValidUntil   *time.Time
	ExpiresIn    time.Duration
	OrderID      string
	UserID       string
}

type WalletAddressLog struct {
	ID           uint64
	AddressID    uint64
	Operation    string
	StatusAfter  string
	StatusBefore string
	OperationAt  time.Time
	ValidUntil   *time.Time
	OrderID      string
	UserID       string
}

func (log *WalletAddressLog) CalculateExpiryTime(duration time.Duration) *time.Time {
	expiryTime := log.OperationAt.Add(duration)
	return &expiryTime
}
