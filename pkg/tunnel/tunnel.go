package tunnel

type Tunnel struct {
	Id         string `json:"id"`
	Protocol   string `json:"protocol"`
	LocalPort uint16 `json:"local_port"`
	ServerPort uint16 `json:"server_port"`
}

type RegisterResult struct {
	Tunnel
	Error string `json:"error"`
}
