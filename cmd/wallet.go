package cmd

import (
	wallet "bank-system-go/internal/app/wallet_service"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	walletCmd = &cobra.Command{
		Use:           "wallet",
		Short:         "Wallet service command",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	walletMigrate = &cobra.Command{
		Use:           "migrate",
		Short:         "Migrate database",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(_ *cobra.Command, _ []string) {
			app, err := wallet.Initialize(cfgFile)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			err = app.Migrate()
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}

	walletService = &cobra.Command{
		Use:           "start",
		Short:         "Start wallet service",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(_ *cobra.Command, _ []string) {
			app, err := wallet.Initialize(cfgFile)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			err = app.Start()
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}
)

func init() {
	walletCmd.AddCommand(walletMigrate, walletService)
}
