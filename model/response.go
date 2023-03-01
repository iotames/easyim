package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
)

type Response struct {
	statusCode     int
	err            error
	conn           net.Conn
	body, response []byte
}

func NewResponse(conn net.Conn) *Response {
	return &Response{conn: conn}
}

func (r Response) getHeader() []string {
	bodyLen := len(r.body)
	return []string{
		fmt.Sprintf("HTTP/1.1 %d %s", r.statusCode, http.StatusText(r.statusCode)),
		fmt.Sprintf("Content-Length: %d", bodyLen),
		"Connection: close", // Connection: keep-alive
		"Server: easyim",
		// ALLOW CORS START
		"Access-Control-Allow-Credentials: true",
		"Access-Control-Allow-Headers: Origin, Content-Length, Content-Type, Accept, Token, Auth-Token, X-Requested-With",
		"Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE, UPDATE",
		"Access-Control-Allow-Origin: *",
		"Access-Control-Expose-Headers: Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type",
		// ALLOW CORS END
		// "Keep-Alive: timeout=4",
		// "Date: Wed, 22 Feb 2023 09:58:51 GMT",
	}
	// https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials
}

func (r Response) httpJson() []byte {
	header := r.getHeader()
	header = append(header, "Content-Type: application/json; charset=utf-8")
	headerText := strings.Join(header, "\n")
	return bytes.Join([][]byte{[]byte(headerText), r.body}, []byte("\n\r\n"))
	// return []byte(headerText + "\n\r\n" + string(r.body))
}

func (r Response) Json(v interface{}) Response {
	var err error
	switch v.(type) {
	case []byte:
		r.body = v.([]byte)
	case string:
		r.body = []byte(v.(string))
	default:
		r.body, err = json.Marshal(v)
		if err != nil {
			r.err = err
		}
	}
	if r.err == nil {
		r.statusCode = http.StatusOK
		r.response = r.httpJson()
	}
	return r
}

func (r Response) OPTIONS() Response {
	return r.Html([]byte{}, 200)
}

func (r Response) Html(body []byte, statusCode int) Response {
	r.statusCode = http.StatusOK
	r.body = body
	header := r.getHeader()
	header = append(header, "Content-Type: text/html; charset=utf-8")
	headerText := strings.Join(header, "\n")
	r.response = []byte(headerText + "\n\r\n" + string(r.body))
	return r
}

func (r Response) Write() error {
	if r.err != nil {
		return r.err
	}
	_, err := r.conn.Write(r.response)
	return err
}

type JsonObject map[string]interface{}
type ResponseApiData struct {
	Code int        `json:"code"`
	Msg  string     `json:"msg"`
	Data JsonObject `json:"data"`
}

func ResponseApi(data JsonObject, msg string, code int) ResponseApiData {
	return ResponseApiData{Data: data, Msg: msg, Code: code}
}

func ResponseOk(msg string) ResponseApiData {
	return ResponseApi(JsonObject{}, msg, http.StatusOK)
}

func ResponseItems(items interface{}) ResponseApiData {
	return ResponseApi(JsonObject{"Items": items}, "success", http.StatusOK)
}

func ResponseFail(msg string, code int) ResponseApiData {
	return ResponseApi(JsonObject{}, msg, code)
}

func ResponseNotFound() ResponseApiData {
	return ResponseFail("NotFound.无法找到请求对象", http.StatusNotFound)
}

func ResponseUnauthorized() ResponseApiData {
	return ResponseFail("Unauthorized.您没有权限访问此页面", http.StatusUnauthorized)
}

func ResponseMethodNotAllowed() ResponseApiData {
	return ResponseFail("MethodNotAllowed.不允许的请求方法", http.StatusMethodNotAllowed)
}

func ResponseServerError() ResponseApiData {
	return ResponseFail("ServerError.服务器内部错误", http.StatusInternalServerError)
}

func ResponseQueryArgsError(msg string) ResponseApiData {
	return ResponseFail("QueryArgsError.请求参数错误:"+msg, http.StatusBadRequest)
}
