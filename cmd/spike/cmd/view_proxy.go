package cmd

import "github.com/spf13/cobra"

var vpCmd = &cobra.Command{
	Use:   "show-proxy",
	Short: "Show proxy of the server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return showProxy()
	},
}

func init()  {
	rootCmd.AddCommand(vpCmd)
}

func showProxy() error{

}