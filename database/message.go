package database

import "time"

const (
	CHAT_SINGLE        uint8 = iota // 发送单聊消息
	CHAT_GROUP                      // 发送群聊消息
	MSG_STATUS_FAIL    uint8 = iota // 消息未送出
	MSG_STATUS_SENT                 // 消息已送出
	MSG_STATUS_ARRIVED              // 消息已送达
	MSG_STATUS_READ                 // 消息已读取
	MSG_TYPE_TEXT      uint8 = iota // 文本消息类型
	MSG_TYPE_IMAGE
	MSG_TYPE_VOICE
	MSG_TYPE_VIDEO
)

// 一对一单聊时，其实只需要保证发出的时序与接收的时序一致，就基本能让用户感觉不到乱序了。
// 多对多的群聊情况下，保证同一群内的所有接收方消息时序一致，也就能让用户感觉不到乱序了，方法有两种，一种单点绝对时序，另一种实现消息id的序列化（也就是实现一种全局递增消息ID）。

type Message struct {
	BaseModel  `xorm:"extends"`
	Seq        int64     // 时序号。由客户端传入的时序参数得出。
	FromUserID int64     `xorm:"notnull default(0) 'from_user_id'"`       // 发送方
	ToUserID   int64     `xorm:"notnull default(0) 'to_user_id'"`         // 接收方
	ChatType   uint8     `xorm:"SMALLINT notnull default(0) 'chat_type'"` // 聊天类型. 单聊0群聊1.默认单聊
	MsgType    uint8     `xorm:"SMALLINT notnull default(0) 'msg_type'"`  // 消息类型。文本0，图片1，语音2，视频3
	Status     uint8     `xorm:"SMALLINT notnull default(0)"`             // 状态。未送出0，已送出1，已送达2，已读取3
	ArrivedAt  time.Time // 送达时间
	ReadAt     time.Time // 已读时间
	Content    string    `xorm:"TEXT notnull 'content'"` // 消息内容。数据库保存为json字符串
}

func NewMsgSingleText(text string, from, to, seq int64) *Message {
	return &Message{
		Seq:        seq,
		FromUserID: from,
		ToUserID:   to,
		ChatType:   CHAT_SINGLE,
		MsgType:    MSG_TYPE_TEXT,
		Status:     MSG_STATUS_FAIL,
		Content:    text,
	}
}

func (msg *Message) Send() {
	msg.Status = MSG_STATUS_SENT
}
