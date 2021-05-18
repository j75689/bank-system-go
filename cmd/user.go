package cmd

import (
	user "bank-system-go/internal/app/user_service"
	"bank-system-go/pkg/util"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	userCmd = &cobra.Command{
		Use:           "user",
		Short:         "User service command",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	userMigrate = &cobra.Command{
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
			util.Launch(app.Migrate, app.Stop, time.Duration(timeout)*time.Second)
		},
	}

	userService = &cobra.Command{
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
			util.Launch(app.Start, app.Stop, time.Duration(timeout)*time.Second)
		},
	}
)

func init() {
	userCmd.AddCommand(userMigrate, userService)
}
