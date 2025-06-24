package wallet

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

const (
	AddressLogOperationPayment = "PAYMENT"
	//
	AddressStatusAvailable   = "AVAILABLE"
	AddressStatusReserved    = "RESERVED"
	AddressStatusBlacklisted = "BLACKLISTED"
)

type Wallet interface {
	InitWalletAddressPools(ctx context.Context, chain string, count, batchSize int) ([]uint64, error)
	AcquireAddress(ctx context.Context, opts ...AcquireOption) (*AddressReservation, error)
	ReleaseAddress(ctx context.Context, address string) error
}

// KeyPair represents a derived public-private key pair
type KeyPair struct {
	PrivateKey []byte
	PublicKey  []byte
	Address    string
}

type AddressReservation struct {
	Address       string    // 錢包地址
	ReservedAt    time.Time // 保留時間
	ExpiresAt     time.Time // 過期時間
	ReservationID string    // 唯一預訂標識符
	CoinType      uint32    // 幣種類型
}

// NewAddressReservationParams 包含創建地址預約所需的所有參數
type NewAddressReservationParams struct {
	Address    string
	AddressID  uint64
	ReservedAt time.Time
	ExpiresAt  time.Time
	CoinType   uint32
	// 可以考慮添加其他元數據，如訂單ID等
}

// NewAddressReservation 創建新的地址預約
func NewAddressReservation(params *NewAddressReservationParams) *AddressReservation {
	// 生成唯一預約ID
	reservationID := fmt.Sprintf("RES-%05d-%03d-%s", params.AddressID, params.CoinType, params.ReservedAt.Format("20060102150405"))
	return &AddressReservation{
		Address:       params.Address,
		ReservedAt:    params.ReservedAt,
		ExpiresAt:     params.ExpiresAt,
		ReservationID: reservationID,
		CoinType:      params.CoinType,
	}
}

// AcquireOption 定義修改選項的函數類型
type AcquireOption func(*AcquireOptions)

type AcquireOptions struct {
	CoinType       uint32
	ExpiresIn      time.Duration
	OrderID        string
	UserID         string
	AmountRequired decimal.Decimal
	MerchantID     string
	AddressType    string
	CallbackURL    string
}

func WithOrderID(orderID string) AcquireOption {
	return func(o *AcquireOptions) {
		o.OrderID = orderID
	}
}
