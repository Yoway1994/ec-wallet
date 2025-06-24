package order

type PaymentToken struct {
	Symbol          string
	Chain           string
	ContractAddress *string
	Decimals        int64
	IsActive        bool
}
