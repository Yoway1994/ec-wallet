package gormRepoImpl

import (
	"context"
	"ec-wallet/internal/domain/order"
	"ec-wallet/internal/infrastructure/logger"
	model "ec-wallet/internal/infrastructure/repository/model/gorm"

	"gorm.io/gorm"
)

func (repo *repository) CreatePaymentOrders(ctx context.Context, tx *gorm.DB, orders []*order.PaymentOrder) ([]uint64, error) {
	zapLogger := logger.FromContext(ctx)
	if tx == nil {
		tx = repo.Db
	}

	if len(orders) == 0 {
		return nil, nil
	}

	modelOrders := model.BatchPaymentOrderDomainToModel(orders)

	if err := tx.Create(&modelOrders).Error; err != nil {
		zapLogger.Error(err.Error())
		return nil, err
	}

	// 收集生成的 ID
	ids := make([]uint64, len(modelOrders))
	for i, order := range modelOrders {
		ids[i] = order.ID
	}

	return ids, nil
}
