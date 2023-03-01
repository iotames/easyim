package model

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/iotames/easyim/config"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type DataPack struct {
	protocol  string
	msgFormat string
}

var (
	once sync.Once
	dp   *DataPack
)

func GetDataPack() *DataPack {
	if dp != nil {
		return dp
	}
	once.Do(func() {
		dp = NewDataPack(config.GetServer().MsgFormat)
	})
	return dp
}

func NewDataPack(msgFormat string) *DataPack {
	return &DataPack{msgFormat: msgFormat}
}

func (dp *DataPack) SetProtocol(protocol string) {
	dp.protocol = protocol
}

func (dp DataPack) Pack(data protoreflect.ProtoMessage) (result []byte, err error) {
	if dp.msgFormat == config.MSG_FORMAT_JSON {
		result, err = json.Marshal(data)
	}
	if dp.msgFormat == config.MSG_FORMAT_PROTOBUF {
		result, err = proto.Marshal(data)
	}
	if dp.protocol == PROTOCOL_WEBSOCKET {
		result = WebSocketPack(result)
	}
	if dp.msgFormat != config.MSG_FORMAT_JSON && dp.msgFormat != config.MSG_FORMAT_PROTOBUF {
		err = fmt.Errorf("msgFormat(%v) can not Pack.", dp.msgFormat)
	}
	return
}

func (dp DataPack) Unpack(data []byte, result protoreflect.ProtoMessage) error {
	if dp.protocol == PROTOCOL_WEBSOCKET {
		data = WebSocketUnpack(data)
	}
	if dp.msgFormat == config.MSG_FORMAT_JSON {
		return json.Unmarshal(data, result)
	}
	if dp.msgFormat == config.MSG_FORMAT_PROTOBUF {
		return proto.Unmarshal(data, result)
	}
	if dp.msgFormat != config.MSG_FORMAT_JSON && dp.msgFormat != config.MSG_FORMAT_PROTOBUF {
		return fmt.Errorf("msgFormat(%v) can not Unpack.", dp.msgFormat)
	}
	return fmt.Errorf("msgFormat(%v) unpack error", dp.msgFormat)
}
