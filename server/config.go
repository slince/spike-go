package server

import (
	"encoding/json"
	"github.com/slince/spike/pkg/auth"
	"github.com/slince/spike/pkg/log"
	"io/ioutil"
)

type Configuration struct {
	Host string `json:"host"`
	Port int
	Log log.Config
	Users []auth.GenericUser
}

func ConfigFromJsonFile(file string) (config Configuration, err error){
	var read []byte
	read,err = ioutil.ReadFile(file)
	if err != nil {
		return
	}
	err = json.Unmarshal(read, &config)
	return
}