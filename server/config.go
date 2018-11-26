package server

import (
	"encoding/json"
	"io/ioutil"
)

// 服务端配置
type Configuration struct {
	Address string `json:"address"`
	Log map[string]string `json:"log"`
	Auth map[string]string `json:"auth"`
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