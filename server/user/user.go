package user

import (
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/iotames/easyim/contract"
	"github.com/iotames/easyim/model"
	"github.com/iotames/miniutils"
)

const ERR_CONNECT_LOST = "connect lost"

// User. 一个TCP连接。ClentSocket
type User struct {
	Name           string
	Addr           string
	protocol       string
	IsClosed       bool
	message        chan []byte
	msgCount       int
	isActive       chan bool
	conn           net.Conn
	server         contract.IServer
	onConnectStart func(u User)
	onConnectLost  func(u User)
}

// 创建一个用户的API
func NewUser(conn net.Conn, s contract.IServer) *User {
	userAddr := conn.RemoteAddr().String()
	u := &User{
		Name:     userAddr,
		Addr:     userAddr,
		message:  make(chan []byte),
		isActive: make(chan bool),
		conn:     conn,
		server:   s,
	}
	//启动监听当前user channel消息的goroutine
	go u.ListenMessage()
	return u
}

func (u User) GetActiveChannel() chan bool {
	return u.isActive
}
func (u *User) KeepActive() {
	u.isActive <- true
}

func (u *User) SetOnConnectLost(f func(u User)) {
	u.onConnectLost = f
}

func (u *User) SetOnConnectStart(f func(u User)) {
	u.onConnectStart = f
}

func (u User) ConnectStart() {
	if u.onConnectStart != nil {
		u.onConnectStart(u)
	}
}

func (u User) ConnectLost() {
	if u.onConnectLost != nil {
		u.onConnectLost(u)
	}
}

// IsHttp 判断本次TCP连接，是否为HTTP协议。
// 因本系统不支持HTTP 的 Keep-Alive。故一个HTTP协议的TCP连接，只能发送一次消息。
// 消息发送超过1次，一定不是HTTP协议。
func (u User) IsHttp(data []byte) bool {
	method := string(data[:7])
	httpMethods := []string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "UPDATE"}
	isHttp := false
	for _, m := range httpMethods {
		if strings.Contains(method, m) {
			isHttp = true
		}
	}
	return isHttp
}

func (u *User) SetProtocol(proto string) {
	u.protocol = proto
}

func (u User) IsWebSocket() bool {
	if u.protocol == model.PROTOCOL_WEBSOCKET {
		return true
	}
	return false
}

func (u User) MsgCount() int {
	return u.msgCount
}

func (u *User) Close() error {
	// s.ListenMessage()方法可能强踢后还在写数据
	// s.mapLock.Lock()
	// delete(s.OnlineMap, user.Name)
	// s.mapLock.Unlock()

	// u.ReceiveDataToSend([]byte("连接长时间不活跃，连接已断开")) // 异步操作消息还没发出去，连接就断开了
	// u.SendData([]byte("连接长时间不活跃，连接已断开")) // OK 给用户发送消息，同步操作
	//销毁用的资源
	close(u.message)
	//关闭连接
	u.IsClosed = true
	return u.conn.Close()
}

// ReceiveDataToSend 接受消息，并通过channel发送给客户端。异步操作。支持并发。
// 当连接断开时，可能会继续发送异步消息。此时须使用同步锁
func (u *User) ReceiveDataToSend(d []byte) {
	u.message <- d
}

// GetConn 获取TCP连接
func (u User) GetConn() net.Conn {
	return u.conn
}

// GetConnData 获取TCP客户端发送的数据
func (u *User) GetConnData() (data []byte, err error) {
	// 最长接受4096长度的信息
	buf := make([]byte, 4096)

	// 如用户未发消息，则代码执行到conn.Read停止
	n, err := u.conn.Read(buf)
	// 主动或被动(网络不好，或长时间未发消息被踢)断开连接，则继续执行

	logger := miniutils.GetLogger("")
	if n == 0 {
		// 客户端主动或意外断开连接
		u.ConnectLost()
		if !u.IsClosed {
			err = u.Close()
			if err != nil {
				logger.Error("--GetConnData--u.Close()--error:", err)
			}
		}
		u.IsClosed = true
		err = fmt.Errorf(ERR_CONNECT_LOST)
		return
	}

	if err != nil {
		if err == io.EOF {
			logger.Debug("----GetConnData--connect Read err(err == io.EOF):", err)
			err = nil
		} else {
			logger.Debug("---GetConnData--connect Read err(err != io.EOF):", err)
			err = fmt.Errorf("connect read err:%v", err)
			return
		}
	}

	// 如果是命令行输入TCP消息，会包含换行符 \n
	data = buf[:n]
	u.msgCount += 1
	return
}

// 监听当前User channel的 方法,一旦有消息，就直接发送给对端客户端
func (u *User) ListenMessage() {
	for {
		msg := <-u.message
		u.SendData(msg)
	}
}

// SendData 发送数据给客户端。同步操作
func (u User) SendData(d []byte) error {
	_, err := u.conn.Write(d)
	return err
}
