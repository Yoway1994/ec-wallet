package domain

import "context"

type Repository interface {
	GetWallets(ctx context.Context, params *GetWalletsParams)
}

type GetWalletsParams struct {
}
