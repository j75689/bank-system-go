package cmd

import (
	user "bank-system-go/internal/app/user_service"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	UserCmd = &cobra.Command{
		Use:           "user",
		Short:         "User service command",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	UserMigrate = &cobra.Command{
		Use:           "migrate",
		Short:         "Migrate database",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(_ *cobra.Command, _ []string) {
			app, err := user.Initialize(cfgFile)
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

	UserService = &cobra.Command{
		Use:           "start",
		Short:         "Start user service",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(_ *cobra.Command, _ []string) {
			app, err := user.Initialize(cfgFile)
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
	UserCmd.AddCommand(UserMigrate, UserService)
}
