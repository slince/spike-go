package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/fs"
	"io/ioutil"
	"os"
)


var content = `
host: 127.0.0.1 # server host
port: 6200
user: 
     username: admin
     password: admin

log:
     console: true  # enable console output
     level: trace  # trace debug info warn error
     file: ./spike.log # generate log file

tunnels:
  - protocol: tcp
    local_port: 3306
    server_port: 6201

  - protocol: udp
    local_host: 8.8.8.8
    local_port: 53
    server_port: 6202

  - protocol: http
    local_port: 80
    server_port: 6203
    headers:
      x-spike: yes
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
	var cfgFile = curDir + "/" + "spike.yaml"
	var _, err = os.Stat(cfgFile)
	if !os.IsNotExist(err) && !force {
		return fmt.Errorf("config file \"%s\" is exists", cfgFile)
	}
	return ioutil.WriteFile(cfgFile, []byte(content), fs.FileMode(0777))
}