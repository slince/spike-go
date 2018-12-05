package server

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	configFile string // 配置文件所在位置
	savePath string // 配置文件保存位置
)

func initConfig() {
	if configFile == "" {
		configFile = "./spiked.json"
	}
	if savePath == "" {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			panic(err)
		}
		savePath = dir
	}
	savePath = strings.TrimRight(savePath, "\\/")
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

// 创建配置文件
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		saveFile := savePath + "/spiked.json"
		srcReader, _ := os.Open("./spiked.json")
		saveWriter, err := os.Create(saveFile)
		if err != nil {
			panic(err.Error())
		}
		_, err = io.Copy(srcReader, saveWriter)
		if err != nil {
			panic(err)
		}
		abSaveFile,_ := filepath.Abs(saveFile)
		fmt.Printf(`"%s" is created`, abSaveFile)
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is /.spiked.json)")
	initCmd.PersistentFlags().StringVar(&savePath, "dir", "", "Save path to storage config file")
	RootCmd.AddCommand(initCmd)
}