package handler

import (
	"fmt"

	"github.com/iotames/easyim/contract"
	"github.com/iotames/easyim/model"
	"github.com/iotames/miniutils"
)

// 用户处理消息的业务 Request
func MainHandler(s contract.IServer, u contract.IUser) error {
	// 通过命令行读取的消息data, 有换行符，转为字符串值为: string(data[:len(data)-1])
	logger := miniutils.GetLogger("")
	data, err := u.GetConnData()
	if err != nil {
		return err
	}

	// 数据过滤
	lendata := len(data)
	if lendata < 10 {
		err = fmt.Errorf("req data too small")
		logger.Debug("---handler.MainHandler--error:", err)
		return err
	}
	isHttp := u.IsHttp(data)
	msgCount := u.MsgCount()

	dp := model.GetDataPack()
	if isHttp && msgCount == 1 {
		// HTTP API 接口业务处理。不支持HTTP 的 Keep-Alive
		req := model.NewRequest(u.GetConn())
		req.SetRawData(data)
		err = req.ParseHttp()
		if err != nil {
			logger.Error(fmt.Sprintf("---ParseHttpError(%v)--RequestRAW(%v)---", err, string(data)))
			return err
		}
		if req.IsWebSocket() {
			// websocket 握手
			dp.SetProtocol(model.PROTOCOL_WEBSOCKET)
			u.SetProtocol(model.PROTOCOL_WEBSOCKET)
			return req.ResponseWebSocket()
		}
		err = HttpHandler(req)
		if err != nil {
			logger.Debug("---handler.MainHandler--HttpHandler--error:", err)
			return err
		}
		// HTTP 一次请求响应后，立即关闭连接。不支持HTTP 的 Keep-Alive
		return u.Close()
	}

	logger.Debug("---TCP---ReceivedMessage--SUCCESS-----u.MsgCount=", u.MsgCount())
	return s.HandlerMsg(u, data)
}
