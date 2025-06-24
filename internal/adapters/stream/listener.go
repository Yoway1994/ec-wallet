package streamadapter

import (
	"context"
	"ec-wallet/internal/adapters/stream/handlers"
	"ec-wallet/internal/domain/stream"
	"ec-wallet/internal/wire"

	"go.uber.org/zap"
)

func StartListeningHandler() {
	ctx := context.Background()
	streamService, err := wire.NewStreamService()
	if err != nil {
		panic(err)
	}
	transferEventChan := make(chan stream.TransferEvent, 100)
	streamService.OnListening(ctx, transferEventChan)

	zapLogger := wire.NewLogger()

	for transferEvent := range transferEventChan {
		err := handlers.HandleTransferEvent(ctx, &transferEvent)
		if err != nil {
			zapLogger.Error("❌ HandleTransferEvent 錯誤", zap.Error(err))
		}
	}
}
