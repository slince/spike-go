package cmd

import (
	"fmt"
	"github.com/slince/spike/client"
	"github.com/slince/spike/pkg/log"
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", curDir + "/spike.yaml" , "Config file (default is Current dir/spike.yaml)")
	rootCmd.PersistentFlags().StringVarP(&host, "host", "H", "127.0.0.1", "Server host")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "P",8808, "Server port")
	rootCmd.PersistentFlags().StringVarP(&username, "username", "u","admin", "User for login")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "admin", "Password for the given user")
}

func start() error{
	var cli, err = createClient()
	if err != nil {
		return err
	}
	return cli.Listen()
}

func createClient() (*client.Client, error){
	var config, err = createConfig()
	if err != nil {
		return nil, err
	}
	return client.NewClient(config)
}

func createConfig() (client.Configuration, error) {
	var _, err = os.Stat(cfgFile)
	var config client.Configuration
	if err != nil {
		config = client.Configuration{
			Log: log.DefaultConfig,
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
		return config, nil
	}
	return client.ConfigFromJsonFile(cfgFile)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
