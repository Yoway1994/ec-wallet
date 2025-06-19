package gormRepoImpl

import (
	"context"
	gormrepo "ec-wallet/internal/domain/gorm_repo"
	"ec-wallet/internal/infrastructure/logger"
	model "ec-wallet/internal/infrastructure/repository/model/gorm"
	"time"

	"gorm.io/gorm"
)

func (repo *repository) QueryWalletAddressPools(ctx context.Context, tx *gorm.DB, params *gormrepo.QueryWalletAddressPoolsParams) ([]*gormrepo.WalletAddressPool, error) {
	zapLogger := logger.FromContext(ctx)
	if tx == nil {
		tx = repo.Db
	}

	var fetched []*model.WalletAddressPool
	tx = repo.applyQueryWalletAddressPoolsParams(tx, params)

	if err := tx.Find(&fetched).Error; err != nil {
		zapLogger.Error(err.Error())
		return nil, err
	}

	return model.BatchWalletAddressPoolModelToDomain(fetched), nil
}

func (repo *repository) GetWalletAddressPool(ctx context.Context, tx *gorm.DB, params *gormrepo.QueryWalletAddressPoolsParams) (*gormrepo.WalletAddressPool, error) {
	zapLogger := logger.FromContext(ctx)
	if tx == nil {
		tx = repo.Db
	}
	var fetched model.WalletAddressPool
	tx = repo.applyQueryWalletAddressPoolsParams(tx, params)
	if err := tx.Take(&fetched).Error; err != nil {
		zapLogger.Error(err.Error())
		return nil, err
	}
	return model.WalletAddressPoolModelToDomain(&fetched), nil
}

func (repo *repository) UpdateWalletAddressPools(ctx context.Context, tx *gorm.DB, updates *gormrepo.UpdateWalletAddressPoolsParams) (int64, error) {
	zapLogger := logger.FromContext(ctx)
	if tx == nil {
		tx = repo.Db
	}

	tx = repo.applyQueryWalletAddressPoolsParams(tx, &updates.Where)

	updateMap := make(map[string]interface{})

	if updates.CurrentStatus != nil {
		updateMap["current_status"] = *updates.CurrentStatus
	}

	if updates.ReservedUntil != nil {
		updateMap["reserved_until"] = *updates.ReservedUntil
	}

	updateMap["updated_at"] = time.Now()

	// 執行更新
	result := tx.Model(&model.WalletAddressPool{}).
		Updates(updateMap)
	err := result.Error
	if err != nil {
		zapLogger.Error(err.Error())
		return 0, err
	}

	return result.RowsAffected, nil
}

func (repo *repository) applyQueryWalletAddressPoolsParams(query *gorm.DB, params *gormrepo.QueryWalletAddressPoolsParams) *gorm.DB {
	if params == nil {
		return query
	}

	if params.ID != nil {
		query = query.Where("id = ?", params.ID)
	}

	if params.Address != nil {
		query = query.Where("address = ?", params.Address)
	}

	if params.Chain != nil {
		query = query.Where("chain = ?", params.Chain)
	}

	if params.CurrentStatus != nil {
		query = query.Where("current_status = ?", params.CurrentStatus)
	}

	return query
}

func (repo *repository) CreateWalletAddressPools(ctx context.Context, tx *gorm.DB, pools []*gormrepo.WalletAddressPool) ([]uint64, error) {
	zapLogger := logger.FromContext(ctx)
	if tx == nil {
		tx = repo.Db
	}

	if len(pools) == 0 {
		return nil, nil
	}

	modelPools := model.BatchWalletAddressPoolDomainToModel(pools)

	if err := tx.Create(&modelPools).Error; err != nil {
		zapLogger.Error(err.Error())
		return nil, err
	}

	// 收集生成的 ID
	ids := make([]uint64, len(modelPools))
	for i, p := range modelPools {
		ids[i] = p.ID
	}

	return ids, nil
}
