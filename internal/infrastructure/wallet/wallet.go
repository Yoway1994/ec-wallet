package walletservice

import (
	"ec-wallet/configs"
	"ec-wallet/internal/domain"
	"ec-wallet/internal/domain/wallet"
)

type walletService struct {
	config *configs.Config
	repo   domain.Repo
}

func NewWalletService(config *configs.Config, repo domain.Repo) wallet.Wallet {
	return &walletService{config: config, repo: repo}
}

var _ wallet.Wallet = (*walletService)(nil)
