package streamadapter

import (
	"context"
	"ec-wallet/internal/domain/stream"
	"ec-wallet/internal/wire"
	"fmt"
)

func StartListeningHandler() {
	ctx := context.Background()
	streamService, err := wire.NewStreamService()
	if err != nil {
		panic(err)
	}
	transferEventChan := make(chan stream.TransferEvent, 100)
	streamService.OnListening(ctx, transferEventChan)

	for transferEvent := range transferEventChan {
		fmt.Println(transferEvent)
	}
}
