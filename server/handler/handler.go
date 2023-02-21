package handler

import (
	"fmt"

	"github.com/iotames/easyim/contract"
)

// 用户处理消息的业务 Request
func Handler(u contract.IUser, data []byte) error {
	// 通过命令行读取的消息data, 有换行符，转为字符串值为: string(data[:len(data)-1])
	fmt.Printf("---Handler---DATA--(%s)\n", string(data))

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
