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

func (repo *repository) QueryPaymentOrders(ctx context.Context, tx *gorm.DB, params *order.QueryPaymentOrdersParams) ([]*order.PaymentOrder, error) {
	zapLogger := logger.FromContext(ctx)
	if tx == nil {
		tx = repo.Db
	}

	tx = repo.applyQueryPaymentOrdersParams(tx, params)

	var fetched []*model.PaymentOrder
	if err := tx.Find(&fetched).Error; err != nil {
		zapLogger.Error(err.Error())
		return nil, err
	}

	return model.BatchPaymentOrderModelToDomain(fetched), nil
}

func (repo *repository) applyQueryPaymentOrdersParams(query *gorm.DB, params *order.QueryPaymentOrdersParams) *gorm.DB {
	if params == nil {
		return query
	}

	if params.ID != nil {
		query = query.Where("id = ?", params.ID)
	}

	if params.Address != nil {
		query = query.Where("address = ?", params.Address)
	}

	return query
}
