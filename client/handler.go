package client

import (
	"errors"
	"fmt"
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

func (hdl *AuthResponseHandler) Handle(message *protol.Protocol) error{
	if message.Headers["code"] != "200" {
		hdl.client.Close() // 关闭客户端
		return errors.New("auth error")
	}
	client := message.Body["client"].(map[string]interface{})
	hdl.client.Id = client["id"].(string)

	hdl.registerTunnel() // 验证成功以后注册隧道
	return nil
}

// 注册所有的隧道到服务器
func (hdl *AuthResponseHandler) registerTunnel() {
	message := &protol.Protocol{
		Action: "register_tunnel",
		Body: map[string]interface{}{
			"tunnels": hdl.client.Tunnels,
		},
	}
	hdl.client.sendMessage(message)
}

// 注册隧道信息返回处理
type RegisterTunnelResponseHandler struct {
	Handler
}

func (hdl *RegisterTunnelResponseHandler) Handle(message *protol.Protocol) error {
	if code,ok := message.Headers["code"]; !ok || code != "200" {
		targetTunnel := message.Body["tunnel"].(map[string]interface{})
		serverPort := targetTunnel["server_port"].(string)
		return fmt.Errorf(`register the tunnle with serverport "%s" error`, serverPort)
	}

	// 注册隧道id
	regTunnels := message.Body["tunnels"].([]interface{})
	for _, regTunnel := range regTunnels{
		regTunnelInfo := regTunnel.(map[string]interface{})
		info := make(map[string]string, len(regTunnelInfo))
		for key, val := range regTunnelInfo {
			info[key] = val.(string)
		}
		for _, tunnel := range hdl.client.Tunnels { //找到本地tunnel，修改tunnel id
			if tunnel.Match(info) {
				tunnel.SetId(info["id"])
			}
		}
	}

	hdl.client.Logger.Info("register tunnel ok")
	return nil
}

// 请求代理
type RequestProxyHandler struct {
	Handler
}

func (hdl *RequestProxyHandler) Handle(message *protol.Protocol) error {
	tunnelId, ok := message.Headers["tunnel-id"]
	publicConnId, publicConnIdOk := message.Headers["pub-conn-id"]
	if !ok || !publicConnIdOk {
		return errors.New("missing tunnel id or public connection id")
	}
	tunnel,err := hdl.client.findTunnelById(tunnelId)
	if err != nil {
		return errors.New("bad tunnel")
	}
	worker := newWorker(hdl.client, publicConnId, tunnel)
	go func() {
		err = worker.Start()
		if err != nil {
			hdl.client.Logger.Error("fail to create worker for request_proxy message")
		}
	}()
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