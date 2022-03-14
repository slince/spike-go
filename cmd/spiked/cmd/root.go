package cmd

import (
	"fmt"
	"github.com/slince/spike/pkg/auth"
	"github.com/slince/spike/pkg/log"
	"github.com/slince/spike/server"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "spiked",
	Short: "Spike is a fast reverse proxy that helps to expose local services to the internet",
	Long: `
 _____   _____   _   _   _    _____   _____  
/  ___/ |  _  \ | | | | / /  | ____| |  _  \ 
| |___  | |_| | | | | |/ /   | |__   | | | | 
\___  \ |  ___/ | | | |\ \   |  __|  | | | | 
 ___| | | |     | | | | \ \  | |___  | |_| | 
/_____/ |_|     |_| |_|  \_\ |_____| |_____/ 
`,
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", curDir + "/spiked.yaml" , "Config file (default is Current dir/spiked.yaml)")
	rootCmd.PersistentFlags().StringVarP(&host, "host", "H", "127.0.0.1", "Bind host")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "P",8808, "Bind port")
	rootCmd.PersistentFlags().StringVarP(&username, "username", "u","", "User")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "Password for the given user")
}

func start() error{
	var config, err = createConfig()
	if err != nil {
		return err
	}
	ser, err := server.NewServer(config)
	if err != nil {
		return err
	}
	return ser.Start()
}

func createConfig() (server.Configuration,error) {
	var _, err = os.Stat(cfgFile)
	var config server.Configuration
	if err != nil {
		config = server.Configuration{
			Users: make([]auth.GenericUser, 0),
			Log:   log.DefaultConfig,
		}
		if len(host) > 0 {
			config.Host = host
		}
		if port > 0 {
			config.Port = port
		}
		if len(username) > 0 {
			config.Users = append(config.Users, auth.GenericUser{
				Username: username, Password: password,
			})
		}
		return config, nil
	}
	return server.ConfigFromJsonFile(cfgFile)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
