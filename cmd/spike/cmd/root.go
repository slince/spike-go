package cmd

import (
	"fmt"
	"github.com/slince/spike/client"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "spike",
	Short: "Spike is a fast reverse proxy that helps to expose local services to the internet",
	RunE: func(cmd *cobra.Command, args []string) error {
		return start()
	},
}

var (
	cfgFile string
	host string
	port int
	username string
	password string
)

func init(){
	//rootCmd.PersistentFlags().Parse()
	var curDir, _ = os.Getwd()
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", curDir + "/.spike.yaml" , "Config file (default is Current dir/.spike.yaml)")
	rootCmd.PersistentFlags().StringVarP(&host, "host", "H", "127.0.0.1", "Server host")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p",8808, "Server port")
	rootCmd.PersistentFlags().StringVarP(&username, "username", "u","admin", "User for login")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "P", "admin", "Password for the given user")
}

func start() error{
	var config, err = createConfig()
	if err != nil {
		return err
	}
	cli, err := client.NewClient(config)
	if err != nil {
		return err
	}
	return cli.Start()
}

func createConfig() (client.Configuration, error){
	var _, err = os.Stat(cfgFile)
	var config client.Configuration
	if err != nil {
	    config = client.Configuration{}
	} else {
		config, err = client.ConfigFromJsonFile(cfgFile)
	}
	if len(host) > 0 {
		config.Host = host
	}
	if port > 0 {
		config.Port = port
	}
	if len(username) > 0 {
		config.User.Username = username
	}
	if len(password) > 0 {
		config.User.Password = password
	}
	return config, err
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
