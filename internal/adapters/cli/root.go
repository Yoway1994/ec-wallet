package cli

import "github.com/spf13/cobra"

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "ec-wallet-cli",
		Short: "EC Wallet CLI Tool",
		Long:  "Command line interface for managing EC Wallet operations",
	}

	// 添加子命令
	rootCmd.AddCommand(NewInitPoolCmd())

	return rootCmd
}
