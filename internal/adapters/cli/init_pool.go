package cli

import (
	"context"
	"ec-wallet/internal/wire"

	"github.com/spf13/cobra"
)

func NewInitPoolCmd() *cobra.Command {
	ctx := context.Background()
	var chain string
	var count int
	var batchSize int

	cmd := &cobra.Command{
		Use:   "init-pool",
		Short: "Initialize wallet address pool",
		RunE: func(cmd *cobra.Command, args []string) error {
			walletService, err := wire.NewWalletService()
			if err != nil {
				return err
			}
			_, err = walletService.InitWalletAddressPools(ctx, chain, count, batchSize)
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&chain, "chain", "ETH", "Blockchain type (ETH, BTC, etc)")
	cmd.Flags().IntVar(&count, "count", 100, "Number of addresses to generate")
	cmd.Flags().IntVar(&batchSize, "batch-size", 20, "Batch size for address generation")

	return cmd
}
