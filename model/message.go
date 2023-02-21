package model

import (
	"fmt"
	"time"

	"github.com/iotames/easyim/database"
)

const MSG_TEXT = "text"
const MSG_IMAGE = "image"
const MSG_VOICE = "voice"
const MSG_VIDEO = "video"

type Message struct {
	ID, FromUser, ToUser, ToGroup, MsgType string
	CreateTime                             int64
	Content                                string
}

func NewMessage(content, msgType, fromUser, toUser, toGroup string) *Message {
	id := database.GetSnowflakeNode().Generate().Int64()
	return &Message{
		ID:       fmt.Sprintf("%d", id),
		FromUser: fromUser, ToUser: toUser, ToGroup: toGroup,
		MsgType:    msgType,
		CreateTime: time.Now().Unix(),
		Content:    content,
	}
}
