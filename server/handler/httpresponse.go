package handler

import (
	"fmt"
)

type JsonObject map[string]interface{}

func Response(data interface{}, msg string, code int) interface{} {
	return struct {
		Code int
		Msg  string
		Data interface{}
	}{Msg: msg, Code: code, Data: data}
}

func ResponseOk(msg string) interface{} {
	return struct {
		Code int
		Msg  string
		Data JsonObject
	}{Msg: msg, Code: 200, Data: JsonObject{}}
}

func ResponseFail(msg string, code int) interface{} {
	return struct {
		Code int
		Msg  string
		Data JsonObject
	}{Msg: msg, Code: code, Data: JsonObject{}}
}

func ResponseItems(items interface{}) interface{} {
	switch items.(type) {
	case string:
		return fmt.Sprintf(`{"Code":200,"Msg":"success","Data":{"Items":%s}}`, items)
	default:
		return struct {
			Code int
			Msg  string
			Data JsonObject
		}{Msg: "success", Code: 200, Data: JsonObject{
			"Items": items,
		}}
	}
}
