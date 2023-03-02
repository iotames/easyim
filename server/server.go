package server

import (
	"fmt"
	"net"
	"sync"

	"github.com/iotames/easyim/config"
	"github.com/iotames/easyim/contract"
	"github.com/iotames/easyim/model"
	"github.com/iotames/miniutils"
)

type Server struct {
	Ip              string
	Port, DropAfter int
	onLineMap       map[string]contract.IUser
	chatRoomsMap    map[string]*ChatRoom

	// //在线用户的列表
	// OnlineMap map[string]*User
	// 对User字典或字典中的user, 进行操作时，要加锁
	lock sync.RWMutex
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

func (s *Server) UserOffline(addr string) {
	s.Lock()
	if _, ok := s.onLineMap[addr]; ok {
		delete(s.onLineMap, addr)
		// TODO 判断从哪个聊天室移除
		fmt.Println("移除onLineMap")
	}
	s.Unlock()
}

func (s *Server) UserOnline(addr string, u contract.IUser) {
	s.Lock()
	if s.onLineMap == nil {
		s.onLineMap = make(map[string]contract.IUser, 10)
	}
	_, ok := s.onLineMap[addr]
	if !ok {
		s.onLineMap[addr] = u
		// TODO 判断加入哪个聊天室
	}
	s.Unlock()
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
// 			u.ReceiveDataToSend([]byte(msg))
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

func (s *Server) HandlerMsg(u contract.IUser, data []byte) error {
	logger := miniutils.GetLogger("")
	dp := model.GetDataPack()
	msg := model.Msg{}
	err := dp.Unpack(data, &msg)
	if err != nil {
		return fmt.Errorf("unpack msg fail:%v", err)
	}
	logger.Debug(fmt.Sprintf("---msg.ChatType(%d)--msg.MsgType(%d)-msg.Seq(%d)--msg.Status(%d)--ReceivedMsg(%v)-", msg.ChatType, msg.MsgType, msg.Seq, msg.Status, msg.String()))
	// 在线调试 http://www.websocket-test.com/, https://websocketking.com/
	msgCount := u.MsgCount()
	if (u.IsWebSocket() && msgCount == 2) || (!u.IsWebSocket() && msgCount == 1) {
		// 接收到的第一条消息
		err = s.firstMsgComeIn(u, data)
		if err != nil {
			return err
		}
		// TODO 发送消息到监听组件
		data, err = dp.Pack(&msg)
		u.ReceiveDataToSend(data)
		return err
	}
	// TODO 发送消息到监听组件
	data, err = dp.Pack(&msg)
	u.ReceiveDataToSend(data)
	return err

	// if msg.ChatType == model.Msg_SINGLE {
	// 	// 单聊。发送给TO_USER
	// 	data, err = dp.Pack(&msg)
	// 	return u.SendData(data)
	// }

	// if msg.ChatType == model.Msg_GROUP {
	// 	// 群聊。发送给群里的每一个成员。
	// 	data, err = dp.Pack(&msg)
	// 	return u.SendData(data)
	// }

	return fmt.Errorf("unknown ChatType")
}

func (s *Server) firstMsgComeIn(u contract.IUser, data []byte) error {
	// TODO 根据access_token进行用户身份鉴权, 再添加到聊天室
	addr := u.GetConn().RemoteAddr().String()
	s.UserOnline(addr, u)
	// TODO 发送消息到监听组件
	return nil
}

// func (s *Server)Stop(){
// 	fmt.Println("[STOP] EasyIM Server")
// }
