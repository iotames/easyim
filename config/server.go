package config

import (
	"os"
	"strconv"
)

const MSG_FORMAT_JSON = "json"
const MSG_FORMAT_PROTOBUF = "protobuf"

type Server struct {
	IP, MsgFormat   string
	Port, DropAfter int
}

func GetServer() Server {
	ip := os.Getenv("SERVER_IP")
	msgFormat := os.Getenv("MSG_FORMAT")
	portStr := os.Getenv("SERVER_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic("Error: Fail To Get SERVER_PORT," + err.Error())
	}
	dropAfterStr := os.Getenv("SERVER_DROP_AFTER")
	dropAfter, err := strconv.Atoi(dropAfterStr)
	if err != nil {
		panic("Error: Fail To Get SERVER_DROP_AFTER," + err.Error())
	}
	return Server{IP: ip, MsgFormat: msgFormat, Port: port, DropAfter: dropAfter}
}
