package model

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

const (
	WebsocketGUID      = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	PROTOCOL_WEBSOCKET = "websocket"
)

type Request struct {
	raw, httpBody []byte
	httpRequest   *http.Request
	conn          net.Conn
}

func NewRequest(conn net.Conn) *Request {
	return &Request{conn: conn}
}

func (r Request) ResponseWebSocket() error {
	reqHeader := r.GetHttpRequest().Header
	accept, err := r.getWebSocketNonceAccept([]byte(reqHeader["Sec-Websocket-Key"][0])) // Sec-WebSocket-Key to Sec-Websocket-Key
	if err != nil {
		return err
	}
	h := []string{
		fmt.Sprintf("HTTP/1.1 %d %s", http.StatusSwitchingProtocols, http.StatusText(http.StatusSwitchingProtocols)),
		"Connection: Upgrade",
		"Upgrade: websocket",
		fmt.Sprintf("Sec-WebSocket-Accept: %s", string(accept)),
	}
	htext := strings.Join(h, "\n") + "\n\r\n"
	data := []byte(htext)
	_, err = r.conn.Write(data)
	return err
}

func (r Request) GetRawData() []byte {
	return r.raw
}
func (r *Request) SetRawData(dt []byte) {
	r.raw = dt
}

func (r Request) GetConn() net.Conn {
	return r.conn
}

func (r Request) GetHttpBodyToJson(v interface{}) error {
	return json.Unmarshal(r.GetHttpBody(), v)
}

func (r *Request) ParseHttp() error {
	var breader *bufio.Reader
	data := r.GetRawData()
	if len(data) > 0 {
		breader = bufio.NewReader(bytes.NewReader(data))
	} else {
		breader = bufio.NewReader(r.GetConn())
	}
	hreq, err := http.ReadRequest(breader)
	if err != nil {
		return fmt.Errorf("err for http.ReadRequest:%v", err)
	}

	var bodyReader io.Reader = hreq.Body
	body, err := io.ReadAll(bodyReader)
	if err != nil {
		return fmt.Errorf("err for io.ReadAll:%v", err)
	}
	r.httpRequest = hreq
	r.httpBody = body
	return nil
}

func (r Request) GetHttpBody() []byte {
	return r.httpBody
}

func (r Request) GetHttpRequest() *http.Request {
	return r.httpRequest
}

func (r Request) RemoteAddr() net.Addr {
	return r.conn.RemoteAddr()
}

func (r Request) LocalAddr() net.Addr {
	return r.conn.LocalAddr()
}

func (r Request) IsWebSocket() bool {
	req := r.httpRequest
	if req.Method != "GET" {
		return false
	}
	h := r.httpRequest.Header
	connVal := h.Get("Connection")
	// if strings.ToUpper(connVal) == "UPGRADE" {
	if connVal == "Upgrade" {
		if h.Get(connVal) == "websocket" {
			return true
		}
	}
	return false
}

// getWebSocketNonceAccept computes the base64-encoded SHA-1 of the concatenation of
// the nonce ("Sec-WebSocket-Key" value) with the websocket GUID string.
func (r Request) getWebSocketNonceAccept(nonce []byte) (expected []byte, err error) {
	h := sha1.New()
	if _, err = h.Write(nonce); err != nil {
		return
	}
	if _, err = h.Write([]byte(WebsocketGUID)); err != nil {
		return
	}
	expected = make([]byte, 28)
	base64.StdEncoding.Encode(expected, h.Sum(nil))
	return
}
