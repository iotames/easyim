package handler

import (
	"fmt"

	"github.com/iotames/easyim/model"
)

func HttpHandler(req *model.Request) error {
	hreq := req.GetHttpRequest()
	body := req.GetHttpBody()
	fmt.Printf("\n--method(%s)--proto(%s)--Header(%+v)--Body(%s)-\n", hreq.Method, hreq.Proto, hreq.Header, string(body))
	return req.ResponseJson(ResponseOk("hello response Json From struct")) //
}
