package walletservice

import (
	"ec-wallet/configs"
	gormrepo "ec-wallet/internal/domain/gorm_repo"
	"ec-wallet/internal/domain/wallet"
)

type walletService struct {
	config *configs.Config
	repo   gormrepo.Repo
}

func NewWalletService(config *configs.Config, repo gormrepo.Repo) wallet.Wallet {
	return &walletService{config: config, repo: repo}
}

var _ wallet.Wallet = (*walletService)(nil)
