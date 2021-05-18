package cmd

import (
	transation "bank-system-go/internal/app/transation_service"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	transationCmd = &cobra.Command{
		Use:           "transation",
		Short:         "Transation service command",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	transationMigrate = &cobra.Command{
		Use:           "migrate",
		Short:         "Migrate database",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(_ *cobra.Command, _ []string) {
			app, err := transation.Initialize(cfgFile)
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

	transationService = &cobra.Command{
		Use:           "start",
		Short:         "Start transation service",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(_ *cobra.Command, _ []string) {
			app, err := transation.Initialize(cfgFile)
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
	transationCmd.AddCommand(transationMigrate, transationService)
}
