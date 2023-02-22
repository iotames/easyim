package contract

import "net"

// 定义服务接口
type IUser interface {
	ReceiveDataToSend([]byte)
	GetConnData() ([]byte, error)
	GetConn() net.Conn
}
