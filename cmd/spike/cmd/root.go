package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "spike",
	Short: "Spike is a fast reverse proxy that helps to expose local services to the internet",
	RunE: func(cmd *cobra.Command, args []string) error {
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
