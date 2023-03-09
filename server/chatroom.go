package server

import (
	"github.com/iotames/easyim/contract"
	"github.com/iotames/easyim/model"
)

// 创建聊天的基本单位: 聊天室
type ChatRoom struct {
	ID          int64
	isGroupChat bool
	msgCount    int
	msg         chan []byte
	usersMap    map[string]contract.IUser
}

func NewChatRoom(u contract.IUser) *ChatRoom {
	usersMap := make(map[string]contract.IUser, 2)
	addr := u.GetConn().RemoteAddr().String()
	usersMap[addr] = u
	cr := &ChatRoom{usersMap: usersMap}
	go cr.ListenMessage()
	return cr
}

func (c *ChatRoom) Join(u contract.IUser) {
	addr := u.GetConn().RemoteAddr().String()
	c.usersMap[addr] = u
}

func (c *ChatRoom) Remove(addr string) {
	delete(c.usersMap, addr)
}

func (c *ChatRoom) ReceiveDataToSend(msg *model.Msg) error {
	// TODO 如果接收消息的用户离线，则保存离线消息。下次用户上线再发送
	c.msgCount += 1
	dp := model.GetDataPack()
	data, err := dp.Pack(msg)
	if err != nil {
		return err
	}
	c.msg <- data
	return nil
}

// 监听本房间是否有消息进来
func (c *ChatRoom) ListenMessage() {
	for {
		msg := <-c.msg
		for _, u := range c.usersMap {
			u.ReceiveDataToSend(msg)
		}
	}
}
