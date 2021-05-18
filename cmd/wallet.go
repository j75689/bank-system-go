package cmd

import (
	wallet "bank-system-go/internal/app/wallet_service"
	"bank-system-go/pkg/util"
	"fmt"
	"os"
	"time"

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
			util.Launch(app.Migrate, app.Stop, time.Duration(timeout)*time.Second)
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
			util.Launch(app.Start, app.Stop, time.Duration(timeout)*time.Second)
		},
	}
)

func init() {
	walletCmd.AddCommand(walletMigrate, walletService)
}
