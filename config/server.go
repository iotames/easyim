package config

import (
	"os"
	"strconv"
)

const MSG_FORMAT_JSON = "json"
const MSG_FORMAT_PROTOBUF = "protobuf"

type Server struct {
	IP, MsgFormat            string
	Port, ApiPort, DropAfter int
}

func GetServer() Server {
	ip := getEnvDefaultStr("SERVER_IP", "0.0.0.0")
	msgFormat := os.Getenv("MSG_FORMAT")
	port := getEnvDefaultInt("SERVER_PORT", 8888)
	apiPort := getEnvDefaultInt("API_SERVER_PORT", 8889)

	dropAfterStr := os.Getenv("SERVER_DROP_AFTER")
	dropAfter, err := strconv.Atoi(dropAfterStr)
	if err != nil {
		panic("Error: Fail To Get SERVER_DROP_AFTER," + err.Error())
	}
	return Server{IP: ip, MsgFormat: msgFormat, Port: port, ApiPort: apiPort, DropAfter: dropAfter}
}
