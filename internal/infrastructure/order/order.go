package orderservice

import (
	"ec-wallet/internal/domain"
	"ec-wallet/internal/domain/order"
)

type orderService struct {
	repo   domain.Repo
	tokens map[string]*order.PaymentToken
}

func NewOrderService(repo domain.Repo, tokens map[string]*order.PaymentToken) order.Order {
	return &orderService{repo: repo, tokens: tokens}
}
