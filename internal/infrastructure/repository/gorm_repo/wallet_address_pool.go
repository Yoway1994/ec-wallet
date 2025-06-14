package gormRepoImpl

import (
	"context"
	gormrepo "ec-wallet/internal/domain/gorm_repo"
	model "ec-wallet/internal/infrastructure/repository/model/gorm"

	"gorm.io/gorm"
)

func (repo *repository) QueryWalletAddressPools(ctx context.Context, tx *gorm.DB, params *gormrepo.QueryWalletAddressPoolsParams) ([]*gormrepo.WalletAddressPool, error) {
	if tx == nil {
		tx = repo.Db
	}

	var fetched []*model.WalletAddressPool
	tx = repo.applyQueryWalletAddressPoolsParams(tx, params)

	if err := tx.Find(&fetched).Error; err != nil {
		return nil, err
	}

	return model.BatchWalletAddressPoolModelToDomain(fetched), nil
}

func (repo *repository) applyQueryWalletAddressPoolsParams(query *gorm.DB, params *gormrepo.QueryWalletAddressPoolsParams) *gorm.DB {
	if params == nil {
		return query
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
	if tx == nil {
		tx = repo.Db
	}

	if len(pools) == 0 {
		return nil, nil
	}

	modelPools := model.BatchWalletAddressPoolDomainToModel(pools)

	if err := tx.Create(&modelPools).Error; err != nil {
		return nil, err
	}

	// 收集生成的 ID
	ids := make([]uint64, len(modelPools))
	for i, p := range modelPools {
		ids[i] = p.ID
	}

	return ids, nil
}
