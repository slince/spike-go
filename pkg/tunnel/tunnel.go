package tunnel

type Tunnel struct {
	Id         string `json:"id"`
	Protocol   string `json:"protocol"`
	ServerPort uint16 `json:"server_port"`
}

type RegisterResult struct {
	Tunnel
	Error string `json:"error"`
}
