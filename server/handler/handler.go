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
	lendata := len(data)
	if lendata < 10 {
		err = fmt.Errorf("req data too small")
		logger.Debug("---handler.MainHandler--error:", err)
		return err
	}
	if u.IsHttp(data) {
		// HTTP API 接口业务处理
		err = HttpHandler(model.NewRequest(data, u.GetConn()))
		if err != nil {
			logger.Debug("---handler.MainHandler--HttpHandler--error:", err)
			return err
		}
		return u.Close()
	}

	// TODO IM即时通讯业务处理
	dataStr := "Response:" + string(data)
	// 接收 FROM_USER 发送给TO_USER

	return u.SendData([]byte(dataStr))

	//提取用户的消息(去除'\n')
	// msg := string(data[:n-1])
	//用户针对msg进行消息处理

	//  len(msg) > 4 && msg[:3] == "to|" {
	// 		//消息格式:  to|张三|消息内容
	// 		remoteUser.SendMsg(u.Name + "对您说:" + content)

	//	} else {
	//		u.server.BroadCast(u, msg)
	//	}
	// return nil
}
