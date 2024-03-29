package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/fs"
	"io/ioutil"
	"os"
)


var content = `
host: 127.0.0.1
port: 6200
users:
  - username: admin
    password: admin
log:
  console: true
  level: trace
  file: "./spiked.log"
`

var force bool
var  initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a configuration file in the current directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		return generateConfig()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "Force to create even is it already exists")
}
func generateConfig() error{
	var curDir, _ = os.Getwd()
	var cfgFile = curDir + "/" + "spiked.yaml"
	var _, err = os.Stat(cfgFile)
	if !os.IsNotExist(err) && !force {
		return fmt.Errorf("config file \"%s\" is exists", cfgFile)
	}
	return ioutil.WriteFile(cfgFile, []byte(content), fs.FileMode(0777))
}