package contract

import (
	"net"
)

// 定义服务接口
type IUser interface {
	ReceiveDataToSend([]byte)
	SendData([]byte) error
	GetConnData() ([]byte, error)
	GetConn() net.Conn
	IsHttp([]byte) bool
	IsWebSocket() bool
	SetProtocol(proto string)
	MsgCount() int
	Close() error
}
