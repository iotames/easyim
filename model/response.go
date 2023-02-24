package model

import (
	"bytes"

	"fmt"
	"net/http"
	"strings"
)

type Response struct {
	statusCode int
	body       []byte
}

func NewResponse(statusCode int, body []byte) *Response {
	return &Response{statusCode: statusCode, body: body}
}

func NewResponseOK(body []byte) *Response {
	return NewResponse(http.StatusOK, body)
}

func (r Response) getHeader() []string {
	bodyLen := len(r.body)
	return []string{
		fmt.Sprintf("HTTP/1.1 %d %s", r.statusCode, http.StatusText(r.statusCode)),
		fmt.Sprintf("Content-Length: %d", bodyLen),
		"Connection: close", // Connection: keep-alive
		"Server: easyim",
		// // ALLOW CORS START
		// "Access-Control-Allow-Credentials: true",
		// "Access-Control-Allow-Headers: Origin, Content-Length, Content-Type, Accept, Token, Auth-Token, X-Requested-With",
		// "Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE, UPDATE",
		// "Access-Control-Allow-Origin: *",
		// "Access-Control-Expose-Headers: Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type",
		// // ALLOW CORS END
		// "Keep-Alive: timeout=4",
		// "Date: Wed, 22 Feb 2023 09:58:51 GMT",
	}
	// https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials
}

func (r Response) HttpJson() []byte {
	header := r.getHeader()
	header = append(header, "Content-Type: application/json; charset=utf-8")
	headerText := strings.Join(header, "\n")
	return bytes.Join([][]byte{[]byte(headerText), r.body}, []byte("\n\r\n"))
	// return []byte(headerText + "\n\r\n" + string(r.body))
}

func (r Response) HttpHtml() []byte {
	header := r.getHeader()
	header = append(header, "Content-Type: text/html; charset=utf-8")
	headerText := strings.Join(header, "\n")
	return []byte(headerText + "\n\r\n" + string(r.body))
}
