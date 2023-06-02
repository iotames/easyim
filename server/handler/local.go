package handler

import (
	"strings"

	"github.com/iotames/easyim/model"
)

var StopSignChan chan string = make(chan string)

func closeListener(req *model.Request, resp *model.Response) model.Response {
	remoteAddr := req.RemoteAddr().String()
	if strings.Contains(remoteAddr, "127.0.0.1") || strings.Contains(remoteAddr, "::1") {
		go func() {
			StopSignChan <- "stop"
		}()
		return resp.Json(model.ResponseOk("操作成功"))
	}
	return resp.Json(model.ResponseFail("仅限内网访问", 400))
}
