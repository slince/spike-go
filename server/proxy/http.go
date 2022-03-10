package proxy

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/slince/spike/pkg/conn"
	"github.com/slince/spike/pkg/log"
	"io"
	"net"
	"net/http"
)

type HttpHandler struct {
	TcpHandler
	headers map[string]string
}

func NewHttpHandler(logger *log.Logger, connPool *conn.Pool, localAddress string, headers map[string]string) *HttpHandler{
	var handler = &HttpHandler{
		TcpHandler: TcpHandler{
			logger: logger,
			proxyConnPool: connPool,
			localAddress: localAddress,
		},
		headers: headers,
	}
	handler.handleConnCallback = handler.handleConn
	return handler
}

func (h *HttpHandler) handleConn(pubConn net.Conn) {
	h.logger.Trace("Accept a public connection[http]:", pubConn.RemoteAddr().String())
	var proxyConn, err = h.proxyConnPool.Get()
	if err != nil {
		h.logger.Error("Failed to get proxy conn from client, error", err)
		pubConn.Close()
		return
	}
	//conn.Combine(proxyConn, pubConn)
	var requests = make(chan *http.Request, 10)

	go func(){
		var b = bufio.NewReader(pubConn)
		var err error
		for {
			var req *http.Request
			req, err = http.ReadRequest(b)
			if err != nil {
				break
			}
			h.modifyRequest(req)
			requests <- req
		}
		_ = pubConn.Close()
		_ = proxyConn.Close()
		close(requests)
		if err != nil {
			h.logger.Error("Failed to read http request from pub conn, error: ", err)
		}
	}()

	go func(){
		var b = bufio.NewReader(proxyConn)
		var err error
		var pubErr bool
		var curReq *http.Request
		for req := range requests{
			curReq = req
			err = req.Write(proxyConn)
			if err != nil {
				break
			}
			var res *http.Response
			res, err = http.ReadResponse(b, req)
			if err != nil {
				break
			}
			res.Header.Set("x-spike-relay", "true")
			err = res.Write(pubConn)
			if err != nil {
				pubErr = true
				break
			}
		}
		// create an error page for the pub conn.
		if err != nil &&!pubErr {
			_ = createErrorResp(curReq, err).Write(pubConn)
		}
		_ = pubConn.Close()
		_ = proxyConn.Close()
		if err != nil {
			h.logger.Error("Failed to transfer http request, error: ", err)
		}
	}()
}

func (h *HttpHandler) modifyRequest(req *http.Request){
	req.Host = h.localAddress //attempt modify request "Host"
	for key, value := range h.headers {
		req.Header.Set(key, value)
	}
}

func createErrorResp(req *http.Request, err error) *http.Response{
	var res = &http.Response{
		Request: req,
		Header: http.Header(map[string][]string{}),
	}
	var errorMsg = fmt.Sprintf(`
<p style="color: red;">Spike Proxy Error: %s</p>
`, err.Error())
	var body =  bytes.NewBufferString(errorMsg)
	res.Body = io.NopCloser(body)
	res.Header.Set("x-spike-relay", "true")
	return res
}