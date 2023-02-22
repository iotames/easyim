package user

import (
	"fmt"
	"io"
	"net"

	"github.com/iotames/easyim/contract"
)

type User struct {
	Name string
	Addr string
	// data    chan []byte
	Message        chan string
	IsActive       chan bool
	conn           net.Conn
	server         contract.IServer
	onConnectStart func(u User)
	onConnectLost  func(u User)
}

// 创建一个用户的API
func NewUser(conn net.Conn, s contract.IServer) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:     userAddr,
		Addr:     userAddr,
		Message:  make(chan string),
		IsActive: make(chan bool),
		conn:     conn,
		server:   s,
	}
	//启动监听当前user channel消息的goroutine
	go user.ListenMessage()

	return user
}

func (u User) GetActiveChannel() chan bool {
	return u.IsActive
}
func (u *User) activeSend() {
	u.IsActive <- true
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

func (u *User) Close() {
	// s.ListenMessage()方法可能强踢后还在写数据
	// s.mapLock.Lock()
	// delete(s.OnlineMap, user.Name)
	// s.mapLock.Unlock()

	// u.ReceiveDataToSend([]byte("连接长时间不活跃，连接已断开")) // 异步操作消息还没发出去，连接就断开了
	u.sendData([]byte("连接长时间不活跃，连接已断开"))
	//销毁用的资源
	close(u.Message)
	//关闭连接
	u.conn.Close()

}

// ReceiveDataToSend 接受消息，并通过channel发送给客户端。异步操作。支持并发。
// 当连接断开时，可能会继续发送异步消息。此时须使用同步锁
func (u *User) ReceiveDataToSend(d []byte) {
	// u.data <- d
	u.Message <- string(d)
}

// GetConn 获取TCP连接
func (u User) GetConn() net.Conn {
	return u.conn
}

// GetConnData 获取TCP客户端发送的数据
func (u User) GetConnData() (data []byte, err error) {
	// 最长接受4096长度的信息
	buf := make([]byte, 4096)

	// 如用户未发消息，则代码执行到conn.Read停止
	n, err := u.conn.Read(buf)
	// 主动或被动(网络不好，或长时间未发消息被踢)断开连接，则继续执行

	if n == 0 {
		// 客户端主动或意外断开连接
		fmt.Println("----connect lost-----")
		u.ConnectLost()
		err = fmt.Errorf("connect lost")
		return
	}

	if err != nil {
		if err == io.EOF {
			fmt.Println("----GetConnData--connect Read err(err == io.EOF):", err)
			err = nil
		} else {
			fmt.Println("---GetConnData--connect Read err(err != io.EOF):", err)
			err = fmt.Errorf("connect read err:%v", err)
			return
		}
	}

	// 如果是命令行输入TCP消息，会包含换行符 \n
	data = buf[:n]

	//用户的任意消息，代表当前用户是一个活跃的
	u.activeSend()
	return
}

// 监听当前User channel的 方法,一旦有消息，就直接发送给对端客户端
func (u *User) ListenMessage() {
	for {
		msg := <-u.Message
		u.sendData([]byte(msg))
	}
}

// sendData 发送数据给客户端。同步操作
func (u User) sendData(d []byte) {
	u.conn.Write(d)
}
