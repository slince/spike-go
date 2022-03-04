package client

import (
	"github.com/slince/spike/pkg/auth"
	"github.com/slince/spike/pkg/log"
	"github.com/slince/spike/pkg/tunnel"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Configuration struct {
	Host string
	Port int
	Log     log.Config
	User    auth.GenericUser
	Tunnels []tunnel.Tunnel
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