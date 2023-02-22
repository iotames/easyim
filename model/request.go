package model

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/http"
)

type Request struct {
	data, httpBody []byte
	httpRequest    *http.Request
	conn           net.Conn
}

func NewRequest(data []byte, conn net.Conn) *Request {
	return &Request{data: data, conn: conn}
}

func (r Request) ResponseBody(body []byte) error {
	data := NewOKResponse(body).Response()
	_, err := r.conn.Write(data)
	return err

}

func (r Request) GetData() []byte {
	return r.data
}

func (r *Request) ParseHttp() error {
	data := r.GetData()
	reader := bytes.NewReader(data)
	hreq, err := http.ReadRequest(bufio.NewReader(reader))
	if err != nil {
		return err
	}
	body, err := io.ReadAll(hreq.Body)
	if err != nil {
		return err
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
