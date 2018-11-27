package server

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var (
	configFile string // 配置文件所在位置
)

func initConfig() {
	if configFile == "" {
		configFile = "./spiked.json"
	}
}

var RootCmd = &cobra.Command{
	Use:   "spiked",
	Short: "A fast reverse proxy written in golang that helps to expose local services to the internet",
	Run: func(cmd *cobra.Command, args []string) {
		var ser *Server
		cfg,err := CreateConfigurationFromFile(configFile)
		if err != nil {
			panic(err)
		}
		ser = NewServer(cfg)
		ser.Run()
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
		}
		fmt.Println(dir)
		fmt.Println(args)
		//srcFile, _ := os.Open("./spiked.json")
		//filepath.Dir()
		//toFile,_ := os.Create("")
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is /.spiked.json)")
	RootCmd.AddCommand(initCmd)
}