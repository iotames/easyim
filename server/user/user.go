package user

import (
	"fmt"
	"io"
	"net"

	"github.com/iotames/easyim/contract"
	"github.com/iotames/easyim/server/handler"
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

	u.SendText("你被踢了")
	//销毁用的资源
	close(u.Message)
	//关闭连接
	u.conn.Close()

}

func (u *User) ReceiveData(d []byte) {
	// u.data <- d
	u.Message <- string(d)
}

// 给当前User对应的客户端发送消息
func (u *User) SendText(msg string) {
	u.SendData([]byte(msg))
}

// 用户处理消息的业务
func (u *User) MsgHandler(conn net.Conn) error {
	fmt.Println("---Begin---MsgHandler---", conn.RemoteAddr().String())

	// 最长接受4096长度的信息
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if n == 0 {
		fmt.Println("----connect lost-----")
		u.ConnectLost()
		return fmt.Errorf("connect lost")
	}

	if err != nil && err != io.EOF {
		fmt.Println("Conn Read err:", err)
		return fmt.Errorf("connect read err:%v", err)
	}

	//用户的任意消息，代表当前用户是一个活跃的
	u.IsActive <- true
	handler.Handler(u, buf[:n])

	return nil
}

// 监听当前User channel的 方法,一旦有消息，就直接发送给对端客户端
func (u *User) ListenMessage() {
	for {
		msg := <-u.Message
		u.SendData([]byte(msg + "\n"))
	}
}

func (u *User) SendData(d []byte) {
	u.conn.Write(d)
}
