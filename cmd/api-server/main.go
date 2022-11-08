package main

import (
	"api-server/internal/busi"

	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// @title api server
// @version 1.0
// @description spacescope block explorer api server
// @termsOfService http://swagger.io/terms/

// @contact.name xueyouchen
// @contact.email xueyou@starboardventures.io

// @host block-explorer-api.spacescope.io
// @BasePath /api/v1

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api-server",
		Short: "as",
		Run: func(cmd *cobra.Command, args []string) {
			if err := entry(); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		},
	}

	cmd.PersistentFlags().StringVar(&busi.Flags.Config, "conf", "", "path of the configuration file")

	return cmd
}

func entry() error {
	busi.Start()
	return nil
}

func main() {
	if err := NewRootCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
