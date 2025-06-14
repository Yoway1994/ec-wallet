package wallet

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type Wallet interface {
	InitWalletAddressPools(ctx context.Context, chain string, count, batchSize int) ([]uint64, error)
	AcquireAddress(ctx context.Context, opts ...AcquireOption) (*AddressReservation, error)
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

// AcquireOption 定義修改選項的函數類型
type AcquireOption func(*AcquireOptions)

type AcquireOptions struct {
	CoinType       uint32
	ExpiresIn      time.Duration
	OrderID        string
	CustomerID     string
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
