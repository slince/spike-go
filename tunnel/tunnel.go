package tunnel

import "fmt"

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

// Tcp 隧道
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
	serverPort := info["server_port"]
	localPort := info["local_port"]
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

// Http 隧道
type HttpTunnel struct {
	TcpTunnel
	ProxyHosts map[string]string
}

func (tn *HttpTunnel) Match(info map[string]string) bool {
	return info["local_port"] == tn.LocalPort && info["server_port"] == tn.ServerPort
}

func (tn *HttpTunnel) MatchTunnel(tunnel Tunnel) bool {
	if tunnel, ok := tunnel.(*HttpTunnel);ok {
		return tn.Match(map[string]string{
			"local_port": tunnel.LocalPort,
			"server_port": tunnel.ServerPort,
		})
	}
	return false
}

// Create many tunnels.
func NewManyTunnels(details []map[string]interface{}) []Tunnel{
	var tunnel Tunnel
	tunnels := make([]Tunnel, len(details))
	for index, info := range details {
		switch info["protocol"] {
		case "tcp":
			tunnel = &TcpTunnel{
				LocalPort: info["local_port"].(string),
				ServerPort: info["server_port"].(string),
				Host: info["host"].(string),
			}
		case "http":
			tunnel = &HttpTunnel{
				TcpTunnel{
					LocalPort: info["local_port"].(string),
					ServerPort: info["server_port"].(string),
					Host: info["host"].(string),
				},
				info["proxy_hosts"].(map[string]string),
			}
		default:
			continue
		}
		fmt.Println(index)
		tunnels[index] = tunnel
	}
	return tunnels
}

