package tunnel

type Tunnel interface {
	//判断是否匹配指定信息
	Match(info map[string]string) bool
}

type TcpTunnel struct {
	LocalPort string
	ServerPort string
}

func (tn *TcpTunnel) Match(info map[string]string) bool {
	serverPort, ok := info["serverPort"]
	return ok && serverPort == tn.ServerPort
}

type HttpTunnel struct {
	TcpTunnel
	Domain string
}

func (tn *HttpTunnel) Match(info map[string]string) bool {
	serverPort, ok := info["serverPort"]
	portMatch := ok && serverPort == tn.ServerPort

	domain, ok := info["domain"]

	return portMatch && ok && domain == tn.Domain
}

// Create many tunnels.
func NewManyTunnels(tunnelInfos []map[string]string) []Tunnel{

	var tunnel Tunnel

	tunnels := make([]Tunnel, 5)
	for _,info := range tunnelInfos {
		tp,_ := info["type"]
		localPort,_ := info["local_port"]
		serverPort,_ := info["server_port"]
		if tp == "tcp" {
			tunnel = &TcpTunnel{
				localPort,
				serverPort,
			}
		} else {
			tunnel = &TcpTunnel{
				localPort,
				serverPort,
			}
		}
		tunnels = append(tunnels, tunnel)
	}
	return tunnels
}

