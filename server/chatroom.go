package server

import (
	"github.com/iotames/easyim/model"
)

// 创建聊天的基本单位: 聊天室
type ChatRoom struct {
	ID       int64
	server   *Server
	msgCount int
	msg      chan []byte
	usersMap map[string]bool
}

func NewChatRoom(addr string, server *Server) *ChatRoom {
	usersMap := make(map[string]bool)
	usersMap[addr] = true
	cr := &ChatRoom{server: server, usersMap: usersMap}
	go cr.ListenMessage()
	return cr
}

func (c *ChatRoom) Join(addr string) {
	c.usersMap[addr] = true
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
	s := c.server
	for {
		msg := <-c.msg
		for k, _ := range c.usersMap {
			u, ok := s.onLineMap[k]
			if ok {
				u.ReceiveDataToSend(msg)
			}
		}
	}
}
