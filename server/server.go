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
	addrToUid       map[string]string
	uidToAddr       map[string]string
	chatRoomsMap    map[string]*ChatRoom
	addrToRoom      map[string]*ChatRoom

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

func (s *Server) Lock() {
	s.lock.Lock()
}

func (s *Server) Unlock() {
	s.lock.Unlock()
}

// 启动服务器的接口
func (s *Server) Start() error {
	//socket listen
	fmt.Printf("[START] EasyIM Server. listenner at IP: %s, Port %d is starting\n", s.Ip, s.Port)
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		return fmt.Errorf("net.Listen err(%v)", err)
	}
	//close listen socket
	defer listener.Close()

	logger := miniutils.GetLogger("")
	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			logger.Error("listener accept err:", err)
			continue
		}

		//do handler
		go Handler(s, conn)
	}
}

func (s *Server) SendMsg(u contract.IUser, msg *model.Msg) error {
	dp := model.GetDataPack()
	data, err := dp.Pack(msg)
	if err != nil {
		return err
	}
	u.ReceiveDataToSend(data)
	return err
}

func (s *Server) HandlerMsg(u contract.IUser, data []byte) error {
	logger := miniutils.GetLogger("")
	dp := model.GetDataPack()
	msg := model.Msg{}
	err := dp.Unpack(data, &msg)
	if err != nil {
		logger.Error(fmt.Sprintf("---unpack msg fail(%v)---dataRAW(%s)--", err, data))
		return fmt.Errorf("unpack msg fail:%v", err)
	}
	logger.Debug(fmt.Sprintf("---msg.ChatType(%d)--msg.MsgType(%d)-msg.Seq(%d)--msg.Status(%d)--ReceivedMsg(%v)-", msg.ChatType, msg.MsgType, msg.Seq, msg.Status, msg.String()))
	// 在线调试 http://www.websocket-test.com/, https://websocketking.com/
	msgCount := u.MsgCount()
	if (u.IsWebSocket() && msgCount == 2) || (!u.IsWebSocket() && msgCount == 1) {
		// 连接建立后客户端主动发送一个心跳事件消息
		return s.FirstMsg(u, data, &msg)
	}
	addr := u.GetConn().RemoteAddr().String()
	room, ok := s.addrToRoom[addr]
	if ok {
		return room.ReceiveDataToSend(&msg)
	}
	// TODO 保存离线消息，下次上线时发送
	return nil
}

func (s *Server) FirstMsg(u contract.IUser, data []byte, msg *model.Msg) error {
	// 建立连接后发送的第一条消息必须为心跳事件消息
	if msg.MsgType != model.Msg_EVENT {
		errMsg := "建立连接后发送的第一条消息必须为心跳事件消息"
		msgFromToUserId := msg.ToUserId
		msg.MsgType = model.Msg_NOTIFY
		msg.ToUserId = msg.FromUserId
		msg.FromUserId = "notify"
		if msgFromToUserId != model.MSG_KEEP_ALIVE {
			msg.Content = fmt.Sprintf("%s(to_user_id=%s)", errMsg, model.MSG_KEEP_ALIVE)
			return s.SendMsg(u, msg)
		}
		msg.Content = errMsg
		return s.SendMsg(u, msg)
	}
	// 处理首次心跳消息，上线用户
	// TODO 根据access_token进行用户身份鉴权, 再添加到聊天室
	return s.UserOnline(u, msg)
}

func (s *Server) getChatRoom(from, to string, chatType model.Msg_ChatType) (roomKey string, room *ChatRoom) {
	if chatType == model.Msg_GROUP {
		roomKey = to
		room, _ = s.chatRoomsMap[roomKey]
		return
	}
	if chatType == model.Msg_SINGLE {
		var ok bool
		roomKey = from + to
		room, ok = s.chatRoomsMap[roomKey]
		if ok {
			return
		}
		roomKey = to + from
		room, ok = s.chatRoomsMap[roomKey]
		return
	}
	return
}

func (s *Server) UserOnline(u contract.IUser, msg *model.Msg) error {
	addr := u.GetConn().RemoteAddr().String()
	s.Lock()
	if s.onLineMap == nil {
		s.onLineMap = make(map[string]contract.IUser, 10)
	}
	_, ok := s.onLineMap[addr]
	if !ok {
		s.onLineMap[addr] = u
		s.addrToUid[addr] = msg.FromUserId
		s.uidToAddr[msg.FromUserId] = addr
		roomKey, room := s.getChatRoom(msg.FromUserId, msg.ToUserId, msg.ChatType)
		// 加入聊天室
		if room == nil {
			room = NewChatRoom(u)
			s.chatRoomsMap[roomKey] = room
		} else {
			room.Join(u)
		}
		_, ok := s.addrToRoom[addr]
		if !ok {
			s.addrToRoom[addr] = room
		}
	}
	s.Unlock()
	msg.Content = "SUCCESS"
	msg.ToUserId = msg.FromUserId
	msg.FromUserId = model.MSG_KEEP_ALIVE
	return s.SendMsg(u, msg)
}

func (s *Server) UserOffline(addr string) {
	s.Lock()
	if _, ok := s.onLineMap[addr]; ok {
		delete(s.onLineMap, addr)
		uid := s.addrToUid[addr]
		delete(s.addrToUid, addr)
		delete(s.uidToAddr, uid)
		// 从聊天室移除
		room, b := s.addrToRoom[addr]
		if b {
			room.Remove(addr)
			delete(s.addrToRoom, addr)
		}
		fmt.Println("移除onLineMap")
	}
	s.Unlock()
}

// func (s *Server)Stop(){
// 	fmt.Println("[STOP] EasyIM Server")
// }
