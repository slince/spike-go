package client

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/slince/spike-go/protol"
)

type MessageHandler interface {
	// Handle the message
	Handle(message *protol.Protocol) error
}

type Handler struct {
	client *Client
}

// 授权响应服务
type AuthResponseHandler struct {
	Handler
}

func (hd *AuthResponseHandler) Handle(message *protol.Protocol) error{
	code, ok := message.Headers["code"]
	if !ok || code != "200" {
		hd.client.Close() // 关闭客户端
		return errors.New("auth error")
	}
	client,_ := message.Body["client"]
	if cl,ok :=client.(map[string]string);ok {
		clientId,_ := cl["id"]
		hd.client.Id = clientId
		hd.registerTunnel() // 验证成功以后注册隧道
		return nil
	}
	return errors.New("bad message body")
}

// 注册所有的隧道到服务器
func (hd *AuthResponseHandler) registerTunnel() {
	message := &protol.Protocol{
		Action: "register_tunnel",
		Body: map[string]interface{}{
			"tunnels": hd.client.Tunnels,
		},
	}
	hd.client.ControlConn.Write(message.ToBytes())
}

// 注册隧道信息返回处理
type RegisterTunnelResponseHandler struct {
	Handler
}

func (hd *RegisterTunnelResponseHandler) Handle(message *protol.Protocol) error {
	if code,ok := message.Headers["code"]; !ok || code != "200" {
		tunnel, _ := message.Body["tunnel"]
		targetTunnel, _ := tunnel.(map[string]string)
		serverPort, _ := targetTunnel["server_port"]
		return fmt.Errorf(`the tunnle with serverport "%s" register tunnel error`, serverPort)
	}

	registeredTunnels, _ := message.Body["tunnels"]
	if regTunnels, ok2 := registeredTunnels.([]map[string]string); ok2 {
		for _, registeredTunnel := range regTunnels{
			targetTunnelId, _ := registeredTunnel["id"]
			for _, tunnel := range hd.client.Tunnels {
				if tunnel.Match(registeredTunnel) {
					tunnel.SetId(targetTunnelId)
				}
			}
		}
	}

	hd.client.Logger.Info("register tunnel ok")
	return nil
}

// 请求代理
type RequestProxyHandler struct {
	Handler
}

func (hd *RequestProxyHandler) Handle(message *protol.Protocol) error {
	tunnelId, ok := message.Headers["tunnel-id"]
	publicConnId, publicConnIdOk := message.Headers["public-connection-id"]
	if !ok || !publicConnIdOk {
		return errors.New("missing tunnel id or public connection id")
	}
	tunnel,err := hd.client.findTunnelById(tunnelId)
	if err != nil {
		return errors.New("bad tunnel")
	}

	worker := &TcpWorker{
		Client: hd.client,
		publicConnId: publicConnId,
		tunnel: tunnel,
	}
	err = worker.Start()
	if err != nil {
		hd.client.Logger.Error("fail to create worker for request_proxy message")
	}
	return nil
}

// handler工厂方法
type MessageHandlerFactory struct {
	client *Client
}

func (factory MessageHandlerFactory) newHandler() Handler{
	return Handler{
		client: factory.client,
	}
}

func (factory MessageHandlerFactory) NewAuthResponseHandler() MessageHandler{
	return &AuthResponseHandler{
		factory.newHandler(),
	}
}
func (factory MessageHandlerFactory) NewRegisterTunnelResponseHandler() MessageHandler{
	return &RegisterTunnelResponseHandler{
		factory.newHandler(),
	}
}

func (factory MessageHandlerFactory) NewRequestProxyHandler() MessageHandler{
	return &RequestProxyHandler{
		factory.newHandler(),
	}
}