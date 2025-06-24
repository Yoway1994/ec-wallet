package orderservice

import (
	"context"
	"ec-wallet/internal/domain/order"
	"ec-wallet/internal/domain/stream"
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

func (s *orderService) ValidateOrder(ctx context.Context, event *stream.TransferEvent) error {
	return nil
}
