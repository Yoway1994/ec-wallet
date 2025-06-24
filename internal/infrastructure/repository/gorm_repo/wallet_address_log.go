package gormRepoImpl

import (
	"context"
	"ec-wallet/internal/domain/wallet"
	"ec-wallet/internal/infrastructure/logger"
	model "ec-wallet/internal/infrastructure/repository/model/gorm"

	"gorm.io/gorm"
)

func (repo *repository) CreateWalletAddressLogs(ctx context.Context, tx *gorm.DB, logs []*wallet.WalletAddressLog) ([]uint64, error) {
	zapLogger := logger.FromContext(ctx)
	if tx == nil {
		tx = repo.Db
	}

	if len(logs) == 0 {
		return nil, nil
	}

	modelLogs := model.BatchWalletAddressLogDomainToModel(logs)

	if err := tx.Create(&modelLogs).Error; err != nil {
		zapLogger.Error(err.Error())
		return nil, err
	}

	// 收集生成的 ID
	ids := make([]uint64, len(modelLogs))
	for i, log := range modelLogs {
		ids[i] = log.ID
	}

	return ids, nil
}
