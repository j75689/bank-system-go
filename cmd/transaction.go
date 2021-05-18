package cmd

import (
	transaction "bank-system-go/internal/app/transaction_service"
	"bank-system-go/pkg/util"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	transactionCmd = &cobra.Command{
		Use:           "transaction",
		Short:         "transaction service command",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	transactionMigrate = &cobra.Command{
		Use:           "migrate",
		Short:         "Migrate database",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(_ *cobra.Command, _ []string) {
			app, err := transaction.Initialize(cfgFile)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			util.Launch(app.Migrate, app.Stop, time.Duration(timeout)*time.Second)
		},
	}

	transactionService = &cobra.Command{
		Use:           "start",
		Short:         "Start transaction service",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(_ *cobra.Command, _ []string) {
			app, err := transaction.Initialize(cfgFile)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			util.Launch(app.Start, app.Stop, time.Duration(timeout)*time.Second)
		},
	}
)

func init() {
	transactionCmd.AddCommand(transactionMigrate, transactionService)
}
