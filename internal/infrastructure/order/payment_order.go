package orderservice

import (
	"context"
	"ec-wallet/internal/domain/order"
	"ec-wallet/internal/domain/stream"
	"ec-wallet/internal/errors"
)

func (s *orderService) CreatePaymentOrder(ctx context.Context, paymentOrder *order.PaymentOrder) error {
	tx := s.repo.Begin()
	defer tx.Rollback()

	_, err := s.repo.CreatePaymentOrders(ctx, tx, []*order.PaymentOrder{paymentOrder})
	if err != nil {
		return err
	}

	_ = tx.Commit()
	return nil
}

func (s *orderService) ValidateOrder(ctx context.Context, event *stream.TransferEvent) (bool, error) {
	tx := s.repo.Begin()
	defer tx.Rollback()
	// 取得order
	params := &order.QueryPaymentOrdersParams{
		Address: &event.To,
	}
	orders, err := s.repo.QueryPaymentOrders(ctx, tx, params)
	if err != nil {
		return false, err
	}

	if len(orders) != 1 {
		return false, errors.ErrOrderNotUnique
	}

	paymentOrder := orders[0]

	// 取得token
	token, ok := s.tokens[paymentOrder.Token]
	if !ok {
		return false, errors.ErrOrderTokenNotFound
	}

	// 檢查合約地址
	if token.ContractAddress != &event.Address {
		return false, errors.ErrOrderContractMismatch
	}

	// 檢查付款數量
	amount, err := event.Amount(token.Decimals)
	if err != nil {
		return false, err
	}
	if amount.LessThan(paymentOrder.AmountUSD) {
		return false, errors.ErrOrderInsufficientAmount
	}

	_ = tx.Commit()
	return true, nil
}
