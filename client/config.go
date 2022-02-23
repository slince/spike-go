package client

import (
	"encoding/json"
	"github.com/slince/spike/pkg/auth"
	"github.com/slince/spike/pkg/log"
	"github.com/slince/spike/pkg/tunnel"
	"io/ioutil"
)

type Configuration struct {
	Host string `json:"host"`
	Port int
	Log log.Config
	Auth auth.GenericUser
	Tunnels []tunnel.Tunnel
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