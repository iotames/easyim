package server

import (
	"fmt"
	"net"
	"sync"

	"github.com/iotames/easyim/config"
	"github.com/iotames/easyim/contract"
	"github.com/iotames/easyim/model"
	"github.com/iotames/easyim/user"
	"github.com/iotames/miniutils"
)

type Server struct {
	Ip              string
	Port, DropAfter int
	addrToUser      map[string]*user.User
	uidToAddr       map[string]string
	chatRoomsMap    map[string]*ChatRoom
	addrToRoom      map[string]*ChatRoom

	// 操作字典时，要加锁
	lock sync.RWMutex
}

// 创建一个server的接口
func NewServer(conf config.Server) *Server {
	server := &Server{
		Ip:           conf.IP,
		Port:         conf.Port,
		DropAfter:    conf.DropAfter,
		addrToUser:   make(map[string]*user.User, 10),
		uidToAddr:    make(map[string]string, 10),
		chatRoomsMap: make(map[string]*ChatRoom, 10),
		addrToRoom:   make(map[string]*ChatRoom, 10),
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
	uu := s.addrToUser[addr]
	// 根据access_token进行用户身份鉴权
	b, err := s.checkToken(uu, u, &msg)
	if !b {
		return err
	}

	room, ok := s.addrToRoom[addr]
	if ok {
		return room.ReceiveDataToSend(&msg)
	}
	// TODO 对未上线的发送对象，保存离线消息，下次上线时发送
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

func (s *Server) checkToken(uu *user.User, u contract.IUser, msg *model.Msg) (bool, error) {
	if uu.GetID() != msg.FromUserId || !uu.CheckToken(msg.AccessToken) {
		errMsg := "access_token 或 from_user_id 不正确"
		msg.MsgType = model.Msg_NOTIFY
		msg.ToUserId = msg.FromUserId
		msg.FromUserId = "notify"
		msg.Content = errMsg
		return false, s.SendMsg(u, msg)
	}
	return true, nil
}

func (s *Server) UserOnline(u contract.IUser, msg *model.Msg) error {
	addr := u.GetConn().RemoteAddr().String()
	uu := user.NewUser(msg.FromUserId)
	b, err := s.checkToken(uu, u, msg)
	if !b {
		return err
	}
	s.Lock()
	_, ok := s.addrToUser[addr]
	if !ok {
		fmt.Println("UserOnline:", addr)
		s.addrToUser[addr] = uu
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
	// TODO 读取离线消息，有则发送
	return s.SendMsg(u, msg)
}

func (s *Server) UserOffline(addr string) {
	s.Lock()
	if _, ok := s.addrToUser[addr]; ok {
		delete(s.addrToUser, addr)
		uid := s.addrToUser[addr].GetID()
		delete(s.uidToAddr, uid)
		// 从聊天室移除
		room, b := s.addrToRoom[addr]
		if b {
			room.Remove(addr)
			delete(s.addrToRoom, addr)
		}
		fmt.Println("移除在线用户")
	}
	s.Unlock()
}

// func (s *Server)Stop(){
// 	fmt.Println("[STOP] EasyIM Server")
// }
