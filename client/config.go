package client

import (
	"encoding/json"
	"io/ioutil"
)

type Configuration struct {
	ServerAddress string `json:"server-address"`
	Log map[string]string
	Auth map[string]string
	Tunnels []TunnelConfiguration
}

type TunnelConfiguration struct {
	Protocol string `json:"protocol"`
	ServerPort string `json:"serverPort"`
	LocalPort string `json:"LocalPort"`
	Host string `json:"host"`
	ProxyHosts map[string]string
}


// 从文件创建一个新的对象
func CreateConfigurationFromFile(file string) (*Configuration, error){
	jsonFile,err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	configuration := &Configuration{}
	err = json.Unmarshal(jsonFile, configuration)
	if err != nil {
		return nil, err
	}
	return configuration, nil
}