package orderservice

import (
	"ec-wallet/internal/domain"
	"ec-wallet/internal/domain/order"
)

type orderService struct {
	repo domain.Repo
}

func NewOrderService(repo domain.Repo) order.Order {
	return &orderService{repo: repo}
}
