package cmd

import (
	"fmt"
	"github.com/slince/spike/pkg/build"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print spiked version",
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func printVersion(){
	fmt.Printf("spiked version: %s; build time: %s; go version: %s\n", build.Version, build.BuildTime, build.GoVersion)
}