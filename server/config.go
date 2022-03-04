package server

import (
	"github.com/slince/spike/pkg/auth"
	"github.com/slince/spike/pkg/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Configuration struct {
	Host string
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
	err = yaml.Unmarshal(read, &config)
	return
}