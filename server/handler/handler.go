package handler

import (
	"fmt"

	"github.com/iotames/easyim/contract"
	"github.com/iotames/easyim/model"
)

// 用户处理消息的业务 Request
func MainHandler(u contract.IUser) error {
	// 通过命令行读取的消息data, 有换行符，转为字符串值为: string(data[:len(data)-1])
	data, err := u.GetConnData()
	if err != nil {
		return err
	}
	lendata := len(data)
	if lendata < 10 {
		return fmt.Errorf("req data too small")
	}
	if u.IsHttp(data) {
		// HTTP API 接口业务处理
		err = HttpHandler(model.NewRequest(data, u.GetConn()))
		if err != nil {
			fmt.Printf("---err-HttpHandler-(%v)----\n", err)
			return err
		}
		return u.Close()
	}
	//用户的任意消息，代表当前用户是一个活跃的
	u.KeepActive()

	// TODO IM即时通讯业务处理

	//提取用户的消息(去除'\n')
	// msg := string(data[:n-1])
	//用户针对msg进行消息处理

	//  len(msg) > 4 && msg[:3] == "to|" {
	// 		//消息格式:  to|张三|消息内容
	// 		remoteUser.SendMsg(u.Name + "对您说:" + content)

	//	} else {
	//		u.server.BroadCast(u, msg)
	//	}
	return nil
}
