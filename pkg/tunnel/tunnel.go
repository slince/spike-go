package tunnel

type Tunnel struct {
	Id         string `yaml:"id"`
	Protocol   string `yaml:"protocol"`
	LocalHost string `yaml:"local_host"`
	LocalPort int `yaml:"local_port"`
	ServerPort int `yaml:"server_port"`
	Headers map[string]string
}

type RegisterResult struct {
	Tunnel
	Error string `json:"error"`
}
