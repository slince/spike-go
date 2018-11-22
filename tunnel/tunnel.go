package tunnel

type Tunnel interface {
	//判断是否匹配指定信息
	Match(info map[string]string) bool
	//匹配两个tunnel是否相同
	MatchTunnel(tunnel Tunnel) bool
	// gets tunnel id
	GetId() string
	// Set id for the tunnel
	SetId(id string)
	// resolve the local address
	ResolveAddress() string
}

type TcpTunnel struct {
	Id string `json:"id"` //由服务端统一分配id
	LocalPort string `json:"local_port"`
	ServerPort string `json:"server_port"`
	Host string `json:"host"`
}

// Set id for tunnel
func (tn *TcpTunnel) SetId(id string) {
	tn.Id = id
}

// get id
func (tn *TcpTunnel) GetId() string {
	return tn.Id
}

// get host
func (tn *TcpTunnel) ResolveAddress() string {
	return tn.Host + ":" + tn.LocalPort
}


func (tn *TcpTunnel) Match(info map[string]string) bool {
	serverPort, _ := info["server_port"]
	localPort, _ := info["local_port"]
	return localPort == tn.LocalPort && serverPort == tn.ServerPort
}

func (tn *TcpTunnel) MatchTunnel(tunnel Tunnel) bool {
	if tunnel, ok := tunnel.(*TcpTunnel);ok {
		return tn.Match(map[string]string{
			"local_port": tunnel.LocalPort,
			"server_port": tunnel.ServerPort,
		})
	}
	return false
}

type HttpTunnel struct {
	TcpTunnel
	Domain string `json:"domain"`
}

func (tn *HttpTunnel) Match(info map[string]string) bool {
	serverPort, _ := info["server_port"]
	localPort, _ := info["local_port"]
	domain, _ := info["domain"]
	return localPort == tn.LocalPort && serverPort == tn.ServerPort && domain == tn.Domain
}

func (tn *HttpTunnel) MatchTunnel(tunnel Tunnel) bool {
	if tunnel, ok := tunnel.(*HttpTunnel);ok {
		return tn.Match(map[string]string{
			"local_port": tunnel.LocalPort,
			"server_port": tunnel.ServerPort,
			"domain": tunnel.Domain,
		})
	}
	return false
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
				"",
				localPort,
				serverPort,
			}
		} else {
			// http tunnel 必须绑定域名
			domain,domainOk := info["domain"]
			if domainOk {
				continue
			}
			tunnel = &HttpTunnel{
				TcpTunnel{
					"",
					localPort,
					serverPort,
				},
				domain,
			}
		}
		tunnels = append(tunnels, tunnel)
	}
	return tunnels
}

