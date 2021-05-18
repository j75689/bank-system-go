package cmd

import (
	"fmt"
	"os"
	"time"

	"bank-system-go/internal/app/http"
	"bank-system-go/pkg/util"

	"github.com/spf13/cobra"
)

var (
	httpCmd = &cobra.Command{
		Use:           "http",
		Short:         "Start http server",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(_ *cobra.Command, _ []string) {
			app, err := http.Initialize(cfgFile)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			util.Launch(app.Start, app.Stop, time.Duration(timeout)*time.Second)
		},
	}
)
