package handler

import (
	"fmt"

	"github.com/iotames/easyim/contract"
	"github.com/iotames/easyim/model"
	"github.com/iotames/miniutils"
)

// 用户处理消息的业务 Request
func MainHandler(u contract.IUser) error {
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
		req := model.NewRequest(data, u.GetConn())
		err = req.ParseHttp()
		if err != nil {
			logger.Error(fmt.Sprintf("---ParseHttpError(%v)---", err))
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

	msg := model.Msg{}
	err = dp.Unpack(data, &msg)
	if err != nil {
		return fmt.Errorf("unpack msg fail:%v", err)
	}
	logger.Debug("-----ReceivedMsg(%v)--msg.ChatType(%d)--", msg.String(), msg.ChatType)
	// 在线调试 http://www.websocket-test.com/, https://websocketking.com/

	if (u.IsWebSocket() && msgCount == 2) || (!u.IsWebSocket() && msgCount == 1) {
		// 接收到的第一条消息
		// TODO 用户身份鉴权, 添加到聊天室

	}

	if msg.ChatType == model.Msg_SINGLE {
		// 单聊。发送给TO_USER
		data, err = dp.Pack(&msg)
		return u.SendData(data)
	}

	if msg.ChatType == model.Msg_GROUP {
		// 群聊。发送给群里的每一个成员。
		data, err = dp.Pack(&msg)
		return u.SendData(data)
	}

	return fmt.Errorf("unknown ChatType")

	//提取用户的消息(去除'\n')
	// msg := string(data[:n-1])
	//用户针对msg进行消息处理

	//  len(msg) > 4 && msg[:3] == "to|" {
	// 		//消息格式:  to|张三|消息内容
	// 		remoteUser.SendMsg(u.Name + "对您说:" + content)

	//	} else {
	//		u.server.BroadCast(u, msg)
	//	}

}
