package handlers

import (
	"context"
	"ec-wallet/internal/domain/stream"
	"ec-wallet/internal/wire"
)

func HandleTransferEvent(ctx context.Context, transferEvent *stream.TransferEvent) error {
	walletService, err := wire.NewWalletService()
	if err != nil {
		return err
	}

	orderService, err := wire.NewOrderService()
	if err != nil {
		return err
	}

	// 檢查訂單
	isValid, err := orderService.ValidateOrder(ctx, transferEvent)
	if !isValid || err != nil {
		return err
	}

	// 釋放地址資源
	err = walletService.ReleaseAddress(ctx, transferEvent.To)
	if err != nil {
		return err
	}

	// 通知前端
	zapLogger := wire.NewLogger()
	zapLogger.Info("成功到帳")
	return nil
}
