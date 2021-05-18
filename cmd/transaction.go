package cmd

import (
	transaction "bank-system-go/internal/app/transaction_service"
	"fmt"
	"os"

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
			err = app.Migrate()
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
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
			err = app.Start()
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}
)

func init() {
	transactionCmd.AddCommand(transactionMigrate, transactionService)
}
