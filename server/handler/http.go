package handler

import (
	"encoding/json"
	"fmt"

	"github.com/iotames/easyim/model"
)

func HttpHandler(req *model.Request) error {
	req.ParseHttp()
	hreq := req.GetHttpRequest()
	body := req.GetHttpBody()
	fmt.Printf("\n--method(%s)--Header(%+v)--Body(%s)-\n", hreq.Method, hreq.Header, string(body))
	data, _ := json.Marshal(ResponseOk("hello this is response Json"))
	return req.ResponseBody(data)
}
