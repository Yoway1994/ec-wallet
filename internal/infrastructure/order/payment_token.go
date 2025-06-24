package orderservice

import (
	"ec-wallet/internal/domain/order"
	model "ec-wallet/internal/infrastructure/repository/model/gorm"

	"gorm.io/gorm"
)

func ProvideTokens(db *gorm.DB) map[string]*order.PaymentToken {
	tokens := make(map[string]*order.PaymentToken)
	var fetched []*model.PaymentToken
	err := db.Find(&fetched).Error
	if err != nil {
		panic(err)
	}

	for _, token := range fetched {
		tokens[token.Symbol] = &order.PaymentToken{
			Symbol:          token.Symbol,
			Chain:           token.Chain,
			ContractAddress: token.ContractAddress,
			Decimals:        token.Decimals,
			IsActive:        token.IsActive,
		}
	}
	return tokens
}
