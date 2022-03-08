package cmd

import (
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/symfony-cli/terminal"
	"strconv"
)

var vpCmd = &cobra.Command{
	Use:   "view-proxy",
	Short: "Show proxy of the server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return showProxy()
	},
}

func init()  {
	rootCmd.AddCommand(vpCmd)
}

func showProxy() error{
	var cli, err = createClient()
	if err != nil {
		return err
	}
	proxies, err := cli.GetProxies()
	if err != nil {
		return err
	}
	table := tablewriter.NewWriter(terminal.Stdout)
	table.SetAutoFormatHeaders(false)
	table.SetHeader([]string{terminal.Format("<header>Protocol</>"), terminal.Format("<header>Server Port</>"), terminal.Format("<header>Client Id</>")})

	for _, proxy := range proxies {
		table.Append([]string{
			proxy.Protocol,
			strconv.Itoa(proxy.ServerPort),
			"client id",
		})
	}
	table.Render()

	return nil
}