package model

import (
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

func NewOKResponse(body []byte) *Response {
	return NewResponse(http.StatusOK, body)
}

func (r Response) GetHeader() []string {
	bodyLen := len(r.body)
	return []string{
		fmt.Sprintf("HTTP/1.1 %d %s", r.statusCode, http.StatusText(r.statusCode)),
		// "Content-Type: text/html; charset=utf-8",
		"Content-Type: application/json",
		fmt.Sprintf("Content-Length: %d", bodyLen),
	}
}

func (r Response) Response() []byte {
	header := r.GetHeader()
	headerText := strings.Join(header, "\n")
	return []byte(headerText + "\n\r\n" + string(r.body))
}
