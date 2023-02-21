package server

import (
	"fmt"
	"net"
	"sync"

	"github.com/iotames/easyim/config"
)

type Server struct {
	Ip              string
	Port, DropAfter int

	// //在线用户的列表
	// OnlineMap map[string]*User
	// 对User字典或字典中的user, 进行操作时，要加锁
	lock sync.RWMutex
	// //消息广播的channel
	// Message chan string
}

// 创建一个server的接口
func NewServer(conf config.Server) *Server {
	server := &Server{
		Ip:        conf.IP,
		Port:      conf.Port,
		DropAfter: conf.DropAfter,
		// OnlineMap: make(map[string]*User),
		// Message:   make(chan string),
	}
	return server
}

func (s *Server) Lock() {
	s.lock.Lock()
}

func (s *Server) Unlock() {
	s.lock.Unlock()
}

// // 监听Message广播消息channel的goroutine，一旦有消息就发送给全部的在线User
// func (s *Server) ListenMessager() {
// 	for {
// 		// fmt.Println("//将msg发送给全部的在线User")
// 		msg := <-s.Message
// 		//将msg发送给全部的在线User
// 		s.mapLock.Lock()
// 		for _, u := range s.OnlineMap {
// 			u.ReceiveData([]byte(msg))
// 		}
// 		s.mapLock.Unlock()
// 	}
// }

// // 广播消息的方法
// func (s *Server) BroadCast(user *User, msg string) {
// 	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
// 	s.Message <- sendMsg
// }

// 启动服务器的接口
func (s *Server) Start() {
	//socket listen
	fmt.Printf("[START] EasyIM Server. listenner at IP: %s, Port %d is starting\n", s.Ip, s.Port)
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	//close listen socket
	defer listener.Close()

	// //启动监听Message的goroutine
	// go s.ListenMessager()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}

		//do handler
		go Handler(s, conn)
	}
}

// func (s *Server)Stop(){
// 	fmt.Println("[STOP] EasyIM Server")
// }
